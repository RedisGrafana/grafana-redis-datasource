import React, { ChangeEvent, PureComponent } from 'react';
import { LegacyForms } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { RedisDataSourceOptions } from './types';

const { FormField } = LegacyForms;

interface Props extends DataSourcePluginOptionsEditorProps<RedisDataSourceOptions> {}

interface State {}

export class ConfigEditor extends PureComponent<Props, State> {
  onURLChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    options.url = event.target.value;
    const jsonData = {
      ...options.jsonData,
    };
    onOptionsChange({ ...options, jsonData });
  };

  render() {
    const { options } = this.props;

    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormField
            label="URL"
            labelWidth={6}
            inputWidth={20}
            onChange={this.onURLChange}
            value={options.url || 'redis://localhost:6379'}
            placeholder="redis://..."
          />
        </div>
      </div>
    );
  }
}
