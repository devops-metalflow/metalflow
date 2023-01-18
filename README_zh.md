# metalflow

[![Build Status](https://github.com/devops-metalflow/metalflow/workflows/ci/badge.svg?branch=main&event=push)](https://github.com/devops-metalflow/metalflow/actions?query=workflow%3Aci)
[![codecov](https://codecov.io/gh/devops-metalflow/metalflow/branch/main/graph/badge.svg?token=El8oiyaIsD)](https://codecov.io/gh/devops-metalflow/metalflow)
[![Go Report Card](https://goreportcard.com/badge/github.com/devops-metalflow/metalflow)](https://goreportcard.com/report/github.com/devops-metalflow/metalflow)
[![License](https://img.shields.io/github/license/devops-metalflow/metalflow.svg)](https://github.com/devops-metalflow/metalflow/blob/main/LICENSE)
[![Tag](https://img.shields.io/github/tag/devops-metalflow/metalflow.svg)](https://github.com/devops-metalflow/metalflow/tags)



> [English](README.md) | 中文



## 介绍

*metalflow* 作为 [metalweb](https://github.com/devops-metalflow/metalweb) 服务端，为其提供 REST API，并通过服务注册中心 `consul` 进行服务发现。



## 前提

- Go >= 1.18.0



## 准备

### [Consul](https://developer.hashicorp.com/consul/)

- **安装**

```bash
wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
apt update
apt install -y consul
```

- **运行**

```bash
consul agent -dev -ui -client=0.0.0.0
```



### MySQL

- **部署**

```bash
docker run -itd --name mysql-test -p 3306:3306 -e MYSQL_ROOT_PASSWORD=db_admin mysql:latest
```

- **初始化**

```bash
mysql -h 10.34.56.78 -u root -p
mysql> CREATE DATABASE metalflow;
```

### Redis

- **部署**

```bash
docker run -itd --name redis-test -p 6379:6379 redis:latest
```
或
```bash
docker run -itd --name redis-test -p 6379:6379 redis:latest --requirepass "123456"
```



## 运行

```bash
make build
./metalflow --config-file=config.yml
```

访问 [http://127.0.0.1:8089/api/ping](http://127.0.0.1:8089/api/ping) 可获取运行状态，例如：

```json
{
  "code": 201,
  "result": "pong",
  "msg": "操作成功"
}
```



## 用法

```
usage: metalflow [<flags>]

MetalBeat

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --version                  Show application version.
  --config-file=CONFIG-FILE  Config file (.yml)
```



## 配置

*metalflow* 相关配置参数见 [conf](https://github.com/devops-metalflow/metalflow/blob/main/initialize/conf)。

配置文件示例见 [config.yml](https://github.com/devops-metalflow/metalflow/blob/main/initialize/conf/config.prod.yml)。



## 协议

本项目协议声明见 [here](LICENSE)。



## 感谢

- [casbin](https://github.com/casbin/casbin): An authorization library that supports access control models like ACL, RBAC, ABAC in Golang.
- [Consul](https://github.com/hashicorp/consul): a distributed, highly available, and data center aware solution to connect and configure applications across dynamic, distributed infrastructure.
- [cronlib](https://github.com/rfyiamcool/cronlib):  golang crontab scheduler.
- [Machinery](https://github.com/RichardKnop/machinery): an asynchronous task queue/job queue based on distributed message passing.
- [Gin](https://github.com/gin-gonic/gin): a web framework written in Go (Golang).
- [gin-jwt](https://github.com/appleboy/gin-jwt): JWT Middleware for Gin framework.
- [gin-web](https://github.com/piupuer/gin-web): Golang version of RBAC permission management scaffolding.
- [Gorm](https://github.com/jinzhu/gorm): The fantastic ORM library for Golang.
- [logrus](https://github.com/sirupsen/logrus):  a structured logger for Go (golang), completely API compatible with the standard library logger.
- [lumberjack](https://github.com/natefinch/lumberjack):  a log rolling package for Go.
- [validator](https://github.com/go-playground/validator): Go Struct and Field validation, including Cross Field, Cross Struct, Map, Slice and Array diving.
- [viper](https://github.com/spf13/viper): Go configuration with fangs.
- [zap](https://github.com/uber-go/zap): Blazing fast, structured, leveled logging in Go.
