# Grafana Redis datasource

[![Grafana 7](https://img.shields.io/badge/Grafana-7-red)](https://www.grafana.com)
[![Radix](https://img.shields.io/badge/Radix-integrated-blue)](https://github.com/mediocregopher/radix)
[![RedisTimeSeries](https://img.shields.io/badge/RedisTimeSeries-inspired-yellowgreen)](https://oss.redislabs.com/redistimeseries/)

## Description

Redis datasource for Grafana 7.0 allows to query data directly from Redis database using [Radix](https://github.com/mediocregopher/radix) client.

To demonstrate it's functionality we included Redis monitoring dashboard:

![Dashboard](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/redis-dashboard.gif)

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

#### Build backend plugin binaries for Linux, Windows and Mac OS

```bash
npm run build:backend
```

## Run

Project consists `docker-compose.yml` to start Redis with RedisTimeSeries module and Grafana 7.0

### Update port for Grafana in `docker-compose.yml`

```
    ports:
      - '3100:3000'
```

### Update datasource `url` in `provisioning/datasources/redis.yaml`

If Redis running on localhost

```
    url: redis://localhost:6379
```

If Redis running in Docker on Mac OS

```
    url: redis://host.docker.internal:6379
```

### Start using `docker-compose`

```bash
npm run start
```

### Open Grafana in your browser [http://localhost:3000](http://localhost:3000)

### Configure Datasource

![Datasource](https://github.com/mikhailredis/grafana-redistimeseries-plugin/blob/master/images/datasource.png)

## Supported Commands

Datasource supports many Redis commands using custom components and provide unified interface to query any commands.

![Query](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/query.png)

## Template variables

Template variables can query any commands and use other variables:

![Variables](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/variables.png)

## Learn more

- [Redis TimeSeries](https://oss.redislabs.com/redistimeseries/)
