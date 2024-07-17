# Backend Blogging REST API

This is a backend REST API built with Golang, using the Chi router and Postgres as the database. It is part of a collection of backend projects provided at [roadmap.sh](https://roadmap.sh/backend/project-ideas#1-personal-blogging-platform-api)

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Server](#running-the-server)
- [Configuration](#configuration)
- [Swagger Documentation](#swagger-documentation)
- [Contributing](#contributing)

## Introduction

This project is a RESTful API that provides various endpoints for managing articles on blog. It leverages the Chi router for handling HTTP requests and Postgres for persistent data storage.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- Golang installed on your machine
- Postgres database set up and running
- Make tool installed for running commands

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/ayo-awe/blogging_api.git
   ```
2. Navigate to the project directory:
   ```sh
   cd blogging_api
   ```
3. Install the project dependencies:
   ```sh
   go mod download
   ```

## Running the Server

To run the server in development mode, use the following command:

```sh
make dev
```

## Configuration

Configuration options for the API, such as database connection strings and server ports, are specified in a `.env` file. Create a `.env` file in the project root and add the necessary configuration variables.

Example `.env` file:

```env
DB_URL=postgresql://<username>:<password>@<host>:<port>/<database>
PORT=8080
```

## Swagger Documentation

The API documentation is available via Swagger. Once the server is running, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes.
