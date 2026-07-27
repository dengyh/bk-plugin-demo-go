# bk-plugin-demo-go

蓝鲸 Go 语言插件 Demo，用于验证 Go 版本插件在多租户环境下的运行。

## 插件说明

这是一个简单的 echo 插件，目前包含四个版本：

- `1.0.0`：基础 echo 版本。
- `1.0.1`：多租户验证版本，从标准运维上下文读取 `tenant_id`，并输出带 `trace_id` 的结构化日志。
- `1.0.2`：同时验证标准运维执行上下文和表单 `data_api` 获取到的租户 ID。
- `1.0.3`：修正 `data_api` 使用的 SaaS `app_code` 为 `bkplugin-go2`。

四个版本的基础参数保持一致：

- **输入**：`hello`（string）
- **上下文输入**：`executor`（string - 任务执行人）
- **输出**：`world`（string）- 返回输入的 `hello` 值

`1.0.1` 额外声明上下文输入 `tenant_id`。执行成功后，可在标准运维的第三方插件日志中查看：

```text
tenant context received
```

对应日志包含 `trace_id` 和 `tenant_id` 字段；如果标准运维没有传入 `tenant_id`，插件会打印
`tenant context missing` 并执行失败。

`1.0.3` 的输入表单会调用：

```text
${SITE_URL}plugin_service/data_api/bkplugin-go2/bk_plugin/plugin_api/tenant_id
```

插件后端从 API 网关注入的 `X-Bkapi-Tenant-Id` 请求头读取租户 ID，返回为下拉选项
“当前租户：<tenant_id>”。选择该项并执行节点后，标准运维节点日志会包含：

```text
tenant context validation succeeded
```

对应结构化字段包含 `trace_id`、`tenant_id` 和 `data_api_tenant_id`。如果两种方式获取的租户
ID 不一致，插件执行失败并打印 `tenant context validation failed`。

SG 环境的 `SITE_URL` 为 `/bk--sops/`，对应的完整验证地址是：

```text
https://apps.sg.bk2game.com/bk--sops/plugin_service/data_api/bkplugin-go2/bk_plugin/plugin_api/tenant_id
```

## 依赖版本

- Go >= 1.23
- `beego-runtime` >= v0.7.5（支持多租户 `X-Bk-Tenant-Id` header，并修复 JSON `plugin_api_dispatch` 空指针）
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
