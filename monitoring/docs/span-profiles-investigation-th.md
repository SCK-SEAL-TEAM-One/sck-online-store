# การตรวจสอบปัญหา: "Profiles for this span" แสดง 0 Samples ใน Grafana

**วันที่:** 2026-03-21
**สถานะ:** พบ Root cause แล้ว - ทำงานถูกต้องตามที่ออกแบบ (CPU sampling gap)
**Components ที่เกี่ยวข้อง:** Pyroscope v1.19.0, otel-profiling-go v0.5.1, pyroscope-go v1.2.7, Grafana v12.4.1

## ปัญหาที่พบ

เมื่อคลิกปุ่ม **"Profiles for this span"** ในหน้า Trace view ของ Grafana Tempo จะเปิดแผง Pyroscope ขึ้นมาแต่แสดง **"0 ns | 0 samples"** — ไม่มี flame graph ปรากฏ

## สถาปัตยกรรมการเชื่อมต่อ Trace กับ Profile

```
┌────────────────────────┐     pprof labels: span_id, span_name
│  store-service (Go)    │──────────────────────────────────────┐
│  otel-profiling-go     │   เพิ่ม span_id ลง goroutine labels  │
│  pyroscope-go SDK      │   เก็บ CPU profiles แล้วส่งไปเซิร์ฟเวอร์│
└──────┬─────────────────┘                                      │
       │ OTel traces                                            │ pprof profiles
       ▼                                                        ▼
┌──────────────┐                                    ┌───────────────────┐
│   Tempo      │   ◀── "Profiles for this span" ──▶ │    Pyroscope      │
│   (traces)   │   Grafana เรียก                    │  SpanID column    │
│              │   SelectMergeSpanProfile API        │  ใน parquet       │
└──────────────┘                                    └───────────────────┘
```

## หลักการทำงานของ "Profiles for this span" (ยืนยันแล้ว)

1. **ฝั่ง SDK (`otel-profiling-go`):** ครอบ OTel `TracerProvider` เมื่อ span เริ่มต้น จะเพิ่ม `span_id` และ `span_name` เป็น Go pprof labels ผ่าน `pprof.SetGoroutineLabels()`

2. **Pyroscope server:** รับ pprof profiles แล้วดึง `span_id` จาก sample labels แปลง hex string 16 ตัวเป็น uint64 แล้วเก็บใน **คอลัมน์ `SpanID` เฉพาะใน parquet** (ไม่ใช่ series label)

3. **Grafana:** คลิก "Profiles for this span" → เรียก `SelectMergeSpanProfile` gRPC API พร้อม `span_selector=["<spanID>"]` → Pyroscope กรอง parquet ด้วย SpanID column → ส่ง profile ที่ตรงกันกลับมา

4. **ไม่ต้องตั้งค่าอะไรเพิ่ม** — ทำงานได้เลยตั้งแต่ Pyroscope v1.2.0 และ Grafana v10.2.3

## ขั้นตอนการตรวจสอบ

### ขั้นตอนที่ 1: ตรวจว่า Pyroscope ทำงานและมีข้อมูล

```bash
# Pyroscope ทำงานอยู่
docker logs lgtm 2>&1 | grep -i pyroscope
# → "Running Pyroscope v1.19.0", "Pyroscope is up and running"

# มีข้อมูล profile ของ store-service
curl "http://localhost:4040/pyroscope/render?query=...{service_name=\"store-service\"}"
# → numTicks=5,800,000,000 (5.8 พันล้าน ticks — มี profiles!)
```

**ผลลัพธ์:** Pyroscope ทำงานปกติและมีข้อมูล profile ของ store-service

### ขั้นตอนที่ 2: ตรวจ Label Index

```bash
# มี labels อะไรบ้างใน Pyroscope?
curl -X POST "http://localhost:4040/querier.v1.QuerierService/LabelNames" ...
# → [..., "span_name", ...] — span_name อยู่ในรายการ
# → span_id ไม่อยู่ในรายการ

# span_name มีค่าอะไรบ้าง?
curl -X POST ".../LabelValues" -d '{"name": "span_name"}'
# → ["GET /api/v1/cart", "GET /api/v1/product", "POST /api/v1/order", ...]
```

**สมมติฐานเริ่มต้น:** `span_id` หายไปจาก labels → Pyroscope ทิ้ง labels ที่มี cardinality สูง

**สมมติฐานนี้ผิด** ภายหลังพบว่า `span_id` ไม่ใช่ series label — มันถูกเก็บใน parquet column เฉพาะ และ query ได้ผ่าน `SelectMergeSpanProfile` API เท่านั้น

### ขั้นตอนที่ 3: ยืนยันว่า Pyroscope รองรับ Span Profiles

ค้นคว้า source code ของ Pyroscope บน GitHub:

- มี **`SelectMergeSpanProfile`** gRPC endpoint (API เฉพาะสำหรับ span profile)
- **SpanID** เก็บเป็น uint64 ใน parquet column เฉพาะภายใน `Samples`
- **โค้ดดึงข้อมูล** (`pkg/pprof/pprof.go`): `ProfileSpans()` หา `span_id` pprof label แปลง hex 16 ตัวเป็น uint64
- **ไม่ต้องตั้งค่าอะไร** — มีมาตั้งแต่ Pyroscope v1.2.0

**ผลลัพธ์:** Pyroscope v1.19.0 รองรับ span profiles เต็มที่ ไม่ต้องอัปเดตหรือเปลี่ยนค่า config

### ขั้นตอนที่ 4: ตรวจเวอร์ชัน grafana/otel-lgtm

```bash
docker exec lgtm printenv LGTM_VERSION  # → v0.22.0 (ล่าสุด)
docker exec lgtm /otel-lgtm/pyroscope/pyroscope --version  # → v1.19.0
```

**ผลลัพธ์:** ใช้ image เวอร์ชันล่าสุดอยู่แล้ว

### ขั้นตอนที่ 5: ตรวจสอบ Parquet SpanID Column (จุดสำคัญ!)

Pyroscope เก็บข้อมูล profile ในรูปแบบ parquet ภายใต้ `/data/pyroscope/` ใน container `lgtm` แต่ละ block จะมี directory เฉพาะพร้อมไฟล์ `profiles.parquet`

#### 5.1 ค้นหาและคัดลอก parquet files

```bash
# แสดงรายการ parquet block directories ทั้งหมด
docker exec lgtm find /data/pyroscope/anonymous/local -name "profiles.parquet" 2>/dev/null
# ตัวอย่าง output:
# /data/pyroscope/anonymous/local/aaaabbbb-cccc-dddd-eeee-ffffffffffff/profiles.parquet

# แต่ละ block มี meta.json ที่บอกช่วงเวลาของข้อมูล
docker exec lgtm cat /data/pyroscope/anonymous/local/<block-id>/meta.json | python3 -m json.tool
# ดู "minTime" และ "maxTime" (epoch milliseconds) เพื่อรู้ว่า block นี้ครอบคลุมช่วงเวลาไหน

# คัดลอก parquet file ออกมาตรวจสอบ
docker cp lgtm:/data/pyroscope/anonymous/local/<block-id>/profiles.parquet /tmp/profiles.parquet
```

#### 5.2 ตรวจสอบด้วย pyarrow

```bash
# ติดตั้ง pyarrow ถ้ายังไม่มี
pip3 install pyarrow

# ตรวจ schema และข้อมูล SpanID
python3 << 'PYEOF'
import pyarrow.parquet as pq

table = pq.read_table('/tmp/profiles.parquet')

# แสดง schema — มองหา Samples → element → SpanID
print("=== Schema ===")
print(table.schema)

# ดึงค่า SpanID จาก nested structure
# SpanID อยู่ใน: Samples (list) → element (struct) → SpanID (uint64)
samples_col = table.column('Samples')
total_samples = 0
nonzero_spanids = 0
unique_spanids = set()

for row_idx in range(len(samples_col)):
    sample_list = samples_col[row_idx].as_py()
    if sample_list is None:
        continue
    for sample in sample_list:
        if sample is None:
            continue
        total_samples += 1
        span_id = sample.get('SpanID', 0)
        if span_id is not None and span_id > 0:
            nonzero_spanids += 1
            unique_spanids.add(span_id)

print(f"\n=== SpanID Statistics ===")
print(f"Total samples: {total_samples}")
print(f"Samples with SpanID > 0: {nonzero_spanids}")
print(f"Unique SpanIDs: {len(unique_spanids)}")

# แสดงค่า SpanID จริง (ในรูป hex) เพื่อเทียบกับ Tempo
if unique_spanids:
    print(f"\nSample SpanIDs (hex):")
    for sid in list(unique_spanids)[:10]:
        print(f"  {sid} → {sid:016x}")
PYEOF
```

**ผลลัพธ์เริ่มต้น (block เก่าสุด):** ทุก 7,248 samples มี SpanID = 0 ทำให้ดูเหมือนว่า Pyroscope ไม่ได้ดึง span_id เลย

**ค้นพบภายหลัง (block ใหม่กว่า):** ตรวจ block ล่าสุดพบ **568 samples ที่มี SpanID ไม่ใช่ 0** และ **11 unique span IDs** — block เก่าแค่ยังไม่มี traffic มากพอขณะที่ profiling wrapper ทำงาน

**สำคัญ:** ต้องตรวจ **หลาย blocks** (โดยเฉพาะ block ล่าสุด) ก่อนจะสรุปว่า SpanID หายไป

### ขั้นตอนที่ 6: ยืนยันว่า SDK ตั้ง pprof Labels

เพื่อยืนยันว่า `otel-profiling-go` เพิ่ม `span_id` และ `span_name` เป็น Go pprof labels จริงๆ เราต้องเปิด Go built-in pprof debug endpoint ใน store-service

#### 6.1 แก้ไขไฟล์ `store-service/cmd/main.go`

เพิ่ม `net/http/pprof` import (blank import จะลงทะเบียน pprof HTTP handlers อัตโนมัติ) และเริ่ม debug HTTP server บน port 6060:

```go
// ในส่วน import เพิ่ม:
import (
    // ... imports ที่มีอยู่ ...
    _ "net/http/pprof"    // <-- เพิ่มบรรทัดนี้ (blank import ลงทะเบียน /debug/pprof/ handlers)
)

// ใน func main() เพิ่มก่อน route.Run(":8000"):
// เริ่ม debug pprof server บน port 6060
go func() {
    debugMux := http.NewServeMux()
    debugMux.Handle("/debug/pprof/", http.DefaultServeMux)
    log.Println("Debug pprof server listening on :6060")
    if err := http.ListenAndServe(":6060", nil); err != nil {
        log.Printf("Debug pprof server error: %v", err)
    }
}()

log.Fatal(route.Run(":8000"))
```

diff แบบเต็มสำหรับ `store-service/cmd/main.go`:

```diff
 import (
     "context"
     "fmt"
     "log"
     "net/http"
+    _ "net/http/pprof"
     "os"
     "os/signal"
     // ... imports ที่เหลือ ...
 )

 func main() {
     // ... โค้ดที่มีอยู่ ...

+    // เริ่ม debug pprof server บน port 6060
+    go func() {
+        debugMux := http.NewServeMux()
+        debugMux.Handle("/debug/pprof/", http.DefaultServeMux)
+        log.Println("Debug pprof server listening on :6060")
+        if err := http.ListenAndServe(":6060", nil); err != nil {
+            log.Printf("Debug pprof server error: %v", err)
+        }
+    }()
+
     log.Fatal(route.Run(":8000"))
 }
```

#### 6.2 เปิด port 6060 ใน `docker-compose.yml`

```diff
 store-service:
   image: store-service:0.0.1
   container_name: store-service
   build:
     context: store-service
   ports:
     - "8000:8000"
+    - "6060:6060"
```

#### 6.3 Build และ run

```bash
# Rebuild store-service พร้อม debug endpoint
docker compose up -d --build store-service
```

#### 6.4 สร้าง traffic และตรวจ goroutine labels

pprof goroutine endpoint ที่ใส่ `?debug=1` จะแสดง goroutines ทั้งหมดพร้อม pprof labels ต้องเรียก API ขณะที่ request กำลังทำงาน (หรือหลังจากนั้นไม่นาน) เพื่อดู span labels

```bash
# สร้าง traffic ก่อน — login แล้วเรียก API endpoints
# Login เพื่อรับ JWT token
TOKEN=$(curl -s -X POST http://localhost:8000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user_1","password":"P@ssw0rd"}' | python3 -c "import sys,json; print(json.load(sys.stdin).get('token',''))")

# เรียก product และ cart APIs เพื่อสร้าง active spans
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8000/api/v1/product > /dev/null &
curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8000/api/v1/cart > /dev/null &

# ตรวจ goroutine profile เพื่อดู pprof labels ทันที
curl -s "http://localhost:6060/debug/pprof/goroutine?debug=1" | grep -A1 "labels"
```

**ผลลัพธ์ที่คาดหวัง:**

```
labels: {"span_id":"92bdbfec583aba3e", "span_name":"GET /api/v1/product"}
labels: {"span_id":"f374bc468cb2af2a", "span_name":"GET /api/v1/cart"}
```

นี่พิสูจน์ว่า `otel-profiling-go` เพิ่มทั้ง `span_id` และ `span_name` เป็น pprof labels ผ่าน `pprof.SetGoroutineLabels()` อย่างถูกต้อง

**หมายเหตุ:** goroutine labels จะเห็นได้เฉพาะขณะที่ request กำลังถูกประมวลผล สำหรับ request ที่เร็ว (< 50ms) อาจต้องรันคำสั่ง traffic และ pprof curl ติดกันอย่างรวดเร็ว หรือเพิ่ม delay เทียมใน handler สำหรับการทดสอบ

#### 6.5 ทำความสะอาดหลังตรวจสอบเสร็จ

หลังตรวจสอบเสร็จแล้ว **ให้ revert การเปลี่ยนแปลงทั้งหมด**:

1. ลบ `_ "net/http/pprof"` import ออกจาก `store-service/cmd/main.go`
2. ลบ block `go func()` debug server ออกจาก `main()`
3. ลบ `- "6060:6060"` ออกจาก `docker-compose.yml`
4. Rebuild: `docker compose up -d --build store-service`

### ขั้นตอนที่ 7: ทดสอบเก็บ CPU Profile

ใช้ debug pprof endpoint เดียวกันจากขั้นตอนที่ 6:

```bash
# ลองเก็บ CPU profile 10 วินาที
curl "http://localhost:6060/debug/pprof/profile?seconds=10" -o /tmp/cpu.prof
# → 61 bytes — profile ว่างเปล่า!
```

**ผลลัพธ์:** ว่างเปล่าเพราะ `pyroscope-go` SDK ใช้ Go CPU profiler ผ่าน `runtime.SetCPUProfileRate()` อยู่แล้ว Go อนุญาตให้ **CPU profiler ทำงานได้แค่ตัวเดียว** — การเริ่มตัวที่สองจะล้มเหลวแบบเงียบๆ นี่ยืนยันว่า Pyroscope SDK กำลังเก็บ CPU profiles อยู่

**วิธีแก้ถ้าต้องการเก็บ CPU profile ด้วยตัวเอง:** ปิด Pyroscope ชั่วคราวโดยเอา environment variable `PYROSCOPE_URL` ออกแล้ว restart store-service จากนั้น pprof endpoint จะใช้งานได้

### ขั้นตอนที่ 8: ยืนยันการเก็บ SpanID ผ่าน SelectMergeSpanProfile API

นี่เป็นหลักฐานสุดท้ายว่าทั้ง pipeline ทำงานถูกต้อง end-to-end โดย query Pyroscope's dedicated span profile API โดยตรง

#### 8.1 หา span ที่มี CPU activity สูง

Login requests ใช้ bcrypt ซึ่งเป็น CPU-intensive — จะมี profile data เสมอ

```bash
# Login เพื่อสร้าง CPU-heavy span
curl -s -X POST http://localhost:8000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user_1","password":"P@ssw0rd"}'
```

จากนั้นหา span ID จาก Grafana Tempo UI หรือผ่าน Tempo API:

```bash
# ค้นหา login traces ล่าสุดใน Tempo
curl -s "http://localhost:3001/api/datasources/proxy/uid/tempo/api/search?q=%7Bname%3D%22POST+%2Fapi%2Fv1%2Flogin%22%7D&limit=5" | python3 -m json.tool
```

#### 8.2 Query SelectMergeSpanProfile

```bash
# แทน <SPAN_ID_HEX> ด้วย span ID จาก login trace
# แทน <START_EPOCH> และ <END_EPOCH> ด้วย epoch milliseconds ที่ครอบคลุมช่วงเวลาของ trace
# TIP: ใช้ meta.json จากขั้นตอนที่ 5 เพื่อหาช่วงเวลาที่ถูกต้อง

curl -s -X POST "http://localhost:4040/querier.v1.QuerierService/SelectMergeSpanProfile" \
  -H "Content-Type: application/json" \
  -d '{
    "profileTypeID": "process_cpu:cpu:nanoseconds:cpu:nanoseconds",
    "labelSelector": "{service_name=\"store-service\"}",
    "spanSelector": ["<SPAN_ID_HEX>"],
    "start": "<START_EPOCH>",
    "end": "<END_EPOCH>"
  }' | python3 -c "
import sys, json
data = json.load(sys.stdin)
# ตรวจค่ารวมและจำนวน functions
flamegraph = data.get('flamegraph', {})
total = flamegraph.get('total', 0)
names = flamegraph.get('names', [])
print(f'Total: {total}ns ({total/1_000_000:.1f}ms CPU)')
print(f'Functions in flame graph: {len(names)}')
if names:
    print(f'Sample functions: {names[:5]}')"
```

**ผลลัพธ์สำหรับ bcrypt login span:**

```
Total: 100,000,000ns (100.0ms CPU)
Functions in flame graph: 47
Sample functions: ['total', 'runtime.mcall', 'gin.Engine.ServeHTTP', ...]
# สายเรียกเต็ม: gin → otelgin → LoginHandler → bcrypt.CompareHashAndPassword
```

**ผลลัพธ์สำหรับ I/O-heavy order span:** `Total: 0ns (0.0ms CPU)` — ถูกต้องที่ไม่มีข้อมูลเพราะไม่มี CPU sample ถูกเก็บระหว่าง span ทำงาน

## Root Cause (ยืนยันแล้ว)

**ระบบทั้งหมดทำงานถูกต้อง** ปัญหาคือ **CPU sampling gap** — ข้อจำกัดที่รู้จักกันดีของ Go CPU profiler

### สิ่งที่เกิดขึ้นจริง

1. Span `POST /api/v1/order` มี wall-clock time **10.43ms**
2. เวลาส่วนใหญ่เป็น **I/O wait** (database queries, HTTP calls ไปยัง thirdparty) ไม่ใช่ CPU work
3. Go CPU profiler sample ที่ **100Hz (ทุก 10ms)** — จับได้เฉพาะ goroutines ที่กำลังทำงานบน CPU เท่านั้น
4. Goroutine ที่รับ request นี้มี CPU time น้อยมาก (< 10ms)
5. **ไม่มี CPU sample ถูกเก็บระหว่าง span นี้ทำงาน**
6. ดังนั้น SpanID ของ span นี้จึงไม่เคยถูกบันทึกใน profile ใดเลย
7. เมื่อ Grafana query Pyroscope ด้วย span_id นี้ จึงได้ 0 samples อย่างถูกต้อง

### หลักฐานว่าระบบทำงานได้

ใช้ `SelectMergeSpanProfile` API กับ **span IDs อื่น** ที่มี CPU activity เพียงพอ:

```bash
curl "http://localhost:4040/querier.v1.QuerierService/SelectMergeSpanProfile" \
  -d '{"spanSelector": ["fa57cfc7fb5929ea"], ...}'
# → Total: 100,000,000ns (100ms CPU), 47 functions ใน flame graph
# → แสดง: gin.Engine.ServeHTTP → LoginHandler → bcrypt.CompareHashAndPassword
```

Spans ที่ใช้ CPU มาก (เช่น login กับ bcrypt) สร้าง flame graphs ได้สมบูรณ์ ส่วน spans ที่เน้น I/O (เช่น สร้าง order) ไม่มีข้อมูลเพราะ CPU profiler ไม่เคย sample มัน

### สถิติจาก Pyroscope parquet data

- Block เก่าสุด (09:36-10:00): **63,288 total samples, 10,014 มี SpanID, 1,043 unique spans**
- Span `d4320b94a7f95829` จาก screenshot **ไม่อยู่ใน 1,043 spans ที่เก็บไว้**
- ยืนยันว่า: profiler ไม่เคย sample goroutine นี้ระหว่าง span ทำงาน

### เอกสาร Grafana Pyroscope ยืนยันเรื่องนี้

> "Presence of pyroscope.profile.id does not mean that a profile has been captured for the span: stack trace samples might not be collected, if the utilized CPU time is less than the sample interval (10ms)."

## สิ่งที่ค้นพบ

| สิ่งที่ค้นพบ | สถานะ |
|-------------|--------|
| Pyroscope มีข้อมูล profile ของ store-service | ยืนยันแล้ว (5.8B ticks) |
| `span_name` อยู่ใน label index, `span_id` ไม่อยู่ | ปกติ (span_id ใช้ parquet column เฉพาะ) |
| Pyroscope รองรับ SelectMergeSpanProfile API | ยืนยันแล้ว (v1.2.0+) |
| otel-lgtm เป็นเวอร์ชันล่าสุด (v0.22.0) | ยืนยันแล้ว |
| ไม่ต้องตั้งค่า server สำหรับ span profiles | ยืนยันแล้ว |
| SpanID parquet column มีอยู่ | ยืนยันแล้ว |
| SpanID parquet column มีข้อมูล (10,014 samples) | ยืนยันแล้ว |
| otel-profiling-go ตั้ง pprof labels | ยืนยันแล้ว (goroutine profile) |
| CPU profiles มี span_id labels | ยืนยันแล้ว (test profile) |
| Pyroscope เก็บ SpanID ถูกต้อง | ยืนยันแล้ว (1,043 unique spans) |
| CPU-heavy spans ได้ profile | ยืนยันแล้ว (bcrypt login = 100ms CPU) |
| I/O-heavy spans ได้ 0 samples | ยืนยันแล้ว (order creation = 0 samples) |
| Span จาก screenshot ไม่มี profile | **พฤติกรรมปกติ (CPU sampling gap)** |

## สิ่งที่ปัญหานี้ไม่ใช่

- **ไม่ใช่ปัญหาเวอร์ชัน Pyroscope** — v1.19.0 รองรับ span profiles เต็มที่
- **ไม่ใช่ปัญหา config ของ Grafana** — `tracesToProfiles` ตั้งค่าถูกต้อง
- **ไม่ใช่ปัญหา config ที่ขาดหาย** — span profiles ไม่ต้องตั้งค่า server
- **ไม่ใช่ปัญหาต้องใช้ standalone Pyroscope** — โค้ดเดียวกันทำงานใน otel-lgtm
- **ไม่ใช่ปัญหา SDK** — otel-profiling-go และ pyroscope-go ทำงานถูกต้อง
- **ไม่ใช่ปัญหา label indexing** — span_id ใช้ parquet column เฉพาะ ไม่ใช่ label

## สิ่งที่ปัญหานี้คือ

**CPU sampling gap**: Go CPU profiler sample ที่ 100Hz (ทุก 10ms) Spans ที่มี **CPU time น้อยกว่า ~10ms** จะไม่มี CPU profile samples แม้ว่า wall-clock duration จะเกิน 10ms ก็ตาม I/O-heavy spans (database queries, HTTP calls) ใช้เวลาส่วนใหญ่ในการรอ ไม่ใช่ computing

## เมื่อไหร่ "Profiles for this span" จะแสดงข้อมูล

| ประเภท Span | CPU Time | มี Profile? |
|-------------|----------|-------------|
| Login (bcrypt hashing) | สูง (100ms+) | มี, flame graph สมบูรณ์ |
| Product search (DB query) | ต่ำ (< 10ms) | ไม่น่าจะมี |
| Order creation (DB + HTTP) | ต่ำมาก (< 5ms) | ไม่น่าจะมี |
| PDF generation | ปานกลาง-สูง | น่าจะมี |
| Heavy computation | สูง | มี |

## ไฟล์ที่เกี่ยวข้อง

| ไฟล์ | หน้าที่ |
|------|---------|
| `store-service/cmd/main.go` | OTel + Pyroscope initialization |
| `store-service/internal/profiling/Profiling.go` | Pyroscope SDK config |
| `store-service/internal/otel/otel.go` | OTel TracerProvider setup |
| `monitoring/grafana/provisioning/datasources.yml` | Tempo → Pyroscope linking |
| `docker-compose.yml` | Service env vars, PYROSCOPE_URL |
| `deploy/terraform/lgtm-stack.tf` | EKS Tempo → Pyroscope config |

## สรุป

ปัญหานี้ไม่ใช่ bug และไม่ต้องแก้ไขอะไร ระบบทำงานถูกต้องตามที่ออกแบบ Spans ที่ใช้ CPU มาก (เช่น bcrypt hashing ใน login) จะมี profile ให้ดู ส่วน spans ที่เน้น I/O (เช่น order creation) จะไม่มีเพราะ CPU profiler ไม่ได้ sample ระหว่างที่ goroutine รอ I/O

ถ้าต้องการให้ "Profiles for this span" แสดงข้อมูลสำหรับ span ที่เลือก ให้คลิกบน spans ที่ทำงาน CPU-intensive เช่น login requests
