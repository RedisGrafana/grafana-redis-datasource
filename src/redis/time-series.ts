import { SelectableValue } from '@grafana/data';

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
