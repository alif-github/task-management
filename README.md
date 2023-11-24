
# Project FITHUB API (With JWT Authentication)

The Project Test API is contains some API for task management. This Project use stacks such HTTP Rest API and JWT Authentication also for access main API.


## Tech Stack

**Server:** Golang, JWT Authentication, gin gonic


## Run Locally

#### 1. Local environment
- Make sure postgresql has installed, else you can download https://www.postgresql.org/download/
- Create new connection on db, you can use pgAdmin or access bin folder postgresql. more on https://www.enterprisedb.com/postgres-tutorials/connecting-postgresql-using-psql-and-pgadmin
- Create new folder for save source repository
- Clone the project
```bash
  git clone https://github.com/alif-github/task-management.git
```
- Go to the project directory
```bash
  cd task-management
```
- Edit location config development on file .env to source folder
- Run server, as the default environment will be "local"
```bash
  go run ./app/main.go local
```
- Server will be running on http://localhost:9078

## API Reference

#### 1. Register User (User)

```http
  POST /v1/fithub/auth/register
```
##### Request Body JSON (Example)
```json
{
    "first_name": "Samuel",
    "last_name": "Wijk",
    "username": "samuel15",
    "password": "Samuel123@",
    "email": "samuel@gmail.com",
    "role_id": 2
}
```
##### Response Body JSON (Example)
```json
{
    "request_id": "ba8fc10b-32ff-49b0-b15a-afd20fb11d32",
    "status": true,
    "message": "Success Register!",
    "data": null
}
```

#### 2. Login (User)
```http
  POST /v1/fithub/auth/login
```
##### Request Body JSON (Example)
```json
{
    "username": "alif15",
    "password": "Alif123@"
}
```
##### Response Body JSON (Example)
```json
{
    "success": true,
    "request_id": "0c287f28-0d76-47c9-b065-5ece28b4bf68",
    "message": "Login Berhasil",
    "data": null
}
```

#### 3. Logout (User)
```http
  GET /v1/fithub/auth/logout
```

#### 4. Get Detail Data (User)
```http
  GET /fithub/users/${ID}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### 5. Get List Data (User)
```http
  GET /fithub/users/${ID}
```

| Query Param   | Example Value | Description                       |
| :----------   | :------------ | :-------------------------------- |
| `page`        | `1`           | Pagination queue from 1 to many   |
| `limit`        | `1`           | Pagination queue from 1 to many  |

#### 6. Add Task (Task)

```http
  POST /fithub/tasks
```
##### Request Body JSON (Example)
```json
{
    "title": "Percobaan 2",
    "description": "Percobaan 2",
    "due_date": "2023-11-24T00:00:00Z",
    "status": "New",
    "user_id": 4
}
```
##### Response Body JSON (Example)
```json
{
    "request_id": "a0d9a453-47c1-4d7f-ab36-439c9efc8784",
    "status": true,
    "message": "Success Store Data!",
    "data": null
}
```

#### 7. Get ID Task (Task)

```http
  GET /fithub/tasks/2
```
##### Response Body JSON (Example)
```json
{
    "request_id": "04528a50-82cf-4d6a-9d04-30bb3eaa7261",
    "status": true,
    "message": "Success Get Detail!",
    "data": {
        "id": 2,
        "title": "Percobaan 1",
        "description": "Percobaan 1",
        "status": "New",
        "due_date": "2023-11-24T00:00:00Z",
        "user": "Samuel Wijk",
        "created_user": "Samuel Wijk",
        "created_at": "2023-11-24T07:35:34.041274Z",
        "updated_user": "Samuel Wijk",
        "updated_at": "2023-11-24T07:35:34.041274Z"
    }
}
```
# ---Next you can review the postman collection