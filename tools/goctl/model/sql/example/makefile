#!/bin/bash

# generate model with cache from ddl
fromDDLWithCache:
	goctl template clean
	goctl model mysql ddl -src="./sql/*.sql" -dir="./sql/model/cache" -cache

fromDDLWithCacheAndIgnoreColumns:
	goctl template clean
	goctl model mysql ddl -src="./sql/*.sql" -dir="./sql/model/ignore_columns/cache" -cache -i 'gmt_create,create_at' -i 'gmt_modified,update_at'

fromDDLWithCacheAndDb:
	goctl template clean
	goctl model mysql ddl -src="./sql/*.sql" -dir="./sql/model/cache_db" -database="1gozero" -cache

fromDDLWithoutCache:
	goctl template clean;
	goctl model mysql ddl -src="./sql/*.sql" -dir="./sql/model/nocache"


# generate model with cache from data source
user=root
password=password
datasource=127.0.0.1:3306
database=gozero

fromDataSource:
	goctl template clean
	goctl model mysql datasource -url="$(user):$(password)@tcp($(datasource))/$(database)" -table="*" -dir ./model/cache -c -style gozero
