import { ConfigEditor, QueryEditor } from 'components';
import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './data-source';
import { RedisQuery } from './redis';
import { RedisDataSourceOptions } from './types';

/**
 * Data Source plugin
 */
export const plugin = new DataSourcePlugin<DataSource, RedisQuery, RedisDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
