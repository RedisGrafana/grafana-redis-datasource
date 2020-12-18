import React from 'react';
import { shallow, ShallowWrapper } from 'enzyme';
import { RadioButtonGroup } from '@grafana/ui';
import { DataSourceSettings } from '@grafana/data';
import { ConfigEditor } from './ConfigEditor';
import { RedisDataSourceOptions, ClientTypeValue } from '../types';

interface OverrideOptions {
  [key: string]: unknown;
  jsonData?: object;
  secureJsonData?: object | null;
}

const getOptions = ({
  jsonData = {},
  secureJsonData = {},
  ...overrideOptions
}: OverrideOptions = {}): DataSourceSettings<RedisDataSourceOptions> => ({
  id: 1,
  orgId: 2,
  name: '',
  typeLogoUrl: '',
  type: '',
  access: '',
  url: '',
  password: '',
  user: '',
  database: '',
  basicAuth: false,
  basicAuthPassword: '',
  basicAuthUser: '',
  isDefault: false,
  secureJsonFields: {},
  readOnly: false,
  withCredentials: false,
  ...overrideOptions,
  jsonData: {
    poolSize: 0,
    timeout: 0,
    pingInterval: 0,
    pipelineWindow: 0,
    tlsAuth: false,
    tlsSkipVerify: false,
    client: ClientTypeValue.CLUSTER,
    sentinelName: '',
    acl: false,
    user: '',
    ...jsonData,
  },
  secureJsonData: {
    password: '',
    tlsClientCert: '',
    tlsClientKey: '',
    tlsCACert: '',
    ...secureJsonData,
  },
});

type ShallowComponent = ShallowWrapper<ConfigEditor['props'], ConfigEditor['state'], ConfigEditor>;

describe('ConfigEditor', () => {
  describe('Type', () => {
    const getTestedComponent = (wrapper: ShallowComponent) => wrapper.find(RadioButtonGroup);
    it('Should pass client value to type field', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(options.jsonData.client);
    });
    it('Should pass standalone as a value if client value is empty', () => {
      const options = getOptions({ jsonData: { client: null } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(ClientTypeValue.STANDALONE);
    });
    it('Should call onOptionsChange function when value was changed', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const newClient = ClientTypeValue.STANDALONE;
      getTestedComponent(wrapper).simulate('change', newClient);
      expect(onOptionsChange).toHaveBeenCalledWith(
        getOptions({
          jsonData: {
            client: newClient,
          },
        })
      );
    });
  });
  describe('MasterName', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'FormField' && node.prop('label') === 'Master Name';
      });
    it('If client is not sentinel should not be shown', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('If client is sentinel should be shown Master Name field', () => {
      const options = getOptions({ jsonData: { client: ClientTypeValue.SENTINEL } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).toBeTruthy();
      expect(testedComponent.prop('value')).toEqual(options.jsonData.sentinelName);
    });
    it('Should call onOptionsChange function when value was changed', () => {
      const options = getOptions({ jsonData: { client: ClientTypeValue.SENTINEL } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      const newValue = '123';
      testedComponent.simulate('change', { target: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith(
        getOptions({
          ...options,
          jsonData: {
            ...options.jsonData,
            sentinelName: newValue,
          },
        })
      );
    });
  });
  describe('Address', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'FormField' && node.prop('label') === 'Address';
      });
    it('Should pass url value to address field', () => {
      const options = getOptions({ url: 'localhost' });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(options.url);
    });
    it('Should call onOptionsChange when value was changed', () => {
      const options = getOptions({ url: 'localhost' });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      const newUrl = 'redis';
      testedComponent.simulate('change', { target: { value: newUrl } });
      expect(onOptionsChange).toHaveBeenCalledWith({ ...options, url: newUrl });
    });
  });
  describe('ACL', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'Switch' && node.prop('label') === 'ACL';
      });
    it('Should pass acl value', () => {
      const options = getOptions({ jsonData: { acl: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('checked')).toEqual(options.jsonData.acl);
    });
    it('Should pass default value if user value is empty', () => {
      const options = getOptions({ jsonData: { acl: null } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('checked')).toEqual(false);
    });
    it('Should call onOptionsChange when value was changed', () => {
      const options = getOptions({ jsonData: { acl: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      const newValue = false;
      testedComponent.simulate('change', { currentTarget: { checked: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        jsonData: {
          ...options.jsonData,
          acl: newValue,
        },
      });
    });
  });
  describe('Username', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'FormField' && node.prop('label') === 'Username';
      });
    it('If acl checked should be shown', () => {
      const options = getOptions({ jsonData: { acl: true, user: 'My user' } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).toBeTruthy();
      expect(testedComponent.prop('value')).toEqual(options.jsonData.user);
    });
    it('If acl not checked should not be shown', () => {
      const options = getOptions({ jsonData: { acl: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should call onOptionsChange when value was changed', () => {
      const options = getOptions({ jsonData: { acl: true, user: 'admin' } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      const newValue = 'guest';
      testedComponent.simulate('change', { target: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        jsonData: {
          ...options.jsonData,
          user: newValue,
        },
      });
    });
  });
  describe('Password', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'SecretFormField' && node.prop('label') === 'Password';
      });
    it('Should pass password value', () => {
      const options = getOptions({ secureJsonData: { password: '123' } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(options.secureJsonData.password);
    });
    it('Should call onResetPassword method when calls onReset prop', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const onResetPasswordMethod = jest.spyOn(wrapper.instance(), 'onResetPassword');
      wrapper.instance().forceUpdate();
      const testedComponent = getTestedComponent(wrapper);
      testedComponent.simulate('reset');
      expect(onResetPasswordMethod).toHaveBeenCalledTimes(1);
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        secureJsonData: {
          ...options.secureJsonData,
          password: '',
        },
        secureJsonFields: {
          ...options.secureJsonFields,
          password: false,
        },
      });
    });
    it('Should call onPasswordChange method when calls onChange prop', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const onPasswordChangeMethod = jest.spyOn(wrapper.instance(), 'onPasswordChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getTestedComponent(wrapper);
      const newValue = '123';
      testedComponent.simulate('change', { target: { value: newValue } });
      expect(onPasswordChangeMethod).toHaveBeenCalledWith({ target: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        secureJsonData: {
          ...options.secureJsonData,
          password: newValue,
        },
      });
    });
  });
  describe('PoolSize', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'FormField' && node.prop('label') === 'Pool Size';
      });
    it('Should pass value from options', () => {
      const options = getOptions({ jsonData: { poolSize: 10 } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(options.jsonData.poolSize);
    });
    it('Should pass default value if poolSize value is empty', () => {
      const options = getOptions({ jsonData: { poolSize: 0 } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(5);
    });
    it('Should call onPoolSizeChange method when calls onChange prop', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const onPoolSizeChangeMethod = jest.spyOn(wrapper.instance(), 'onPoolSizeChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getTestedComponent(wrapper);
      const newValue = 15;
      testedComponent.simulate('change', { target: { value: newValue } });
      expect(onPoolSizeChangeMethod).toHaveBeenCalledWith({ target: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        jsonData: {
          ...options.jsonData,
          poolSize: newValue,
        },
      });
    });
  });
  describe('Timeout', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'FormField' && node.prop('label') === 'Timeout, sec';
      });
    it('Should pass value from options', () => {
      const options = getOptions({ jsonData: { timeout: 10 } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(options.jsonData.timeout);
    });
    it('Should pass default value if timeout value is empty', () => {
      const options = getOptions({ jsonData: { timeout: '' } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(10);
    });
    it('Should call onTimeoutChange method when calls onChange prop', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const onTimeoutChangeMethod = jest.spyOn(wrapper.instance(), 'onTimeoutChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getTestedComponent(wrapper);
      const newValue = '15';
      testedComponent.simulate('change', { target: { value: newValue } });
      expect(onTimeoutChangeMethod).toHaveBeenCalledWith({ target: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        jsonData: {
          ...options.jsonData,
          timeout: parseInt(newValue, 10),
        },
      });
    });
  });
  describe('PingInterval', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'FormField' && node.prop('label') === 'Ping Interval, sec';
      });
    it('Should pass value from options', () => {
      const options = getOptions({ jsonData: { pingInterval: 10 } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(options.jsonData.pingInterval);
    });
    it('Should pass default value if pingInterval value is empty', () => {
      const options = getOptions({ jsonData: { pingInterval: '' } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(0);
    });
    it('Should call onPingIntervalChange method when calls onChange prop', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedMethod = jest.spyOn(wrapper.instance(), 'onPingIntervalChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getTestedComponent(wrapper);
      const newValue = '15';
      testedComponent.simulate('change', { target: { value: newValue } });
      expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        jsonData: {
          ...options.jsonData,
          pingInterval: parseInt(newValue, 10),
        },
      });
    });
  });
  describe('PipelineWindow', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'FormField' && node.prop('label') === 'Pipeline Window, μs';
      });
    it('Should pass value from options', () => {
      const options = getOptions({ jsonData: { pipelineWindow: 10 } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(options.jsonData.pipelineWindow);
    });
    it('Should pass default value if pipelineWindow value is empty', () => {
      const options = getOptions({ jsonData: { pipelineWindow: '' } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(0);
    });
    it('Should call onPipelineWindowChange method when calls onChange prop', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedMethod = jest.spyOn(wrapper.instance(), 'onPipelineWindowChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getTestedComponent(wrapper);
      const newValue = '15';
      testedComponent.simulate('change', { target: { value: newValue } });
      expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        jsonData: {
          ...options.jsonData,
          pipelineWindow: parseInt(newValue, 10),
        },
      });
    });
  });
  describe('ClientAuthentication', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'Switch' && node.prop('label') === 'Client Authentication';
      });
    it('Should pass value from options', () => {
      const options = getOptions({ jsonData: { tlsAuth: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('checked')).toEqual(options.jsonData.tlsAuth);
    });
    it('Should pass default value if tlsAuth value is empty', () => {
      const options = getOptions({ jsonData: { tlsAuth: '' } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('checked')).toEqual(false);
    });
    it('Should call onChangeOptions', () => {
      const options = getOptions();
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      const newValue = true;
      testedComponent.simulate('change', { currentTarget: { checked: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        jsonData: {
          ...options.jsonData,
          tlsAuth: newValue,
        },
      });
    });
  });
  describe('SkipVerify', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'Switch' && node.prop('label') === 'Skip Verify';
      });
    it('Should be shown if tlsAuth=true', () => {
      const options = getOptions({ jsonData: { tlsAuth: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).toBeTruthy();
    });
    it('Should not be shown if tlsAuth=false', () => {
      const options = getOptions({ jsonData: { tlsAuth: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should pass value from options', () => {
      const options = getOptions({ jsonData: { tlsAuth: true, tlsSkipVerify: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('checked')).toEqual(options.jsonData.tlsSkipVerify);
    });
    it('Should pass default value if tlsSkipVerify value is empty', () => {
      const options = getOptions({ jsonData: { tlsAuth: true, tlsSkipVerify: '' } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.prop('checked')).toEqual(false);
    });
    it('Should call onChangeOptions', () => {
      const options = getOptions({ jsonData: { tlsAuth: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      const newValue = true;
      testedComponent.simulate('change', { currentTarget: { checked: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        jsonData: {
          ...options.jsonData,
          tlsSkipVerify: newValue,
        },
      });
    });
  });
  describe('ClientCertificate', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'TextArea' && node.prop('onChange') === wrapper.instance().onTlsClientCertificateChange;
      });
    it('Should be shown if tlsAuth=true and tlsClientCert=false', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsClientCert: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).toBeTruthy();
    });
    it('Should not be shown if tlsAuth=false', () => {
      const options = getOptions({ jsonData: { tlsAuth: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should not be shown if tlsClientCert=true', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsClientCert: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should call onTlsClientCertificateChange when onChange prop was called', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsClientCert: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedMethod = jest.spyOn(wrapper.instance(), 'onTlsClientCertificateChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getTestedComponent(wrapper);
      const newValue = '123';
      testedComponent.simulate('change', { currentTarget: { value: newValue } });
      expect(testedMethod).toHaveBeenCalledWith({ currentTarget: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        secureJsonData: {
          ...options.secureJsonData,
          tlsClientCert: newValue,
        },
      });
    });
    it('Should call onResetTlsClientCertificate when reset button was clicked', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsClientCert: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedMethod = jest.spyOn(wrapper.instance(), 'onResetTlsClientCertificate');
      wrapper.instance().forceUpdate();
      const testedComponent = wrapper.findWhere((node) => {
        return node.name() === 'Button' && node.prop('onClick') === wrapper.instance().onResetTlsClientCertificate;
      });
      expect(testedComponent.exists()).toBeTruthy();
      testedComponent.simulate('click');
      expect(testedMethod).toHaveBeenCalled();
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        secureJsonFields: {
          ...options.secureJsonFields,
          tlsClientCert: false,
        },
        secureJsonData: {
          ...options.secureJsonData,
          tlsClientCert: '',
        },
      });
    });
  });
  describe('ClientKey', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'TextArea' && node.prop('onChange') === wrapper.instance().onTlsClientKeyChange;
      });
    it('Should be shown if tlsAuth=true and tlsClientKey=false', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsClientKey: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).toBeTruthy();
    });
    it('Should not be shown if tlsAuth=false', () => {
      const options = getOptions({ jsonData: { tlsAuth: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should not be shown if tlsClientKey=true', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsClientKey: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should call onTlsClientKeyChange when onChange prop was called', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsClientKey: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedMethod = jest.spyOn(wrapper.instance(), 'onTlsClientKeyChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getTestedComponent(wrapper);
      const newValue = '123';
      testedComponent.simulate('change', { currentTarget: { value: newValue } });
      expect(testedMethod).toHaveBeenCalledWith({ currentTarget: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        secureJsonData: {
          ...options.secureJsonData,
          tlsClientKey: newValue,
        },
      });
    });
    it('Should call onResetTlsClientKey when reset button was clicked', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsClientKey: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedMethod = jest.spyOn(wrapper.instance(), 'onResetTlsClientKey');
      wrapper.instance().forceUpdate();
      const testedComponent = wrapper.findWhere((node) => {
        return node.name() === 'Button' && node.prop('onClick') === wrapper.instance().onResetTlsClientKey;
      });
      expect(testedComponent.exists()).toBeTruthy();
      testedComponent.simulate('click');
      expect(testedMethod).toHaveBeenCalled();
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        secureJsonFields: {
          ...options.secureJsonFields,
          tlsClientKey: false,
        },
        secureJsonData: {
          ...options.secureJsonData,
          tlsClientKey: '',
        },
      });
    });
  });
  describe('CertificationAuthority', () => {
    const getTestedComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'TextArea' && node.prop('onChange') === wrapper.instance().onTlsCACertificateChange;
      });
    it('Should be shown if tlsAuth=true and tlsCACert=false', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsCACert: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).toBeTruthy();
    });
    it('Should not be shown if tlsAuth=false', () => {
      const options = getOptions({ jsonData: { tlsAuth: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should not be shown if tlsClientKey=true', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsCACert: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedComponent = getTestedComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should call onTlsCACertificateChange when onChange prop was called', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsCACert: false } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedMethod = jest.spyOn(wrapper.instance(), 'onTlsCACertificateChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getTestedComponent(wrapper);
      const newValue = '123';
      testedComponent.simulate('change', { currentTarget: { value: newValue } });
      expect(testedMethod).toHaveBeenCalledWith({ currentTarget: { value: newValue } });
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        secureJsonData: {
          ...options.secureJsonData,
          tlsCACert: newValue,
        },
      });
    });
    it('Should call onResetTlsCACertificate when reset button was clicked', () => {
      const options = getOptions({ jsonData: { tlsAuth: true }, secureJsonFields: { tlsCACert: true } });
      const onOptionsChange = jest.fn();
      const wrapper = shallow<ConfigEditor>(<ConfigEditor options={options} onOptionsChange={onOptionsChange} />);
      const testedMethod = jest.spyOn(wrapper.instance(), 'onResetTlsCACertificate');
      wrapper.instance().forceUpdate();
      const testedComponent = wrapper.findWhere((node) => {
        return node.name() === 'Button' && node.prop('onClick') === wrapper.instance().onResetTlsCACertificate;
      });
      expect(testedComponent.exists()).toBeTruthy();
      testedComponent.simulate('click');
      expect(testedMethod).toHaveBeenCalled();
      expect(onOptionsChange).toHaveBeenCalledWith({
        ...options,
        secureJsonFields: {
          ...options.secureJsonFields,
          tlsCACert: false,
        },
        secureJsonData: {
          ...options.secureJsonData,
          tlsCACert: '',
        },
      });
    });
  });
});
