GO - Simple bank
# GO - Simple Bank Service
This service implements these main features:

1. Create and manage users and its bank accounts.
2. Record all balance changes to each of the accounts.
3. Perform money transactions between the accounts.

Infrastructure:
- Database migrations: `golang-migrate`
- Generate CRUD code from SQL: `sqlc`
- Database: `postgreSQL`
- Mock database for testing: `gomock`
- Database documentation: `dbdocs`
- Generate SQL schema: `dbml2sql`
- Queue: `Redis`
- Web framework: `Gin`
- RPC framework: `gRPC`
- API gateway supports both HTTP and gRPC request: `grpc-gateway`
- API documentation: `Swagger`
- CI/CD: `github-action`
- Containerize: `docker`, `docker-compose`
- Deployment: `Kubernetes`
- AWS services: `IAM - ECR - RDS - Secrets Manager - EKS - Route53`
- TLS: `Let's Encrypt`

## Setup local development

### Install tools
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    brew install golang-migrate
    ```

- [DB Docs](https://dbdocs.io/docs)

    ```bash
    npm install -g dbdocs
    dbdocs login

- [DBML CLI](https://www.dbml.org/cli/#installation)

    ```bash
    npm install -g @dbml/cli
    dbml2sql --version
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

- Create the bank-network

    ``` bash
    make network
    ```

- Start postgres container:

    ```bash
    make postgres
    ```

- Create database:

    ```bash
    make createdb
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```

- Run db migration up 1 version:

    ```bash
    make migrateup1
    ```

- Run db migration down all versions:

    ```bash
    make migratedown
    ```

- Run db migration down 1 version:

    ```bash
    make migratedown1
    ```

### Documentation

- Generate DB documentation:

    ```bash
    make db_docs
    ```

- Access the DB documentation at [this address](https://dbdocs.io/hoangtk.0100/simple_bank). Password: `hoangtk`

### How to generate code

- Generate schema SQL file with DBML:

    ```bash
    make db_schema
    ```

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

- Generate DB mock with gomock:

    ```bash
    make mock
    ```

- Create a new db migration:

    ```bash
    make new_migration name=<migration_name>
    ```

### How to run

- Run server:

    ```bash
    make server
    ```

- Run test:

    ```bash
    make test
    ```

- Access database:

    ```bash
    make db
    ```

- Run containers:

    ```bash
    make up
    ```

- Remove containers:

    ```bash
    make down
    ```

## Deploy to kubernetes cluster

- [Install nginx ingress controller](https://kubernetes.github.io/ingress-nginx/deploy/#aws):

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.48.1/deploy/static/provider/aws/deploy.yaml
    ```

- [Install cert-manager](https://cert-manager.io/docs/installation/kubernetes/):

    ```bash
    kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml
    ```
