import { SelectableValue } from '@grafana/data';

/**
 * Query Type Values
 */
export enum QueryTypeValue {
  REDIS = 'command',
  TIMESERIES = 'timeSeries',
  SEARCH = 'search',
  CLI = 'cli',
  GEARS = 'gears',
  GRAPH = 'graph',
  JSON = 'json',
}

/**
 * Query Type
 */
export const QueryType: Array<SelectableValue<QueryTypeValue>> = [
  {
    label: 'Redis',
    description: 'Hashes, Sets, Lists, Strings, Streams, etc.',
    value: QueryTypeValue.REDIS,
  },
  {
    label: 'RedisGears',
    description: 'Dynamic framework for data processing',
    value: QueryTypeValue.GEARS,
  },
  {
    label: 'RedisJSON',
    description: 'JSON data type for Redis',
    value: QueryTypeValue.JSON,
  },
  {
    label: 'RedisGraph',
    description: 'Graph database',
    value: QueryTypeValue.GRAPH,
  },
  {
    label: 'RediSearch',
    description: 'Secondary Index & Query Engine',
    value: QueryTypeValue.SEARCH,
  },
  {
    label: 'RedisTimeSeries',
    description: 'Time Series data structure',
    value: QueryTypeValue.TIMESERIES,
  },
];

/**
 * Query Type for Command-line interface
 */
export const QueryTypeCli: SelectableValue<QueryTypeValue> = {
  label: 'Command-line interface (CLI)',
  description: 'Be mindful, not all commands are supported',
  value: QueryTypeValue.CLI,
};
