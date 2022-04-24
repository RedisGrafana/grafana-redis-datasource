import { SelectableValue } from '@grafana/data';

/**
 * Supported Commands
 */
export enum RedisTimeSeries {
  GET = 'ts.get',
  INFO = 'ts.info',
  MGET = 'ts.mget',
  MRANGE = 'ts.mrange',
  QUERYINDEX = 'ts.queryindex',
  RANGE = 'ts.range',
}

/**
 * Commands List
 */
export const RedisTimeSeriesCommands = [
  {
    label: RedisTimeSeries.GET.toUpperCase(),
    description: 'Returns the last sample',
    value: RedisTimeSeries.GET,
  },
  {
    label: RedisTimeSeries.INFO.toUpperCase(),
    description: 'Returns information and statistics on the time-series',
    value: RedisTimeSeries.INFO,
  },
  {
    label: RedisTimeSeries.MGET.toUpperCase(),
    description: 'Returns the last samples matching the specific filter',
    value: RedisTimeSeries.MGET,
  },
  {
    label: RedisTimeSeries.MRANGE.toUpperCase(),
    description: 'Query a timestamp range across multiple time-series by filters',
    value: RedisTimeSeries.MRANGE,
  },
  {
    label: RedisTimeSeries.QUERYINDEX.toUpperCase(),
    description: 'Query all the keys matching the filter list',
    value: RedisTimeSeries.QUERYINDEX,
  },
  { label: RedisTimeSeries.RANGE.toUpperCase(), description: 'Query a range', value: RedisTimeSeries.RANGE },
];

/**
 * Aggregation Values
 */
export enum AggregationValue {
  NONE = '',
  AVG = 'avg',
  COUNT = 'count',
  MAX = 'max',
  MIN = 'min',
  RANGE = 'range',
  SUM = 'sum',
  FIRST = 'first',
  LAST = 'last',
  STDP = 'std.p',
  STDS = 'std.s',
  VARP = 'var.p',
  VARS = 'var.s',
}

/**
 * Aggregations
 */
export const Aggregations: Array<SelectableValue<AggregationValue>> = [
  { label: 'None', description: 'No aggregation', value: AggregationValue.NONE },
  { label: 'Avg', description: 'Average', value: AggregationValue.AVG },
  { label: 'Count', description: 'Count number of samples', value: AggregationValue.COUNT },
  { label: 'Max', description: 'Maximum', value: AggregationValue.MAX },
  { label: 'Min', description: 'Minimum', value: AggregationValue.MIN },
  { label: 'Range', description: 'Diff between maximum and minimum in the bucket', value: AggregationValue.RANGE },
  { label: 'Sum', description: 'Summation', value: AggregationValue.SUM },
  { label: 'First', description: 'The value with the lowest timestamp in the bucket', value: AggregationValue.FIRST },
  { label: 'Last', description: 'The value with the highest timestamp in the bucket', value: AggregationValue.LAST },
  { label: 'Std.p', description: 'Population standard deviation of the values', value: AggregationValue.STDP },
  { label: 'Std.s', description: 'Sample standard deviation of the values', value: AggregationValue.STDS },
  { label: 'Var.p', description: 'Population variance of the values', value: AggregationValue.VARP },
  { label: 'Var.s', description: 'Sample variance of the values', value: AggregationValue.VARS },
];
