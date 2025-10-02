# MPC Connector Frankfurt

This is a connector to call the Polytune MPC server.

## Setup (Docker)

Before proceeding, please make sure that the MPC server is running and is reachable

1. Build the docker container:
   ```bash
   docker build .
   ```
3. You can configure the connector by passing values in the environment variables. An [example.env](example.env) if
   provided.
4. Multiple instances of the connector can be run connecting to different MPC servers by changing the values in the env
   file.
   For example:
   ```bash
   docker --name connector0 run --env-file 1.env -p 3000:3000 connector
   ```
   ```bash
   docker --name connector0 run --env-file 2.env -p 3001:3000 connector
   ```
5. Launch the task on both connectors, ensuring the task on leader is launched in the end.
   ```bash
   curl -X POST -v http://localhost:3001/launch
   ```
   ```bash
   curl -X POST -v http://localhost:3000/launch
   ```
7. The garble program file used in the connector can be modified by changing the contents of
   file [program.garble](mpc/program.garble)

## Setup (Local)

1. Create an environment variable file and add appropriate values
   ```bash
   cp example.env .env
   ```
2. Run the program
   ```bash
   go run cmd/*.go
   ```
