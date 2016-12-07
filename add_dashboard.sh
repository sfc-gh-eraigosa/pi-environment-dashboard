#!/bin/bash

export dashboard_json=$(cat dashboard.json|tr "\n" " ")
dashboard_json='{"dashboard":'"$dashboard_json"'}'

curl http://admin:admin@16.214.19.214:3000/api/dashboards/db -X POST -H 'Content-Type: application/json;charset=UTF8' --data-binary ''"$dashboard_json"''
