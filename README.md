# GORM Relationships with Gin Framework

This project provides a REST API for managing users, posts, and tags with functionalities such as creating users, associating posts with users, and associating tags with posts.

## Table of Contents

1. [Technologies Used](#technologies-used)
2. [Overview](#overview)
3. [Setup Instructions](#setup-instructions)
4. [Access Swagger UI](#access-swagger-ui)
5. [Using Swagger UI](#using-swagger-ui)

## Technologies Used
![Go](https://img.shields.io/badge/Language-Go-blue) ![Gin](https://img.shields.io/badge/Framework-Gin-brightgreen) ![GORM](https://img.shields.io/badge/ORM-GORM-blue) ![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue) ![Swagger](https://img.shields.io/badge/API-Swagger-orange) ![Docker](https://img.shields.io/badge/Docker-Enabled-blue) ![Docker Compose](https://img.shields.io/badge/Docker%20Compose-Used-blueviolet)

## Overview

This project demonstrates the use of GORM to establish relationships between database models and exposes RESTful APIs for interacting with these relationships. Specifically, it covers:

1. **One-to-Many Relationship**: Between `User` and `Post`. A user can create multiple posts.

2. **Many-to-Many Relationship**: Between `Post` and `Tag`. A post can have multiple tags, and a tag can belong to multiple posts.

## Setup Instructions
**Step 1: Clone the Repository**

Clone the repository to your local machine:

```bash
git clone https://github.com/Fahad-I-Khan/go-gin-gorm.git
cd go-gin-gorm
```

**Step 2: Install Go Dependencies**

Run the following commands to install the necessary Go dependencies for the project:

- Install Gin (the web framework used in this project):

```bash
go get github.com/gin-gonic/gin
```
- Install PostgreSQL driver for Go:

```bash
go get github.com/lib/pq
```
- Install Swagger dependencies for API documentation generation:
  - Install Swag CLI (this command-line tool generates Swagger docs from Go comments):

    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    ```
  - Install Swagger UI for Gin:

    ```bash
    go get github.com/swaggo/gin-swagger
    go get github.com/swaggo/files
    ```
- Install CORS middleware for Gin (to handle cross-origin requests):

```bash
go get github.com/gin-contrib/cors
```
- Optional: Run `go mod tidy` to clean up the `go.mod` file and download any missing dependencies:

```bash
go mod tidy
```
- Generate Swagger documentation:

```bash
swag init
```

This command scans your Go code and generates the Swagger documentation in the `docs` folder. You need to run `swag init` every time you modify your API routes or comments.

## Stopping the Application
To stop the application, use the following command:

```bash
docker-compose down
```

This stops and removes the containers, but the data in the PostgreSQL container persists due to the pgdata volume.

## Access Swagger UI
Once the server is running, you can access the Swagger UI at:

```bash
http://localhost:8000/swagger/index.html
```

## Using Swagger UI

Swagger UI allows you to test the API endpoints. Each endpoint will display a default example value in the "Try it out" section. Replace the default example values with your request JSON when sending API requests.

**Example Instructions**

### 1. Create a User

**Endpoint:** `POST /api/v1/users`

**Example Request:**

Replace the default JSON with:

```json
{
  "name": "Alice",
  "email": "alice@example.com"
}
```

### 2. Create a Post for a User

**Endpoint:** `POST /api/v1/posts`

**Example Request:**

Replace the default JSON with:

```json
{
  "title": "First Post",
  "body": "This is my first post content.",
  "user_id": 1
}
```

### 3. Add a Tag to a Post

**Step 1: Create a Tag**

**Endpoint:** `POST /api/v1/tags`

**Example Request:**

Replace the default JSON with:

```json
`{
  "name": "GoLang"
}
```

**Step 2: Associate the Tag with a Post**

**Endpoint:** `POST /api/v1/posts/{postID}/tags/{tagID}`

**Example Request:**

Use the post ID and tag ID created earlier.

```
POST /api/v1/posts/1/tags/1
```

**Example Response:**

```json
{
  "message": "Tag added to post"
}
```

## Notes

- Always replace the default values in Swagger UI with your actual request data.

- Use GET endpoints to verify data after performing POST requests.

- For debugging invalid inputs, check field names and the request structure.