#!/bin/bash

# generate model with cache from ddl
goctl model mysql ddl -src="./sql/user.sql" -dir="./sql/model" -c

# generate model with cache from data source
goctl model mysql datasource -url="user:password@tcp(127.0.0.1:3306)/database" -table="table1,table2"  -dir="./model"