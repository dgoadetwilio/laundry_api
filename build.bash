#!/usr/bin/env bash

API_NAME=laundry_api
echo "building ${API_NAME} for linux..."
set -x
env GOOS=linux GOARCH=arm go build -o "bin/${API_NAME}" -v main.go
EXIT_STATUS=${?}
set +x

if [[ "${EXIT_STATUS}" -eq 0 ]]; then 
  echo "${API_NAME} built successfully"
else
  echo "error building {API_NAME}"
fi
exit "${EXIT_STATUS}"
