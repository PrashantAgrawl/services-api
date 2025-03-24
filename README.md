# Services API

A RESTful API built with Go to manage services and their versions, using PostgreSQL as the database. This project demonstrates a microservices-like architecture with endpoints to create, retrieve, and list services, leveraging raw SQL for database interactions to avoid GORM relationship issues.

## Features
- Create a new service with a name, description, and list of versions.
- Retrieve all services with pagination, filtering, and sorting.
- Fetch a specific service by ID, including its versions.
- Uses PostgreSQL for persistent storage.
- Raw SQL queries for database operations, bypassing GORM ORM limitations.
- Logging with request IDs for traceability.
- Built with Gin for fast HTTP routing.

## Prerequisites
- **Go**: Version 1.21 or higher (check with `go version`).
- **PostgreSQL**: Version 13 or higher (check with `psql --version`).

## Installation
- **Clone the Repository**:
  - Using Git:
    ```bash
    git clone git@github.com:prashantagrawal/services-api.git
    OR 
    gh repo clone PrashantAgrawl/services-api
    cd services-api
- **Install Dependencies**:
    ```bash
    go mod download

## Database Setup
- **Ensure PostgreSQL is Running**:
  - macOS:
    ```bash
    brew services start postgresql
- **Verify User Access**:
    ```bash
    psql -U workflow -h localhost
- **Create Database services**: 
    ```bash
    CREATE DATABASE services;
- **Create Role pertaining to it**:
    ```bash
    CREATE ROLE workflow WITH PASSWORD 'workflow';
    CREATE ROLE workflow WITH LOGIN SUPERUSER INHERIT CREATEDB CREATEROLE NOREPLICATION BYPASSRLS;
- **Create Tables post that**: 
    SQL is already there in ddl.sql

## Running the code:
- **SETUP**
    ```bash
    go mod tidy
    go mod vendor
    go build
    go run main.go

## NOTE : Seed the data before running the build or uncomment the db.SeedData call in main.go

## Post the installations below are the curl : 
- **Get Services API (To list all the Services)**:
    ```bash
    curl -X GET "http://localhost:8080/services?page=1&page_size=10"
- **Get Service API (To fetch each service data using service id)**: 
    ```bash
    curl -X GET "http://localhost:8080/services/1"
