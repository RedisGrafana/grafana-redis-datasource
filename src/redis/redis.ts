import { SelectableValue } from '@grafana/data';

/**
 * Supported Commands
 */
export enum Redis {
  CLIENT_LIST = 'clientList',
  CLUSTER_INFO = 'clusterInfo',
  CLUSTER_NODES = 'clusterNodes',
  GET = 'get',
  HGET = 'hget',
  HGETALL = 'hgetall',
  HKEYS = 'hkeys',
  HLEN = 'hlen',
  HMGET = 'hmget',
  INFO = 'info',
  LLEN = 'llen',
  TMSCAN = 'tmscan',
  SCARD = 'scard',
  SLOWLOG_GET = 'slowlogGet',
  SMEMBERS = 'smembers',
  TTL = 'ttl',
  TYPE = 'type',
  ZRANGE = 'zrange',
  XINFO_STREAM = 'xinfoStream',
  XLEN = 'xlen',
  XRANGE = 'xrange',
  XREVRANGE = 'xrevrange',
}

/**
 * Commands List
 */
export const RedisCommands = [
  {
    label: 'CLIENT LIST',
    description: 'Returns information and statistics about the client connections server',
    value: Redis.CLIENT_LIST,
  },
  {
    label: 'CLUSTER INFO',
    description: 'Provides INFO style information about Redis Cluster vital parameters',
    value: Redis.CLUSTER_INFO,
  },
  {
    label: 'CLUSTER NODES',
    description: 'Provides current cluster configuration, given by the set of known nodes',
    value: Redis.CLUSTER_NODES,
  },
  {
    label: Redis.GET.toUpperCase(),
    description: 'Returns the value of key',
    value: Redis.GET,
  },
  {
    label: Redis.HGET.toUpperCase(),
    description: 'Returns the value associated with field in the hash stored at key',
    value: Redis.HGET,
  },
  {
    label: Redis.HGETALL.toUpperCase(),
    description: 'Returns all fields and values of the hash stored at key',
    value: Redis.HGETALL,
  },
  {
    label: Redis.HKEYS.toUpperCase(),
    description: 'Returns all field names in the hash stored at key',
    value: Redis.HKEYS,
  },
  {
    label: Redis.HLEN.toUpperCase(),
    description: 'Returns the number of fields contained in the hash stored at key',
    value: Redis.HLEN,
  },
  {
    label: Redis.HMGET.toUpperCase(),
    description: 'Returns the values associated with the specified fields in the hash stored at key',
    value: Redis.HMGET,
  },
  {
    label: Redis.INFO.toUpperCase(),
    description: 'Returns information and statistics about the server',
    value: Redis.INFO,
  },
  { label: Redis.LLEN.toUpperCase(), description: 'Returns the length of the list stored at key', value: Redis.LLEN },
  {
    label: Redis.TMSCAN.toUpperCase(),
    description: 'Returns keys with types and memory usage (CAUSE LATENCY)',
    value: Redis.TMSCAN,
  },
  {
    label: Redis.SCARD.toUpperCase(),
    description: 'Returns the set cardinality (number of elements) of the set stored at key',
    value: Redis.SCARD,
  },
  {
    label: 'SLOWLOG GET',
    description: 'Returns the Redis slow queries log',
    value: Redis.SLOWLOG_GET,
  },
  {
    label: Redis.SMEMBERS.toUpperCase(),
    description: 'Returns all the members of the set value stored at key',
    value: Redis.SMEMBERS,
  },
  {
    label: Redis.TTL.toUpperCase(),
    description: 'Returns the string representation of the type of the value stored at key',
    value: Redis.TTL,
  },
  {
    label: Redis.TYPE.toUpperCase(),
    description: 'Returns the string representation of the type of the value stored at key',
    value: Redis.TYPE,
  },
  {
    label: Redis.ZRANGE.toUpperCase(),
    description: 'Returns the specified range of elements in the sorted set at key',
    value: Redis.ZRANGE,
  },
  {
    label: 'XINFO STREAM',
    description: 'Returns general information about the stream stored at the specified key',
    value: Redis.XINFO_STREAM,
  },
  {
    label: Redis.XLEN.toUpperCase(),
    description: 'Returns the number of entries inside a stream',
    value: Redis.XLEN,
  },
  {
    label: Redis.XRANGE.toUpperCase(),
    description: 'Returns the stream entries matching a given range of IDs',
    value: Redis.XRANGE,
  },
  {
    label: Redis.XREVRANGE.toUpperCase(),
    description: 'Returns the stream entries matching a given range of IDs in reverse order',
    value: Redis.XREVRANGE,
  },
];

/**
 * ZRANGE Query Values
 */
export enum ZRangeQueryValue {
  BYINDEX = '',
  BYSCORE = 'BYSCORE',
  BYLEX = 'BYLEX',
}

/**
 * Aggregations
 */
export const ZRangeQuery: Array<SelectableValue<ZRangeQueryValue>> = [
  {
    label: 'Index range',
    description:
      'The <min> and <max> arguments represent zero-based indexes, where 0 is the first element, 1 is the next element, and so on.',
    value: ZRangeQueryValue.BYINDEX,
  },
  {
    label: 'Score range',
    description: 'Returns the range of elements from the sorted set having scores equal or between <min> and <max>',
    value: ZRangeQueryValue.BYSCORE,
  },
];
