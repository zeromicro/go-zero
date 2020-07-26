# cpu监控准确度测试

1. 启动测试pod

   `make deploy`

2. 通过`kubectl get po -n adhoc`确认`sheeding` pod已经成功运行，通过如下命令进入pod

   `kubectl exec -it -n adhoc shedding -- sh`

3. 启动负载

   `/app # go-cpu-load -p 50 -c 1`

   默认`go-cpu-load`是对每个core加上负载的，所以测试里指定了`1000m`，等同于1 core，我们指定`-c 1`让测试更具有可读性

   `-p`可以多换几个值测试

4. 验证测试准确性

   `kubectl logs -f -n adhoc shedding`

   可以看到日志中的`CPU`报告，`1000m`表示`100%`，如果看到`500m`则表示`50%`，每分钟输出一次

   `watch -n 5 kubectl top pod -n adhoc`

   可以看到`kubectl`报告的`CPU`使用率，两者进行对比，即可知道是否准确

