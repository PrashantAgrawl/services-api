Services API
A RESTful API built with Go to manage services and their versions, using PostgreSQL as the database. This project demonstrates a microservices-like architecture with endpoints to create, retrieve, and list services, leveraging raw SQL for database interactions to avoid GORM relationship issues.

Features
Create a new service with a name, description, and list of versions.
Retrieve all services with pagination, filtering, and sorting.
Fetch a specific service by ID, including its versions.
Uses PostgreSQL for persistent storage.
Raw SQL queries for database operations, bypassing GORM ORM limitations.
Logging with request IDs for traceability.
Built with Gin for fast HTTP routing.

Prerequisites
Go: Version 1.21 or higher (go version to check).
PostgreSQL: Version 13 or higher (psql --version to check).
Fork the Repository:
    Click "Fork" on GitHub.
Clone Your Fork:
    git@github.com:PrashantAgrawl/services-api.git
Installation
1. Clone the Repository:
    git clone git@github.com:prashantagrawal/services-api.git
    cd services-api
2. Install Dependencies:
    go mod download

Database Setup
1. Ensure PostgreSQL is Running:
    macOS: brew services start postgresql
2. Verify User Access:
    psql -U workflow -h localhost
3. Create Database services : 
    CREATE DATABASE services;
4. Create Role pertaining to it.
    CREATE ROLE workflow WITH PASSWORD 'workflow';
    CREATE ROLE workflow WITH LOGIN SUPERUSER INHERIT CREATEDB CREATEROLE NOREPLICATION BYPASSRLS
5. Create Tables post that : 
    SQL is already there in ddl.sql