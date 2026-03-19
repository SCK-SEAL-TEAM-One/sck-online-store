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
	docker compose up liquibase

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
	cp -f store-web/.env_local store-web/.env
	docker compose up -d thirdparty point-service db store-service store-web nginx seed liquibase --build

start_test_suite_grid:
	cp -f store-web/.env_grid store-web/.env
	docker compose up -d thirdparty point-service db store-service store-web nginx seed liquibase --build
	docker compose up selenium-hub chrome -d

start_test_suite_sck:
	docker compose up selenium-hub chrome -d

stop_test_suite:
	docker compose down

run_robot: URL ?= http://localhost/product/list
run_robot:
	cd atdd/ui \
	&& python3 -m venv .venv \
	&& . .venv/bin/activate \
	&& pip install -r requirements.txt \
	&& robot -v URL:$(URL) -v REMOTE_HUB_URL:$(REMOTE_HUB_URL) -x ./reports/authen.xml ./001-Authentication \
	&& robot -v URL:$(URL) -v REMOTE_HUB_URL:$(REMOTE_HUB_URL) -x ./reports/pdf.xml ./002-Order-Summary-PDF \
	&& deactivate

run_robot_authentication: URL ?= http://localhost/product/list
run_robot_authentication:
	cd atdd/ui \
	&& python3 -m venv .venv \
	&& . .venv/bin/activate \
	&& pip install -r requirements.txt \
	&& robot -v URL:$(URL) -v REMOTE_HUB_URL:${REMOTE_HUB_URL} -x ./reports/authen.xml ./001-Authentication \
	&& deactivate

run_robot_order_summary_pdf: URL ?= http://localhost/product/list
run_robot_order_summary_pdf:
	cd atdd/ui \
	&& python3 -m venv .venv \
	&& . .venv/bin/activate \
	&& pip install -r requirements.txt \
	&& robot -v URL:$(URL) -v REMOTE_HUB_URL:${REMOTE_HUB_URL} -x ./reports/pdf.xml ./002-Order-Summary-PDF \
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

# --- Development workflow: run all tests before commit ---

unit_test_all:
	cd store-service && go test -v ./...
	cd point-service && npm test
	cd store-web && npm run test:component

code_analysis_all: code_analysis_backend code_analysis_frontend

test_all: code_analysis_all unit_test_all start_test_suite run_newman run_robot stop_test_suite
	@echo "All tests passed!"

gen-swagger:
	cd store-service && swag init -g cmd/main.go -o cmd/docs

# --- EKS Build & Deploy ---
# Image tag format: eks-YYMMDD-HHMM (e.g., eks-260319-1045)
EKS_TAG := eks-$(shell date +%y%m%d-%H%M)
DOCKER_REPO := siamchamnankit

eks_build_store:
	docker build --platform linux/amd64 -t $(DOCKER_REPO)/store-service:$(EKS_TAG) store-service/
	@echo "Built $(DOCKER_REPO)/store-service:$(EKS_TAG)"

eks_build_point:
	docker build --platform linux/amd64 -t $(DOCKER_REPO)/point-service:$(EKS_TAG) point-service/
	@echo "Built $(DOCKER_REPO)/point-service:$(EKS_TAG)"

eks_build_all: eks_build_store eks_build_point

eks_push_store: eks_build_store
	docker push $(DOCKER_REPO)/store-service:$(EKS_TAG)

eks_push_point: eks_build_point
	docker push $(DOCKER_REPO)/point-service:$(EKS_TAG)

eks_push_all: eks_push_store eks_push_point

eks_deploy_store: eks_push_store
	sed -i '' 's|image: $(DOCKER_REPO)/store-service:.*|image: $(DOCKER_REPO)/store-service:$(EKS_TAG)|' deploy/k8s/store-service/service.yml
	kubectl apply -f deploy/k8s/store-service/service.yml
	kubectl rollout status deployment/store-service-deployment --timeout=120s
	@echo "Deployed store-service:$(EKS_TAG)"

eks_deploy_point: eks_push_point
	sed -i '' 's|image: $(DOCKER_REPO)/point-service:.*|image: $(DOCKER_REPO)/point-service:$(EKS_TAG)|' deploy/k8s/point-service/service.yml
	kubectl apply -f deploy/k8s/point-service/service.yml
	kubectl rollout status deployment/point-service-deployment --timeout=120s
	@echo "Deployed point-service:$(EKS_TAG)"

eks_deploy_all: eks_deploy_store eks_deploy_point
