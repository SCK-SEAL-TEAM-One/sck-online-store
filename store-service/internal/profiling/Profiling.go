package profiling

import (
	"log"
	"os"
	"runtime"

	pyroscope "github.com/grafana/pyroscope-go"
)

func StartPyroscope(pyroscopeURL, serviceName string) (*pyroscope.Profiler, error) {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	hostname, _ := os.Hostname()

	profiler, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: serviceName,
		ServerAddress:   pyroscopeURL,
		Tags:            map[string]string{"hostname": hostname},
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	if err != nil {
		return nil, err
	}

	log.Printf("Pyroscope profiler started, sending to %s", pyroscopeURL)
	return profiler, nil
}
