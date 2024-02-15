# run all systems
all: backend_start store_web

# run backend api include arrange system
backend_start: store_service point_service bank shipping

# run all arrange systems of backend
backend_setup: store_db point_service bank shipping

# run all test of backend
backend_test_all: backend_unit_test backend_integration_test all_done

backend_unit_test:
	cd store-service && go test ./...

backend_integration_test:
	docker compose up -d store-db-test bank-gateway shipping-gateway
	sleep 20
	cat tearup/store/init.sql | docker exec -i store-db-test /usr/bin/mysql -u user --password=password --default-character-set=utf8  store
	cd store-service && go test -tags=integration ./...
	docker compose down store-db-test 

store_db:
	docker compose up -d store-db 

store_service:
	docker compose up -d store-service --build

store_service_dev_mode:
	cd ./store-service/cmd && DBCONNECTION=user:password@\(localhost:3306\)/store POINT_GATEWAY=localhost:8001 BANK_GATEWAY=localhost:8882 SHIPPING_GATEWAY=localhost:8883 go run main.go

point_service:
	docker compose up -d point-service

store_web:
	docker compose up -d store-web

bank:
	docker compose up -d bank-gateway --build

shipping:
	docker compose up -d shipping-gateway --build

down:
	docker compose down

all_done:
	echo "All done"