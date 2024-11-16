# 基于 Gin 的管理系统后端接口服务

[![star](https://gitee.com/TsMask/mask_api_gin/badge/star.svg?theme=dark)](https://gitee.com/TsMask/mask_api_gin/stargazers)
![Build Gin](https://img.shields.io/badge/Build-Gin-green.svg)
![Release V0.4.0](https://img.shields.io/badge/Release-V0.4.0-orange.svg)
![License MIT](https://img.shields.io/badge/License-MIT-blue.svg)

## 介绍

`Node-Midwayjs` 项目代码迁移，相同配置，性能翻倍。

基于 `Golang` 的 `Gin` Web框架进行开发，如需进一步了解框架，参见 [Gin Web Framework](https://gin-gonic.com/zh-cn/docs/quickstart/) 官方文档。

## 项目文档

- 项目代码进行服务器部署的网站 => [在线预览](http://124.223.91.248:8102/#/)
- `Apifox` 提供mock服务和程序接口文档 => [接口文档](https://mask-api-midwayjs.apifox.cn/)
- `Mask管理系统` 文档专栏，相关使用和开发升级文档。 => [专栏](https://juejin.cn/column/7188761626017792056/)

### 前端

- `Vue3-Element` 仓库地址 => [mask_vue3_element](https://gitee.com/TsMask/mask_vue3_element)
- `Vue3-Antd` 仓库地址 => [mask_vue3_antd](https://gitee.com/TsMask/mask_vue3_antd)

### 后端

- `Node-Midwayjs` 仓库地址 => [mask_api_midwayjs](https://gitee.com/TsMask/mask_api_midwayjs)
- `Go-Gin` 仓库地址 => [mask_api_gin](https://gitee.com/TsMask/mask_api_gin)

## 内置功能

1. 用户管理：用户是系统操作者，该功能主要完成系统用户配置。
2. 部门管理：配置系统组织机构，树结构展现支持数据权限。
3. 岗位管理：配置系统用户所属担任职务。
4. 菜单管理：配置系统菜单，操作权限，按钮权限标识等。
5. 角色管理：角色菜单权限分配、设置角色按机构进行数据范围权限划分。
6. 字典管理：对系统中经常使用的一些较为固定的数据进行维护。
7. 参数管理：对系统动态配置常用参数。
8. 通知公告：系统通知公告信息发布维护。
9. 系统日志：系统正常操作日志记录和查询；系统异常信息日志记录和查询。
10. 登录日志：系统登录日志记录查询包含登录异常。
11. 在线用户：当前系统中活跃用户状态监控。
12. 调度任务：在线（添加、修改、删除）任务调度包含执行结果日志。
13. 服务信息：监视当前系统CPU、内存、磁盘、堆栈等相关信息。
14. 缓存信息：对系统的缓存信息查询，命令统计。

## 项目结构

```text
mask_api_gin
├── script                      目录-程序可用脚本
├── src                         目录-源代码
├   ├── assets                  目录-程序内部静态资源文件
├   ├── config                  目录-程序相关运行参数配置
├   ├── framework               目录-程序核心通用代码
├   ├── modules                 目录-业务模块
├   ├   ├── system              目录-系统内部接口模块
├   ├   ├   ├── controller      目录-接口路由控制层
├   ├   ├   ├── model           目录-数据对象模型层
├   ├   ├   ├── repository      目录-CURD数据存储层
├   ├   ├   ├── service         目录-业务逻辑服务层
├   ├   ├   └── system.go       文件-模块路由注册
├   ├   └── ...
├   ├── app.go                  文件-程序框架启动核心
├   └── configuration.go        文件-程序框架配置加载
├── .gitignore
├── go.mod                      文件-程序模块依赖
├── go.sum
├── LICENSE
├── main.go                     文件-程序主入口
└── README.md                   文件-程序说明
```

## 快速启动

你需要先安装准备 **开发环境** 后使用 **程序命令**

### 开发环境

| 技术 | 说明 | 版本 |
| ---- | ---- | ---- |
| Go | go项目的运行环境 | 1.21+ |
| Redis | 缓存存储程序 | 6+ |
| Mysql | 数据存储程序 | 8+ |

程序可用脚本 `script` 目录内含初始化数据库SQL脚本文件

### 程序命令

#### 本地开发-Window/Liunx

```bash
# 修改配置参数 /src/connfig/config.local.yaml
go mod download              # 下载项目所需依赖
go run main.go               # 开发模式启动项目
open http://localhost:6275   # 启动成功后得到服务访问地址

go run main.go --env prod    # 生产模式启动项目
```

#### 生产部署-Liunx

```bash
# 修改配置参数 /src/connfig/config.prod.yaml
go mod download             # 下载项目所需依赖
go build -o api -ldflags="-s -w"   # 构建项目二进制文件
./api --env prod            # 二进制文件启动时指定生产模式
```

> 有任何问题或者建议，可以在 [_Issues_](https://gitee.com/TsMask/mask_api_gin/issues) 或通过QQ群：[_57242844_](https://jq.qq.com/?_wv=1027&k=z6Y4YQcB) 提出想法。  
> 如果觉得项目对您有帮助，可以来个Star ⭐
