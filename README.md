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


-----------------------------------------------------------------------------------------------------
Requests 

1.patient register request :

curl -X POST -H "Content-Type: application/json" -d '{
    "name": "Test User",
    "email": "testuser1@example.com",
    "password": "password123",
    "role": "user"
}' http://localhost:8080/auth/register/user

patient register response :

{
  "message": "User created successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYjM3ZDAyOTYtMjgzYy00MjkzLTkwYzgtMjZhNDliNWQwNzFhIiwicm9sZSI6InVzZXIiLCJlbWFpbCI6InRlc3R1c2VyMUBleGFtcGxlLmNvbSIsImV4cCI6MTc0MTU0MTI3NX0.3nWLXXutxH2q7PaDXatyUKtnHkDoQolTao9Y8wKI9ro",
  "user": {
    "email": "testuser1@example.com",
    "id": "b37d0296-283c-4293-90c8-26a49b5d071a",
    "name": "Test User"
  }
}

----------------------------------------------------------------------------------------------------------------------------------------

doctor register request :

curl -X POST -H "Content-Type: application/json" -d '{
    "name": "Test Doctor",
    "email": "testdoctor1@example.com",
    "password": "password123",
    "specialization": "Cardiology",
    "experience": 10,
    "qualification": "MD",
    "hospital_name": "General Hospital",
    "consultation_fee": "100",
    "role": "doctor"
}' http://localhost:8080/auth/register/doctor


doctor register response :

{
  "doctor": {
    "consultation_fee": 100,
    "email": "testdoctor1@example.com",
    "experience": 10,
    "hospital_name": "General Hospital",
    "id": "224f3c9b-2b3c-4f75-9453-d2264015fdff",
    "name": "Test Doctor",
    "qualification": "MD",
    "specialization": "Cardiology"
  },
  "message": "Doctor created successfully",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMjI0ZjNjOWItMmIzYy00Zjc1LTk0NTMtZDIyNjQwMTVmZGZmIiwicm9sZSI6ImRvY3RvciIsImVtYWlsIjoidGVzdGRvY3RvcjFAZXhhbXBsZS5jb20iLCJleHAiOjE3NDE1NDEzNzl9.ScgvbCFh1HtsxRVlJ-ghphVgLeiGcdrSQ6TiVYUH1s8"
}

----------------------------------------------------------------------------------------------------------------------------------------

user login request :

curl -X POST -H "Content-Type: application/json" -d '{
    "email": "testuser1@example.com",
    "password": "password123",
    "role": "user"
}' http://localhost:8080/auth/login/user

user login response :

{
  "message": "User login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYjM3ZDAyOTYtMjgzYy00MjkzLTkwYzgtMjZhNDliNWQwNzFhIiwicm9sZSI6InVzZXIiLCJlbWFpbCI6InRlc3R1c2VyMUBleGFtcGxlLmNvbSIsImV4cCI6MTc0MTU0MTUwM30.KvjN2jGvTjcRpGITfZ_9YDv8RqGrtyEoiPsGrGH4dP0",
  "user": {
    "email": "testuser1@example.com",
    "id": "b37d0296-283c-4293-90c8-26a49b5d071a",
    "name": "Test User"
  }
}

----------------------------------------------------------------------------------------------------------------------------------------

doctor login request :

curl -X POST -H "Content-Type: application/json" -d '{
    "email": "testdoctor1@example.com",
    "password": "password123"
}' http://localhost:8080/auth/login/doctor


doctor login response :

{
  "doctor": {
    "consultation_fee": 100,
    "email": "testdoctor1@example.com",
    "experience": 10,
    "hospital_name": "General Hospital",
    "id": "224f3c9b-2b3c-4f75-9453-d2264015fdff",
    "name": "Test Doctor",
    "qualification": "MD",
    "specialization": "Cardiology"
  },
  "message": "Doctor login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMjI0ZjNjOWItMmIzYy00Zjc1LTk0NTMtZDIyNjQwMTVmZGZmIiwicm9sZSI6ImRvY3RvciIsImVtYWlsIjoidGVzdGRvY3RvcjFAZXhhbXBsZS5jb20iLCJleHAiOjE3NDE1NDE2MTR9.Ep_gu-eXNrXDijLWPvEcJxM5S-r-Vi3URl15-kjDTmk"
}


----------------------------------------------------------------------------------------------------------------------------------------

OAuth login request :
https://accounts.google.com/o/oauth2/auth?client_id=532745854641-tkcf9fepfa26rbafmo0900k1r108via4.apps.googleusercontent.com&redirect_uri=http://localhost:8080/auth/google/callback&response_type=code&scope=email%20profile


OAuth new user response :
 
{
  "message": "new user login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTRmODliMzQtM2E2OS00MTdhLWIyZWYtYjgxMjlmN2ZhYzgzIiwicm9sZSI6InVzZXIiLCJlbWFpbCI6ImNob2NvYm95ODMyNUBnbWFpbC5jb20iLCJleHAiOjE3NDE1NDE3ODB9.rsaAj8Xt713KX7T8yb6RT3Gyt-cqO8kj-K8U7SPw_KY",
  "user": {
    "email": "chocoboy8325@gmail.com",
    "id": "94f89b34-3a69-417a-b2ef-b8129f7fac83",
    "name": "Sriram Janardhanan"
  }
}


OAuth existing user login response :

{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOTRmODliMzQtM2E2OS00MTdhLWIyZWYtYjgxMjlmN2ZhYzgzIiwicm9sZSI6InVzZXIiLCJlbWFpbCI6ImNob2NvYm95ODMyNUBnbWFpbC5jb20iLCJleHAiOjE3NDE1NDE4MTh9.F5o_IGakVGjl0SfXh2AgBp2YoW3X69LGiM7gFYNlxcw",
  "user": {
    "email": "chocoboy8325@gmail.com",
    "id": "94f89b34-3a69-417a-b2ef-b8129f7fac83",
    "name": "Sriram Janardhanan"
  }
}

----------------------------------------------------------------------------------------------------------------------------------------

create doctor availability request : (use the token and doctorid ressponse generated by register or login request )

curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMjI0ZjNjOWItMmIzYy00Zjc1LTk0NTMtZDIyNjQwMTVmZGZmIiwicm9sZSI6ImRvY3RvciIsImVtYWlsIjoidGVzdGRvY3RvcjFAZXhhbXBsZS5jb20iLCJleHAiOjE3NDE1NDE2MTR9.Ep_gu-eXNrXDijLWPvEcJxM5S-r-Vi3URl15-kjDTmk" -d '{
    "start_time": "2024-03-08T10:00:00Z",
    "end_time": "2024-03-08T11:00:00Z"
}' http://localhost:8080/doctors/224f3c9b-2b3c-4f75-9453-d2264015fdff/availability