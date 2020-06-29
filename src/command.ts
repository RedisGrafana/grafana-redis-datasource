import { SelectableValue } from '@grafana/data';

/**
 * Query Type Values
 */
export enum QueryTypeValue {
  COMMAND = 'command',
  TIMESERIES = 'timeSeries',
  CLI = 'cli',
}

/**
 * Query Type
 */
export const QueryType: Array<SelectableValue<string>> = [
  {
    label: 'Redis commands',
    description: 'Hashes, Sets, Lists, Strings, Streams, etc.',
    value: QueryTypeValue.COMMAND,
  },
  {
    label: 'RedisTimeSeries commands',
    description: 'Redis Module adding a Time Series data structure to Redis',
    value: QueryTypeValue.TIMESERIES,
  },
  {
    label: 'Command-line interface (CLI)',
    description: 'Be mindful, not all commands are supported',
    value: QueryTypeValue.CLI,
  },
];

/**
 * Commands
 */
export const Commands = {
  command: [
    {
      label: 'CLIENT LIST',
      description: 'Returns information and statistics about the client connections server',
      value: 'clientList',
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
    { label: 'INFO', description: 'Returns information and statistics about the server', value: 'info' },
    { label: 'LLEN', description: 'Returns the length of the list stored at key', value: 'llen' },
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
  ],
  timeSeries: [
    {
      label: 'TS.MRANGE',
      description: 'Query a timestamp range across multiple time-series by filters',
      value: 'ts.mrange',
    },
    { label: 'TS.RANGE', description: 'Query a range', value: 'ts.range' },
  ],
};

/**
 * Input for Commands
 */
export const CommandParameters = {
  aggregation: ['ts.range', 'ts.mrange'],
  field: ['hget'],
  filter: ['ts.mrange'],
  key: [
    'get',
    'hget',
    'hgetall',
    'hkeys',
    'hlen',
    'llen',
    'scard',
    'smembers',
    'ts.range',
    'ttl',
    'type',
    'xinfoStream',
    'xlen',
  ],
  legend: ['ts.range'],
  legendLabel: ['ts.mrange'],
  section: ['info'],
  valueLabel: ['ts.mrange'],
};

/**
 * Aggregations
 */
export const Aggregations: Array<SelectableValue<string>> = [
  { label: 'None', description: 'No aggregation', value: '' },
  { label: 'Avg', description: 'Average', value: 'avg' },
  { label: 'Count', description: 'Count number of samples', value: 'count' },
  { label: 'Max', description: 'Maximum', value: 'max' },
  { label: 'Min', description: 'Minimum', value: 'min' },
  { label: 'Range', description: 'Diff between maximum and minimum in the bucket', value: 'range' },
  { label: 'Sum', description: 'Sum', value: 'sum' },
];

/**
 * Info sections
 */
export const InfoSections: Array<SelectableValue<string>> = [
  { label: 'Server', description: 'General information about the Redis server', value: 'server' },
  { label: 'Clients', description: 'Client connections section', value: 'clients' },
  { label: 'Memory', description: 'Memory consumption related information', value: 'memory' },
  { label: 'Persistence', description: 'RDB and AOF related information', value: 'persistence' },
  { label: 'Stats', description: 'General statistics', value: 'stats' },
  { label: 'Replication', description: 'Master/replica replication information', value: 'replication' },
  { label: 'CPU', description: 'CPU consumption statistics', value: 'cpu' },
  { label: 'Command Stats', description: 'Redis command statistics', value: 'commandstats' },
  { label: 'Cluster', description: 'Redis Cluster section', value: 'cluster' },
  { label: 'Keyspace', description: 'Database related statistics', value: 'keyspace' },
];
