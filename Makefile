backend_all: backend_unit_test backend_integration_test all_done

backend_unit_test:
	cd store-service && go test ./...

backend_integration_test:
	docker compose up -d store-db bank-gateway shipping-gateway
	sleep 5
	cat tearup/store/init.sql | docker exec -i store-db /usr/bin/mysql -u user --password=password --default-character-set=utf8  store
	cd store-service && go test -tags=integration ./...
	# docker-compose down

store_db:
	docker compose up -d store-db 

point_db:
	docker compose up -d point-db 

down:
	docker compose down

all_done:
	echo "All done"