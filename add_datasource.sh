#!/bin/bash

curl 'http://admin:admin@16.214.19.214:3000/api/datasources' -X POST -H 'Content-Type: \
application/json;charset=UTF8' --data-binary \
'{"orgId":1,"name":"flux","type":"influxdb","typeLogoUrl":"public/app/plugins/datasource/influxdb/img/influxdb_logo.svg","access":"direct","url":"http://influxdb:8086","password":"root","user":"root","database":"environment","basicAuth":false,"basicAuthUser":"","basicAuthPassword":"","withCredentials":false,"isDefault":false}'
