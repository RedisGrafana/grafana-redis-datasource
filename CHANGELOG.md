# Change Log

## v1.3.0

- Update description and GitHub issues #83
- Add RediSearch FT.INFO command #97
- Add HMGET Command #98
- Update release workflow #99
- Update Grafana dependencies to 7.3.5 #100
- Update Grafana SDK 0.80.0 #101
- Update data source icon and refactoring #102
- Update field's name for HGET command to align with HMGET #103
- Update HGETALL command to return fields and support streaming similar to HGET, HMGET #104
- Add tests for React Config and Query editors #105

## v1.2.1

- Support Connecting to Redis via Unix Socket #58
- Support Redis 6 ACL authentication #60
- Update Grafana dependencies to 7.2.0 #66
- Update and optimize dashboards for Grafana 7.2.0 #67
- Add Streaming for Command Statistics #68
- Add Size parameter for SLOWLOG GET #79
- Update GitHub org to RedisGrafana #80
- Plugin health check failed for ARM on Linux #61
- Timeseries data time stamp truncated to seconds #64

## v1.2.0

- Added docker cmd line option to start in README #31
- How to query a specific database inside the same Redis single node #34
- Update docker-compose to load datasource from the repository and add development file #39
- Use "ScopedVars" when applying template variables #37
- Refactoring to support new commands and modules #42
- Add support for TS.GET, TS.INFO, and TS.QUERYINDEX commands #45
- Add Redis dashboard to support multiple Redis instances #49
- Plugin executable missing for arm64 architecture #48 (Grafana SDK: grafana/grafana-plugin-sdk-go#221)
- Return 0 for all buckets with 0 counts on time-series TS.RANGE queries #50
- Add Redis Cluster support and update monitoring dashboard #52
- Connection issue to Redis deployed in k8s (Sentinel) #38
- MRANGE: add fill zero option #53
- Add Streaming capabilities to visualize INFO command #57
- Slowlog returns 'No data' for Redis 3.0.6 #33
- Fix backend lint issues #41
- ts.mrange returns no data when label has spaces within #44

## v1.1.2

- Remove developer jargon from README #30
- Redis Datasource is Unsigned. K8S+Helm installation #29

## v1.1.1

- Screenshots added to plugin.json and updated in the README
- CHANGELOG added to display on the Plugin page

## v1.1.0

- Updated to Grafana 7.1.0 and the latest version of Radix #27
- Add dashboard as a part of datasource #25
- Add Field config units to the response #26

## v1.0.0

- Initial release based on Grafana 7.0.5.
- Allows configuring password, TLS, and advanced settings.
- Supports Redis commands: CLIENT LIST, GET, HGET, HGETALL, HKEYS, HLEN, INFO, LLEN, SCARD, SLOWLOG GET, SMEMBERS, TTL, TYPE, XLEN.
- Supports RedisTimeSeries commands: TS.MRANGE, TS.RANGE.
- Provides Redis monitoring dashboard.
