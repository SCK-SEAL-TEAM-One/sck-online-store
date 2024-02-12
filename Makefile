
backend_start: store_service point_service bank shipping
backend_setup: store_db point_service bank shipping
backend_test_all: backend_unit_test backend_integration_test all_done

backend_unit_test:
	cd store-service && go test ./...

backend_integration_test:
	docker compose up -d store-db-test bank-gateway shipping-gateway
	sleep 10
	cat tearup/store/init.sql | docker exec -i store-db-test /usr/bin/mysql -u user --password=password --default-character-set=utf8  store
	cd store-service && go test -tags=integration ./...
	docker compose down store-db-test 

store_db:
	docker compose up -d store-db 

store_service:
	docker compose up -d store-service --build

point_service:
	docker compose up -d point-service

bank:
	docker compose up -d bank-gateway --build

shipping:
	docker compose up -d shipping-gateway --build

down:
	docker compose down

all_done:
	echo "All done"