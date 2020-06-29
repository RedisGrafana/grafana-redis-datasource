import { SelectableValue } from '@grafana/data';

/**
 * Commands
 */
export const Commands: Array<SelectableValue<string>> = [
  { label: 'HGET', description: 'Returns the value associated with field in the hash stored at key', value: 'hget' },
  { label: 'HGETALL', description: 'Returns all fields and values of the hash stored at key', value: 'hgetall' },
  { label: 'SMEMBERS', description: 'Returns all the members of the set value stored at key', value: 'smembers' },
  { label: 'TS.RANGE', description: 'Query a range', value: 'tsrange' },
  {
    label: 'TS.MRANGE',
    description: 'Query a timestamp range across multiple time-series by filters',
    value: 'tsmrange',
  },
];

/**
 * Input for Commands
 */
export const CommandParameters = {
  aggregation: ['tsrange', 'tsmrange'],
  field: ['hget'],
  filter: ['tsmrange'],
  key: ['tsrange', 'hgetall', 'hget', 'smembers'],
  legend: ['tsrange'],
  legendLabel: ['tsmrange'],
  valueLabel: ['tsmrange'],
};

/**
 * Aggregations
 */
export const Aggregations: Array<SelectableValue<string>> = [
  { label: 'None', description: 'no aggregation', value: '' },
  { label: 'Max', description: 'max', value: 'max' },
  { label: 'Min', description: 'min', value: 'min' },
  { label: 'Rate', description: 'rate', value: 'rate' },
  { label: 'Count', description: 'count number of samples', value: 'count' },
  { label: 'Range', description: 'Diff between max and min in the bucket', value: 'range' },
];
