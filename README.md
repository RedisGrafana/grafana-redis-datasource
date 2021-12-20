# Redis Data Source for Grafana

![Dashboard](https://raw.githubusercontent.com/RedisGrafana/grafana-redis-datasource/master/src/img/redis-dashboard.png)

[![Grafana 8](https://img.shields.io/badge/Grafana-8-orange)](https://www.grafana.com)
[![Redis Data Source](https://img.shields.io/badge/dynamic/json?color=blue&label=Redis%20Data%20Source&query=%24.version&url=https%3A%2F%2Fgrafana.com%2Fapi%2Fplugins%2Fredis-datasource)](https://grafana.com/grafana/plugins/redis-datasource)
[![Redis Application plugin](https://img.shields.io/badge/dynamic/json?color=blue&label=Redis%20Application%20plugin&query=%24.version&url=https%3A%2F%2Fgrafana.com%2Fapi%2Fplugins%2Fredis-app)](https://grafana.com/grafana/plugins/redis-app)
[![Redis Explorer plugin](https://img.shields.io/badge/dynamic/json?color=blue&label=Redis%20Explorer%20plugin&query=%24.version&url=https%3A%2F%2Fgrafana.com%2Fapi%2Fplugins%2Fredis-explorer-app)](https://grafana.com/grafana/plugins/redis-explorer-app)
[![Go Report Card](https://goreportcard.com/badge/github.com/RedisGrafana/grafana-redis-datasource)](https://goreportcard.com/report/github.com/RedisGrafana/grafana-redis-datasource)
![CI](https://github.com/RedisGrafana/grafana-redis-datasource/workflows/CI/badge.svg)
[![codecov](https://codecov.io/gh/RedisGrafana/grafana-redis-datasource/branch/master/graph/badge.svg?token=YX7995RPCF)](https://codecov.io/gh/RedisGrafana/grafana-redis-datasource)
[![Language grade: JavaScript](https://img.shields.io/lgtm/grade/javascript/g/RedisGrafana/grafana-redis-datasource.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/RedisGrafana/grafana-redis-datasource/context:javascript)

## Introduction

The Redis Data Source for Grafana is a plugin that allows users to connect to any Redis database On-Premises and in the Cloud. It provides out-of-the-box predefined dashboards and lets you build customized dashboards to monitor Redis and application data.

### Demo

Demo is available on [demo.volkovlabs.io](https://demo.volkovlabs.io):

- [Redis Overview dashboard](https://demo.volkovlabs.io/d/TgibHBv7z/redis-overview?orgId=1&refresh=1h)
- [Projects](https://demo.volkovlabs.io)

### Requirements

- **Grafana 8.0+** is required for Redis Data Source 2.X.
- **Grafana 7.1+** is required for Redis Data Source 1.X.

### Redis Application plugin

You can add as many data sources as you want to support multiple Redis databases. [Redis Application plugin](https://grafana.com/grafana/plugins/redis-app) helps manage various Redis Data Sources and provides Custom panels.

### Redis Explorer plugin

[The Redis Explorer plugin](https://grafana.com/grafana/plugins/redis-explorer-app) connects to Redis Enterprise software clusters using REST API. It provides application pages to add Redis Data Sources for managed databases and dashboards to see cluster configuration.

## Getting Started

Redis Data Source can be installed from the Grafana Marketplace or use the `grafana-cli` tool to install from the command line:

```bash
grafana-cli plugins install redis-datasource
```

![Grafana Marketplace](https://raw.githubusercontent.com/RedisGrafana/grafana-redis-datasource/master/src/img/grafana-marketplace.png)

For Docker instructions and installation without Internet access, follow the [Quickstart](https://redisgrafana.github.io/quickstart/) page.

### Configuration

Data Source allows to connect to Redis using TCP port, Unix socket, Cluster, Sentinel and supports SSL/TLS authentication. For detailed information, take a look at the [Configuration](https://redisgrafana.github.io/redis-datasource/configuration/) page.

![Datasource](https://raw.githubusercontent.com/RedisGrafana/grafana-redis-datasource/master/src/img/datasource.png)

## Documentation

Please take a look at the [Documentation](https://redisgrafana.github.io/redis-datasource/overview/) to learn more about plugin and features.

### Supported commands

List of all supported commands and how to use them with examples you can find in the [Commands](https://redisgrafana.github.io/redis-datasource/commands/) section.

![Query](https://raw.githubusercontent.com/RedisGrafana/grafana-redis-datasource/master/src/img/query.png)

## Development

[Developing Redis Data Source](https://redisgrafana.github.io/development/redis-datasource/) page provides instructions on building the data source.

Are you interested in the latest features and updates? Start nightly built [Docker image for Redis Application plugin](https://redisgrafana.github.io/development/images/), including Redis Data Source.

## Feedback

We love to hear from users, developers, and the whole community interested in this plugin. These are various ways to get in touch with us:

- Ask a question, request a new feature, and file a bug with [GitHub issues](https://github.com/RedisGrafana/grafana-redis-datasource/issues/new/choose).
- Star the repository to show your support.

## Contributing

- Fork the repository.
- Find an issue to work on and submit a pull request.
- Could not find an issue? Look for documentation, bugs, typos, and missing features.

## License

- Apache License Version 2.0, see [LICENSE](https://github.com/RedisGrafana/grafana-redis-datasource/blob/master/LICENSE).
