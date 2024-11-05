# Social API

## Overview
The **Social API** provides a set of functionalities  for adding blog posts and engaging in conversations around them. This API enables user authentication (registration and login), account activation with confirmation emails, blog post creation and management, following other users, and adding comments. It is designed to be a backend system that integrates with any frontend framework.

## Features
- **Authentication & Authorization**
  - User registration and login
  - JWT-based authentication
  - Email verification
- **Content Management**
  - Blog post CRUD operations
  - Comment system
  - User following system
  - Feed system
- **Performance**
  - Redis caching for user profiles
  - Rate limiting
  - Connection pooling
- **Security**
  - Basic authentication for admin routes
  - Role based authorization for post management
  - CORS configuration
  - Environment-based configuration
## API Documentation
API documentation is accessible on Swagger when the application is running at `http://<host>:8080/v1/swagger/index.html`, where `<host>` is the server’s hostname.

## Installation
Follow these steps to set up the Social platform API locally:
### Prerequisites
Ensure that the following tools are installed on your machine:
- **Go** (Version 1.22.5 or later)
- **Git** for cloning the repository
- **PostgreSQL** (or use Docker for database setup)
- **Redis** for user profile caching (optional)
### Cloning the Repository

Clone the repository to your local machine and navigate to the project directory:

```bash
git clone https://github.com/CP-Payne/social.git
cd social
```
### Installing Dependencies
The project uses `go.mod` for managing dependencies. Install the required dependencies:

```bash
go mod download
```
This will fetch and install all the necessary Go packages and modules.

### Environment Setup
Configure environment variables before running the API. You can manually set the variables or, if you have direnv installed, create a .envrc file in the root directory with the following values:

```plaintext
# Server Configuration
export ADDR=":8080"
export EXTERNAL_URL="localhost:8080"
export FRONTEND_URL="http://localhost:5173"
export ENV="development"

# Database Configuration
export DB_ADDR="postgres://postgres:postgres@localhost:5432/socialnetwork?sslmode=disable"
export DB_MAX_OPEN_CONNS=30
export DB_MAX_IDLE_CONNS=30
export DB_MAX_IDLE_TIME="15m"

# Email Configuration
export FROM_EMAIL=""
export MAILTRAP_USERNAME="api"
export MAILTRAP_PASSWORD=""
export SENDGRID_API_KEY=""

# Redis Configuration
export REDIS_ADDR="localhost:6379"
export REDIS_PW=""
export REDIS_DB=0
export REDIS_ENABLED=false

# Security Configuration
export RATELIMITER_REQUESTS_COUNT=20
export RATE_LIMITER_ENABLED=true
export AUTH_BASIC_USER="admin"
export AUTH_BASIC_PASS="admin"
export AUTH_TOKEN_SECRET="secret-example"
export CORS_ALLOWED_ORIGIN="http://localhost:5173"
```
> Note: The environment variables listed above include default values that will be used if no environment variables are set. Either Mailtrap or SendGrid can be used for SMTP account activation emails.

### Database Setup

#### Using Docker for PostgreSQL and Redis
If PostgreSQL and Redis are already installed, add the relevant details in the environment variables. Otherwise, if you have Docker installed, use Docker Compose to set up PostgreSQL and Redis:
```bash
docker-compose up -d
```

#### Setting Up the Database Schema
Once the database is up, create the necessary schema. You can do this manually or via migration tools like `migrate`. To automate the migration process, use the Makefile:
```bash
make migrate-up
```
Alternatively, you can execute the SQL files located in `/cmd/migrate/migrations` manually.

#### Populating Test Data
To seed the database with test data, run:
```bash
go run cmd/migrate/seed/*
```

OR use the Makefile to run the seed files:
```bash
make seed
```

### Running the Server

After setup, start the API server by running the following command from the root of the project:
```bash
go run ./cmd/api
```

The API server will run on the port specified in the environment variables (default is `8080`). You can access it via `http://localhost:<port>`.

### Frontend (Proof of Concept)
This project includes a proof-of-concept frontend (created with React) to test user activation. Once a user registers, an email is sent with an activation link. Clicking this link opens a confirmation page.

To start the frontend:
```bash
cd ./web
npm run dev
```

## Tools and Technologies Used
- **Golang** (v1.22.5)
- **Redis** for profile caching
- **PostgreSQL** as the database
- **Mailtrap** and **SendGrid** for smtp emails

## Future Enhancements
Potential future enhancements for the API include:

- **Redis for Rate Limiting**: Improve performance by using Redis instead of an in-memory map for tracking rate limits.
- **Additional Tests**: Expand test coverage across the application.
- **Full Frontend Application**: Develop a comprehensive frontend application for the project.