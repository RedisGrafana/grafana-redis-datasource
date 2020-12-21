import { DataSourcePlugin } from '@grafana/data';
import { plugin } from './module';

/**
 * Plugin
 */
describe('Plugin', () => {
  it('Should be an instance of DataSourcePlugin', () => {
    expect(plugin).toBeInstanceOf(DataSourcePlugin);
  });
});
