# https 测试证书

## 1. 生成证书

```bash
openssl req -x509 -newkey rsa:4096 -keyout www.mask.cn.key -out www.mask.cn_chain.crt -days 3650 -nodes -subj "/C=CN/ST=Beijing/L=Beijing/O=Mask/OU=Mask/CN=www.mask.cn"
```

## 2. 配置文件

```yaml
# 应用服务配置
server:
  port: 6275
  # 服务端私钥
  key: "./script/https/www.mask.cn.key"
  # 服务端证书
  cert: "./script/https/www.mask.cn_chain.crt"

```
