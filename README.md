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

The database is not mapped to a host port, you can enter in the container to run the integration tests using the following
```
make up
```

For helping with debug you can access phpadmin instance to check your database at address
> http://localhost:3006
> username: root
> password: root

You always change the docker-compose file to map the database in your host if you preffer.


# Integration tests

Having your local environment running locally, you can enter the container to run the tests with all dependencies already set.
```
make exec

make integration-test // integration tests

make test // unit tests
```

You can find more options running the make help target

```
make help
```

# Considerations

## pkg module


## Logging
## Migration
## Dual write problem
## Pagination
## Password encrypt

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



