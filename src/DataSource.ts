import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { RedisDataSourceOptions, RedisQuery } from './types';

export class DataSource extends DataSourceWithBackend<RedisQuery, RedisDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<RedisDataSourceOptions>) {
    super(instanceSettings);
  }
}
