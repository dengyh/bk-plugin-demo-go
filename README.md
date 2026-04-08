# bk-plugin-demo-go

蓝鲸 Go 语言插件 Demo，用于验证 Go 版本插件在多租户环境下的运行。

## 插件说明

这是一个简单的 echo 插件：

- **输入**: `hello` (string)
- **上下文输入**: `executor` (string - 任务执行人)
- **输出**: `world` (string) - 返回输入的 `hello` 值

## 依赖版本

- Go >= 1.23
- `beego-runtime` >= v0.7.0（支持多租户 `X-Bk-Tenant-Id` header）
- `bk-plugin-framework-go` v0.5.0

## 本地开发

```bash
# 确保已设置 Redis 环境变量
export REDIS_HOST="127.0.0.1"
export REDIS_PORT="6379"
export REDIS_PASSWORD=""

# 编译
go build -o bin/bk-plugin-demo-go .

# 启动 web 服务
./bin/bk-plugin-demo-go server

# 启动 worker
./bin/bk-plugin-demo-go worker
```

## 多租户配置

部署到 PaaS 平台时，平台会自动注入 `BKPAAS_APP_TENANT_ID` 环境变量。
手动部署时，需设置 `BK_APP_TENANT_ID` 环境变量。
