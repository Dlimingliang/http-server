# 模块十二作业

### 作业要求

1. httpserver以istioIngressGateway形式发布出来
2. 安全保证
3. 7层路由规则
4. open tracing


### 完成过程

1. 安装istio,并且在我们想要运行的环境启用envoy的sidecar注入.
```
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.12.0
cp bin/istioctl /usr/local/bin
istioctl install --set profile=demo -y
kubectl create ns httpserver
kubectl label ns httpserver istio-injection=enabled
```
2. 把我们之前的httpserver.yaml改到httpser的namespace.然后部署。会发现一个pod内会出现俩个容器
```
使用module12/httpser/目录下的httpserver.yaml
```
3. 将httpserver以istioIngressGateway形式发布




