# RPC Companion (Work In Progress)

This is the repository for the RPC Companion solution based on the Data Companion ADR

## Starting the Postgres Database

Access the Docker folder:

`cd ./database/docker`

Run the services

`docker-compose -f docker-compose.yml up`

## Accessing the database with pgAdmin

Open a browser and navigate to
http://127.0.0.1:5050

Use the following credentials:

```
User: pgadmin@pgadmin.com
Password: pgadmin123
```