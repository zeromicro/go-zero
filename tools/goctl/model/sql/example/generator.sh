#!/bin/bash

# generate model with cache from ddl
goctl model mysql ddl -src="./sql/*.sql" -dir="./sql/model/user" -c

# generate model with cache from data source
#user=root
#password=password
#datasource=127.0.0.1:3306
#database=test
#goctl model mysql datasource -url="${user}:${password}@tcp(${datasource})/${database}" -table="*" -dir ./model