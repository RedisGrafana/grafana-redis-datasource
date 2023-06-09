#!/bin/bash

go test ./pkg/... -tags=clusterIntegration -v
last=$?
if [[ $last != 0 ]]; then
  echo "exit code for test: " $last
  exit $last
fi

exit 0
