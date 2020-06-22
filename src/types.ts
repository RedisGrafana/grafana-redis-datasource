import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface RedisQuery extends DataQuery {
  keyname?: string;
  cmd?: string;
  aggregation?: string;
  bucket?: string;
  legend?: string;
}

export const defaultQuery: Partial<RedisQuery> = {};

/**
 * These are options configured for each DataSource instance
 */
export interface RedisDataSourceOptions extends DataSourceJsonData {}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface RedisSecureJsonData {}
