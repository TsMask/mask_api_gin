# 程序可用脚本

## 模块代理

中国最可靠的 Go 模块代理  

```shell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct
```

打包时改变编译平台

```shell
go env -w GOOS=linux
go env -w GOOS=windows
```

## 初始化数据库

- `db_init.sql` 初始化MySQL数据库数据

> **账号/密码**  
> **管理员**：maskAdmin/Admin@1234  
> **普通人员**：maskUser/User@1234  

## Docker 部署

- `Dockerfile` 构建Docker镜像脚本文件
- `Dockerfile分阶段` 构建Docker最小镜像
