import { DataQuery, DataSourceJsonData } from '@grafana/data';

/**
 * Redis Query
 */
export interface RedisQuery extends DataQuery {
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
  filter: string;

  /**
   * Value label
   *
   * @type {string}
   */
  value: string;

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
   * Legend label
   *
   * @type {string}
   */
  legend?: string;
}

/**
 * These are options configured for each DataSource instance
 */
export interface RedisDataSourceOptions extends DataSourceJsonData {
  /**
   * Pool Size
   *
   * @type {number}
   */
  size: number;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface RedisSecureJsonData {
  /**
   * Database password
   *
   * @type {string}
   */
  password?: string;
}
