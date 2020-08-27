# Redis Data Source for Grafana

![Dashboard](https://raw.githubusercontent.com/RedisTimeSeries/grafana-redis-datasource/master/src/img/redis-dashboard.png)

[![Grafana 7](https://img.shields.io/badge/Grafana-7-red)](https://www.grafana.com)
[![Radix](https://img.shields.io/badge/Radix-powered-blue)](https://github.com/mediocregopher/radix)
[![RedisTimeSeries](https://img.shields.io/badge/RedisTimeSeries-inspired-yellowgreen)](https://oss.redislabs.com/redistimeseries/)
[![Redis Enterprise](https://img.shields.io/badge/Redis%20Enterprise-supported-orange)](https://redislabs.com/redis-enterprise/)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedisTimeSeries/grafana-redis-datasource)](https://goreportcard.com/report/github.com/RedisTimeSeries/grafana-redis-datasource)
[![CircleCI](https://circleci.com/gh/RedisTimeSeries/grafana-redis-datasource.svg?style=svg)](https://circleci.com/gh/RedisTimeSeries/grafana-redis-datasource)
[![Downloads](https://img.shields.io/badge/dynamic/json?color=green&label=downloads&query=%24.downloads&url=https%3A%2F%2Fgrafana.com%2Fapi%2Fplugins%2Fredis-datasource)](https://grafana.com/grafana/plugins/redis-datasource)

## Summary

- [**Introduction**](#introduction)
- [**Getting Started**](#getting-started)
- [**Supported Commands**](#supported-commands)
- [**Template variables**](#templates-variables)
- [**Feedback**](#feedback)
- [**Contributing**](#contributing)
- [**License**](#license)

## Introduction

### What is the Redis Data Source for Grafana?

If you’re not familiar with Grafana, it’s a very popular tool used to build dashboards to monitor applications, infrastructures, and software components. The Redis Data Source for Grafana is a plug-in that allows users to connect to the Redis database and build dashboards in Grafana to easily monitor Redis data. It provides an out-of-the-box predefined dashboard, but also lets you build customized dashboards tuned to your specific needs.

### What Grafana version is supported?

Grafana 7.1 and later with a new plug-in platform supported.

### Does this Data Source require anything special configured on the Redis databases?

Data Source can connect to any Redis database. No special configuration is required.

### Does this Data Source support [Redis Cluster](https://redis.io/topics/cluster-tutorial) and [Sentinel](https://redis.io/topics/sentinel)?

Yes, Redis Cluster and Sentinel supported since version 1.2.

### How to connect to Redis logical database?

Please use `/db-number` or `?db=db-number` in the Data Source URL to specify the database number as defined in the [Schema](https://www.iana.org/assignments/uri-schemes/prov/redis).

```
redis://redis-server:6379/0
```

### Build Data Source

To learn step by step how to build Redis Data Source from scratch and register in new or existing Grafana please take a look at [BUILD](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/BUILD.md) instructions.

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

Project provides `docker-compose.yml` to start Redis with RedisTimeSeries module and Grafana.

```bash
docker-compose up
```

Open Grafana in your browser and configure Redis Data Source. You can add as many data sources as you want to support multiple Redis databases.

![Datasource](https://raw.githubusercontent.com/RedisTimeSeries/grafana-redis-datasource/master/src/img/datasource.png)

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

Data Source have to be built following [BUILD](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/BUILD.md) instructions before starting using `docker-compose-dev.yml` file.

```bash
docker-compose -f docker-compose-dev.yml up
```

## Supported Commands

Data Source supports many Redis commands using custom components and provide unified interface to query any command.

![Query](https://raw.githubusercontent.com/RedisTimeSeries/grafana-redis-datasource/master/src/img/query.png)

## Template variables

Template variables can query any command and use other variables as parameters.

![Variables](https://raw.githubusercontent.com/RedisTimeSeries/grafana-redis-datasource/master/src/img/variables.png)

## Feedback

We love to hear from users, developers and the whole community interested by this plugin. These are various ways to get in touch with us:

- Ask a question, request a new feature and file a bug with [GitHub issues](https://github.com/RedisTimeSeries/grafana-redis-datasource/issues/new/choose).
- Star the repository to show your support.

## Contributing

- Fork the repository.
- Find an issue to work on and submit a pull request
  - Pick a [good first issue](https://github.com/RedisTimeSeries/grafana-redis-datasource/labels/good%20first%20issue).
- Could not find an issue? Look for documentation, bugs, typos, and missing features.

## Other interesting resources

- [RedisTimeSeries](https://oss.redislabs.com/redistimeseries/)
- [Redis Pop-up store](https://github.com/RedisTimeSeries/redis-pop-up-store)

## License

- Apache License Version 2.0, see [LICENSE](https://github.com/RedisTimeSeries/grafana-redis-datasource/blob/master/LICENSE)
