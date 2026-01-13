# dnsctl

> 管理 etcd 中 DNS 主机记录的命令行工具.

## 安装

### 从源码编译

```sh
git clone https://github.com/etcdhosts/dnsctl.git
cd dnsctl
task build

# 二进制文件在 dist/ 目录
./dist/dnsctl --help
```

### 安装到 GOPATH

```sh
task install
dnsctl --help
```

## 配置

创建 `~/.dnsctl.yaml`:

```yaml
endpoints:
  - https://172.16.1.21:2379
  - https://172.16.1.22:2379
  - https://172.16.1.23:2379
key: /etcdhosts
ca: /etc/etcd/ssl/etcd-ca.pem
cert: /etc/etcd/ssl/etcd.pem
cert_key: /etc/etcd/ssl/etcd-key.pem
```

或生成示例配置:

```sh
dnsctl example > ~/.dnsctl.yaml
```

### 配置选项

| 选项 | 说明 |
|------|------|
| `endpoints` | etcd 端点列表 |
| `key` | etcd 键或前缀 (默认: `/etcdhosts`) |
| `dial_timeout` | 连接超时 (默认: `5s`) |
| `req_timeout` | 请求超时 (默认: `5s`) |
| `ca` | CA 证书文件 |
| `cert` | 客户端证书文件 |
| `cert_key` | 客户端密钥文件 |
| `username` | etcd 用户名 |
| `password` | etcd 密码 |

## 使用方法

### 列出记录

```sh
# 默认 hosts 格式
dnsctl list

# JSON 输出
dnsctl list -o json

# YAML 输出
dnsctl list -o yaml

# 读取指定版本
dnsctl list -r 12345
```

### 编辑记录

使用系统编辑器编辑 DNS 记录:

```sh
# 使用默认编辑器 ($EDITOR 或 vi)
dnsctl edit

# 使用指定编辑器
EDITOR=nano dnsctl edit
EDITOR=vim dnsctl edit
```

功能:
- 自动去重 (移除重复记录)
- 保存前验证 hosts 格式
- 保留扩展属性 (权重, TTL, 健康检查)

### 查看历史

```sh
# 列出所有版本 (单键模式)
dnsctl history

# 限制显示数量
dnsctl history -n 10

# 按主机分键模式: 列出所有域名
dnsctl history

# 按主机分键模式: 查看域名历史
dnsctl history example.com
```

输出示例:
```
REVISION        VERSION     RECORDS   MODIFIED
--------------  ----------  --------  --------------------
12350           3           5         2024-01-12 10:30:00 (latest)
12340           2           4         2024-01-12 09:15:00
12330           1           3         -
```

### 对比版本

```sh
# 对比两个版本
dnsctl diff 12340 12350

# 对比旧版本与当前版本 (0 = 当前)
dnsctl diff 12340 0
```

输出带颜色:
- 红色 (`-`): 删除的行
- 绿色 (`+`): 新增的行

### 清除主机名

```sh
# 删除主机名的所有记录
dnsctl purge example.com
```

### 其他命令

```sh
# 显示版本
dnsctl --version

# 显示帮助
dnsctl --help

# 使用自定义配置文件
dnsctl -c /path/to/config.yaml list

# 生成示例配置
dnsctl example > ~/.dnsctl.yaml
```

## 示例

### 配置负载均衡服务

```sh
# 使用编辑器配置多个后端并设置权重
dnsctl edit

# 在编辑器中添加:
# 192.168.1.1 api.example.com # +etcdhosts weight=3
# 192.168.1.2 api.example.com # +etcdhosts weight=2
# 192.168.1.3 api.example.com # +etcdhosts weight=1

# 验证
dnsctl list
```

### 迁移后端

```sh
# 使用编辑器修改记录
dnsctl edit

# 在编辑器中:
# - 添加新后端: 192.168.1.10 api.example.com # +etcdhosts weight=1
# - 删除旧后端行: 192.168.1.1 api.example.com
```

### 导出和备份

```sh
# 导出为 JSON
dnsctl list -o json > backup.json

# 导出为 YAML
dnsctl list -o yaml > backup.yaml
```

## 构建

```sh
# 为当前平台构建
task build

# 运行测试
task test

# 运行集成测试 (需要 Docker)
task test-integration

# 构建所有发布版本
task release

# 创建 GitHub 发布
task gh-release
```

## 相关项目

- [etcdhosts](https://github.com/etcdhosts/etcdhosts) - CoreDNS 插件
- [client-go](https://github.com/etcdhosts/client-go) - Go 客户端库

## 许可证

Apache 2.0
