import { DataQuery, SelectableValue } from '@grafana/data';

/**
 * Query Type Values
 */
export enum QueryTypeValue {
  COMMAND = 'command',
  TIMESERIES = 'timeSeries',
  CLI = 'cli',
}

/**
 * Query Type
 */
export const QueryType: Array<SelectableValue<string>> = [
  {
    label: 'Redis commands',
    description: 'Hashes, Sets, Lists, Strings, Streams, etc.',
    value: QueryTypeValue.COMMAND,
  },
  {
    label: 'RedisTimeSeries commands',
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
  key?: string;

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
  aggregation?: string;

  /**
   * Bucket
   *
   * @type {string}
   */
  bucket?: string;

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
  section?: string;

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
}
