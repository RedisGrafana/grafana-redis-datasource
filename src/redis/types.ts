import { ZRangeQueryValue } from 'redis';
import { DataQuery } from '@grafana/data';
import { StreamingDataType } from '../constants';
import { InfoSectionValue } from './info';
import { QueryTypeValue } from './query';
import { AggregationValue } from './time-series';
import { SortDirectionValue } from './search';

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
   * Query for Search Command
   *
   * @type {string}
   */
  searchQuery?: string;

  /**
   * search return fields
   *
   * @type {string[]}
   */
  returnFields?: string[];

  /**
   * offset into result set to start at
   */
  offset?: number;

  /**
   * The direction to sort.
   */
  sortDirection?: SortDirectionValue;

  /**
   * The value to sort by
   */
  sortBy?: string;

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
   * @type {AggregationValue}
   */
  aggregation?: AggregationValue;

  /**
   * ZRANGE Query
   *
   * @see https://redis.io/commands/zrange
   * @type {ZRangeQueryValue}
   */
  zrangeQuery?: ZRangeQueryValue;

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
   * Streaming data type
   * @type {StreamingDataType}
   */
  streamingDataType?: StreamingDataType;

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

  /**
   * Samples for MEMORY USAGE command
   *
   * @type {number}
   */
  samples?: number;

  /**
   * Start for Streams
   *
   * @type {string}
   */
  start?: string;

  /**
   * Stop for Streams
   *
   * @type {string}
   */
  end?: string;

  /**
   * Minimum for ZSet
   *
   * @type {string}
   */
  min?: string;

  /**
   * Maximum for ZSet
   *
   * @type {string}
   */
  max?: string;

  /**
   * Cypher
   *
   * @type {string}
   */
  cypher?: string;

  /**
   * Path
   *
   * @type {string}
   */
  path?: string;
}
