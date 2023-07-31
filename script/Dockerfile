FROM golang:alpine

## 禁用cgo  
ENV CGO_ENABLED 0

## 设置构建目标
ENV GOOS linux

## 依赖下载代理
ENV GOPROXY https://goproxy.cn,direct

## 设置本地区时
ENV TZ="Asia/Shanghai"

## 工作目录存放程序源码
WORKDIR /home/mask_api_gin

## 复制实际需要的文件到工作目录
COPY ./src ./src
COPY ./go.sum ./
COPY ./go.mod ./
COPY ./main.go ./

## 安装程序依赖，需要编译
RUN go mod download

## 进行源码编译，生产文件 app
RUN go build -ldflags="-s -w" -o app

## 暴露端口要与程序端口一致
EXPOSE 6275

## 程序启动命令
CMD ["./app", "--env", "prod"]
