import React from 'react';
import { shallow, ShallowWrapper } from 'enzyme';
import { QueryEditor } from './query-editor';
import { QueryTypeValue } from '../../redis';
import { getQuery } from '../../tests/utils';

type ShallowComponent = ShallowWrapper<QueryEditor['props'], QueryEditor['state'], QueryEditor>;

/**
 * Query Editor
 */
describe('QueryEditor', () => {
  const onRunQuery = jest.fn();
  const onChange = jest.fn();

  beforeEach(() => {
    onRunQuery.mockReset();
    onChange.mockReset();
  });

  /**
   * Query Type
   */
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
      const newValue = QueryTypeValue.REDIS;
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

  /**
   * Query
   */
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
      const query = getQuery({ type: QueryTypeValue.REDIS });
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

  /**
   * Command
   */
  describe('Command', () => {
    const getComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.prop('onChange') === wrapper.instance().onCommandChange;
      });

    it('Should be shown if type!=cli', () => {
      const query = getQuery({ type: QueryTypeValue.REDIS });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.exists()).toBeTruthy();
    });

    it('Should not not be shown if type=cli', () => {
      const query = getQuery({ type: QueryTypeValue.CLI });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.exists()).not.toBeTruthy();
    });

    it('Should set value from query', () => {
      const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(query.command);
    });

    it('Should call onCommandChange when onChange prop was called', () => {
      const query = getQuery({ type: QueryTypeValue.REDIS });
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

  /**
   * Command properties
   */
  describe('Command fields', () => {
    /**
     * Key name
     */
    describe('Key', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Key';
        });

      it('Should be shown when command exists in commands.key', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'get' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.key', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'gettt' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'get', keyName: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.keyName);
      });

      it('Should call onKeyChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'get' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onKeyNameChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          keyName: newValue,
        });
      });
    });

    /**
     * Filter
     */
    describe('LabelFilter', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Label Filter';
        });

      it('Should be shown when command exists in commands.filter', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.mrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.filter', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.mrange', filter: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.filter);
      });

      it('Should call onFilterChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.mrange' });
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

    /**
     * Field
     */
    describe('Field', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Field';
        });

      it('Should be shown when command exists in commands.field', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'hget' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.field', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'hget', field: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.field);
      });

      it('Should call onFieldChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'hget' });
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

    /**
     * Legend
     */
    describe('Legend', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Legend';
        });

      it('Should be shown when command exists in commands.legend', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.range' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.legend', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.range', legend: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.legend);
      });

      it('Should call onLegendChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.range' });
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

    /**
     * Legend Label
     */
    describe('LegendLabel', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Legend Label';
        });

      it('Should be shown when command exists in commands.legendLabel', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.mrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.legendLabel', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.mrange', legend: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.legend);
      });

      it('Should call onLegendChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.mrange' });
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

    /**
     * Value Label
     */
    describe('ValueLabel', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Value Label';
        });

      it('Should be shown when command exists in commands.valueLabel', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.mrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.valueLabel', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.mrange', value: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.value);
      });

      it('Should call onValueChange method when onChange prop was called', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'ts.mrange' });
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

    /**
     * Size
     */
    describe('Size', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Size';
        });

      it('Should be shown when command exists in commands.size', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'slowlogGet' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.size', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'slowlogGet', size: 123 });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.size);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'slowlogGet' });
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

    /**
     * Cursor
     */
    describe('Cursor', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Cursor';
        });

      it('Should be shown when command exists in commands.cursor', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.cursor', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan', cursor: '2' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.cursor);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          cursor: newValue,
        });
      });
    });

    /**
     * Match
     */
    describe('Match Pattern', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Match pattern';
        });

      it('Should be shown when command exists in commands.match', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.match', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan', match: '**' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.match);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          match: newValue,
        });
      });
    });

    /**
     * Start
     */
    describe('Start', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Start';
        });

      it('Should be shown when command exists in commands.start', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'xrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.start', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'xrange', start: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.start);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'xrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          start: newValue,
        });
      });
    });

    /**
     * End
     */
    describe('End', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'End';
        });

      it('Should be shown when command exists in commands.end', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'xrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.end', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'xrange', end: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.end);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'xrange' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          end: newValue,
        });
      });
    });

    /**
     * Count
     */
    describe('Count', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Count';
        });

      it('Should be shown when command exists in commands.count', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.count', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan', count: 123 });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.count);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          count: parseInt(newValue, 10),
        });
      });
    });

    /**
     * Samples
     */
    describe('Samples', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField', node.prop('label') === 'Samples';
        });

      it('Should be shown when command exists in commands.samples', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.samples', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan', count: 123 });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.samples);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'tmscan' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          samples: parseInt(newValue, 10),
        });
      });
    });

    /**
     * Section
     */
    describe('Section', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.prop('onChange') === wrapper.instance().onInfoSectionTextChange;
        });

      it('Should be shown when command exists in commands.section', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'info' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.section', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'info', section: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.section);
      });

      it('Should call onInfoSectionTextChange method when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.REDIS, command: 'info' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onInfoSectionTextChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { value: newValue });
        expect(testedMethod).toHaveBeenCalledWith({ value: newValue });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          section: newValue,
        });
      });
    });

    /**
     * Aggregation
     */
    describe('Aggregation', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.prop('onChange') === wrapper.instance().onAggregationTextChange;
        });

      it('Should be shown when command exists in commands.aggregation', () => {
        const query = getQuery({ type: QueryTypeValue.TIMESERIES, command: 'ts.range' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when command is not exists in commands.aggregation', () => {
        const query = getQuery({ type: QueryTypeValue.TIMESERIES, command: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ type: QueryTypeValue.TIMESERIES, command: 'ts.range', aggregation: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.aggregation);
      });

      it('Should call onAggregationTextChange method when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.TIMESERIES, command: 'ts.range' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onAggregationTextChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = '1234';
        testedComponent.simulate('change', { value: newValue });
        expect(testedMethod).toHaveBeenCalledWith({ value: newValue });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          aggregation: newValue,
        });
      });
    });

    /**
     * Bucket
     */
    describe('Bucket', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.prop('onChange') === wrapper.instance().onBucketTextChange;
        });

      it('Should be shown when command exists in commands.aggregation and aggregation value is selected', () => {
        const query = getQuery({ type: QueryTypeValue.TIMESERIES, command: 'ts.range', aggregation: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when aggregation is empty', () => {
        const query = getQuery({ type: QueryTypeValue.TIMESERIES, command: 'ts.range', aggregation: '' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({
          type: QueryTypeValue.TIMESERIES,
          command: 'ts.range',
          aggregation: '123',
          bucket: 123,
        });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.bucket);
      });

      it('Should call onBucketTextChange method when value was changed', () => {
        const query = getQuery({ type: QueryTypeValue.TIMESERIES, command: 'ts.range', aggregation: '123' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedMethod = jest.spyOn(wrapper.instance(), 'onBucketTextChange');
        wrapper.instance().forceUpdate();
        const testedComponent = getComponent(wrapper);
        const newValue = 1234;
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(testedMethod).toHaveBeenCalledWith({ target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          bucket: newValue,
        });
      });
    });

    /**
     * Fill Missing
     */
    describe('FillMissing', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'Switch' && node.prop('label') === 'Fill Missing';
        });

      it('Should be shown when command exists in commands.aggregation and commands.fill, aggregation value is selected and bucket is filled', () => {
        const query = getQuery({
          type: QueryTypeValue.TIMESERIES,
          command: 'ts.range',
          aggregation: '123',
          bucket: 123,
        });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when bucket is not filled', () => {
        const query = getQuery({
          type: QueryTypeValue.TIMESERIES,
          command: 'ts.range',
          aggregation: '123',
          bucket: 0,
        });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({
          type: QueryTypeValue.TIMESERIES,
          command: 'ts.range',
          aggregation: '123',
          bucket: 123,
          fill: true,
        });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('checked')).toEqual(query.fill);
      });

      it('Should set default value if query does not have fill value', () => {
        const query = getQuery({
          type: QueryTypeValue.TIMESERIES,
          command: 'ts.range',
          aggregation: '123',
          bucket: 123,
          fill: null,
        });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('checked')).toEqual(false);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({
          type: QueryTypeValue.TIMESERIES,
          command: 'ts.range',
          aggregation: '123',
          bucket: 123,
          fill: true,
        });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = false;
        testedComponent.simulate('change', { currentTarget: { checked: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          fill: newValue,
        });
      });
    });
  });

  /**
   * Streaming options
   */
  describe('Streaming fields', () => {
    /**
     * Streaming enabled
     */
    describe('Streaming', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'Switch' && node.prop('label') === 'Streaming';
        });

      it('Should be shown when refId=A', () => {
        const query = getQuery({ refId: 'A' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when refId!=A', () => {
        const query = getQuery({ refId: 'B' });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ refId: 'A', streaming: true });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('checked')).toEqual(query.streaming);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ refId: 'A', streaming: false });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = true;
        testedComponent.simulate('change', { currentTarget: { checked: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          streaming: newValue,
        });
      });
    });

    /**
     * Streaming interval
     */
    describe('Interval', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField' && node.prop('label') === 'Interval';
        });

      it('Should be shown when refId=A and streaming is checked', () => {
        const query = getQuery({ refId: 'A', streaming: true });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when refId=A and streaming is not checked', () => {
        const query = getQuery({ refId: 'A', streaming: false });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ refId: 'A', streaming: true, streamingInterval: 100 });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.streamingInterval);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ refId: 'A', streaming: true });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = '123';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          streamingInterval: parseInt(newValue, 10),
        });
      });
    });

    /**
     * Streaming capacity
     */
    describe('Capacity', () => {
      const getComponent = (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.name() === 'FormField' && node.prop('label') === 'Capacity';
        });

      it('Should be shown when refId=A and streaming is checked', () => {
        const query = getQuery({ refId: 'A', streaming: true });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).toBeTruthy();
      });

      it('Should not be shown when refId=A and streaming is not checked', () => {
        const query = getQuery({ refId: 'A', streaming: false });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.exists()).not.toBeTruthy();
      });

      it('Should set value from query', () => {
        const query = getQuery({ refId: 'A', streaming: true, streamingCapacity: 100 });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        expect(testedComponent.prop('value')).toEqual(query.streamingCapacity);
      });

      it('Should call onChange prop when value was changed', () => {
        const query = getQuery({ refId: 'A', streaming: true });
        const wrapper = shallow<QueryEditor>(
          <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
        );
        const testedComponent = getComponent(wrapper);
        const newValue = '123';
        testedComponent.simulate('change', { target: { value: newValue } });
        expect(onChange).toHaveBeenCalledWith({
          ...query,
          streamingCapacity: parseInt(newValue, 10),
        });
      });
    });
  });
});
