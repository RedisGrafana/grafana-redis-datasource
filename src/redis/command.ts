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
    { label: 'HMGET', description: 'Returns the values associated with the specified fields in the hash stored at key', value: 'hmget' },
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
    /**
     * {
     *   label: 'XINFO STREAM',
     *   description: 'Returns general information about the stream stored at the specified key',
     *   value: 'xinfoStream',
     * },
     */
    {
      label: 'XLEN',
      description: 'Returns the number of entries inside a stream',
      value: 'xlen',
    },
  ],
  timeSeries: [
    {
      label: 'TS.GET',
      description: 'Returns the last sample',
      value: 'ts.get',
    },
    {
      label: 'TS.INFO',
      description: 'Returns information and statistics on the time-series',
      value: 'ts.info',
    },
    {
      label: 'TS.MRANGE',
      description: 'Query a timestamp range across multiple time-series by filters',
      value: 'ts.mrange',
    },
    {
      label: 'TS.QUERYINDEX',
      description: 'Query all the keys matching the filter list',
      value: 'ts.queryindex',
    },
    { label: 'TS.RANGE', description: 'Query a range', value: 'ts.range' },
  ],
  search: [
    {
      label: 'FT.INFO',
      description: 'Returns information and statistics on the index',
      value: 'ft.info',
    },
  ],
  gears: [
    {
      label: 'RG.PYSTATS',
      description: 'Returns memory usage statistics from the Python interpreter',
      value: 'rg.pystats',
    },
    {
      label: 'RG.DUMPREGISTRATIONS',
      description: 'Outputs the list of function registrations',
      value: 'rg.dumpregistrations',
    },
  ]
};

/**
 * Input for Commands
 */
export const CommandParameters = {
  aggregation: ['ts.range', 'ts.mrange'],
  field: ['hget','hmget'],
  filter: ['ts.mrange', 'ts.queryindex'],
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
    'ts.range',
    'ts.get',
    'ts.info',
    'ttl',
    'type',
    'xinfoStream',
    'xlen',
    'ft.info'
  ],
  legend: ['ts.range'],
  legendLabel: ['ts.mrange'],
  section: ['info'],
  value: ['ts.range'],
  valueLabel: ['ts.mrange'],
  fill: ['ts.range', 'ts.mrange'],
  size: ['slowlogGet', 'tmscan'],
  cursor: ['tmscan'],
  match: ['tmscan'],
  count: ['tmscan'],
  samples: ['tmscan']
};
