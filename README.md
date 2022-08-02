user
-----------------------

# Run locally

docker network create cadicallegari_network


# HTTP example requests

> http GET localhost:8080/v1/users

> echo '{"first_name": "Alice", "last_name": "Bob", "nickname": "AB123", "password": "supersecurepassword", "email": "alice@bob.com", "country": "UK"}' | \
    http POST http://localhost:8080/v1/users

> http GET localhost:8080/v1/users/1

> echo '{"first_name": "Alice edite", "last_name": "Bob", "nickname": "AB123", "password": "supersecurepassword", "email": "alice@bob.com", "country": "UK"}' | \
    http PUT http://localhost:8080/v1/users/1

> http DELETE http://localhost:8080/v1/users/1


## TODO
- encode passwd
- integration tests
- implement grpc interface??
- add health check
- delete user on post message failures and add text on the readme about it

- midlewares
    https://pkg.go.dev/github.com/go-chi/chi/middleware#CleanPath
    https://pkg.go.dev/github.com/go-chi/chi/middleware#Heartbeat


