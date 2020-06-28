import React, { ChangeEvent, PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { LegacyForms } from '@grafana/ui';
import { RedisDataSourceOptions, RedisSecureJsonData } from './types';

/**
 * Form Field
 */
const { SecretFormField, FormField } = LegacyForms;

/**
 * Editor Property
 */
interface Props extends DataSourcePluginOptionsEditorProps<RedisDataSourceOptions> {}

/**
 * State
 */
interface State {}

/**
 * Config Editor
 */
export class ConfigEditor extends PureComponent<Props, State> {
  /**
   * URL change
   *
   * @param event Event
   */
  onURLChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    options.url = event.target.value;

    const jsonData = {
      ...options.jsonData,
    };

    onOptionsChange({ ...options, jsonData });
  };

  /**
   * Pool Size change
   *
   * @param event Event
   */
  onPoolSizeChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = { ...options.jsonData, poolSize: Number(event.target.value) };

    onOptionsChange({ ...options, jsonData });
  };

  /**
   * Timeout change
   *
   * @param event Event
   */
  onTimeoutChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = { ...options.jsonData, timeout: Number(event.target.value) };

    onOptionsChange({ ...options, jsonData });
  };

  /**
   * Ping interval change
   *
   * @param event Event
   */
  onPingIntervalChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = { ...options.jsonData, pingInterval: Number(event.target.value) };

    onOptionsChange({ ...options, jsonData });
  };

  /**
   * Pipeline window change
   *
   * @param event Event
   */
  onPipelineWindowChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    const jsonData = { ...options.jsonData, pipelineWindow: Number(event.target.value) };

    onOptionsChange({ ...options, jsonData });
  };

  /**
   * Password Secure field (only sent to the backend)
   *
   * @param event Event
   */
  onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({ ...options, secureJsonData: { password: event.target.value } });
  };

  /**
   * Password Reset
   */
  onResetPassword = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: { ...options.secureJsonFields, password: false },
      secureJsonData: { ...options.secureJsonData, password: '' },
    });
  };

  /**
   * Render Editor
   */
  render() {
    const { options } = this.props;
    const { url, jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as RedisSecureJsonData;

    /**
     * Return
     */
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
            tooltip="Accepts host:port address or a URI, as defined in https://www.iana.org/assignments/uri-schemes/prov/redis"
            placeholder="redis://..."
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Pool Size"
            labelWidth={10}
            inputWidth={10}
            onChange={this.onPoolSizeChange}
            value={jsonData.poolSize || 5}
            tooltip="Will keep open at least the given number of connections to the redis instance at the given address.
            The recommended size of the pool depends on the number of concurrent goroutines that will use the pool and
            whether implicit pipelining is enabled or not."
          />

          <FormField
            label="Ping Interval, sec"
            labelWidth={10}
            inputWidth={10}
            onChange={this.onPingIntervalChange}
            value={jsonData.pingInterval || 0}
            tooltip="Specifies the interval in seconds at which a ping event happens.
            A shorter interval means connections are pinged more frequently, but also means more traffic with the server.
            If interval is zero then ping will be disabled."
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Timeout, sec"
            labelWidth={10}
            inputWidth={10}
            onChange={this.onTimeoutChange}
            value={jsonData.timeout || 10}
            tooltip="Sets the duration in seconds for connect, read and write timeouts."
          />
        </div>

        <div className="gf-form">
          <FormField
            label="Pipeline Window, μs"
            labelWidth={10}
            inputWidth={10}
            onChange={this.onPipelineWindowChange}
            value={jsonData.pipelineWindow || 0}
            tooltip="Sets the duration in microseconds after which internal pipelines will be flushed.
            If window is zero then implicit pipelining will be disabled."
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
