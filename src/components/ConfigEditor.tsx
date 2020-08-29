import React, { ChangeEvent, PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { Button, InlineFormLabel, LegacyForms, RadioButtonGroup, TextArea } from '@grafana/ui';
import { ClientType, ClientTypeValue, RedisDataSourceOptions, RedisSecureJsonData } from '../types';

/**
 * Form Field
 */
const { SecretFormField, FormField, Switch } = LegacyForms;

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
   * Pool Size change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onPoolSizeChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({ ...options, jsonData: { ...options.jsonData, poolSize: Number(event.target.value) } });
  };

  /**
   * Timeout change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onTimeoutChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({ ...options, jsonData: { ...options.jsonData, timeout: Number(event.target.value) } });
  };

  /**
   * Ping interval change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onPingIntervalChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({ ...options, jsonData: { ...options.jsonData, pingInterval: Number(event.target.value) } });
  };

  /**
   * Pipeline window change
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onPipelineWindowChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({ ...options, jsonData: { ...options.jsonData, pipelineWindow: Number(event.target.value) } });
  };

  /**
   * Password Secure field (only sent to the backend)
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  onPasswordChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({ ...options, secureJsonData: { ...options.secureJsonData, password: event.target.value } });
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
   * TLS Client Certificate
   *
   * @param {ChangeEvent<HTMLTextAreaElement>} event Event
   */
  onTlsClientCertificateChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: { ...options.secureJsonData, tlsClientCert: event.currentTarget.value },
    });
  };

  /**
   * TLS Client Certificate Reset
   */
  onResetTlsClientCertificate = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: { ...options.secureJsonFields, tlsClientCert: false },
      secureJsonData: { ...options.secureJsonData, tlsClientCert: '' },
    });
  };

  /**
   * TLS Certification Authority
   *
   * @param {ChangeEvent<HTMLTextAreaElement>} event Event
   */
  onTlsCACertificateChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: { ...options.secureJsonData, tlsCACert: event.currentTarget.value },
    });
  };

  /**
   * TLS CA Certificate Reset
   */
  onResetTlsCACertificate = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: { ...options.secureJsonFields, tlsCACert: false },
      secureJsonData: { ...options.secureJsonData, tlsCACert: '' },
    });
  };

  /**
   * TLS Client key
   *
   * @param {ChangeEvent<HTMLTextAreaElement>} event Event
   */
  onTlsClientKeyChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonData: { ...options.secureJsonData, tlsClientKey: event.currentTarget.value },
    });
  };

  /**
   * TLS Client Key Reset
   */
  onResetTlsClientKey = () => {
    const { onOptionsChange, options } = this.props;
    onOptionsChange({
      ...options,
      secureJsonFields: { ...options.secureJsonFields, tlsClientKey: false },
      secureJsonData: { ...options.secureJsonData, tlsClientKey: '' },
    });
  };

  /**
   * Render Editor
   */
  render() {
    const { options, onOptionsChange } = this.props;
    const { url, jsonData, secureJsonFields } = options;
    const secureJsonData = (options.secureJsonData || {}) as RedisSecureJsonData;

    /**
     * Return
     */
    return (
      <div className="gf-form-group">
        <h3 className="page-heading">Redis</h3>

        <div className="gf-form">
          <InlineFormLabel width={10} tooltip="">
            Type
          </InlineFormLabel>
          <RadioButtonGroup
            options={ClientType}
            value={jsonData.client || ClientTypeValue.STANDALONE}
            onChange={(value) => {
              const jsonData = { ...options.jsonData, client: value as ClientTypeValue };
              onOptionsChange({ ...options, jsonData });
            }}
          />
        </div>

        {jsonData.client === ClientTypeValue.SENTINEL && (
          <div className="gf-form">
            <FormField
              label="Master Name"
              labelWidth={10}
              inputWidth={10}
              value={jsonData.sentinelName}
              tooltip="Provide Master Name to connect to."
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                onOptionsChange({ ...options, jsonData: { ...options.jsonData, sentinelName: event.target.value } });
              }}
            />
          </div>
        )}

        <div className="gf-form">
          <FormField
            label="Address"
            labelWidth={10}
            inputWidth={20}
            onChange={(event: ChangeEvent<HTMLInputElement>) => {
              onOptionsChange({ ...options, url: event.target.value });
            }}
            value={url || ''}
            tooltip="Accepts host:port address or a URI, as defined in https://www.iana.org/assignments/uri-schemes/prov/redis.
            For Redis Cluster and Sentinel can contain multiple values with comma.
            For a Unix Socket accepts path to socket file."
            placeholder="redis://..."
          />
        </div>

        <div className="gf-form">
          <SecretFormField
            isConfigured={(secureJsonFields && secureJsonFields.password) as boolean}
            value={secureJsonData.password || ''}
            label="Password"
            placeholder="Database password"
            labelWidth={10}
            inputWidth={20}
            tooltip="When specified AUTH command will be used to authenticate with the provided password."
            onReset={this.onResetPassword}
            onChange={this.onPasswordChange}
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
        </div>

        <br />
        <h3 className="page-heading">Advanced settings</h3>

        <div className="gf-form">
          <Switch
            label="ACL"
            labelClass="width-10"
            tooltip="Allows certain connections to be limited in terms of the commands that can be executed and the keys that can be accessed"
            checked={jsonData.acl || false}
            onChange={(event) => {
              const jsonData = { ...options.jsonData, acl: event.currentTarget.checked };
              onOptionsChange({ ...options, jsonData });
            }}
          />

          {jsonData.acl && (
            <FormField
              label="Username"
              labelWidth={10}
              inputWidth={10}
              value={jsonData.user}
              tooltip="Provide ACL Username to authenticate."
              onChange={(event: ChangeEvent<HTMLInputElement>) => {
                onOptionsChange({ ...options, jsonData: { ...options.jsonData, user: event.target.value } });
              }}
            />
          )}
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
        <h3 className="page-heading">TLS</h3>

        <div className="gf-form-inline">
          <Switch
            label="Client Authentication"
            labelClass="width-10"
            checked={jsonData.tlsAuth || false}
            onChange={(event) => {
              const jsonData = { ...options.jsonData, tlsAuth: event.currentTarget.checked };
              onOptionsChange({ ...options, jsonData });
            }}
          />

          {jsonData.tlsAuth && (
            <Switch
              label="Skip Verify"
              labelClass="width-10"
              tooltip="If checked, the server's certificate will not be checked for validity."
              checked={jsonData.tlsSkipVerify || false}
              onChange={(event) => {
                const jsonData = { ...options.jsonData, tlsSkipVerify: event.currentTarget.checked };
                onOptionsChange({ ...options, jsonData });
              }}
            />
          )}
        </div>

        {jsonData.tlsAuth && (
          <>
            <div className="gf-form-inline">
              <div className="gf-form gf-form--v-stretch">
                <label className="gf-form-label width-10">Client Certificate</label>
              </div>

              {secureJsonFields && secureJsonFields.tlsClientCert ? (
                <Button type="reset" variant="secondary" onClick={this.onResetTlsClientCertificate}>
                  Reset
                </Button>
              ) : (
                <div className="gf-form gf-form--grow">
                  <TextArea
                    css=""
                    rows={7}
                    className="gf-form-input gf-form-textarea"
                    placeholder="Begins with -----BEGIN CERTIFICATE-----"
                    onChange={this.onTlsClientCertificateChange}
                  />
                </div>
              )}
            </div>

            <div className="gf-form-inline">
              <div className="gf-form gf-form--v-stretch">
                <label className="gf-form-label width-10">Client Key</label>
              </div>
              <div className="gf-form gf-form--grow">
                {secureJsonFields && secureJsonFields.tlsClientKey ? (
                  <Button type="reset" variant="secondary" onClick={this.onResetTlsClientKey}>
                    Reset
                  </Button>
                ) : (
                  <TextArea
                    css=""
                    rows={7}
                    className="gf-form-input gf-form-textarea"
                    placeholder="Begins with -----BEGIN PRIVATE KEY-----"
                    onChange={this.onTlsClientKeyChange}
                  />
                )}
              </div>
            </div>

            <div className="gf-form-inline">
              <div className="gf-form gf-form--v-stretch">
                <label className="gf-form-label width-10">Certification Authority</label>
              </div>
              {secureJsonFields && secureJsonFields.tlsCACert ? (
                <Button type="reset" variant="secondary" onClick={this.onResetTlsCACertificate}>
                  Reset
                </Button>
              ) : (
                <div className="gf-form gf-form--grow">
                  <TextArea
                    css=""
                    rows={7}
                    className="gf-form-input gf-form-textarea"
                    placeholder="Begins with -----BEGIN CERTIFICATE-----"
                    onChange={this.onTlsCACertificateChange}
                  />
                </div>
              )}
            </div>
          </>
        )}
      </div>
    );
  }
}
