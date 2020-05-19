import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { RedisTimeSeriesQuery, RedisTimeSeriesDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, RedisTimeSeriesQuery, RedisTimeSeriesDataSourceOptions>(
  DataSource
)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
