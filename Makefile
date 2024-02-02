
store_db:
	docker compose up store-db -d

point_db:
	docker compose up point-db -d

down:
	docker compose down