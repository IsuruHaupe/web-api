# Web-Api

This is a simple API where an user can get a quick overview over all contacts resources like person, skills...

The following use cases are implemented : 

* [X] **UC1** : CRUD endpoint for managing contacts. A contact is composed of the following attributes : 
    * Firstname
    * Lastname
    * Fullname
    * Address
    * Email
    * Mobile phone number

* [X] **UC2** : CRUD endpoint for managing skills. A skill is composed of the following attributes : 
    * Name
    * Level (expertise)

* [X] **UC3** : Documentation using Swagger. More on that in [Swagger](#swagger) section.

* [X] **UC4** : Authentification/Authorization, where users can only change their contacts and skills.

# Makefile 

We use a makefile to easily setup the local environment. The following commands are available : 

    * make postgres : run a postgres12 instance
    * make createdb : create a database named web-api
    * make dropdb : remove the db in the docker instance
    * make migrateforce : force database migration from sql files (used when migration failed using migrate CLI)
    * make migrateup : migrate database schema to the docker instance
    * make migratedown : remove database from the docker instance
    * make sqlc : generate type-safe interfaces to query database
    * make tests : run all the tests in the module
    * make server : run the server 
    * make mock : generate mocks for testing the API
    * make swagger : create the swagger documentation for the API

## Setup local environment

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    brew install golang-migrate
    ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

### Setup infrastructure

- Start a postgres container:

    ```bash
    make postgres
    ```

- Create the web-api database:

    ```bash
    make createdb
    ```

- Run database migration:

    ```bash
    make migrateup
    ```
### Running the application 

- Run server:

    ```bash
    make server
    ```

### Running the tests 

- Generate DB mock with gomock:

    ```bash
    make mock
    ```

- Run test:

    ```bash
    make tests
    ```


# Docker 

If you don't want to setup manually everything you can just use the docker-compose provided. 

Before that, **change this line in `app.env` file**: 

```
DB_SOURCE=postgresql://root:secret@localhost:5432/web-api?sslmode=disable
```

by this one : 

```
DB_SOURCE=postgresql://root:secret@postgres:5432/web-api?sslmode=disable
```

* Run the application in a container:

    ```bash
    make docker-up
    ```

* Remove the application :

    ```bash
    make docker-down
    ```
# Swagger 

You can follow and use the application via the swagger. When the app is running go to [this link](http://localhost:8080/swagger/index.html) and you will will access the swagger documentation of the API. All the routes except for creating a new user are restricted. When you created a new user and logged in via the specific route, remember to add the token in the Authorize so that all the requests used the token to be identified.

* Swagger documentation

![Swagger](./images/swagger.png)

* Create an user 

![Swagger](./images/CreateUser.png)

* Login to retrieve a token

![Swagger](./images/LoginUser.png)

![Swagger](./images/RetrieveToken.png)

* Use that token for the following requests

![Swagger](./images/Authorize.png)

![Swagger](./images/AddTokerBearer.png)

