# Project Name

This project is about building simple Supply Chain Management App that help User to manage item movement from one entity to other entities.  The stack consist of Backend(Golang) and Frontend(React). Besides, I use Postgres as the storage and Redis for caching.

## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Server](#running-the-server)
- [API Documentation](#api-documentation)
- [License](#license)

## Features
- Transfer Items from 1 entity to other entity. Current entity available are Customer, Technician and Warehouse. 
- Tracking each transaction process for each item.
- Inform the detail information of each item in 1 transaction.
- [to be designed]

## Prerequisites
List required software/tools and versions:
- Go >= 1.21
- PostgreSQL >= 15
- Redis >= 7 (if used)
- Docker & Docker Compose

## Installation

Clone the repository:
\`\`\`execute **bash command

> git clone https://github.com/milkywoosh/scm-simple.git

\`\`\`
make sure docker is installed in your system, if docker not installed yet then it must be.

Then spin up all of them together with docker-compose:

>docker compose up -d --build


# then

## Configuration

Copy the example environment file and fill in the values at .env:


| Variable      | Description              | Default |
|---------------|---------------------------|---------|
| `PORT`        | Server port                | 8000    |
| `PG_CONNSTRING`| Database connection string | 'postgresql://<username>:<password>@localhost:<port-postgres>/<datatabase-name>?sslmode=disable'      |
| `REDIS_PASSWORD`   | Redis connection string    | your-redis-pswrd      |

## Running the Server

### Development
**bash command

before this, make sure all the stack component is up, for exp: postgres

>go run ./...

or
>go run main.go



### Using Docker
**bash command
go to root workdir where docker-compose.yaml exists, then:

>docker-compose start

>docker-compose stop

if you have changed and want to refresh it, then execute:
>docker compose down -v

then
>docker compose up -d --build


The server will be available at `http://localhost:8000/api/`.

## API Documentation
I use Swagger docs

to do a renewal after update docs at handler, you can exec the following command. Then your docs will be automatically  updated

>swag init --parseDependency --parseDependencyLevel 1

you can see swagger docs here at development phase:
http://localhost:8000/swagger/index.html

## License
MIT