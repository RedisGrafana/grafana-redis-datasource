import { DataSourcePlugin } from '@grafana/data';
import { ConfigEditor, QueryEditor } from './components';
import { DataSource } from './datasource';
import { RedisQuery } from './redis';
import { RedisDataSourceOptions } from './types';

/**
 * Data Source plugin
 */
export const plugin = new DataSourcePlugin<DataSource, RedisQuery, RedisDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
