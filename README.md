# BANK-API
A SAAS banking application built with Go, and PostgreSQL, To manage your bank account, and make transactions.

## How to run

* using docker

If you have docker installed in your machine, then to build and run app inside docker container, follow the steps

- Build docker image

```shell

$docker build -t bank-api:latest .

```

after build process is completed, then the next step is running bank-api image

- Run bank-api image

```shell

$docker run -p 8000:8000 bank-api:latest

```

* using podman

If you are using podman instead of docker, the command slightly the same but all you need is change `docker` with `podman`

- Build podman image

```shell

podman build -t bank-api:latest .

```

after build process is completed, then the next step is running bank-api image

- Run bank-api image

```shell

$podman run -p 8000:8000 bank-api:latest

```

* Running natively

- Run app

you can run app directly without build as docker or podman image, but to do that you need to have Go SDK installe in 
your machine. Please refer to [download and install instruction](https://go.dev/doc/install) to install Go into your machine.

If you have installed Go SDK into your machine, to run follow below step

- Run as REST API

```shell

$go run main.go rest

```

- Run as gRPC API

```shell

$go run main.go gapi

```

- Build app

As Go is compiled programming language so we can choose either run directly or compile the app. The pros for compile the app
is that after compile, you can deliver or run compiled app without need to have Go SDK in your machine.

```shell
$go build -o bank-api main.go

```

and after build success, we can execute the app using command

```shell

$./bank-api rest

```

you can see complete command here 

```shell

$go run main.go

```

result 

```shell
Bank API made with Go

Usage:
  BANK [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  gapi        Run Banking API HTTP server (gRPC)
  gateway     Run Banking API HTTP server (gRPC Gateway)
  help        Help about any command
  migrate     Run Banking API migration
  rest        Run Banking API HTTP server (rest-API)

```



_Note : if you want to change application configuration, open `config.yaml` file and change the value match with your system._ 



## Features

The applicatoin uses `paseto` for authentication.

### User

- Create a user
- Login
- Renew access token
- Upadte user

### Account
- Create an account
- Get all accounts (Of the logged in user)
- Delete an account (Soft delete, Can be restored)
- Restore an account (After has been deleted)

### Transaction
- Create a transaction
- Get all transactions (Of a specific account)

## Tech Stack

- Gin
- PostgresSQL
- gomock

## GRPC Services

The project uses GRPC besides the REST API, to communicate with db. But the GRPC are not implemented yet fully as the API.
