 docker compose up --build   -> this will run the postgress image 

 go run cmd/main.go  -> this will run the Health-Sync backend server 

 migrate create -ext sql -dir db/migrations -seq create_users_table    -> thhi will create migration table (optional)


 migrate -database "postgres://admin:adminofheal@localhost:5432/healthsync_db?sslmode=disable" -path db/migrations up  -> this will create table mentions on the migration to the postgres image 

 sqlc generate  -> to generate the sql queries  

 curl -X POST http://localhost:8080/auth/register \
     -H "Content-Type: application/json" \
     -d '{
           "email": "test@example.com",
           "password": "securepassword",
           "name": "Test User"
         }'                                                  -> sample request



docker ps ->   take the postgres container 

docker exec -it ac7435b99763 bash -> enter the terminal of the image using the container id 

psql -U admin -d healthsync_db -> now enter into psql , now query as needed



https://accounts.google.com/o/oauth2/auth?client_id=532745854641-tkcf9fepfa26rbafmo0900k1r108via4.apps.googleusercontent.com&redirect_uri=http://localhost:8080/auth/google/callback&response_type=code&scope=email%20profile

