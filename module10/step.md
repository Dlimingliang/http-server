# 模块十作业

### 作业要求

1. 为HTTPServer添加0-2秒的随机延时
2. 为HTTPServer项目添加延时Metric
3. 将HTTPServer部署至测试集群，并且完成Prometheus配置
4. 从Promethus界面中查询延时指标数据
5. 创建一个Grafana Dashboard展现延时分配情况


### 完成过程

1. 为HTTPServer添加0-2秒的随机延时
```
在defaultHandler中添加随机的timeSleep
```

2. 为HTTPServer项目添加延时Metric
```
添加对应的指标包，并且注册metrics
```

3. 将HTTPServer部署至测试集群，并且完成Prometheus配置
```
重新构建镜像，并且修改yaml中使用最新的镜像,并且添加prometheus的抓取标签
```
![aaa](./prmetheus-target.png)
