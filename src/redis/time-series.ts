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
}

export enum ReducerValue {
  AVG = 'avg',
  SUM = 'sum',
  MIN = 'min',
  MAX = 'max',
  RANGE = 'range',
  COUNT = 'count',
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
];

/**
 * Reducers
 */
export const Reducers: Array<SelectableValue<ReducerValue>> = [
  { label: 'Avg', description: 'Arithmetic mean of all non-NaN values', value: ReducerValue.AVG },
  { label: 'Sum', description: 'Sum of all non-NaN values', value: ReducerValue.SUM },
  { label: 'Min', description: 'Minimum non-NaN value', value: ReducerValue.MIN },
  { label: 'Max', description: 'Maximum non-NaN value', value: ReducerValue.MAX },
  {
    label: 'Range',
    description: 'Difference between maximum non-Nan value and minimum non-NaN value',
    value: ReducerValue.RANGE,
  },
  { label: 'Count', description: 'Number of non-NaN values', value: ReducerValue.COUNT },
  {
    label: 'Std Population',
    description: 'Population standard deviation of all non-NaN values',
    value: ReducerValue.STDP,
  },
  { label: 'Std Sample', description: 'Sample standard deviation of all non-NaN values', value: ReducerValue.STDS },
  { label: 'Var Population', description: 'Population variance of all non-NaN values', value: ReducerValue.VARP },
  { label: 'Var Sample', description: 'Sample variance of all non-NaN values', value: ReducerValue.VARS },
];
