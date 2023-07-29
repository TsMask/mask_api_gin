## 第一阶段 ====> 打包编译输出可执行文件
FROM golang:alpine AS builder

## 禁用cgo  
ENV CGO_ENABLED 0

## 设置构建目标
ENV GOOS linux

## 依赖下载代理
ENV GOPROXY https://goproxy.cn,direct

## 工作目录存放程序源码
WORKDIR /home/build

## 复制实际需要的文件到工作目录
COPY ./src ./src
COPY ./go.sum ./
COPY ./go.mod ./

## 安装程序依赖，需要编译
RUN go mod download

## 进行源码编译，生产文件 /home/mask_api
RUN go build -ldflags="-s -w" -o /home/mask_api ./main.go

## 第二阶段 ====> 构建可执行文件镜像
FROM alpine

## 工作目录存放程序源码
WORKDIR /home

## 安装时区工具
RUN apk add --no-cache tzdata

## 设置本地区时
ENV TZ="Asia/Shanghai"

## 将第一阶段文件复制到
COPY --from=builder /home/mask_api /home/mask_api

## 暴露端口要与程序端口一致
EXPOSE 6275

## 程序启动命令
CMD ["./mask_api", "--env", "prod"]
