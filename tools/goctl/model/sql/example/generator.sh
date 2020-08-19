#!/bin/bash

# generate model with cache from ddl
goctl model -src ./sql/user.sql -dir ./model -c true

# generate model with cache from data source
goctl model datasource -url="user:password@tcp(127.0.0.1:3306)/database" -table="table1,table2"  -dir="./model"