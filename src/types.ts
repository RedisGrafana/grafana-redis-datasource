import { DataQuery, DataSourceJsonData } from '@grafana/data';

/**
 * Redis Query
 */
export interface RedisQuery extends DataQuery {
  /**
   * Key name
   *
   * @type {string}
   */
  keyname?: string;

  /**
   * Command
   *
   * @type {string}
   */
  cmd?: string;

  /**
   * Aggregation
   *
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
