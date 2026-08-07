# 子路径部署

应用可通过运行时环境变量 `APP_BASE_PATH` 部署到自定义 URL 子路径。镜像中不包含固定路径，留空或设置为 `/` 时仍使用根路径。

`docker-compose.yml` 使用公开镜像 `ghcr.io/qrzbing/new-api:subpath`，并将 `.env` 中的配置传入容器：

```dotenv
APP_BASE_PATH=/new-api
```

```shell
docker compose up -d
```

对应的 Compose 配置为：

```yaml
services:
  new-api:
    image: ghcr.io/qrzbing/new-api:subpath
    environment:
      - APP_BASE_PATH=${APP_BASE_PATH:-}
```

配置后，可直接通过 `http://localhost:3000/new-api/` 访问。修改环境变量后需要重启容器。

`APP_BASE_PATH` 必须以 `/` 开头，支持 `/new-api`、`/tools/new_api2` 等单级或多级路径。每段仅可包含 ASCII 字母、数字、短横线和下划线；尾部 `/` 会被自动移除。

系统设置中的 `ServerAddress` 应填写包含子路径的完整公开 URL，例如 `https://example.com/new-api`。

## Caddy

应用会自行识别并移除 URL 前缀，Caddy 必须原样转发请求，不要使用会移除前缀的 `handle_path`。

独占域名时，可代理该域名的全部请求：

```caddyfile
example.com {
    reverse_proxy new-api:3000
}
```

与其他应用共用域名时，使用 `handle` 保留请求路径：

```caddyfile
example.com {
    @new_api path /new-api /new-api/*
    handle @new_api {
        reverse_proxy new-api:3000
    }

    handle {
        # 其他应用
    }
}
```
