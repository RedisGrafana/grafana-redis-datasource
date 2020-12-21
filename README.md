# Redis Data Source for Grafana

![Dashboard](https://raw.githubusercontent.com/RedisGrafana/grafana-redis-datasource/master/src/img/redis-dashboard.png)

[![Grafana 7](https://img.shields.io/badge/Grafana-7-red)](https://www.grafana.com)
[![Radix](https://img.shields.io/badge/Radix-powered-blue)](https://github.com/mediocregopher/radix)
[![Redis Enterprise](https://img.shields.io/badge/Redis%20Enterprise-supported-orange)](https://redislabs.com/redis-enterprise/)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedisGrafana/grafana-redis-datasource)](https://goreportcard.com/report/github.com/RedisGrafana/grafana-redis-datasource)
![CI](https://github.com/RedisGrafana/grafana-redis-datasource/workflows/CI/badge.svg)
[![Downloads](https://img.shields.io/badge/dynamic/json?color=green&label=downloads&query=%24.downloads&url=https%3A%2F%2Fgrafana.com%2Fapi%2Fplugins%2Fredis-datasource)](https://grafana.com/grafana/plugins/redis-datasource)

## Summary

- [**Introduction**](#introduction)
- [**Getting Started**](#getting-started)
- [**Supported commands**](#supported-commands)
- [**Template variables**](#template-variables)
- [**Learn more**](#learn-more)
- [**Feedback**](#feedback)
- [**Contributing**](#contributing)
- [**License**](#license)

## Introduction

### What is the Redis Data Source for Grafana?

The Redis Data Source for Grafana is a plug-in that allows users to connect to the Redis database and build dashboards in Grafana to easily monitor Redis and application data. It provides an out-of-the-box predefined dashboard, but also lets you build customized dashboards tuned to your specific needs.

### What Grafana version is supported?

Grafana 7.1 and later with a new plug-in platform supported.

### Does this Data Source require anything special configured on the Redis databases?

Data Source can connect to any Redis database. No special configuration is required.

### Does this Data Source support [Redis Cluster](https://redis.io/topics/cluster-tutorial) and [Sentinel](https://redis.io/topics/sentinel)?

Redis Cluster and Sentinel supported since version 1.2.

### Does this Data Source support Redis modules?

Data Source supports [RedisTimeSeries](https://oss.redislabs.com/redistimeseries/). We are working on supporting additional modules.

### How to connect to Redis logical database

Please use `/db-number` or `?db=db-number` in the Data Source URL to specify the database number as defined in the [Schema](https://www.iana.org/assignments/uri-schemes/prov/redis).

```
redis://redis-server:6379/0
```

### How to build Data Source

To learn how to build Redis Data Source from scratch and register in new or existing Grafana please take a look at [BUILD](https://github.com/RedisGrafana/grafana-redis-datasource/blob/master/BUILD.md) instructions.

## Getting Started

### Install using `grafana-cli`

Use the `grafana-cli` tool to install from the commandline:

```bash
grafana-cli plugins install redis-datasource
```

### Run using `docker`

```bash
docker run -d -p 3000:3000 --name=grafana -e "GF_INSTALL_PLUGINS=redis-datasource" grafana/grafana
```

### Run using `docker-compose`

Project provides `docker-compose.yml` to start Redis with all Redis Labs modules and Grafana.

```bash
docker-compose up
```

Open Grafana in your browser and configure Redis Data Source. You can add as many data sources as you want to support multiple Redis databases.

![Datasource](https://raw.githubusercontent.com/RedisGrafana/grafana-redis-datasource/master/src/img/datasource.png)

There are certain settings that can be configured based on your own setup:

- Grafana port
- Data Source URL

#### Configure Grafana port in `docker-compose.yml`

If standard port 3000 is occupied by another application update the port to bind Grafana to

```
    ports:
      - '3000:3000'
```

#### Configure Data Source URL in `provisioning/datasources/redis.yaml`

If Redis is running and listening on localhost:6379 no changes are required

```
    url: redis://localhost:6379
```

If Redis is running as Docker container on MacOS, please update host to `host.docker.internal`

```
    url: redis://host.docker.internal:6379
```

### Run using `docker-compose` for development

Data Source have to be built following [BUILD](https://github.com/RedisGrafana/grafana-redis-datasource/blob/master/BUILD.md) instructions before starting using `docker-compose-dev.yml` file.

```bash
docker-compose -f docker-compose-dev.yml up
```

## Supported commands

Data Source supports various Redis commands using custom components and provides a unified interface to query any command.

![Query](https://raw.githubusercontent.com/RedisGrafana/grafana-redis-datasource/master/src/img/query.png)

## Template variables

Template variables can query any command and use other variables as parameters.

![Variables](https://raw.githubusercontent.com/RedisGrafana/grafana-redis-datasource/master/src/img/variables.png)

## Learn more

- [Introducing the Redis Data Source Plug-in for Grafana](https://redislabs.com/blog/introducing-the-redis-data-source-plug-in-for-grafana/)
- [How to Use the New Redis Data Source for Grafana Plug-in](https://redislabs.com/blog/how-to-use-the-new-redis-data-source-for-grafana-plug-in/)
- [3 Real-Life Apps Built with Redis Data Source for Grafana](https://redislabs.com/blog/3-real-life-apps-built-with-redis-data-source-for-grafana/)

## Feedback

We love to hear from users, developers and the whole community interested by this plugin. These are various ways to get in touch with us:

- Ask a question, request a new feature and file a bug with [GitHub issues](https://github.com/RedisGrafana/grafana-redis-datasource/issues/new/choose).
- Star the repository to show your support.

## Contributing

- Fork the repository.
- Find an issue to work on and submit a pull request.
- Could not find an issue? Look for documentation, bugs, typos, and missing features.

## License

- Apache License Version 2.0, see [LICENSE](https://github.com/RedisGrafana/grafana-redis-datasource/blob/master/LICENSE).
