# Change Log

## v1.2.0

- Added docker cmd line option to start in README #31
- How to query a specific database inside the same Redis single node #34
- Update docker-compose to load datasource from the repository and add development file #39
- Use "ScopedVars" when applying template variables #37 (fix for #36)
- Slowlog returns 'No data' for Redis 3.0.6 #33
- Fix backend lint issues #41
- Refactoring to support new commands and modules #42

## v1.1.2

- Remove developer jargon from README #30
- Redis Datasource is Unsigned. K8S+Helm installation #29

## v1.1.1

- Screenshots added to plugin.json and updated in the README
- CHANGELOG added to display on the Plugin page

## v1.1.0

- Updated to Grafana 7.1.0 and the latest version of Radix [#27](https://github.com/RedisTimeSeries/grafana-redis-datasource/pull/27)
- Add dashboard as a part of datasource [#25](https://github.com/RedisTimeSeries/grafana-redis-datasource/pull/25)
- Add Field config units to the response [#26](https://github.com/RedisTimeSeries/grafana-redis-datasource/pull/26)

## v1.0.0

- Initial release based on Grafana 7.0.5.
- Allows configuring password, TLS, and advanced settings.
- Supports Redis commands: CLIENT LIST, GET, HGET, HGETALL, HKEYS, HLEN, INFO, LLEN, SCARD, SLOWLOG GET, SMEMBERS, TTL, TYPE, XLEN.
- Supports RedisTimeSeries commands: TS.MRANGE, TS.RANGE.
- Provides Redis monitoring dashboard.
