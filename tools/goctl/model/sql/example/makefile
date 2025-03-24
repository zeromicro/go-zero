#!/bin/bash

# generate model with cache from ddl
fromDDLWithCache:
	goctl template clean
	goctl model mysql ddl -src="./sql/*.sql" -dir="./sql/model/cache" -cache --prefix gozero

fromDDLWithCacheAndIgnoreColumns:
	goctl template clean
	goctl model mysql ddl -src="./sql/*.sql" -dir="./sql/model/ignore_columns/cache" -cache -i 'gmt_create,create_at' -i 'gmt_modified,update_at' --prefix gozero

fromDDLWithCacheAndDb:
	goctl template clean
	goctl model mysql ddl -src="./sql/*.sql" -dir="./sql/model/cache_db" -database="1gozero" -cache --prefix gozero

fromDDLWithoutCache:
	goctl template clean;
	goctl model mysql ddl -src="./sql/*.sql" -dir="./sql/model/nocache"


# generate model with cache from data source
user=app_user
password=PLO75FbcfmFYRuQEGmygZ9PyQCQbmgeD5
datasource=k8s-istiosys-unifydev-2cdcebd306-b64e09be84220820.elb.ap-southeast-1.amazonaws.com:3306
database=db_fiat

fromDataSource:
	goctl template clean
	goctl model mysql datasource -url="$(user):$(password)@tcp($(datasource))/$(database)" -table="biz_switch" -dir ./model/cache -c -style gozero --prefix gozero
