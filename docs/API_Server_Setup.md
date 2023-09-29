# Running API Server

This document explains how you can setup a development environment for `ZWEB Builder` API server.

## Pre-requisites

- [Go](https://go.dev/doc/install)

- [PostgreSQL](https://www.postgresql.org/download/)

- [zweb-supervisor-backend](https://github.com/zilliangroup/zweb-supervisor-backend)

## Local Setup

1. Setup the PostgreSQL database

    - Running the [script](../scripts/postgres-init.sh) to create the database and tables

2. Setup the `zweb-supervisor-backend`

    - Following the setup steps in [zweb-supervisor-backend](https://github.com/zilliangroup/deploy-zweb-manually/tree/main/build-by-yourself#build-zweb-supervisor-backend)

3. Change the default env config

   Change the default env config in `pkg/db/connection.go` to the PostgreSQL config.

   Change the default env config in `internal/util/supervisior/token_validator.go` to the `zweb-supervisor-backend` config.

4. Running the ZWEB Builder API server

    ```bash
    go run github.com/zilliangroup/builder-backend/cmd/http-server
    ```

   This will start the ZWEB Builder API server on  `http://127.0.0.1:8001`.

5. Extract the JWT token for the user `root`

    ```bash
    curl 'http://{{zweb-supervisor-backend-addr}}/api/v1/auth/signin' --data-raw '{"email":"root","password":"password"}' -v
    ```

   Get the value of response header `zweb-token` as the next API call's `Authorization` header value.

6. Test the API server

    ```bash
    curl 'http://127.0.0.1:8001/api/v1/teams/:teamID/apps' -H 'Authorization: {{Value of response header `zweb-token`}}'
    ```

   The value of `:teamID` is `ILAfx4p1C7d0`.

## Need Assistance

- If you are unable to resolve any issue while doing the setup, please feel free to ask questions on our [Discord channel](https://discord.com/invite/zilliangroup) or initiate a [Github discussion](https://github.com/orgs/zilliangroup/discussions). We'll be happy to help you.
- In case you notice any discrepancy, please raise an issue on Github.