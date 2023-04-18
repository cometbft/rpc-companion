# RPC Companion 

### WIP (Work In Progress)

This is the repository for the RPC Companion solution based on the Data Companion ADR

## Running the demo

In order to run the demo, you need `docker-compose` and `docker` installed on your machine.

You will also need a local instance of [CometBFT](https://github.com/cometbft/cometbft) running in your machine.

#### 1. Running CometBFT

Once you install CometBFT locally, please run the following commands:

Initialize CometBFT:
```
cometbft init
```

Run the kvstore app:
```
cometbft start --proxy_app kvstore
```

To ensure the service is running, open a browser and navigate to: http://localhost:26657/block?height=1

If you see a JSON response for the block at height 1 then it's working.

#### 2. Running the database

The database (Postgres) and its admin interface (pgadmin) runs inside docker containers. 

In order to run them locally please run the following command (this assumes the terminal is in the `rpc-companion` forlder:

```
docker-compose -f ./database/docker/docker-compose.yml up
```

To ensure the database was started properly, open a browser and navigate to: http://localhost:5050

If you get a prompt to login, please use:

```
Email: pgadmin@pgadmin.com
Password: pgadmin123
```

If everything works, you should see a `comet_postgres_group` entry on the left side navigation browser window.

#### 3. Running the RPC Companion REST service

In order to run the demo, please run the following command:

```
go run main.go
```
The program will fetch blocks from height 1 to 100 and insert them into the database.

To test the REST service, please open a browser and navigate to: http://localhost:8080/v1/block?height=1, the service
will return a JSON that should be equivalent to the one returned by the CometBFT RPC service.


