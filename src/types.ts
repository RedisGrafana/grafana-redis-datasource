import { DataQuery, DataSourceJsonData } from '@grafana/data';

export interface RedisTimeSeriesQuery extends DataQuery {
  keyname?: string;
  cmd?: string;
  aggregation?: string;
  bucket?: string;
  legend?: string;
}

export const defaultQuery: Partial<RedisTimeSeriesQuery> = {};

/**
 * These are options configured for each DataSource instance
 */
export interface RedisTimeSeriesDataSourceOptions extends DataSourceJsonData {}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface RedisTimeSeriesSecureJsonData {}
