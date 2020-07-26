# mapreduce用法

## Map

> channel是Map的返回值

由于Map是个并发操作，如果不用range或drain的方式，那么在使用返回值的时候，可能Map里面的代码还在读写这个返回值，可能导致数据不全或者`concurrent read write错误`

* 如果需要收集Map生成的结果，那么使用如下方式

	```
	for v := range channel {
		// v is with type interface{}
	}
	```

* 如果不需要收集结果，那么就需要显式的调用mapreduce.Drain，如

	```
	mapreduce.Drain(channel)
	```
	
## MapReduce

* mapper和reducer方法里可以调用cancel，调用了cancel之后返回值会是`nil, false`
* mapper里面如果有item不写入writer，那么这个item就不会被reduce收集
* mapper里面如果有处理item时panic，那么这个item也不会被reduce收集
* reduce是单线程，所有mapper出来的结果在这里串行处理
* reduce里面不写writer，或者panic，会导致返回`nil, false`