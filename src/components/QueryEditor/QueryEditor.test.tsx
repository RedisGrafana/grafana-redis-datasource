import { configure, shallow, ShallowWrapper, mount } from 'enzyme';
import React from 'react';
import { RedisGraph } from 'redis/graph';
import { SelectableValue } from '@grafana/data';
import { Input } from '@grafana/ui';
import {
  AggregationValue,
  QueryTypeCli,
  QueryTypeValue,
  Redis,
  RedisGears,
  RedisJson,
  RedisQuery,
  RedisTimeSeries,
} from '../../redis';
import { getQuery } from '../../tests/utils';
import { QueryEditor } from './QueryEditor';
import { RediSearch } from '../../redis/search';
import Adapter from "@wojtekmaj/enzyme-adapter-react-17"
import sinon from 'sinon';


configure({adapter: new Adapter()});



type ShallowComponent = ShallowWrapper<QueryEditor['props'], QueryEditor['state'], QueryEditor>;

/**
 * Data Source
 */
const dataSourceMock = {
  name: 'datasource',
};
const dataSourceInstanceSettingsMock = {
  jsonData: { cliDisabled: false },
};

const dataSourceSrvGetMock = jest.fn().mockImplementation(() => Promise.resolve(dataSourceMock));
const dataSourceSrvGetInstanceSettingsMock = jest.fn().mockImplementation(() => dataSourceInstanceSettingsMock);

jest.mock('@grafana/runtime', () => ({
  getDataSourceSrv: () => ({
    get: dataSourceSrvGetMock,
    getInstanceSettings: dataSourceSrvGetInstanceSettingsMock,
  }),
}));

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

    it('Should set value from query to CLI', () => {
      const query = getQuery({ type: QueryTypeValue.CLI });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(query.type);
      const cli = testedComponent
        .prop('options')
        .filter((o: SelectableValue<QueryTypeValue>) => o.value === QueryTypeValue.CLI);
      expect(cli).toEqual([QueryTypeCli]);
    });

    it('CLI should not be seen if disabled', () => {
      dataSourceInstanceSettingsMock.jsonData.cliDisabled = true;
      const query = getQuery({ type: QueryTypeValue.REDIS });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.prop('value')).toEqual(query.type);
      const cli = testedComponent
        .prop('options')
        .filter((o: SelectableValue<QueryTypeValue>) => o.value === QueryTypeValue.CLI);
      expect(cli).toEqual([]);
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
   * Streaming
   */
  describe('Streaming', () => {
    const getComponent = (wrapper: ShallowComponent) =>
      wrapper.findWhere((node) => {
        return node.name() === 'Switch' && node.prop('label') === 'Streaming';
      });

    it('Should set value from query', () => {
      const query = getQuery({ streaming: true });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedComponent = getComponent(wrapper);
      expect(testedComponent.prop('checked')).toEqual(true);
    });

    it('Should call onStreamingChange when onChange prop was called', () => {
      const query = getQuery({ streaming: true });
      const wrapper = shallow<QueryEditor>(
        <QueryEditor datasource={{} as any} query={query} onRunQuery={onRunQuery} onChange={onChange} />
      );
      const testedMethod = jest.spyOn(wrapper.instance(), 'onStreamingChange');
      wrapper.instance().forceUpdate();
      const testedComponent = getComponent(wrapper);
      const newValue = true;
      testedComponent.simulate('change', { currentTarget: { checked: newValue } });
      expect(testedMethod).toHaveBeenCalledWith({ currentTarget: { checked: newValue } });
      expect(onChange).toHaveBeenCalledWith({
        ...query,
        streaming: newValue,
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
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.GET },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'keyName',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Function';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: RedisGears.PYEXECUTE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'filter',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Label Filter';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: RedisTimeSeries.MRANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'field',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Field';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.HGET },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'legend',
        testName: 'Legend',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Legend';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: RedisTimeSeries.RANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'legend',
        testName: 'Legend Label',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Legend Label';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: RedisTimeSeries.MRANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'value',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Value Label';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: RedisTimeSeries.MRANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'cypher',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onCypherChange;
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.GRAPH, command: RedisGraph.QUERY },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'offset',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onOffsetChange;
          }),
        type: 'number',
        queryWhenShown: { refId: '', type: QueryTypeValue.SEARCH, command: RediSearch.SEARCH },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'searchQuery',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onSearchQueryChange;
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.SEARCH, command: RediSearch.SEARCH },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'path',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onPathChange;
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.JSON, command: RedisJson.GET },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'size',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Size';
          }),
        type: 'number',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.SLOWLOG_GET },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'cursor',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Cursor';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.TMSCAN },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'match',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Match pattern';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.TMSCAN },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'start',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Start';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.XRANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'end',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'End';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.XRANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'min',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Minimum';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.ZRANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'max',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Maximum';
          }),
        type: 'string',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.ZRANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'count',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Count';
          }),
        type: 'number',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.TMSCAN },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'samples',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.name() === 'FormField' && node.prop('label') === 'Samples';
          }),
        type: 'number',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.TMSCAN },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'section',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onInfoSectionChange;
          }),
        type: 'select',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.GET },
      },
      {
        name: 'aggregation',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onAggregationChange;
          }),
        type: 'select',
        queryWhenShown: { refId: '', type: QueryTypeValue.TIMESERIES, command: RedisTimeSeries.RANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
      },
      {
        name: 'zrangeQuery',
        getComponent: (wrapper: ShallowComponent) =>
          wrapper.findWhere((node) => {
            return node.prop('onChange') === wrapper.instance().onZRangeQueryChange;
          }),
        type: 'select',
        queryWhenShown: { refId: '', type: QueryTypeValue.REDIS, command: Redis.ZRANGE },
        queryWhenHidden: { refId: '', type: QueryTypeValue.REDIS, command: Redis.INFO },
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
          command: RedisTimeSeries.RANGE,
          aggregation: AggregationValue.AVG,
        },
        queryWhenHidden: {
          refId: '',
          type: QueryTypeValue.TIMESERIES,
          command: RedisTimeSeries.RANGE,
          aggregation: undefined,
        },
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
          command: RedisTimeSeries.RANGE,
          aggregation: AggregationValue.AVG,
          bucket: 123,
          fill: false,
        },
        queryWhenHidden: {
          refId: '',
          type: QueryTypeValue.TIMESERIES,
          command: RedisTimeSeries.RANGE,
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
