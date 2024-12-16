FROM golang:1.22-alpine

## 编译参数 项目名和版本号
ARG NAME
ARG VERSION

## 禁用cgo  
ENV CGO_ENABLED 0
## 设置构建目标
ENV GOOS linux
## 依赖下载代理
ENV GOPROXY https://goproxy.cn,direct

# 设置时区和语言环境
ENV TZ="Asia/Shanghai"
ENV LANG="en_US.UTF-8"

## 工作目录存放程序源码
WORKDIR /app

## 复制实际需要的文件到工作目录
COPY ./src ./src
COPY ./go.sum ./
COPY ./go.mod ./
COPY ./main.go ./
COPY ./src/config/config.prod.yaml ./config.yaml

## 安装程序依赖，需要编译
RUN go mod download

## 进行源码编译，生产文件 app
RUN go build -o app -ldflags="-s -w \ 
    -X 'mask_api_gin/src/framework/config.Name=${NAME}' \ 
    -X 'mask_api_gin/src/framework/config.Version=${VERSION}' \ 
    -X 'mask_api_gin/src/framework/config.BuildTime=$(date)' \ 
    -X 'mask_api_gin/src/framework/config.GoVer=$(go version)'"

## 暴露端口要与程序端口一致
EXPOSE 6275

## 程序启动命令
CMD ["./app", "--config", "./config.yaml"]

# 构建
# docker build --build-arg NAME=mask_api --build-arg VERSION=0.0.1 -t mask_api:0.0.1 .
# docker run -it mask_api:0.0.1 sh
# ./app -V