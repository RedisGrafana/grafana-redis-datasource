import { RedisGears, RedisGearsCommands } from './gears';
import { RedisGraph, RedisGraphCommands } from './graph';
import { RedisJson, RedisJsonCommands } from './json';
import { QueryTypeValue } from './query';
import { Redis, RedisCommands } from './redis';
import { RediSearch, RediSearchCommands } from './search';
import { RedisTimeSeries, RedisTimeSeriesCommands } from './time-series';

/**
 * Commands
 */
export const Commands = {
  [QueryTypeValue.REDIS]: RedisCommands,
  [QueryTypeValue.TIMESERIES]: RedisTimeSeriesCommands,
  [QueryTypeValue.SEARCH]: RediSearchCommands,
  [QueryTypeValue.GEARS]: RedisGearsCommands,
  [QueryTypeValue.GRAPH]: RedisGraphCommands,
  [QueryTypeValue.JSON]: RedisJsonCommands,
};

/**
 * Input for Commands
 */
export const CommandParameters = {
  aggregation: [RedisTimeSeries.RANGE, RedisTimeSeries.MRANGE],
  field: [Redis.HGET, Redis.HMGET],
  filter: [RedisTimeSeries.MRANGE, RedisTimeSeries.QUERYINDEX, RedisTimeSeries.MGET],
  keyName: [
    Redis.GET,
    Redis.HGET,
    Redis.HGETALL,
    Redis.HKEYS,
    Redis.HLEN,
    Redis.HMGET,
    Redis.LLEN,
    Redis.SCARD,
    Redis.SMEMBERS,
    RedisTimeSeries.RANGE,
    RedisTimeSeries.GET,
    RedisTimeSeries.INFO,
    Redis.TTL,
    Redis.TYPE,
    Redis.XINFO_STREAM,
    Redis.XLEN,
    RediSearch.INFO,
    Redis.XRANGE,
    Redis.XREVRANGE,
    RedisGraph.QUERY,
    RedisGraph.SLOWLOG,
    RedisGraph.EXPLAIN,
    RedisGraph.PROFILE,
    Redis.ZRANGE,
    RedisJson.TYPE,
    RedisJson.GET,
    RedisJson.OBJKEYS,
    RedisJson.OBJLEN,
    RedisJson.ARRLEN,
  ],
  legend: [RedisTimeSeries.RANGE],
  legendLabel: [RedisTimeSeries.MRANGE, RedisTimeSeries.MGET],
  section: [Redis.INFO],
  value: [RedisTimeSeries.RANGE],
  valueLabel: [RedisTimeSeries.MRANGE, RedisTimeSeries.MGET],
  fill: [RedisTimeSeries.RANGE, RedisTimeSeries.MRANGE],
  size: [Redis.SLOWLOG_GET, Redis.TMSCAN],
  cursor: [Redis.TMSCAN],
  match: [Redis.TMSCAN],
  count: [Redis.TMSCAN, Redis.XRANGE, Redis.XREVRANGE],
  samples: [Redis.TMSCAN],
  min: [Redis.ZRANGE],
  max: [Redis.ZRANGE],
  start: [Redis.XRANGE, Redis.XREVRANGE],
  end: [Redis.XRANGE, Redis.XREVRANGE],
  cypher: [RedisGraph.EXPLAIN, RedisGraph.QUERY, RedisGraph.PROFILE],
  zrangeQuery: [Redis.ZRANGE],
  path: [RedisJson.TYPE, RedisJson.OBJKEYS, RedisJson.GET, RedisJson.OBJLEN, RedisJson.ARRLEN],
  pyFunction: [RedisGears.PYEXECUTE],
};
