# MPC Connector Frankfurt

This is a connector to call the Polytune MPC server.

## Setup (Docker)

Before proceeding, please make sure that the MPC server is running and is reachable

1. Build the docker container: `docker build .`
2. You can configure the connector by passing values in the environment variables. An [example.env](example.env) if
   provided.
3. Multiple instances of the connector can be run connecting to different MPC servers by changing the values in the env
   file.
   For example:
   `docker --name connector0 run --env-file 1.env -p 3000:3000 connector`
   `docker --name connector0 run --env-file 2.env -p 3001:3000 connector`
4. Launch the task on both connectors, ensuring the task on leader is launched in the end.
   `curl -X POST -v http://localhost:3001/launch`
   `curl -X POST -v http://localhost:3000/launch`
5. The garble program file used in the connector can be modified by changing the contents of
   file [program.garble](mpc/program.garble)

## Setup (Local)

1. Create an environment variable file and add appropriate values
  `cp example.env .env`
2. Run the program
   `go run cmd/*.go`