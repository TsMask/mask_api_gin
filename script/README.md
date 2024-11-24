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
- `db_init_demo.sql` 初始化demo模块数据表，不需要可不导入

> **账号/密码**  
> **系统管理员**：system/Abcd@1234..  
> **管理员**：admin/Abcd@1234..  
> **普通人员**：user/Abcd@1234..  

## Docker 部署

- `Dockerfile` 构建Docker镜像脚本文件
- `Dockerfile分阶段` 构建Docker最小镜像

```shell
# 构建
docker build --build-arg NAME=mask_api --build-arg VERSION=0.0.1 -t mask_api:0.0.1 .

# 启动
docker run -d \
--privileged=true \
--restart=always \
-p 8080:6275 \
-e TZ="Asia/Shanghai" \
-m 512M \
--name mask_api_000400 \
mask_api:0.0.1

```
