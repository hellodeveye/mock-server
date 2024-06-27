## 项目名称   
### 简介
**Mock Server** 是一个用于mock 接口的一个本地客户端。它通过本地编写yaml文件使用CLI快速启动来帮助用户实现接口mock，此项目适用于前端开发者和后端开发者。
### 特性
- 能够注册到本地开发环境的Nacos上, 使用泳道标签进行访问；
- 能够通过配置文件同时mock多个服务, 形成一个mock suite；
- 本地调试依赖外部服务, 开发环境没有可用数据, 使用MockServer能够快速构建mock接口；
- 前后端接口约定好后, 能够快速进行mock；
### 快速开始
#### 环境要求
- Go 1.20.5+
#### 安装指南
  为了开始使用项目名称，请按照以下步骤安装：
```
# 指定tap
brew tap hellodeveye/Homebrew-tap
# 使用Homebrew安装
brew install mock-server
```
#### 运行一个示例
```
# 启动
./mock-server -config example.yml
```
#### 配置文件示例
```
nacos:
enabled: true  # 控制是否启用 Nacos 注册
server_addr: "192.168.2.128"
server_port: 31146
namespace_id: "public"
server:
- name: service-01
  port: 8081
  endpoints:
    - url: /api/v2/hello
      method: GET
      response: Hello world!
      status: 200
      delay: 1000
      headers:
      Content-Type: html/text
```
### ToDo

- [x] Endpoint 支持启用enable开关
- [x] Mock 服务随机端口
- [x] 路由转发
- [x] 开启路由转发config
- [x] 支持转发都特定环境
- [ ] 配置文件沉淀
- [ ] gitlab
- [ ] suite
- [ ] 可视化页面
