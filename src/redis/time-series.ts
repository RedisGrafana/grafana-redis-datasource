import { SelectableValue } from '@grafana/data';

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
  { label: 'Sum', description: 'Sum', value: AggregationValue.SUM },
];
