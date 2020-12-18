import React from 'react';
import { shallow, ShallowWrapper } from 'enzyme';
import { QueryEditor } from './QueryEditor';
import { QueryTypeValue, RedisQuery } from '../redis';

const getQuery = (overrideQuery: object = {}): RedisQuery => ({
  key: '',
  aggregation: '',
  bucket: '',
  legend: '',
  command: '',
  field: '',
  filter: '',
  value: '',
  query: '',
  type: QueryTypeValue.CLI,
  section: '',
  size: 1,
  fill: true,
  streaming: true,
  streamingInterval: 1,
  streamingCapacity: 1,
  refId: '',
  ...overrideQuery,
});

type ShallowComponent = ShallowWrapper<QueryEditor['props'], QueryEditor['state'], QueryEditor>;

describe('QueryEditor', () => {
  const onRunQuery = jest.fn();
  const onChange = jest.fn();

  beforeEach(() => {
    onRunQuery.mockReset();
    onChange.mockReset();
  });

  describe('Type', () => {
    const getComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.prop('onChange') === wrapper.instance().onTypeChange;
      });
    it('Should set value from query', () => {
      const query = getQuery({ type: QueryTypeValue.CLI });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(query.type);
    });
    it('Should call onTypeChange when onChange prop was called', () => {
      const query = getQuery({ type: QueryTypeValue.CLI });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedMethod = jest.spyOn(wrapper.instance(), 'onTypeChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getComponent(wrapper);
      const newValue = QueryTypeValue.COMMAND;
      testedComponent.simulate('change', { value: newValue });
      expect(testedMethod).toHaveBeenCalledWith({ value: newValue });
      expect(onChange).toHaveBeenCalledWith({
        ...query,
        type: newValue,
        query: '',
        command: '',
      });
    });
  });

  describe('Query', () => {
    const getComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.prop('onChange') === wrapper.instance().onQueryChange;
      });
    it('Should be shown if type=cli', () => {
      const query = getQuery({ type: QueryTypeValue.CLI });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.exists()).toBeTruthy();
    });
    it('Should not be shown if type!=cli', () => {
      const query = getQuery({ type: QueryTypeValue.COMMAND });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should set value from query', () => {
      const query = getQuery({ type: QueryTypeValue.CLI });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(query.query);
    });
    it('Should call onQueryChange when onChange prop was called', () => {
      const query = getQuery({ type: QueryTypeValue.CLI });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedMethod = jest.spyOn(wrapper.instance(), 'onQueryChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getComponent(wrapper);
      const newValue = '123';
      testedComponent.simulate('change', { target: { value: newValue } });
      expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
      expect(onChange).toHaveBeenCalledWith({
        ...query,
        query: newValue,
      });
    });
  });

  describe('Command', () => {
    const getComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.prop('onChange') === wrapper.instance().onCommandChange;
      });
    it('Should be shown if type!=cli', () => {
      const query = getQuery({ type: QueryTypeValue.COMMAND });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.exists()).toBeTruthy();
    });
    it('Should not be shown if type=cli', () => {
      const query = getQuery({ type: QueryTypeValue.CLI });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });
    it('Should set value from query', () => {
      const query = getQuery({ type: QueryTypeValue.COMMAND, command: '123' });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(query.command);
    });
    it('Should call onCommandChange when onChange prop was called', () => {
      const query = getQuery({ type: QueryTypeValue.COMMAND });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedMethod = jest.spyOn(wrapper.instance(), 'onCommandChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getComponent(wrapper);
      const newValue = '123';
      testedComponent.simulate('change', { value: newValue });
      expect(testedMethod).toHaveBeenCalledWith({ value: newValue });
      expect(onChange).toHaveBeenCalledWith({
        ...query,
        command: newValue,
      });
    });
  });
  describe('Command fields', () => {
    describe('Key', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Key';
        });

      it('Should be shown when command exists in commands.key', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'get' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should be shown when command is not exists in commands.key', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'gettt' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'get', key: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.key);
      });

      it('Should call onKeyChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'get' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onKeyChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          key: newValue,
        });
      });
    });

    describe('LabelFilter', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Label Filter';
        });

      it('Should be shown when command exists in commands.filter', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.mrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should be shown when command is not exists in commands.filter', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.mrange', filter: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.filter);
      });

      it('Should call onFilterChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.mrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onFilterChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          filter: newValue,
        });
      });
    });

    describe('Field', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Field';
        });

      it('Should be shown when command exists in commands.field', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'hget' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should be shown when command is not exists in commands.field', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'hget', field: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.field);
      });

      it('Should call onFieldChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'hget' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onFieldChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          field: newValue,
        });
      });
    });

    describe('Legend', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Legend';
        });

      it('Should be shown when command exists in commands.legend', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.range' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should be shown when command is not exists in commands.legend', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.range', legend: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.legend);
      });

      it('Should call onLegendChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.range' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onLegendChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          legend: newValue,
        });
      });
    });

    describe('LegendLabel', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Legend Label';
        });

      it('Should be shown when command exists in commands.legendLabel', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.mrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should be shown when command is not exists in commands.legendLabel', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.mrange', legend: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.legend);
      });

      it('Should call onLegendChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.mrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onLegendChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          legend: newValue,
        });
      });
    });

    describe('ValueLabel', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Value Label';
        });

      it('Should be shown when command exists in commands.valueLabel', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.mrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should be shown when command is not exists in commands.valueLabel', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.mrange', value: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.value);
      });

      it('Should call onValueChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'ts.mrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onValueChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          value: newValue,
        });
      });
    });

    describe('Size', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Size';
        });

      it('Should be shown when command exists in commands.size', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'slowlogGet' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should be shown when command is not exists in commands.size', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'slowlogGet', size: 123 });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.size);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.COMMAND, command: 'slowlogGet' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          size: parseInt(newValue, 10),
        });
      });
    });
  });
});
