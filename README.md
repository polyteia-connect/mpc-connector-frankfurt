# MPC Connector Frankfurt

This is a connector to call the Polytune MPC server for Frankfurt use case.
Connector is divided into 2 modules: Measles and ESU. Both modules run as separate services.
Both modules can be configured via environment variables. See below the configuration options for each module.

## Measles Module

### Configuration Options (Environment Variables)

- `PORT`: **(Optional, Default: 3000)** The port to listen on.
- `HOST`: **(Optional, Default: 0.0.0.0)** The host to listen on.
- `DEBUG`: **(Optional, Default: false)** Whether to enable debug mode. This flag logs all requests and responses to the console. Only recommended for development as in production it could generate a lot of log noise.
- `JWT_KEY_FILE`: **(Required)** The path to the JWT key file in PEM format. This key is used to sign the JWT token that is used to authenticate requests to the gateway server.
- `JWT_ISSUER`: **(Required)** The issuer of the JWT token. This will set the `iss` claim in the JWT token.
- `MPC_BASE_URL`: **(Required)** The base URL of the MPC server. This is the URL of the polytune MPC server responsible for the measles data that will be used to launch tasks.
- `MPC_LEADER_ID`: **(Optional, Default: 0)** The leader ID of the MPC server. This is the ID of the party that will be the leader of the MPC task. In case of measles, the leader is the party that will be responsible for the measles task.
- `MPC_PARTY_ID`: **(Optional, Default: 0)** The party ID of the MPC server. This is the ID of the party that will be the participant of the MPC task. In case of measles, the party ID is the ID of the party that will be the participant of the measles task which in our case is the same as the leader ID.
- `MPC_PARTICIPANTS`: **(Required)** The participants of the MPC server. This is a comma separated list of URLs of the MPC servers that will be the participants of the MPC task.
- `CALLBACK_BASE_URL`: **(Required)** The base URL of the callback server. This is the URL of the server that will be called back when the MPC task is completed by the Polytune MPC server.
- `ESU_BASE_URL`: **(Required)** The base URL of the ESU server. This is the URL of the ESU gateway server that will be used to call the ESU connector module.

## ESU Module

### Configuration Options (Environment Variables)

- `PORT`: **(Optional, Default: 3000)** The port to listen on.
- `HOST`: **(Optional, Default: 0.0.0.0)** The host to listen on.
- `DEBUG`: **(Optional, Default: false)** Whether to enable debug mode. This flag logs all requests and responses to the console. Only recommended for development as in production it could generate a lot of log noise.
- `JWT_KEY_FILE`: **(Required)** The path to the JWT key file in PEM format. This key is used to sign the JWT token that is used to authenticate requests to the gateway server.
- `JWT_ISSUER`: **(Required)** The issuer of the JWT token. This will set the `iss` claim in the JWT token.
- `MPC_BASE_URL`: **(Required)** The base URL of the MPC server. This is the URL of the polytune MPC server responsible for the esu vaccination data that will be used to launch tasks.
- `MPC_LEADER_ID`: **(Optional, Default: 0)** The leader ID of the MPC server. This is the ID of the party that will be the leader of the MPC task. In our, the leader is the party that will be responsible for the measles task.
- `MPC_PARTY_ID`: **(Optional, Default: 1)** The party ID of the MPC server. This is the ID of the party that will be the participant of the MPC task. In our use case, the party ID is the ID of the party that will be the participant of the esu task.
- `MPC_PARTICIPANTS`: **(Required)** The participants of the MPC server. This is a comma separated list of URLs of the MPC servers that will be the participants of the MPC task.
- `VACCINATION_BASE_URL`: **(Required)** The base URL of the vaccination server. This is the URL of the vaccination server that will be used to get the vaccination data.

# Testing

In order the test the whole setup, a mock vaccination server is provided.
A `docker-compose.yml` file is provided to run the whole setup.
```bash
docker compose up -d
```

Once the setup is running, you can test the whole flow by sending a request to the measles server.
```bash
curl -X POST http://localhost:3001/measles-vaccination-check/schedule -H "Content-Type: application/json" -d '{"requestId": "123e4567-e89b-12d3-a456-426614174000", "fileStateIds": ["123e4567-e89b-12d3-a456-426614174000"]}'
```

Once the request is sent, you can check the result by sending a request to the measles server.
```bash
curl http://localhost:3001/measles-vaccination-check/result/123e4567-e89b-12d3-a456-426614174000
```
