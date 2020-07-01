# Grafana Redis Datasource

[![Grafana 7](https://img.shields.io/badge/Grafana-7-red)](https://www.grafana.com)
[![Radix](https://img.shields.io/badge/Radix-powered-blue)](https://github.com/mediocregopher/radix)
[![RedisTimeSeries](https://img.shields.io/badge/RedisTimeSeries-inspired-yellowgreen)](https://oss.redislabs.com/redistimeseries/)
[![Redis Enterprise](https://img.shields.io/badge/Redis%20Enterprise-supported-orange)](https://redislabs.com/redis-enterprise/)

## What is the Grafana Redis Datasource?

The Grafana Redis Datasource, is a plugin that allows users to connect to Redis database and build dashboards in Grafana to easily monitor Redis data. It provides out-of-the box predefined dashboards - but the plugin allows to build entirely customized dashboards, tuned to your needs.

## What is Grafana?

If you are not familiar with Grafana yet, it is a very popular tool used to build dashboards allowing to monitor applications, infrastructures and any kind of software components.

## What Grafana version is supported?

Only Grafana 7.0 and later with a new plugin platform supported.

## Does this datasource require anything special configured on the Redis databases?

Datasource can connect to any Redis database.

## Redis Monitoring dashboard

To demonstrate datasource functionality we included Redis monitoring dashboard.

![Dashboard](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/images/redis-dashboard.png)

## Build and register plugin

To learn how to build Redis Datasource and register in new or existing Grafana please take a look at [BUILD](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/BUILD.md)

## Run using `docker-compose`

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

### Start Redis and Grafana

```bash
docker-compose up
```

### Open Grafana in your browser [http://localhost:3000](http://localhost:3000) and configure datasource

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
