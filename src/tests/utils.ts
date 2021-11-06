import { AggregationValue, InfoSectionValue, QueryTypeValue, RedisQuery } from '../redis';

/**
 * Query
 */
export const getQuery = (overrideQuery: object = {}): RedisQuery => ({
  keyName: '',
  aggregation: AggregationValue.NONE,
  bucket: 0,
  legend: '',
  command: '',
  field: '',
  path: '',
  cypher: '',
  filter: '',
  value: '',
  query: '',
  type: QueryTypeValue.CLI,
  section: InfoSectionValue.STATS,
  size: 1,
  fill: true,
  streaming: true,
  streamingInterval: 1,
  streamingCapacity: 1,
  refId: '',
  ...overrideQuery,
});
