# bk-plugin-demo-go

蓝鲸 Go 语言插件 Demo，用于验证 Go 版本插件在多租户环境下的运行。

## 插件说明

这是一个简单的 echo 插件，目前包含两个版本：

- `1.0.0`：基础 echo 版本。
- `1.0.1`：多租户验证版本，从标准运维上下文读取 `tenant_id`，并输出带 `trace_id` 的结构化日志。

两个版本的基础参数保持一致：

- **输入**：`hello`（string）
- **上下文输入**：`executor`（string - 任务执行人）
- **输出**：`world`（string）- 返回输入的 `hello` 值

`1.0.1` 额外声明上下文输入 `tenant_id`。执行成功后，可在标准运维的第三方插件日志中查看：

```text
tenant context received
```

对应日志包含 `trace_id` 和 `tenant_id` 字段；如果标准运维没有传入 `tenant_id`，插件会打印
`tenant context missing` 并执行失败。

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

上述环境变量表示插件应用自身所属租户。`1.0.1` 验证的是标准运维任务所属租户，取值来自
invoke 请求的 `context.tenant_id`，两者不能混用。
