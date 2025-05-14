#!/bin/bash

URL="http://localhost:8080/healthz"

STATUS=$(curl -s -o /dev/null -w "%{http_code}" -m 2 $URL)
if [ "$STATUS" -ne 200 ]; then
  echo "Service unhealthy!"
  exit 1
fi
