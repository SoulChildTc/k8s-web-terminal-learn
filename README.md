### Kubernetes Dashboard Web Terminal 学习
核心代码来自于Kubernetes Dashboard

![image](https://raw.githubusercontent.com/SoulChildTc/k8s-web-terminal-learn/master/image/1.png)

相关文章


#### 安装依赖
```go
go mod tidy
```
#### 配置修改
app.js
```bash
# Namespace、Pod、Container、shell需要修改
/api/v1/pod/default/busybox-deployment-dcb89bc87-hrspk/shell/busybox?shell=sh
```
k8s_client.go
```go
// 修改kube config文件路径
```

#### 访问地址
```go
127.0.0.1:8090/web
```

