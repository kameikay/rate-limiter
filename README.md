# Rate limiter

## Description
The Rate Limiter is a project that provides a mechanism to limit the rate of incoming requests to a system or API, using as a middlware. It helps prevent abuse, protect resources, and ensure fair usage. This project uses Redis to store the rate limit data, but you can easily modify it to use another storage system, respecting the same interface.

## How It Works
The Rate Limiter works by tracking the number of requests made within a specific time window. It enforces a predefined limit on the number of requests that can be made within that window. If the limit is exceeded, the Rate Limiter will respond a HTTP status code 429, with a message "you have reached the maximum number of requests or actions allowed within a certain time frame". It use Redis sorted sets to store the rate limit data, and it uses a sliding window algorithm to track the number of requests made within the time window.

## Configuration
To configure the Rate Limiter, you can modify the following parameters on the .env file:

- **REDIS_ADDR**: The address of the Redis server used to store the rate limit data. (e.g. localhost:6379) - string
- **REDIS_DB**: The database number used to store the rate limit data. (e.g. 0) - integer
- **REDIS_PASSWORD**: The password used to authenticate with the Redis server. - string
- **RATE_LIMITER_DEFAULT_MAX_REQUESTS**: The default maximum number of requests allowed within the time window. - integer
- **RATE_LIMITER_DEFAULT_BLOCK_TIME**: The default block time in milliseconds for users who exceed the rate limit. - integer
- **RATE_LIMITER_CUSTOM_TOKENS**: A list of custom tokens with their own rate limit and block time (JSON format). Each token should include the following properties:
  - **Name**: The name of the token used to identify the user or client. (e.g. "user1") - string
  - **Max Requests**: The maximum number of requests allowed within one second. - integer
  - **Block Time**: The block time in milliseconds for users who exceed the rate limit. - integer

## How to play with it
1. Clone the repository
2. Run `docker-compose up` to start the Redis server and the Rate Limiter with air
3. Make a request to the Rate Limiter using the following endpoint: `http://localhost:8080/`
4. You can modify the rate limit configuration by editing the .env file and restarting the Rate Limiter

## Stress Test
To stress test the Rate Limiter, you can use [fortio](https://fortio.org/) to send a large number of requests to the Rate Limiter and observe how it handles the rate limit.

```bash

# Getting blocked for IP access when RATE_LIMITER_DEFAULT_MAX_REQUESTS=7
docker run --rm --network=rate-limiter fortio/fortio load -qps 8 -c 1 -t 1s http://rate-limiter:8080

# Getting blocked for random token access
docker run --rm --network=rate-limiter fortio/fortio load -qps 10 -c 1 -t 1s -H "API_KEY: random" http://rate-limiter:8080

# Getting blocked for "DEF321" token access when RATE_LIMITER_CUSTOM_TOKENS=[{"name":"DEF321","max_requests_per_second":15,"block_time_in_milliseconds":500}]
docker run --rm --network=rate-limiter fortio/fortio load -qps 16 -c 1 -t 1s -H "API_KEY: DEF321" http://rate-limiter:8080

```