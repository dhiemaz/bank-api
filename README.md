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
