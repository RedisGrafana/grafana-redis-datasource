# Grafana Redis Datasource

[![Grafana 7](https://img.shields.io/badge/Grafana-7-red)](https://www.grafana.com)
[![Radix](https://img.shields.io/badge/Radix-powered-blue)](https://github.com/mediocregopher/radix)
[![RedisTimeSeries](https://img.shields.io/badge/RedisTimeSeries-inspired-yellowgreen)](https://oss.redislabs.com/redistimeseries/)
[![Redis Enterprise](https://img.shields.io/badge/Redis%20Enterprise-supported-orange)](https://redislabs.com/redis-enterprise/)

## Description

Redis datasource for Grafana 7.0 allows to query data directly from Redis database using [Radix](https://github.com/mediocregopher/radix) client. No additional adapters is required.

### Redis Monitoring dashboard

To demonstrate datasource functionality we included Redis monitoring dashboard.

![Dashboard](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/redis-dashboard.png)

## Build

Redis datasource consists of both frontend and backend components.

### React Frontend

#### Install dependencies

```bash
yarn install
```

#### Build frontend

```bash
npm run build
```

### Golang Backend

#### Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency

```bash
go get -u github.com/grafana/grafana-plugin-sdk-go
```

#### Build backend plugin binaries for Linux, Windows and MacOS

```bash
npm run build:backend
```

## Run

Project provides `docker-compose.yml` to start Redis with RedisTimeSeries module and Grafana 7.0

### Grafana port in `docker-compose.yml`

If standard port 3000 is occupied by another application update the port to bind Grafana to

```
    ports:
      - '3000:3000'
```

### Datasource url in `provisioning/datasources/redis.yaml`

If Redis is running and listening on localhost:6379 no changes are required

```
    url: redis://localhost:6379
```

If Redis is running as Docker container on MacOS, please update host to `host.docker.internal`

```
    url: redis://host.docker.internal:6379
```

### Start using `docker-compose`

```bash
npm run start
```

## Open Grafana in your browser [http://localhost:3000](http://localhost:3000) and configure datasource

You can add as many datasources as you want to support multiple Redis databases.

![Datasource](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/datasource.png)

## Supported Commands

Datasource supports many Redis commands using custom components and provide unified interface to query any command.

![Query](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/query.png)

## Template variables

Template variables can query any command and use other variables as parameters.

![Variables](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/variables.png)

## Learn more

- [Redis TimeSeries](https://oss.redislabs.com/redistimeseries/)

## License

Apache License Version 2.0, see [LICENSE](LICENSE)
