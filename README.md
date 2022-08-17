# blogs service

a sample of CRUD service for blogs with gRPC

## run mongodb

1. create `/data/db` directory

```sh
mkdir -p /data/db
```

2. run mongodb

```sh
mongod --dbpath /data/db
```

## run server

```sh
make server
```

## run client

```sh
make client
```

## (optional) run evans CLI

```sh
evans -p 50051 -r
```

Basically with [evans](https://github.com/ktr0731/evans) CLI you can talk to the service directly and get info or even make rpc calls within the REPL.
