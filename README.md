# RPC Companion (Work In Progress)

This is the repository for the RPC Companion solution based on the [ADR-101 - Data Companion Pull API](https://github.com/cometbft/cometbft/blob/main/docs/architecture/adr-101-data-companion-pull-api.md)

Please see ADR-102 - RPC Companion for more information in regards this implementation and architecture solution.

## Starting the Postgres Database

> NOTE: This assumes you have [Docker](https://www.docker.com/) and Docker Compose already installed in your machine

The reference implementation of RPC Companion utilizes Postgres as its ingest service storage system to save 
the data acquired from the node. It employs a simple schema to store the node data in a relational database. 
However, it does not normalize the data or utilize a schema that can store structured data such as a Block. 
Essentially, the data is saved as a byte array.

In order to run the database, access the Docker folder:

`cd ./database/docker`

And run the docker services that will host the database (Postgres) and the database IDE (pgAdmin)

`docker-compose -f docker-compose.yml up`

## Accessing the database with pgAdmin

Open a browser and navigate to
http://127.0.0.1:5050

Use the following credentials:

```
User: pgadmin@pgadmin.com
Password: pgadmin123
```

## Run CometBT (with gRPC services support)

Checkout the `main` branch from [cometbft](https://github.com/cometbft/cometbft) repository:

```
git checkout https://github.com/cometbft/cometbft.git

cd cometbft
make install
```

Configure the gRPC services. Modify the `$HOME/.cometbft/config/config.toml` file to enable the gRPC services as per
[these instructions](https://github.com/cometbft/cometbft/blob/main/docs/data-companion/grpc.md#enabling-the-grpc-services)

Start an instance of the kvstore app, e.g.:

```
cometbft init
cometbft start --proxy_app kvstore
```

## Start the ingest service

The ingest service has a crucial role in monitoring new blocks generated on the CometBFT node. After detecting
the newly created block, it retrieves the necessary information and inserts it into the database. 
Once the data is safely stored, the service uses the data companion API to notify the node that the 
information can be pruned.

In order to run the ingest service please make sure you follow this steps outlined below.

### Configuration

Open a new terminal tab. Create a new file in `$HOME/.rpc-companion` named `config.toml` and add the addresses
for the gRPC endpoints (regular and privileged) and the database connection information, e.g.:

```
[grpc_client]
address = "0.0.0.0:8080"
privileged_address = "0.0.0.0:8088"

[storage]
connection = "postgres://postgres:postgres@0.0.0.0:15432/postgres?sslmode=disable"
```

Save the file.

### Run the rpc-companion ingest service

Build the `rpc-companion` binary and run the ingest service

```
go build

./rpc-companion ingest service
```

If everything is compiled and configured correctly, you will see logs displaying the ingest service fetching 
new blocks. The service then sets the retain height information so CometBFT can prune them from its storage.