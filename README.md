user
-----------------------

# Run locally

docker network create cadicallegari_network


# HTTP example requests

> http get localhost:8080/v1/users


## TODO
- encode passwd
- integration tests
- implement grpc interface??
- add health check
- delete user on post message failures and add text on the readme about it

- midlewares
    https://pkg.go.dev/github.com/go-chi/chi/middleware#CleanPath
    https://pkg.go.dev/github.com/go-chi/chi/middleware#Heartbeat


