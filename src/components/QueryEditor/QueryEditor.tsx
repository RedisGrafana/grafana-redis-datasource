import React, { ChangeEvent, PureComponent } from 'react';
import { RedisGraph } from 'redis/graph';
import { css } from '@emotion/css';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { getDataSourceSrv } from '@grafana/runtime';
import {
  Button,
  FieldArray,
  Form,
  InlineFormLabel,
  Input,
  LegacyForms,
  RadioButtonGroup,
  Select,
  TextArea,
} from '@grafana/ui';
import { StreamingDataType, StreamingDataTypes } from '../../constants';
import { DataSource } from '../../datasource';
import {
  Aggregations,
  AggregationValue,
  CommandParameters,
  Commands,
  InfoSections,
  InfoSectionValue,
  QueryType,
  QueryTypeCli,
  QueryTypeValue,
  Redis,
  RedisGears,
  RedisJson,
  RedisQuery,
  RedisTimeSeries,
  Reducers,
  ReducerValue,
  ZRangeQuery,
  ZRangeQueryValue,
} from '../../redis';
import { RedisDataSourceOptions } from '../../types';
import { RediSearch, SortDirection, SortDirectionValue } from '../../redis/search';
import { FieldValuesContainer } from '../../redis/fieldValuesContainer';

/**
 * Form Field
 */
const { FormField, Switch } = LegacyForms;

/**
 * Editor Property
 */
type Props = QueryEditorProps<DataSource, RedisQuery, RedisDataSourceOptions>;

/**
 * Query Editor
 */
export class QueryEditor extends PureComponent<Props> {
  /**
   * Change handler for number field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createNumberFieldHandler = (name: keyof RedisQuery) => (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({ ...this.props.query, [name]: Number(event.target.value) });
  };

  /**
   * Change handler for text field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createTextFieldHandler = (name: keyof RedisQuery) => (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({ ...this.props.query, [name]: event.target.value });
  };

  /**
   * Change handler for textarea field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createTextareaFieldHandler = (name: keyof RedisQuery) => (event: ChangeEvent<HTMLTextAreaElement>) => {
    this.props.onChange({ ...this.props.query, [name]: event.target.value });
  };

  /**
   * Change handler for select field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createSelectFieldHandler<ValueType>(name: keyof RedisQuery) {
    return (val: SelectableValue<ValueType>) => {
      this.props.onChange({ ...this.props.query, [name]: val.value });
    };
  }

  /**
   * Change handler for radio button field
   *
   * @param {value: ValueType}
   */
  createRedioButtonFieldHandler<ValueType>(name: keyof RedisQuery) {
    return (value?: ValueType) => {
      this.props.onChange({ ...this.props.query, [name]: value });
    };
  }

  /**
   * Change handler for switch field
   *
   * @param {ChangeEvent<HTMLInputElement>} event Event
   */
  createSwitchFieldHandler = (name: keyof RedisQuery) => (event: React.SyntheticEvent<HTMLInputElement>) => {
    this.props.onChange({ ...this.props.query, [name]: event.currentTarget.checked });
  };

  createFieldArrayHandler = (name: 'returnFields') => (event: React.SyntheticEvent<HTMLInputElement>) => {
    const index = Number(event.currentTarget.name.split(':', 2)[1]);
    if (!this.props.query[name]) {
      this.props.query[name] = [];
    }

    this.props.query[name]![index] = event.currentTarget.value;
  };

  /**
   * Key name change
   */
  onKeyNameChange = this.createTextFieldHandler('keyName');

  /**
   * Query change
   */
  onQueryChange = this.createTextareaFieldHandler('query');

  /**
   * searchQuery change
   */
  onSearchQueryChange = this.createTextareaFieldHandler('searchQuery');

  /**
   * Filter change
   */
  onFilterChange = this.createTextFieldHandler('filter');

  /**
   * Cypher change
   */
  onCypherChange = this.createTextareaFieldHandler('cypher');

  /**
   * Path change
   */
  onPathChange = this.createTextareaFieldHandler('path');

  /**
   * Field change
   */
  onFieldChange = this.createTextFieldHandler('field');

  /**
   * Legend change
   */
  onLegendChange = this.createTextFieldHandler('legend');

  /**
   * Value change
   */
  onValueChange = this.createTextFieldHandler('value');

  /**
   * sortBy change
   */
  onSortByChange = this.createTextFieldHandler('sortBy');

  /**
   * Command change
   */
  onCommandChange = this.createSelectFieldHandler<string>('command');

  /**
   * Type change
   *
   * @param val Value
   */
  onTypeChange = (val: SelectableValue<string>) => {
    const { onChange, query } = this.props;
    onChange({
      ...query,
      type: val.value as QueryTypeValue,
      query: '',
      command: '',
    });
  };

  /**
   * Ts Reducer Change
   */
  onTsReducerChange = this.createSelectFieldHandler<ReducerValue>('tsReducer');

  /**
   * Group By Change
   */
  onTsGroupByLabelChange = this.createTextFieldHandler('tsGroupByLabel');

  /**
   * Aggregation change
   */
  onAggregationChange = this.createSelectFieldHandler<AggregationValue>('aggregation');

  /**
   * LATEST change
   */

  onLatestChange = this.createSwitchFieldHandler('tsLatest');

  /**
   * ZRANGE Query change
   */
  onZRangeQueryChange = this.createSelectFieldHandler<ZRangeQueryValue>('zrangeQuery');

  /**
   * FT.SEARCH Sort By
   */
  onSortDirectionChange = this.createSelectFieldHandler<SortDirectionValue>('sortDirection');

  /**
   * Info section change
   */
  onInfoSectionChange = this.createSelectFieldHandler<InfoSectionValue>('section');

  /**
   * Bucket change
   */
  onBucketChange = this.createNumberFieldHandler('bucket');

  /**
   * Size change
   */
  onSizeChange = this.createNumberFieldHandler('size');

  /**
   * Count change
   */
  onCountChange = this.createNumberFieldHandler('count');

  /**
   * Count change
   */
  onOffsetChange = this.createNumberFieldHandler('offset');

  /**
   * Samples change
   */
  onSamplesChange = this.createNumberFieldHandler('samples');

  /**
   * Cursor change
   */
  onCursorChange = this.createTextFieldHandler('cursor');

  /**
   * Match change
   */
  onMatchChange = this.createTextFieldHandler('match');

  /**
   * Start change
   */
  onStartChange = this.createTextFieldHandler('start');

  /**
   * End change
   */
  onEndChange = this.createTextFieldHandler('end');

  /**
   * Min change
   */
  onMinChange = this.createTextFieldHandler('min');

  /**
   * Max change
   */
  onMaxChange = this.createTextFieldHandler('max');

  /**
   * Fill change
   */
  onFillChange = this.createSwitchFieldHandler('fill');

  /**
   * Streaming change
   */
  onStreamingChange = this.createSwitchFieldHandler('streaming');

  /**
   * Streaming interval change
   */
  onStreamingIntervalChange = this.createNumberFieldHandler('streamingInterval');

  /**
   * Streaming capacity change
   */
  onStreamingCapacityChange = this.createNumberFieldHandler('streamingCapacity');

  /**
   * Streaming data type change
   */
  onStreamingDataTypeChange = this.createRedioButtonFieldHandler<StreamingDataType>('streamingDataType');

  onReturnFieldChange = this.createFieldArrayHandler('returnFields');

  /**
   * Render Editor
   */
  render() {
    const defaultValues: FieldValuesContainer = {
      fieldArray: [''],
    };

    const {
      keyName,
      aggregation,
      zrangeQuery,
      bucket,
      legend,
      command,
      field,
      filter,
      cypher,
      path,
      value,
      query,
      searchQuery,
      offset,
      sortDirection,
      sortBy,
      type,
      section,
      size,
      fill,
      cursor,
      count,
      match,
      samples,
      start,
      end,
      min,
      max,
      streaming,
      streamingInterval,
      streamingCapacity,
      streamingDataType,
      tsGroupByLabel,
      tsReducer,
      tsLatest,
    } = this.props.query;
    const { onRunQuery, datasource } = this.props;

    /**
     * Check if CLI disabled
     */
    const jsonData = getDataSourceSrv().getInstanceSettings(datasource.uid)?.jsonData as RedisDataSourceOptions;
    const cliDisabled = jsonData?.cliDisabled;

    /**
     * Return
     */
    return (
      <div>
        <div className="gf-form">
          <InlineFormLabel width={8}>Type</InlineFormLabel>
          <Select
            className={css`
              margin-right: 5px;
            `}
            width={40}
            options={cliDisabled ? QueryType : [...QueryType, QueryTypeCli]}
            menuPlacement="bottom"
            value={type}
            onChange={this.onTypeChange}
          />

          {type === QueryTypeValue.CLI && (
            <>
              <InlineFormLabel width={8}>Command</InlineFormLabel>
              <TextArea value={query} className="gf-form-input" onChange={this.onQueryChange} />
            </>
          )}
          {type && type !== QueryTypeValue.CLI && (
            <>
              <InlineFormLabel width={8}>Command</InlineFormLabel>
              <Select options={Commands[type]} menuPlacement="bottom" value={command} onChange={this.onCommandChange} />
            </>
          )}
        </div>

        {type !== QueryTypeValue.CLI && command && (
          <div className="gf-form">
            {CommandParameters.keyName.includes(command as Redis) && (
              <FormField
                labelWidth={8}
                inputWidth={30}
                value={keyName}
                onChange={this.onKeyNameChange}
                label="Key"
                tooltip="Key name"
              />
            )}

            {CommandParameters.filter.includes(command as RedisTimeSeries) && (
              <FormField
                labelWidth={8}
                inputWidth={30}
                value={filter}
                onChange={this.onFilterChange}
                label="Label Filter"
                tooltip="Whenever filters need to be provided, a minimum of one l=v filter must be applied.
                The list of possible filters:
                https://oss.redislabs.com/redistimeseries/commands/#filtering"
              />
            )}

            {CommandParameters.field.includes(command as Redis) && (
              <FormField labelWidth={8} inputWidth={30} value={field} onChange={this.onFieldChange} label="Field" />
            )}

            {CommandParameters.legend.includes(command as RedisTimeSeries) && (
              <FormField
                labelWidth={8}
                inputWidth={20}
                value={legend}
                onChange={this.onLegendChange}
                label="Legend"
                tooltip="Legend override"
              />
            )}

            {CommandParameters.value.includes(command as RedisTimeSeries) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={value}
                onChange={this.onValueChange}
                label="Value"
                tooltip="Value override"
              />
            )}

            {CommandParameters.legendLabel.includes(command as RedisTimeSeries) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={legend}
                onChange={this.onLegendChange}
                label="Legend Label"
              />
            )}

            {CommandParameters.valueLabel.includes(command as RedisTimeSeries) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={value}
                onChange={this.onValueChange}
                label="Value Label"
              />
            )}

            {CommandParameters.size.includes(command as Redis) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={size}
                type="number"
                onChange={this.onSizeChange}
                label="Size"
              />
            )}

            {CommandParameters.cursor.includes(command as Redis) && (
              <FormField labelWidth={8} inputWidth={10} value={cursor} onChange={this.onCursorChange} label="Cursor" />
            )}

            {CommandParameters.match.includes(command as Redis) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={match}
                onChange={this.onMatchChange}
                placeholder="*"
                label="Match pattern"
              />
            )}
          </div>
        )}

        {command && CommandParameters.cypher.includes(command as RedisGraph) && (
          <div className="gf-form">
            <InlineFormLabel
              tooltip="The syntax is based on Cypher, and only a subset of the language currently supported: \
                https://oss.redislabs.com/redisgraph/commands/#query-language"
              width={8}
            >
              Cypher
            </InlineFormLabel>
            <TextArea value={cypher} className="gf-form-input" onChange={this.onCypherChange} />
          </div>
        )}

        {command && CommandParameters.pyFunction.includes(command as RedisGears) && (
          <div className="gf-form">
            <FormField
              labelWidth={8}
              inputWidth={30}
              value={keyName}
              onChange={this.onKeyNameChange}
              label="Function"
            />
          </div>
        )}

        {command && CommandParameters.path.includes(command as RedisJson) && (
          <div className="gf-form">
            <InlineFormLabel
              tooltip="RedisJSON's syntax is a subset of common best practices and resembles JSONPath not by accident."
              width={8}
            >
              Path
            </InlineFormLabel>
            <TextArea value={path} className="gf-form-input" onChange={this.onPathChange} />
          </div>
        )}

        {command && CommandParameters.searchQuery.includes(command as RediSearch) && (
          <div className="gf-form">
            <InlineFormLabel tooltip="The RediSearch Query to issue to the index." width={10}>
              Query
            </InlineFormLabel>
            <TextArea value={searchQuery} className="gf-form-input" onChange={this.onSearchQueryChange} />
          </div>
        )}

        {command && CommandParameters.returnFields.includes(command as RediSearch) && (
          <Form id="returnFieldsForm" onSubmit={() => true} defaultValues={defaultValues}>
            {({ control }) => (
              <div className="gf-form">
                <InlineFormLabel
                  tooltip="Add return Fields to query to minimize what needs to be pulled back from Redis."
                  width={10}
                >
                  Return Fields
                </InlineFormLabel>
                <FieldArray name="returnFields" control={control}>
                  {({ fields, append }) => (
                    <>
                      <Button id="addReturnFieldButton" onClick={() => append({})}>
                        Add Return Field
                      </Button>
                      <div id="returnFieldInputs">
                        {fields.map((field, index) => (
                          <Input name={`returnField:${index}`} key={field.id} onChange={this.onReturnFieldChange} />
                        ))}
                      </div>
                    </>
                  )}
                </FieldArray>
              </div>
            )}
          </Form>
        )}

        {command && CommandParameters.offset.includes(command as RediSearch) && (
          <FormField
            labelWidth={8}
            inputWidth={10}
            value={offset}
            type="number"
            onChange={this.onOffsetChange}
            label="Offset"
            defaultValue="0"
          />
        )}

        {command && CommandParameters.limit.includes(command as RediSearch) && (
          <FormField
            labelWidth={8}
            inputWidth={10}
            value={count}
            type="number"
            onChange={this.onCountChange}
            label="Limit"
            defaultValue="10"
          />
        )}

        {command && CommandParameters.sortBy.includes(command as RediSearch) && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Sort Direction</InlineFormLabel>
            <Select
              onChange={this.onSortDirectionChange}
              options={SortDirection}
              width={20}
              value={sortDirection}
              defaultValue={SortDirectionValue.NONE}
            />
            {sortDirection && sortDirection !== SortDirectionValue.NONE && (
              <div>
                <FormField
                  labelWidth={8}
                  inputWidth={10}
                  value={sortBy}
                  onChange={this.onSortByChange}
                  type="string"
                  label={'Sort By'}
                />
              </div>
            )}
          </div>
        )}

        {type === QueryTypeValue.REDIS && command && CommandParameters.zrangeQuery.includes(command as Redis) && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Range Query</InlineFormLabel>
            <Select
              className={css`
                margin-right: 5px;
              `}
              options={ZRangeQuery}
              width={30}
              onChange={this.onZRangeQueryChange}
              value={zrangeQuery}
              menuPlacement="bottom"
            />
          </div>
        )}

        {type !== QueryTypeValue.CLI && command && (
          <div className="gf-form">
            {CommandParameters.start.includes(command as Redis) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={start}
                onChange={this.onStartChange}
                placeholder="time range"
                label="Start"
                tooltip="Based on the selected time range, if not specified"
              />
            )}

            {CommandParameters.end.includes(command as Redis) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={end}
                onChange={this.onEndChange}
                placeholder="time range"
                label="End"
                tooltip="Based on the selected time range, if not specified"
              />
            )}

            {CommandParameters.min.includes(command as Redis) && (
              <FormField labelWidth={8} inputWidth={10} value={min} onChange={this.onMinChange} label="Minimum" />
            )}

            {CommandParameters.max.includes(command as Redis) && (
              <FormField labelWidth={8} inputWidth={10} value={max} onChange={this.onMaxChange} label="Maximum" />
            )}

            {CommandParameters.count.includes(command as Redis) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={count}
                type="number"
                onChange={this.onCountChange}
                label="Count"
              />
            )}

            {CommandParameters.samples.includes(command as Redis) && (
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={samples}
                type="number"
                onChange={this.onSamplesChange}
                label="Samples"
                placeholder="5"
                tooltip="Number of sampled nested values. 0 (all values) is not supported."
              />
            )}
          </div>
        )}

        {type === QueryTypeValue.REDIS && command && CommandParameters.section.includes(command as Redis) && (
          <div className="gf-form">
            <InlineFormLabel width={8}>Section</InlineFormLabel>
            <Select options={InfoSections} onChange={this.onInfoSectionChange} value={section} menuPlacement="bottom" />
          </div>
        )}

        {type === QueryTypeValue.TIMESERIES &&
          command &&
          CommandParameters.aggregation.includes(command as RedisTimeSeries) && (
            <div className="gf-form">
              <InlineFormLabel width={8}>Aggregation</InlineFormLabel>
              <Select
                className={css`
                  margin-right: 5px;
                `}
                options={Aggregations}
                width={30}
                onChange={this.onAggregationChange}
                value={aggregation}
                menuPlacement="bottom"
              />
              {aggregation && (
                <FormField
                  labelWidth={8}
                  value={bucket}
                  type="number"
                  onChange={this.onBucketChange}
                  label="Time Bucket"
                  tooltip="Time bucket for aggregation in milliseconds"
                />
              )}
              {aggregation && bucket && CommandParameters.fill.includes(command as RedisTimeSeries) && (
                <Switch
                  label="Fill Missing"
                  labelClass="width-10"
                  tooltip="If checked, the datasource will fill missing intervals."
                  checked={fill || false}
                  onChange={this.onFillChange}
                />
              )}
            </div>
          )}

        {type === QueryTypeValue.TIMESERIES &&
          command &&
          CommandParameters.tsGroupBy.includes(command as RedisTimeSeries) && (
            <div className="gf-form">
              <FormField
                labelWidth={8}
                inputWidth={10}
                value={tsGroupByLabel}
                onChange={this.onTsGroupByLabelChange}
                label="Group By"
                tooltip="The label to group your time-series by for your reduction"
              />
              {tsGroupByLabel && (
                <Select
                  options={Reducers}
                  width={30}
                  onChange={this.onTsReducerChange}
                  value={tsReducer}
                  menuPlacement="bottom"
                />
              )}
            </div>
          )}

        {type === QueryTypeValue.TIMESERIES &&
          command &&
          CommandParameters.tsLatest.includes(command as RedisTimeSeries) && (
            <div className="gf-form">
              <Switch
                label="Latest"
                labelClass="width-8"
                tooltip="If checked, will return the latest (incomplete bucket)."
                checked={tsLatest || false}
                onChange={this.onLatestChange}
              />
            </div>
          )}

        <div className="gf-form">
          <Switch
            label="Streaming"
            labelClass="width-8"
            tooltip="If checked, the datasource will stream datas."
            checked={streaming || false}
            onChange={this.onStreamingChange}
          />
          {streaming && (
            <>
              <FormField
                labelWidth={8}
                value={streamingInterval}
                type="number"
                onChange={this.onStreamingIntervalChange}
                label="Interval"
                tooltip="Streaming interval in milliseconds. Default is 1000ms. For multiple Streaming targets minimum value will be taken."
                placeholder="1000"
              />
              <FormField
                labelWidth={8}
                value={streamingCapacity}
                type="number"
                onChange={this.onStreamingCapacityChange}
                label="Capacity"
                tooltip="Values will be constantly added and will never exceed the given capacity. Default is 1000."
                placeholder="1000"
              />
            </>
          )}
        </div>

        {streaming && (
          <div className="gf-form">
            <InlineFormLabel width={8} tooltip="If checked Time series, the last line of data will be applied.">
              Data type
            </InlineFormLabel>
            <RadioButtonGroup
              options={StreamingDataTypes}
              value={streamingDataType || StreamingDataType.TIMESERIES}
              onChange={this.onStreamingDataTypeChange}
            />
            {streamingDataType !== StreamingDataType.DATAFRAME && (
              <FormField
                className={css`
                  margin-left: 5px;
                `}
                labelWidth={8}
                inputWidth={30}
                value={field}
                tooltip="Specify field(s) to return from backend. Required for Alerting to trigger on specific fields."
                onChange={this.onFieldChange}
                label="Filter Field"
              />
            )}
          </div>
        )}

        <Button onClick={onRunQuery}>Run</Button>
      </div>
    );
  }
}
