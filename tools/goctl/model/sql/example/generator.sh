#!/bin/bash

# generate usermodel with cache
goctl model -src ./sql/user.sql -dir ./model -c true

# generate usercoursemodel without cache
goctl model -src ./sql/course.sql -dir ./model
