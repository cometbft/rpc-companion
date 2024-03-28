# RPC Companion (Work In Progress)

This is the repository for the RPC Companion solution based on the [ADR-101 - Data Companion Pull API](https://github.com/cometbft/cometbft/blob/main/docs/architecture/adr-101-data-companion-pull-api.md)

Please see [ADR-102 - RPC Companion](https://github.com/cometbft/cometbft/blob/main/docs/references/architecture/adr-102-rpc-companion.md) for more information in regards this implementation and architecture solution.

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

```
time=2023-10-18T17:23:55.311-04:00 level=INFO msg="New block" service=Ingest module=Fetcher method=WatchNewBlock height=1549
time=2023-10-18T17:23:55.312-04:00 level=INFO msg="Get block" service=Ingest module=Fetcher method=GetBlock height=1549
time=2023-10-18T17:23:55.312-04:00 level=INFO msg="Processing job" service=Ingest module=Fetcher height=1549
time=2023-10-18T17:23:55.314-04:00 level=INFO msg="Get block retain height" service=Ingest module=Fetcher method=GetBlockRetainHeight retain_height=1548 app_retain_height=0
time=2023-10-18T17:23:55.316-04:00 level=INFO msg="Set block retain height" service=Ingest module=Fetcher method=SetBlockRetainHeight height=1549
time=2023-10-18T17:23:55.316-04:00 level=INFO msg="Processed block job" service=Ingest module=Fetcher method=ProcessBlockJob height=1549
time=2023-10-18T17:23:56.327-04:00 level=INFO msg="New block" service=Ingest module=Fetcher method=WatchNewBlock height=1550
time=2023-10-18T17:23:56.328-04:00 level=INFO msg="Get block" service=Ingest module=Fetcher method=GetBlock height=1550
time=2023-10-18T17:23:56.328-04:00 level=INFO msg="Processing job" service=Ingest module=Fetcher height=1550
time=2023-10-18T17:23:56.331-04:00 level=INFO msg="Get block retain height" service=Ingest module=Fetcher method=GetBlockRetainHeight retain_height=1549 app_retain_height=0
time=2023-10-18T17:23:56.332-04:00 level=INFO msg="Set block retain height" service=Ingest module=Fetcher method=SetBlockRetainHeight height=1550
time=2023-10-18T17:23:56.332-04:00 level=INFO msg="Processed block job" service=Ingest module=Fetcher method=ProcessBlockJob height=1550

```
