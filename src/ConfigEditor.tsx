import React, { ChangeEvent, PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { LegacyForms } from '@grafana/ui';
import { RedisDataSourceOptions, RedisSecureJsonData } from './types';

const { SecretFormField, FormField } = LegacyForms;

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

  onSizeChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = {
      ...options.jsonData,
      size: Number(event.target.value),
    };
    onOptionsChange({ ...options, jsonData });
  };

  // Secure field (only sent to the backend)
  onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: {
        password: event.target.value,
      },
    });
  };

  onResetPassword = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: {
        ...options.secureJsonFields,
        password: false,
      },
      secureJsonData: {
        ...options.secureJsonData,
        password: '',
      },
    });
  };

  render() {
    const { options } = this.props;
    const { url, jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as RedisSecureJsonData;

    return (
      <div className="gf-form-group">
        <h3 className="page-heading">Redis</h3>
        <div className="gf-form">
          <FormField
            label="URL"
            labelWidth={10}
            inputWidth={20}
            onChange={this.onURLChange}
            value={url || ''}
            placeholder="redis://..."
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Pool Size"
            labelWidth={10}
            inputWidth={10}
            onChange={this.onSizeChange}
            value={jsonData.size || 1}
            placeholder="1"
          />
        </div>

        <br />
        <h3 className="page-heading">Auth</h3>
        <div className="gf-form-inline">
          <div className="gf-form">
            <SecretFormField
              isConfigured={(secureJsonFields && secureJsonFields.password) as boolean}
              value={secureJsonData.password || ''}
              label="Password"
              placeholder="Database password"
              labelWidth={10}
              inputWidth={20}
              onReset={this.onResetPassword}
              onChange={this.onPasswordChange}
            />
          </div>
        </div>
      </div>
    );
  }
}
