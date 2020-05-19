import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { RedisTimeSeriesDataSourceOptions, RedisTimeSeriesQuery } from './types';

export class DataSource extends DataSourceWithBackend<RedisTimeSeriesQuery, RedisTimeSeriesDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedisTimeSeriesDataSourceOptions>) {
    super(instanceSettings);
  }
}
