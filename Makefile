migrate-up:
	docker exec -it go-server migrate -path db/migrations -database "postgres://admin:adminofheal@postgres:5432/healthsync_db?sslmode=disable" up

migrate-down:
	docker exec -it go-server migrate -path db/migrations -database "postgres://admin:adminofheal@postgres:5432/healthsync_db?sslmode=disable" down

migrate-new:
	@read -p "Enter migration name: " name; \
	docker exec -it go-server migrate create -ext sql -dir db/migrations -seq $$name
