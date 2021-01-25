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
    label: 'RediSearch',
    description: 'Secondary Index & Query Engine',
    value: QueryTypeValue.SEARCH,
  },
  {
    label: 'RedisTimeSeries',
    description: 'Time Series data structure',
    value: QueryTypeValue.TIMESERIES,
  },
  {
    label: 'Command-line interface (CLI)',
    description: 'Be mindful, not all commands are supported',
    value: QueryTypeValue.CLI,
  },
];
