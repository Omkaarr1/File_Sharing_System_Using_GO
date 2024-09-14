# File Sharing System

## Introduction

The **File Sharing System** is a backend application built with Go that allows users to upload, share, and manage files. It includes features such as user registration and login, file uploading, and file sharing via unique URLs. This system is designed to be secure and scalable, utilizing JWT (JSON Web Token) authentication for user sessions, PostgreSQL as the primary database for storing user and file metadata, and optionally AWS S3 or local storage for storing files.

This project is ideal for learning how to build a scalable file-sharing platform with authentication and database interactions in Go.

---

## Functionality

### 1. **User Registration and Authentication**
- Users can register by providing an email and password.
- Passwords are securely hashed using the `bcrypt` library.
- Users can log in by providing their credentials, which will return a JWT token.
- JWT tokens are used to authenticate further requests to the API.

### 2. **File Upload**
- Authenticated users can upload files to the system.
- The metadata of the files (such as file name, size, and URL) is stored in PostgreSQL.
- Files can be stored locally or on cloud storage like AWS S3.

### 3. **File Sharing**
- Users can share files by generating unique URLs for each uploaded file.
- The file can then be accessed via the shared URL for public download.

### 4. **File Management**
- Users can retrieve the list of files they have uploaded.
- Files can be deleted after a certain period or manually by the user.

---

## How the Project Works

### 1. **Authentication and Authorization**
- Users register by providing an email and password. The password is hashed using `bcrypt` before being stored in the database.
- JWT (JSON Web Tokens) are used to authenticate users after they log in.
- Every request requiring authentication must include a valid JWT token in the header.

### 2. **File Uploading**
- Authenticated users can upload files via a `/upload` endpoint.
- Upon file upload, the file's metadata (file name, size, etc.) is saved in PostgreSQL, and the file is uploaded either locally or to an S3 bucket (depending on configuration).

### 3. **File Sharing**
- Files can be shared using a unique URL. When a user uploads a file, they receive a URL which they can share with others.
- The file can be accessed via this URL for public download.

### 4. **Database Integration**
- PostgreSQL is used as the primary database.
- `pgx/v4` is used for interacting with PostgreSQL to store user and file information.

---

## How to Run the Project

Follow these steps to set up and run the project:

### 1. **Clone the Repository**

First, clone this repository to your local machine:

```bash
git clone https://github.com/yourusername/file-sharing-system.git
cd file-sharing-system
```
## Set Up Environment Variables

You need to set environment variables to configure the database connection, JWT secret, and optionally, AWS S3 for file storage.

### Create a `.env` file in the project root directory:

```bash
touch .env
```

## Add the following environment variables in the .env file:

# Database configuration
DATABASE_URL="postgres://username:password@localhost:5432/file_sharing_db?sslmode=disable"

# JWT Secret Key
JWT_SECRET="your_jwt_secret_key"

# AWS S3 configuration (if applicable)
AWS_ACCESS_KEY_ID="your_aws_access_key_id"
AWS_SECRET_ACCESS_KEY="your_aws_secret_access_key"
AWS_REGION="your_aws_region"
AWS_S3_BUCKET="your_s3_bucket_name"

# Install Dependencies
Ensure Go is installed on your system. You can install Go from here.
Next, install the project dependencies:

``` go
    go mod tidy
```

### This will fetch all the dependencies required for the project.

# Set Up PostgreSQL Database
Make sure PostgreSQL is running on your local machine or cloud provider.
Create a new PostgreSQL database:

``` sql
    CREATE DATABASE file_sharing_db;
```

Run any necessary migrations to set up the database tables (these can be defined manually or with a migration tool).

# **Run the Project**
Once the environment variables and database are set up, you can run the project using:
``` go
    go run main.go
```
The server will start on the default port 8080. You can change this by modifying the server configuration in main.go.

# **Testing the API**
You can use a tool like Postman or curl to interact with the API. Below are some example requests:

User Registration:

``` bash
    curl -X POST http://localhost:8080/register -d '{"email":"test@example.com","password":"password123"}' -H "Content-Type: application/json"
```

User Login:
``` bash
    curl -X POST http://localhost:8080/login -d '{"email":"test@example.com","password":"password123"}' -H "Content-Type: application/json"
```

File Upload (requires JWT token):
``` bash
    curl -X POST http://localhost:8080/upload -H "Authorization: Bearer <JWT_TOKEN>" -F "file=@path/to/your/file.txt"
```

# **Run Tests**
You can run tests to validate the functionality of the project:
``` bash
    go test -v ./...
```
This will run all the unit tests defined in the *_test.go files.

# **Project Structure**

``` bash
file-sharing-system/
├── handlers/          # Contains the HTTP handlers (e.g., register, login, upload, etc.)
│   ├── auth.go
│   ├── file.go
├── models/            # Contains the data models and database interaction code
│   ├── user.go
├── utils/             # Contains utility functions like database connections
│   ├── db.go
├── main.go            # The main entry point for the application
├── .env               # Environment variables file
├── go.mod             # Go module file
├── go.sum             # Go dependencies file
└── README.md          # Project documentation
```

# **Technologies Used**
    Go: The programming language used for the backend.
    PostgreSQL: For user and file metadata storage.
    AWS S3: (Optional) For file storage.
    JWT: For authentication and authorization.
    bcrypt: For password hashing.