# Esther

**Esther** is the *events store*. Its purpose is to store some callback events, in order to execute them later. These callback events can be created, then modified if needed, and at the end you can apply them and/or delete them.

## TL;DR

Set the needed environment variables (see examples in the [.env](./.env) file) and launch the server:

```bash
export PORT=8080
export MONGODB_SERVICE_HOST=mongodb://<your_mongodb_host>
export MONGODB_PORT=<your_mongodb_port>                     # usually 27017
export MONGODB_DATABASE_NAME=esther
go run .
```

Then open a new console and start using the endpoints:

- `/ready`: TODO

The OpenAPI specification and the associated tools are available here:

- <http://localhost:8080/openapi/>

## Persistence

TODO

## Testing

You can run the tests locally and see the test coverage:

```bash
go test -v -race -coverprofile=coverage.out $(go list ./... | grep -v /vendor/)
go tool cover -func=coverage.out
```

## Using Docker

You can use `docker-compose` to run **Esther** in a *Docker* environment:

```bash
cp docker-compose.override.yml.dist docker-compose.override.yml
docker-compose up --build
```

The server will be listening on the port defined in your `docker-compose.override.yml` file.
