# FFR4MS-FrDocker

## 介绍

Fr-Docker是FFR4MS平台的容器监控模块，对微服务系统中各个微服务实例容器的通信消息和性能指标进行监控，实现快速检测定位故障微服务实例。

## 前提条件

### 微服务系统

微服务系统的扩展和部署请参考:

FFR4MS: [Gitee](https://gitee.com/zengtao321/ffr4ms) [GitHub](https://github.com/ztibeike/ffr4ms)

FFR4MS-Demo [Gitee](https://gitee.com/zengtao321/ffr4ms-demo) [GitHub](https://github.com/ztibeike/ffr4ms-demo)

### 环境配置

1. Golang v1.20
2. Pcap
```bash
apt install libpcap-dev
```
3. MongoDB
```bash
docker run --name frdocker-mongo --restart always -p 27017:27017 -d mongo --auth
```

## 安装

1. 配置MongoDB的用户名密码
```go
// config/db_config.go
MONGO_HOST = "localhost"
MONGO_PORT = 27017
MONGO_USER = "frdocker"
MONGO_PASS = "frdocker"
```
2. 编译安装
```bash
make && make install
```

## 使用
1. 查询微服务系统使用的网卡
```bash
ifconfig | grep br
export network="br-xxxxxxxxxxxx"
```
2. 指定Fr-Eureka注册中心地址
```bash
export registry="host:port"
```
3. 开启Fr-Docker
```bash
# 无日志颜色
frdocker frecovery -n ${network} -r ${registry}
# 启用日志颜色
frdocker frecovery -n ${network} -r ${registry} --color
```

