import { DataQuery, SelectableValue } from '@grafana/data';
import { InfoSectionValue } from './info';
import { AggregationValue } from './time-series';

/**
 * Query Type Values
 */
export enum QueryTypeValue {
  COMMAND = 'command',
  TIMESERIES = 'timeSeries',
  SEARCH = 'search',
  CLI = 'cli',
}

/**
 * Query Type
 */
export const QueryType: Array<SelectableValue<QueryTypeValue>> = [
  {
    label: 'Redis',
    description: 'Hashes, Sets, Lists, Strings, Streams, etc.',
    value: QueryTypeValue.COMMAND,
  },
  {
    label: 'RediSearch',
    description: 'Redis Secondary Index & Query Engine',
    value: QueryTypeValue.SEARCH,
  },
  {
    label: 'RedisTimeSeries',
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
 * Redis Query
 */
export interface RedisQuery extends DataQuery {
  /**
   * Type
   *
   * @type {string}
   */
  type: QueryTypeValue;

  /**
   * Query command
   *
   * @type {string}
   */
  query?: string;

  /**
   * Field
   *
   * @type {string}
   */
  field?: string;

  /**
   * Redis TimeSeries filter
   *
   * @see https://oss.redislabs.com/redistimeseries/commands/#filtering
   * @type {string}
   */
  filter?: string;

  /**
   * Command
   *
   * @type {string}
   */
  command?: string;

  /**
   * Key name
   *
   * @type {string}
   */
  keyName?: string;

  /**
   * Value label
   *
   * @type {string}
   */
  value?: string;

  /**
   * Aggregation
   *
   * @see https://oss.redislabs.com/redistimeseries/commands/#aggregation-compaction-downsampling
   * @type {string}
   */
  aggregation?: AggregationValue;

  /**
   * Bucket
   *
   * @type {number}
   */
  bucket?: number;

  /**
   * Fill
   *
   * @type {boolean}
   */
  fill?: boolean;

  /**
   * Legend label
   *
   * @type {string}
   */
  legend?: string;

  /**
   * Info Section
   *
   * @type {string}
   */
  section?: InfoSectionValue;

  /**
   * Size
   *
   * @type {number}
   */
  size?: number;

  /**
   * Support streaming
   *
   * @type {boolean}
   */
  streaming?: boolean;

  /**
   * Streaming interval in milliseconds
   *
   * @type {number}
   */
  streamingInterval?: number;

  /**
   * Streaming capacity
   *
   * @type {number}
   */
  streamingCapacity?: number;

  /**
   * Cursor for SCAN command
   *
   * @type {string}
   */
  cursor?: string;

  /**
   * Match for SCAN command
   *
   * @type {string}
   */
  match?: string;

  /**
   * Count for SCAN command
   *
   * @type {number}
   */
  count?: number;
}
