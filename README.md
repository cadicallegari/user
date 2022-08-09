user
===============

# Run locally

docker network create cadicallegari_network

target=builder make build

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


## TODO
- start the service from scratch and also run tests

- readme
    - dual write problem
    - mention migrations dir
    - pagination
    - encrypt algorithm
    - logs
    - pkg

