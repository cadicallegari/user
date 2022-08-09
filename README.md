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

TL;DR

First
```
make create-network dev-build up
```

Mysql container might take a minute or so to get ready

Enter into the container to run the tests
```
make exec
```

Then run the integration tests
```
make integration-test
```


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

Having your local environment running locally, you can enter the container to run the tests with all dependencies already set.
```
make exec
make test // unit tests
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
it should be moved to its repo and used as a dependency module.

## Logging

Logging is included in the `pkg` directory.
The implementation extends the [logrus](https://github.com/sirupsen/logrus) package to add context logging features.

You can configure log format and level through `USER_LOG_FORMATTER` and `USER_LOG_LEVEL` respectively.

## Migration

Migrations is also included in the `pkg` directory.
It makes easy to run integration tests, and they are also checked and run if needed when the service starts.

You can configure the directory where the migration is located through `USER_MYSQL_MIGRATIONS_DIR` env var.

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

To navigate forward and backward into the data, the user must use the `next_page` and `prev_page` fields from the response.

Example:

> next page: /v1/users?page=3
> prev page: /v1/users?page=1

## Password encrypt

The password received in the request body is encrypted using bcrypt algorithm.
The plain text password is discarded and only the encrypted value is stored in the database.

Currently, the encrypted password is not returned in the GET responses but can be easily added, dependent on the use cases.

The cost to generate the encrypted password can be configured using `USER_PASSWORD_GENERATION_COST` env var.

## Events
### Dual write problem

# Deploy

These make targets are normally called by CI to build the final image and publish to the docker registry.

```
make build // build production ready image
make push  // push the image to repository
```


# Request examples

# HTTP example requests

> http GET localhost:8080/v1/users

> echo '{"first_name": "Alice", "last_name": "Bob", "nickname": "AB123", "password": "supersecurepassword", "email": "alice@bob.com", "country": "UK"}' | \
    http POST http://localhost:8080/v1/users

> http GET localhost:8080/v1/users/1

> echo '{"first_name": "Alice edite", "last_name": "Bob", "nickname": "AB123", "password": "supersecurepassword", "email": "alice@bob.com", "country": "UK"}' | \
    http PUT http://localhost:8080/v1/users/1

> http DELETE http://localhost:8080/v1/users/1


http GET 'localhost:8080/v1/users?search=alice'
http GET 'localhost:8080/v1/users?country=BR'
http GET 'localhost:8080/v1/users?per_page=1&page=1'



