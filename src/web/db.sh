#!/usr/bin/env bash
HOST=${1:-"localhost:8080"}

curl -XPOST -H "Content-Type: application/json" http://"${HOST}"/api/v1/user -d '{"name":"marius","alias":"Marius Magureanu","email":"marius@emailc.om"}'
curl -XPOST -H "Content-Type: application/json" http://"${HOST}"/api/v1/user -d '{"name":"anca","alias":"Anca Molodet","email":"anca@emailc.om"}'
curl -XPOST -H "Content-Type: application/json" http://"${HOST}"/api/v1/user -d '{"name":"marin","alias":"Marian Zeic","email":"marin@emailc.om"}'

