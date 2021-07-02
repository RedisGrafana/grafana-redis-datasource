import { RedisGearsCommands } from './gears';
import { RedisGraph, RedisGraphCommands } from './graph';
import { QueryTypeValue } from './query';
import { RedisTimeSeries, RedisTimeSeriesCommands } from './time-series';

/**
 * Commands
 */
export const Commands = {
  [QueryTypeValue.REDIS]: [
    {
      label: 'CLIENT LIST',
      description: 'Returns information and statistics about the client connections server',
      value: 'clientList',
    },
    {
      label: 'CLUSTER INFO',
      description: 'Provides INFO style information about Redis Cluster vital parameters',
      value: 'clusterInfo',
    },
    {
      label: 'CLUSTER NODES',
      description: 'Provides current cluster configuration, given by the set of known nodes',
      value: 'clusterNodes',
    },
    {
      label: 'GET',
      description: 'Returns the value of key',
      value: 'get',
    },
    { label: 'HGET', description: 'Returns the value associated with field in the hash stored at key', value: 'hget' },
    { label: 'HGETALL', description: 'Returns all fields and values of the hash stored at key', value: 'hgetall' },
    { label: 'HKEYS', description: 'Returns all field names in the hash stored at key', value: 'hkeys' },
    { label: 'HLEN', description: 'Returns the number of fields contained in the hash stored at key', value: 'hlen' },
    {
      label: 'HMGET',
      description: 'Returns the values associated with the specified fields in the hash stored at key',
      value: 'hmget',
    },
    { label: 'INFO', description: 'Returns information and statistics about the server', value: 'info' },
    { label: 'LLEN', description: 'Returns the length of the list stored at key', value: 'llen' },
    {
      label: 'SCAN with Type and Memory Usage',
      description: 'Returns keys with types and memory usage (CAUSE LATENCY)',
      value: 'tmscan',
    },
    {
      label: 'SCARD',
      description: 'Returns the set cardinality (number of elements) of the set stored at key',
      value: 'scard',
    },
    {
      label: 'SLOWLOG GET',
      description: 'Returns the Redis slow queries log',
      value: 'slowlogGet',
    },
    { label: 'SMEMBERS', description: 'Returns all the members of the set value stored at key', value: 'smembers' },
    {
      label: 'TTL',
      description: 'Returns the string representation of the type of the value stored at key',
      value: 'ttl',
    },
    {
      label: 'TYPE',
      description: 'Returns the string representation of the type of the value stored at key',
      value: 'type',
    },
    {
      label: 'XINFO STREAM',
      description: 'Returns general information about the stream stored at the specified key',
      value: 'xinfoStream',
    },
    {
      label: 'XLEN',
      description: 'Returns the number of entries inside a stream',
      value: 'xlen',
    },
    {
      label: 'XRANGE',
      description: 'Returns the stream entries matching a given range of IDs',
      value: 'xrange',
    },
    {
      label: 'XREVRANGE',
      description: 'Returns the stream entries matching a given range of IDs in reverse order',
      value: 'xrevrange',
    },
  ],
  [QueryTypeValue.TIMESERIES]: RedisTimeSeriesCommands,
  [QueryTypeValue.SEARCH]: [
    {
      label: 'FT.INFO',
      description: 'Returns information and statistics on the index',
      value: 'ft.info',
    },
  ],
  [QueryTypeValue.GEARS]: RedisGearsCommands,
  [QueryTypeValue.GRAPH]: RedisGraphCommands,
};

/**
 * Input for Commands
 */
export const CommandParameters = {
  aggregation: [RedisTimeSeries.RANGE, RedisTimeSeries.MRANGE],
  field: ['hget', 'hmget'],
  filter: [RedisTimeSeries.MRANGE, RedisTimeSeries.QUERYINDEX, RedisTimeSeries.MGET],
  keyName: [
    'get',
    'hget',
    'hgetall',
    'hkeys',
    'hlen',
    'hmget',
    'llen',
    'scard',
    'smembers',
    RedisTimeSeries.RANGE,
    RedisTimeSeries.GET,
    RedisTimeSeries.INFO,
    'ttl',
    'type',
    'xinfoStream',
    'xlen',
    'ft.info',
    'xrange',
    'xrevrange',
    RedisGraph.QUERY,
    RedisGraph.SLOWLOG,
    RedisGraph.EXPLAIN,
    RedisGraph.PROFILE,
  ],
  legend: [RedisTimeSeries.RANGE],
  legendLabel: [RedisTimeSeries.MRANGE, RedisTimeSeries.MGET],
  section: ['info'],
  value: [RedisTimeSeries.RANGE],
  valueLabel: [RedisTimeSeries.MRANGE, RedisTimeSeries.MGET],
  fill: [RedisTimeSeries.RANGE, RedisTimeSeries.MRANGE],
  size: ['slowlogGet', 'tmscan'],
  cursor: ['tmscan'],
  match: ['tmscan'],
  count: ['tmscan', 'xrange', 'xrevrange'],
  samples: ['tmscan'],
  start: ['xrange', 'xrevrange'],
  end: ['xrange', 'xrevrange'],
  cypher: [RedisGraph.EXPLAIN, RedisGraph.QUERY, RedisGraph.PROFILE],
};
