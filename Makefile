run:
	docker build -t app --no-cache=true .
	docker run --name postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=tehnodom --network host -p '5432:5432' -d postgres
	docker run -v migrations/ --network host migrate/migrate -path=/migrations/ -database postgres://postgres:postgres@localhost:5432/tehnodom?sslmode=disable up
	docker run --name app --network host -p 8080:8080 -d app

migrate_down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/tehnodom?sslmode=disable" down


.PHONY: run migrate_down