import { DataSourceJsonData } from '@grafana/data';
import { ClientTypeValue } from './constants';

/**
 * Options configured for each DataSource instance
 */
export interface RedisDataSourceOptions extends DataSourceJsonData {
  /**
   * Pool Size
   *
   * @type {number}
   */
  poolSize: number;

  /**
   * Timeout
   *
   * @type {number}
   */
  timeout: number;

  /**
   * Pool Ping Interval
   *
   * @type {number}
   */
  pingInterval: number;

  /**
   * Pool Pipeline Window
   *
   * @type {number}
   */
  pipelineWindow: number;

  /**
   * TLS Authentication
   *
   * @type {boolean}
   */
  tlsAuth: boolean;

  /**
   * TLS Skip Verify
   *
   * @type {boolean}
   */
  tlsSkipVerify: boolean;

  /**
   * Client Type
   *
   * @type {ClientTypeValue}
   */
  client: ClientTypeValue;

  /**
   * Sentinel Master group name
   *
   * @type {string}
   */
  sentinelName: string;

  /**
   * ACL enabled
   *
   * @type {boolean}
   */
  acl: boolean;

  /**
   * ACL Username
   *
   * @type {string}
   */
  user: string;
}

/**
 * Value that is used in the backend, but never sent over HTTP to the frontend
 */
export interface RedisSecureJsonData {
  /**
   * Database password
   *
   * @type {string}
   */
  password?: string;

  /**
   * TLS Client Certificate
   *
   * @type {string}
   */
  tlsClientCert?: string;

  /**
   * TLS Client Key
   *
   * @type {string}
   */
  tlsClientKey?: string;

  /**
   * TLS Authority Certificate
   *
   * @type {string}
   */
  tlsCACert?: string;
}
