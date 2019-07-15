# Simple RESTful API

## Routes
* `GET /ping`: health check<br/>
* `GET /v1/users/verify/:code`: verify email<br/>
* `POST /v1/users/new`: register user<br/>
* `POST /v1/users/login`: login<br/>
* `POST /v1/users/forgot_password`: forgot password<br/>
* `POST /v1/users/reset/:code`: reset password<br/>

## Configuration
Create config/app.yaml from config/template.yaml<br/>

## Test
```shell
# register new user
curl -X POST -H 'Content-Type: application/json' -d '{"business_name": "demo", "full_name": "demo", "business_email": "example@mail.ru", "business_phone": "123123", "password": "Qwerty1!", "confirm_password": "Qwerty1!"}' 'http://localhost:8080/v1/users/new'

# login
curl -X POST -H 'Content-Type: application/json' -d '{"business_email": "example@mail.ru", "password": "Qwerty1!"}' 'http://localhost:8080/v1/users/login'

# verify email
curl -X GET 'http://localhost:8080/v1/users/verify/{code}?user_id=1'

# forgot password
curl -X POST -H 'Content-Type: application/json' -d '{"business_email": "example@mail.ru"}' http://localhost:8080/v1/users/forgot_password

#reset password
curl -X POST -H 'Content-Type: application/json' -d '{"password": "asd123ASD!", "confirm_password": "asd123ASD!"}' 'http://localhost:8080/v1/users/reset/{code}?user_id=1'
```

## Project Structure

This project divided into four main packages:

* `models`: contains the data structures used for communication between different layers.
* `services`: contains the main business logic of the application.
* `daos`: contains the DAO (Data Access Object) layer that interacts with persistent storage.
* `apis`: contains the API layer that wires up the HTTP routes with the corresponding service APIs.

[Dependency inversion principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)
is followed to make these packages independent of each other and thus easier to test and maintain.

The rest of the packages are used globally:
 
* `app`: contains routing middlewares and application-level configurations
* `utils`: contains utility code

The main entry of the application is in the `main.go` file. It does the following work:

* load external configuration
* establish database connection
* instantiate components and inject dependencies
* start the HTTP server
