import React from 'react';
import { shallow, ShallowWrapper } from 'enzyme';
import { QueryEditor } from './query-editor';
import { AggregationValue, QueryTypeValue, RedisQuery } from '../../redis';
import { getQuery } from '../../tests/utils';

type ShallowComponent = ShallowWrapper<QueryEditor['props'], QueryEditor['state'], QueryEditor>;

/**
 * Query Field
 */
interface QueryFieldTest {
  name: keyof RedisQuery;
  testName?: string;
  getComponent: (wrapper: ShallowComponent) => ShallowWrapper;
  type: 'number' | 'string' | 'select' | 'switch' | 'radioButton';
  queryWhenShown: RedisQuery;
  queryWhenHidden: RedisQuery;
}

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
   * Run tests for query fields
   *
   * @param tests
   */
  const runQueryFieldsTest = (tests: QueryFieldTest[]) =>
    tests.forEach(({ name, getComponent, queryWhenShown, queryWhenHidden, type, testName = name }) => {
      describe(testName, () => {
        it('Should be shown', () => {
          const query = getQuery(queryWhenShown);
          const wrapper = shallow<QueryEditor>(
            <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
          );
          const testedComponent = getComponent(wrapper);
          expect(testedComponent.exists()).toBeTruthy();
        });

        it('Should not be shown', () => {
          const query = getQuery(queryWhenHidden);
          const wrapper = shallow<QueryEditor>(
            <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
          );
          const testedComponent = getComponent(wrapper);
          expect(testedComponent.exists()).not.toBeTruthy();
        });

        it('Should set value from query', () => {
          const query = getQuery({ [name]: 123, ...queryWhenShown });
          const wrapper = shallow<QueryEditor>(
            <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
          );
          const testedComponent = getComponent(wrapper);
          expect(testedComponent.prop(type === 'switch' ? 'checked' : 'value')).toEqual(query[name]);
        });

        it('Should call onChange prop when value was changed', () => {
          const query = getQuery(queryWhenShown);
          const wrapper = shallow<QueryEditor>(
            <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
          );
          const testedComponent = getComponent(wrapper);

          let newValue: any = '1234';
          if (type === 'number' || type === 'string') {
            testedComponent.simulate('change', { target: { value: newValue } });
          } else if (type === 'select') {
            testedComponent.simulate('change', { value: newValue });
          } else if (type === 'switch') {
            newValue = true;
            testedComponent.simulate('change', { currentTarget: { checked: newValue } });
          } else if (type === 'radioButton') {
            testedComponent.simulate('change', newValue);
          }
          expect(onChange).toHaveBeenCalledWith({
            ...query,
            [name]: type === 'number' ? parseInt(newValue, 10) : newValue,
          });
        });
      });
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

  runQueryFieldsTest([
    {
      name: 'query',
      getComponent: (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.prop('onChange') === wrapper.instance().onQueryChange;
        }),
      type: 'string',
      queryWhenShown: { refId: '', type: QueryTypeValue.CLI },
      queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS },
    },
    {
      name: 'command',
      getComponent: (wrapper: ShallowComponent) =>
        wrapper.findWhere((node) => {
          return node.prop('onChange') === wrapper.instance().onCommandChange;
        }),
      type: 'select',
      queryWhenShown: { refId: '', type: QueryTypeValue.REDIS },
      queryWhenHidden: { refId: '', type: QueryTypeValue.CLI },
    },
  ]);

  /**
   * Command properties
   */
  describe('Command fields', () => {
    runQueryFieldsTest([
      {
        name: 'keyName',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Key';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'get' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'get123' },
      },
      {
        name: 'filter',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Label Filter';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'ts.mrange' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'ts.mrange123' },
      },
      {
        name: 'field',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Field';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'hget' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'hget123' },
      },
      {
        name: 'legend',
        testName: 'Legend',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Legend';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'ts.range' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'ts.range123' },
      },
      {
        name: 'legend',
        testName: 'Legend Label',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Legend Label';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'ts.mrange' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'ts.mrange123' },
      },
      {
        name: 'value',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Value Label';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'ts.mrange' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'ts.mrange123' },
      },
      {
        name: 'size',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Size';
          }),
        type: 'number',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'slowlogGet' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'slowlogGet123' },
      },
      {
        name: 'cursor',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Cursor';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'tmscan' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'tmscan123' },
      },
      {
        name: 'match',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Match pattern';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'tmscan' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'tmscan123' },
      },
      {
        name: 'start',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Start';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'xrange' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'xrange123' },
      },
      {
        name: 'end',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'End';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'xrange' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'xrange123' },
      },
      {
        name: 'count',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Count';
          }),
        type: 'number',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'tmscan' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'tmscan123' },
      },
      {
        name: 'samples',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Samples';
          }),
        type: 'number',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'tmscan' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'tmscan123' },
      },
      {
        name: 'section',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onInfoSectionChange;
          }),
        type: 'select',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: 'info' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: 'info123' },
      },
      {
        name: 'aggregation',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onAggregationChange;
          }),
        type: 'select',
        queryWhenShown: { refId: '', type: QueryTypeValue.TIMESERIES, command: 'ts.range' },
        queryWhenHidden: { refId: '', type: QueryTypeValue.TIMESERIES, command: 'ts.range123' },
      },
      {
        name: 'bucket',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onBucketChange;
          }),
        type: 'number',
        queryWhenShown: {
          refId: '',
          type: QueryTypeValue.TIMESERIES,
          command: 'ts.range',
          aggregation: AggregationValue.AVG,
        },
        queryWhenHidden: { refId: '', type: QueryTypeValue.TIMESERIES, command: 'ts.range', aggregation: undefined },
      },
      {
        name: 'fill',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'Switch' && node.prop('label') === 'Fill Missing';
          }),
        type: 'switch',
        queryWhenShown: {
          refId: '',
          type: QueryTypeValue.TIMESERIES,
          command: 'ts.range',
          aggregation: AggregationValue.AVG,
          bucket: 123,
          fill: false,
        },
        queryWhenHidden: {
          refId: '',
          type: QueryTypeValue.TIMESERIES,
          command: 'ts.range',
          aggregation: AggregationValue.AVG,
          bucket: 0,
        },
      },
    ]);
  });

  /**
   * Streaming options
   */
  describe('Streaming fields', () => {
    runQueryFieldsTest([
      {
        name: 'streaming',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'Switch' && node.prop('label') === 'Streaming';
          }),
        type: 'switch',
        queryWhenShown: {
          refId: 'A',
          type: QueryTypeValue.TIMESERIES,
        },
        queryWhenHidden: {
          refId: 'B',
          type: QueryTypeValue.TIMESERIES,
        },
      },
      {
        name: 'streamingInterval',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Interval';
          }),
        type: 'number',
        queryWhenShown: {
          refId: 'A',
          type: QueryTypeValue.TIMESERIES,
          streaming: true,
        },
        queryWhenHidden: {
          refId: 'A',
          type: QueryTypeValue.TIMESERIES,
          streaming: false,
        },
      },
      {
        name: 'streamingCapacity',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Capacity';
          }),
        type: 'number',
        queryWhenShown: {
          refId: 'A',
          type: QueryTypeValue.TIMESERIES,
          streaming: true,
        },
        queryWhenHidden: {
          refId: 'A',
          type: QueryTypeValue.TIMESERIES,
          streaming: false,
        },
      },
      {
        name: 'streamingDataType',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onStreamingDataTypeChange;
          }),
        type: 'radioButton',
        queryWhenShown: {
          refId: 'A',
          type: QueryTypeValue.TIMESERIES,
          streaming: true,
        },
        queryWhenHidden: {
          refId: 'A',
          type: QueryTypeValue.TIMESERIES,
          streaming: false,
        },
      },
    ]);
  });
});
