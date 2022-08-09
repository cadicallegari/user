user
===============

User service is a simple REST API service, with CRUD operations for user entity.


# Code structure

Apart from files like docker-compose.yml, Dockerfile, Makefile, etc.
The project structure is composed of the following structure

```
├── cmd (service binaries/entry points)
├── http (http related code)
├── mem (mem related code)
├── mock (mocks for tests)
├── mysql (mysql related code)
├── pkg (code to support service implementation, normally is a external dep)
├── user.go (service domain definitions)
└── service.go (service implementation)
```

# Run locally

### TL;DR

First
```
make create-network dev-build up
```

After the previous command, the service is running in the docker containers and
accessible in the port 8080.
> http://localhost:8080


Entering into the container to run the tests
```
make exec
```

Then run the integration tests
```
make integration-test
```

### Explanation

You can run the project locally and also run integration tests using docker-compose.
You will need to create the network and build the dev image if it is your first time running it
running the following make targets.

```
make create-network
make dev-build
```

To start the dependencies like database you can use target
```
make up
```

For helping with debug you can access phpadmin instance to check your database at address
> http://localhost:3006
> username: root
> password: root

You can always change the docker-compose file to map the database in your host if you preffer,
and run the service or tests setting `USER_MYSQL_URL` env var.

# Integration and unit tests

Having your local environment running, you can enter the container to run the tests with all dependencies already set.
```
make exec
make test // unit tests also can run locally without any aditional configuration
make integration-test // integration tests
```

You can find more options running the make help target

```
make help
```

# Considerations

## pkg module

The `pkg` directory contains a bunch of packages that help build services and are not related to the
services domain. For reusability in multi-repo architecture,
it should be moved to its repo and used as a dependency.

## Logging

Logging is included in the `pkg` directory.
The implementation extends logrus to add context logging features.

You can configure log format and level through `USER_LOG_FORMATTER` and `USER_LOG_LEVEL` respectively.

## Migration

Migrations is also included in the `pkg` directory along with more support for handling databases.
The migrations can run automatically and are helpfull during integration tests.

You can configure the directory where the migration is located through `USER_MYSQL_MIGRATIONS_DIR` env var.
There are other configurations that can be set via env vars, check out `pkg/xdatabase/xsql/config.go` for more details.

## Pagination

The pagination uses simple `page` and `per_page` method.

The response for `GET /v1/users` is an object with the following fields

```
{
    "next_page": 3,
    "prev_page": 1,
    "total": 200,
    "users": [...]
}
```

To navigate forward and backward into the data, the client must use the `next_page` and `prev_page` fields from the response.

Example:

```
next page: /v1/users?page=3
prev page: /v1/users?page=1
```

## Password encrypt

The password received in the request body is encrypted using bcrypt algorithm.
The plain text password is discarded and only the encrypted value is stored in the database.

Currently, the encrypted password is not returned in the GET responses but can be easily added, dependent on the use cases.

The cost to generate the encrypted password can be configured using `USER_PASSWORD_GENERATION_COST` env var.

## Events

The system is ready to publish events after state changes in the users.

For simplicity, none real event broker was used. There is only a dummy implementation of the `EventService` under the mem module.

To give an example of a possible solution, each method in the `EventService` can publish in topics like
`user.created`, `user.updated`, `user.deleted` respectivily, in kafka, pubsub, nats or other solution.

### Dual write problem

For simplicity, the current solution for publishing events has the dual write problem.

Among other solutions for this problem, we can highlight `Listen yourself`and `Outbox pattern`.

# Deploy

The artifact generated in the project is a docker image. It can be deployed in any container management system,
like Kubernetes, Mesos, etc.

The CI usually calls these makefile targets to build the final image and publish it to the docker registry.

```
make build // build production ready image
make push  // push the image to repository
```

After that, you can use the image in your container management system.
This part is not covered in this project.


# HTTP request examples

## Create user

```
curl -H "Content-Type: application/json" -X POST localhost:8080/v1/users \
    -d '{"first_name": "Alice", "last_name": "Chains", "nickname": "AB123", "password": "supersecurepassword", "email": "alice@chains.com", "country": "UK"}'
```

## Get users

```
curl -X GET 'localhost:8080/v1/users'
curl -X GET 'localhost:8080/v1/users?search=alice'
curl -X GET 'localhost:8080/v1/users?country=BR'
curl -X GET 'localhost:8080/v1/users?per_page=1&page=1'
```

## Get user

```
curl -X GET localhost:8080/v1/users/{user_id}
```

## Update user

```
curl -H "Content-Type: application/json" -X PUT localhost:8080/v1/users/{user_id} \
    -d '{"first_name": "Alice", "last_name": "Bob", "nickname": "AB123", "password": "supersecurepassword", "email": "alice@bob.com", "country": "UK"}'
```

## Delete users

```
curl -X DELETE localhost:8080/v1/users
```


