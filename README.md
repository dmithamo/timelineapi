# timelineapi

## What this is

An api that exposes CRUD functionality - built with Go and a sprinkle of blood and a sprinkle of sweat and a sprinkle of tears.

### Run, app, run

To test, execute the following commands in a Terminal (bash or similar):

```bash
#!/bin/bash
git clone <this-repo>
cd timelineapi
git checkout develop
go run cmd/api/* -addr=:3001 -dsn='dbuser:dbpassword@[host]/dbname' -rdb=true -cdsn='[redis-server host]:[port]'
```

If all went well, the api should serve available endpoints on localhost

### Flags definition

|Flag|Description|Default|
|:------|:-----|----:|
`addr`| port at which the app will run | `:3001`
`dsn`| DSN of the database | `REQUIRED`
`rdb`| setting this to true will create the database tables, dropping them if they already exist | `false`
`cdsn`| DSN of the `redis` database for persisting sessions | `REQUIRED`

### The Stack

This project uses the following open source technologies

- [The Go programming Language](https://golang.org/)
- [MySQL Database](https://www.mysql.com/)
- [Gorilla Mux Router](https://github.com/gorilla/mux)
