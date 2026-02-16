# run all systems
include store-web/.env
export $(shell sed 's/=.*//' 'store-web/.env')


all: backend_start store_web

test:
	echo $(STORE_SERVICE_URL)

pv:
	python3 --version
	robot --help

# install dependency
install_dependency_frontend:
	cd store-web && npm install

install_dependency_backend:
	cd store-service && go mod tidy

# code analysis
code_analysis_frontend:
	cd store-web && npm run lint

code_analysis_backend:
	cd store-service && go vet ./...

# run backend api include arrange system
backend_start: store_service point_service thirdparty

# run all arrange systems of backend
backend_setup: store_db point_service thirdparty

backend_clear_test_cache:
	cd store-service && go clean --testcache

# run all test of backend
backend_test_all: backend_unit_test setup_test_fixtures backend_integration_test all_done

backend_unit_test:
	# cd store-service && go test -v ./...
	cd store-service && go test -v 2>&1 ./... | go-junit-report -set-exit-code > report.xml

setup_test_fixtures:
	docker compose up -d db thirdparty
	sleep 7

backend_integration_test: setup_test_fixtures
	cd store-service && go test -tags=integration ./...
	docker compose down 

store_db:
	docker compose up -d db 

store_service_dev_mode:
	cd ./store-service/cmd && \
	DB_CONNECTION="user:password@tcp(localhost:3306)/store?parseTime=true" \
	POINT_GATEWAY=localhost:8001 \
	BANK_GATEWAY=localhost:8882 \
	SHIPPING_GATEWAY=localhost:8883 \
	JWT_SECRET=my-secret-key \
	go run main.go

point_service:
	docker compose up -d point-service

store_web:
	docker compose up -d store-web --build

thirdparty:
	docker compose up -d thirdparty --build

start_all:
	 docker compose up -d db adminer seed liquibase thirdparty point-service store-service store-web nginx --build

down:
	docker compose down

all_done:
	echo "All done"

build_backend:
	docker compose build store-service

build_frontend:
	docker compose build store-web

build_nginx:
	docker compose build nginx

start_test_suite:
	docker compose up -d thirdparty point-service db store-service store-web nginx seed liquibase --build

stop_test_suite:
	docker compose down

run_robot:
	cd atdd/ui \
	&& python3 -m venv .venv \
	&& source .venv/bin/activate \
	&& pip install -r requirements.txt \
	&& robot -v URL:http://localhost/product/list ./001-Authentication \
	&& robot -v URL:http://localhost/product/list ./002-Order-Summary-PDF \
	&& deactivate

run_robot_authentication:
	cd atdd/ui \
	&& python3 -m venv .venv \
	&& source .venv/bin/activate \
	&& pip install -r requirements.txt \
	&& robot -v URL:http://localhost/product/list ./001-Authentication \
	&& deactivate

run_robot_order_summary_pdf:
	cd atdd/ui \
	&& python3 -m venv .venv \
	&& source .venv/bin/activate \
	&& pip install -r requirements.txt \
	&& robot -v URL:http://localhost/product/list ./002-Order-Summary-PDF \
	&& deactivate

# run_newman: 
# 	cd atdd/api \
# 	&& newman run sck-online-store.postman_collection.json \
# 	 -e sck-online-store.local.postman_environment.json \
# 	 -r cli,junit,htmlextra

run_newman:
	$(MAKE) run_newman_authentication
	$(MAKE) run_newman_order_summary_pdf

run_newman_authentication:
	cd atdd/api \
	&& newman run collections/001-Authentication.postman_collection.json \
	  --folder "TSS-AUTH-001" \
		-e sck-online-store.local.postman_environment.json \
		-d data/001-Authentication/TSS-AUTH-001.json \
		-r cli,junit,htmlextra \
	&& newman run collections/001-Authentication.postman_collection.json \
	  --folder "TSS-AUTH-002" \
		-e sck-online-store.local.postman_environment.json \
		-d data/001-Authentication/TSS-AUTH-002.json \
		-r cli,junit,htmlextra \
	&& newman run collections/001-Authentication.postman_collection.json \
	  --folder "TSS-AUTH-003" \
		-e sck-online-store.local.postman_environment.json \
		-d data/001-Authentication/TSS-AUTH-003.json \
		-r cli,junit,htmlextra \
	&& newman run collections/001-Authentication.postman_collection.json \
	  --folder "TSA-AUTH-001" \
		-e sck-online-store.local.postman_environment.json \
		-d data/001-Authentication/TSA-AUTH-001.json \
		-r cli,junit,htmlextra \
	&& newman run collections/001-Authentication.postman_collection.json \
	  --folder "TSA-AUTH-002" \
		-e sck-online-store.local.postman_environment.json \
		-d data/001-Authentication/TSA-AUTH-002.json \
		-r cli,junit,htmlextra \
	&& newman run collections/001-Authentication.postman_collection.json \
	  --folder "TSA-AUTH-003" \
		-e sck-online-store.local.postman_environment.json \
		-d data/001-Authentication/TSA-AUTH-003.json \
		-r cli,junit,htmlextra

run_newman_order_summary_pdf:
	cd atdd/api \
	&& newman run collections/002-Order-Summary-PDF.postman_collection.json \
	  --folder "TSS-OSP-001" \
		-e sck-online-store.local.postman_environment.json \
		-d data/002-Order-Summary-PDF/TSS-OSP-001.json \
		-r cli,junit,htmlextra \
	&& newman run collections/002-Order-Summary-PDF.postman_collection.json \
	  --folder "TSS-OSP-002" \
		-e sck-online-store.local.postman_environment.json \
		-d data/002-Order-Summary-PDF/TSS-OSP-002.json \
		-r cli,junit,htmlextra

code-coverage:
	cd store-service && go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

gen-swagger:
	cd store-service && swag init -g cmd/main.go -o cmd/docs 
