#!/usr/bin/env bash

curl --silent -XPOST http://localhost:8081/bigfoo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d '{"name": "Anomandaris", "iterations": 547}' > request16K.json
sleep 1
curl --silent -XPOST http://localhost:8081/bigfoo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d '{"name": "Anomandaris", "iterations": 4300}' > request128K.json
sleep 1
curl --silent -XPOST http://localhost:8081/bigfoo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d '{"name": "Anomandaris", "iterations": 8600}' > request256K.json
sleep 1
curl --silent -XPOST http://localhost:8081/bigfoo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d '{"name": "Anomandaris", "iterations": 12900}' > request384K.json
sleep 1
curl --silent -XPOST http://localhost:8081/bigfoo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d '{"name": "Anomandaris", "iterations": 17200}' > request512K.json
sleep 1
curl --silent -XPOST http://localhost:8081/bigfoo -H 'Accept: application/json, */*;q=0.5' -H 'Content-type: application/json'  -d '{"name": "Anomandaris", "iterations": 34400}' > request1024K.json
sleep 1