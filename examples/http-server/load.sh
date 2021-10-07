#!/usr/bin/env bash

while true
do
#curl -XPOST http://localhost:8081/foo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d '{"name": "QuickBen"}'
#curl --silent -XPOST http://localhost:8081/bigfoo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d '{"name": "Anomandaris"}' > /dev/null
curl --silent -XPOST http://localhost:8081/bigfoorequest -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d @request.json > /dev/null
done