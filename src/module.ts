import { ConfigEditor, QueryEditor } from 'components';
import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { RedisDataSourceOptions, RedisQuery } from './types';

/**
 * Data Source plugin
 */
export const plugin = new DataSourcePlugin<DataSource, RedisQuery, RedisDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
