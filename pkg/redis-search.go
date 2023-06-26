package main

import (
  "strconv"

  "github.com/grafana/grafana-plugin-sdk-go/backend"
  "github.com/grafana/grafana-plugin-sdk-go/backend/log"
  "github.com/grafana/grafana-plugin-sdk-go/data"
  "github.com/redisgrafana/grafana-redis-datasource/pkg/models"
)

func queryFtSearch(qm queryModel, client redisClient) backend.DataResponse {
  response := backend.DataResponse{}

  var result interface{}
  args := []string{qm.Key}
  if qm.SearchQuery == "" {
    args = append(args, "*")
  } else {
    args = append(args, qm.SearchQuery)
  }

  if qm.ReturnFields != nil && len(qm.ReturnFields) > 0 {
    args = append(args, "RETURN")
    args = append(args, strconv.Itoa(len(qm.ReturnFields)))
    args = append(args, qm.ReturnFields...)
  }

  if qm.Count != 0 || qm.Offset > 0 {
    var count int
    if qm.Count == 0 {
      count = 10
    } else {
      count = qm.Count
    }
    args = append(args, "LIMIT", strconv.Itoa(qm.Offset), strconv.Itoa(count))
  }

  if qm.SortBy != "" {
    args = append(args, "SORTBY", qm.SortBy, qm.SortDirection)
  }

  err := client.RunCmd(&result, qm.Command, args...)

  if err != nil {
    return errorHandler(response, err)
  }

  frame := data.NewFrame("Results")
  fieldValuesMap := make(map[string][]string)

  fieldValuesMap["keyName"] = make([]string, len(result.([]interface{}))/2)

  for i := 1; i < len(result.([]interface{})); i += 2 {
    keyName := string((result.([]interface{}))[i].([]uint8))
    fieldValuesMap["keyName"][i/2] = keyName
    fieldValueArr := (result.([]interface{}))[i+1].([]interface{})

    for j := 0; j < len(fieldValueArr); j += 2 {
      fieldName := string(fieldValueArr[j].([]uint8))

      if _, ok := fieldValuesMap[fieldName]; !ok {
        fieldValuesMap[fieldName] = make([]string, len(result.([]interface{}))/2)
      }

      fieldValue := string(fieldValueArr[j+1].([]uint8))
      fieldValuesMap[fieldName][i/2] = fieldValue
    }
  }

  for fieldName, slice := range fieldValuesMap {
    frame.Fields = append(frame.Fields, data.NewField(fieldName, nil, slice))
  }

  response.Frames = append(response.Frames, frame)

  return response
}

/**
 * FT.INFO {index}
 *
 * @see https://oss.redislabs.com/redisearch/Commands/#ftinfo
 */
func queryFtInfo(qm queryModel, client redisClient) backend.DataResponse {
  response := backend.DataResponse{}

  // Execute command
  var result map[string]interface{}
  err := client.RunCmd(&result, qm.Command, qm.Key)

  // Check error
  if err != nil {
    return errorHandler(response, err)
  }

  // Create data frame response
  frame := data.NewFrame(qm.Key)

  // Add fields and values
  for key := range result {
    // Value
    switch value := result[key].(type) {
    case int64:
      // Add field
      field := data.NewField(key, nil, []int64{value})
      frame.Fields = append(frame.Fields, field)
    case []byte:
      // Parse Float
      if floatValue, err := strconv.ParseFloat(string(value), 64); err == nil {
        field := data.NewField(key, nil, []float64{floatValue})

        // Set unit
        if models.SearchInfoConfig[key] != "" {
          field.Config = &data.FieldConfig{Unit: models.SearchInfoConfig[key]}
        }

        frame.Fields = append(frame.Fields, field)
      } else {
        frame.Fields = append(frame.Fields, data.NewField(key, nil, []string{string(value)}))
      }
    case string:
      frame.Fields = append(frame.Fields, data.NewField(key, nil, []string{string(value)}))
    case []interface{}:
    default:
      log.DefaultLogger.Error(models.SearchInfo, "Conversion Error", "Unsupported Value type")
    }
  }

  // Add the frame to the response
  response.Frames = append(response.Frames, frame)

  // Return Response
  return response
}
