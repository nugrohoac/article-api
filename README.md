# Kumparan assessment
Article Repository

## Prerequisites
Make sure you have installed all of the following prerequisites on your development machine:
* go - [Download & Install](https://go.dev/dl/)
* golang-migrate - [Install golang-migrate](https://github.com/golang-migrate/migrate)
* PostgreSQL
* redis
* mockery - [Install golang-migrate](https://github.com/vektra/mockery)


## Integration Test
This is will running integration test and unit test.
#### 1. Create file migrate
```bash
$ migrate create -ext sql -dir destination/directory name_migration
# Example
$ migrate create -ext sql -dir internal/postgresql/migrations create_table_applicant_score
```

Each migration has an up and down migration.
```bash
1481574547_create_article_table.up.sql
1481574547_create_article_table.down.sql
```
* write up ddl at file.up.sql
* write reverse of ddl at file.down.sql

#### 2. Create file migrate
Make sure you have modified configuration at [postgre suite](./internal/postgresql/postgre_suite.go)
* driver     = "defauult postgres"
* host       = "default is localhost if running local"
* dbname     = "make sure you have been created database for testing"
* sslMode    = "deafult disable"
* userName   = "default postgres"
* password   = ""
* searchPath = "default is public"
* port       = "default is 5432"


#### 3. Adjust redis connection
Please adjust connection base on your machine at [redis suite](./internal/cache/redis_suite.go)

#### Running Integration Test
```bash
$ make test
```

## Unit Test
This is just running unit test without integration test. Make sure mocks is up-to-date.
* generate or update mock base on name of [interface](src/business/contract.go)
#### generate or update mocks
```bash
$ mockery -name=name-of-interface
$ mockery --dir=source/directory --name=nameInterface --output=destination/directory
```
#### Running Integration Test
```bash
$ make unittest
```

## Running Apps
#### Running Integration Test
* download all dependencies base on go.mod
```bash
$ go mod vendor
```
* running command
```bash
$ make run
```

if your config env is custom
```bash
$ go run cmd/main/main.go -config-path=your-custom-env
```