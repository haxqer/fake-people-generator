[English](README.md) | [简体中文](README.zh-CN.md)

# Fake People Generator

一个使用 Go 编写的合成个人资料生成 API。当前版本只支持美国资料生成，目标是输出“高拟真但明确为合成”的个人信息，适合测试、演示、联调、压测和假数据填充场景。

## 功能特性

- 支持 `REST JSON` API
- 内置一个工具型前端页面，打开首页即可自动生成示例资料
- 支持单条和批量生成，单次最多 `100` 条
- 支持按 `gender`、`age_min`、`age_max`、`state` 过滤
- 支持通过 `include_fields` 精确请求增强字段和敏感字段
- 生成的姓名、年龄、生日、州/城市/邮编、区号、电话彼此一致
- 返回更完整的合成资料字段：
  - 姓名：`first_name`、`middle_name`、`last_name`、`full_name`
  - 基本信息：`gender`、`age`、`birth_date`、`marital_status`、`occupation`
  - 就业信息：`employment.status`、`employment.company_name`、`employment.job_title`、`employment.company_domain`、`employment.work_email`
  - 教育信息：`education.level`、`education.degree_type`、`education.school_name`、`education.major`、`education.graduation_year`
  - 敏感信息：`sensitive.ssn`、`sensitive.driver_license.*`、`sensitive.passport.*`、`sensitive.bank_card.*`
  - 联系方式：`phone`、`username`、`email`
  - 地址信息：`street_line1`、`street_line2`、`city`、`state`、`state_code`、`postal_code`
  - 地域衍生：`time_zone`
- 提供 OpenAPI 文档和 Swagger UI

## 说明

- 当前仅支持 `country = "US"`
- 输出数据是合成资料，不保证对应真实存在的人、住址或电话号码
- 默认返回基础字段；教育增强和敏感字段只有在 `include_fields` 显式声明时才返回

## 运行方式

要求：

- Go `1.25+`

启动服务：

```bash
go run ./cmd/server
```

默认监听：

```text
http://127.0.0.1:8080
```

首页工具页面：

```text
http://127.0.0.1:8080/
```

也可以通过环境变量修改端口：

```bash
ADDR=:9090 go run ./cmd/server
```

或者使用更适合容器平台的 `PORT`：

```bash
PORT=9090 go run ./cmd/server
```

## Docker 部署

构建镜像：

```bash
docker build -t fake-people-generator:local .
```

启动容器：

```bash
docker run --rm -p 8080:8080 fake-people-generator:local
```

如果需要自定义容器内监听端口：

```bash
docker run --rm -e PORT=9090 -p 9090:9090 fake-people-generator:local
```

当前 `Dockerfile` 使用多阶段构建：

- 构建阶段：`golang:1.25-bookworm`
- 运行阶段：`distroless static-debian12:nonroot`
- 默认容器端口：`8080`
- 支持 `linux/amd64` 和 `linux/arm64` 多架构构建

## GitHub Actions 自动构建并推送 Docker Hub

仓库已添加工作流：

```text
.github/workflows/docker.yml
```

触发方式：

- push 到 `main`
- push 版本标签，例如 `v0.1.0`
- 手动触发 `workflow_dispatch`

在 GitHub 仓库里需要先配置以下 secrets：

- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

可选配置一个仓库变量：

- `DOCKERHUB_REPOSITORY`

说明：

- 如果没有设置 `DOCKERHUB_REPOSITORY`，默认使用当前 GitHub 仓库名作为 Docker Hub 镜像名
- 最终镜像名默认是 `docker.io/<DOCKERHUB_USERNAME>/fake-people-generator`
- `main` 分支会推送 `latest`、`main` 和 `sha-<commit>` 标签
- `v*` 标签会推送对应版本标签，例如 `v0.1.0`
- 工作流会先执行 `go test ./...`，通过后再构建并推送多架构镜像

发布示例：

```bash
git tag v0.1.0
git push origin v0.1.0
```

拉取示例：

```bash
docker pull docker.io/<DOCKERHUB_USERNAME>/fake-people-generator:latest
```

## API

### 健康检查

```bash
curl http://127.0.0.1:8080/healthz
```

### 生成资料

请求：

```bash
curl -X POST http://127.0.0.1:8080/v1/profiles/generate \
  -H 'Content-Type: application/json' \
  -d '{
    "country": "US",
    "count": 2,
    "filters": {
      "gender": "female",
      "age_min": 28,
      "age_max": 35,
      "state": "CA"
    },
    "include_fields": [
      "education.level",
      "employment.industry",
      "sensitive.ssn"
    ]
  }'
```

响应示例：

```json
{
  "meta": {
    "country": "US",
    "count": 1,
    "generated_at": "2026-03-20T05:35:07Z"
  },
  "profiles": [
    {
      "id": "b88418b7-3fc4-4405-bfe3-9dbe558a75ea",
      "country": "US",
      "gender": "female",
      "first_name": "Emily",
      "middle_name": "Grace",
      "last_name": "Smith",
      "full_name": "Emily Grace Smith",
      "age": 29,
      "birth_date": "1996-04-20",
      "marital_status": "married",
      "occupation": "Software Engineer",
      "username": "emily.smith42",
      "email": "emily.smith@gmail.com",
      "time_zone": "America/Los_Angeles",
      "employment": {
        "status": "employed",
        "company_name": "Pacific Systems Inc.",
        "company_domain": "pacificsystems.com",
        "job_title": "Senior Software Engineer",
        "department": "Engineering",
        "work_email": "emily.smith@pacificsystems.com",
        "industry": "Technology",
        "seniority": "senior",
        "start_date": "2020-09-14",
        "work_mode": "hybrid",
        "manager_name": "Nancy Wilson",
        "office_city": "San Jose",
        "office_state": "California"
      },
      "education": {
        "level": "bachelor",
        "school_name": "California University",
        "major": "Computer Science",
        "graduation_year": 2018,
        "student_status": "alumni",
        "school_email": "emily.smith@california.edu"
      },
      "sensitive": {
        "ssn": "428-15-9012"
      },
      "phone": {
        "e164": "+14086827484",
        "national": "(408) 682-7484",
        "area_code": "408"
      },
      "address": {
        "street_line1": "9535 Highland Boulevard",
        "street_line2": "Apt 217",
        "city": "San Jose",
        "state": "California",
        "state_code": "CA",
        "postal_code": "95113"
      }
    }
  ]
}
```

### 请求参数

请求体：

```json
{
  "country": "US",
  "count": 1,
  "filters": {
    "gender": "male|female|random",
    "age_min": 18,
    "age_max": 65,
    "state": "CA"
  },
  "include_fields": [
    "education.level",
    "employment.industry",
    "sensitive.ssn"
  ]
}
```

规则：

- `country` 必填，当前只能是 `US`
- `count` 选填，默认 `1`，范围 `1..100`
- `filters.gender` 选填，支持 `male`、`female`、`random`
- `filters.age_min` 和 `filters.age_max` 必须一起传
- `filters.state` 为两位美国州代码，例如 `CA`、`NY`、`TX`
- `include_fields` 选填；未声明的增强字段不会返回

常用 `include_fields`：

- 教育：`education.level`、`education.degree_type`、`education.school_name`、`education.major`、`education.graduation_year`、`education.student_status`、`education.school_city`、`education.school_state`、`education.school_email`
- 就业增强：`employment.industry`、`employment.seniority`、`employment.start_date`、`employment.work_mode`、`employment.manager_name`、`employment.office_city`、`employment.office_state`
- 敏感：`sensitive.ssn`、`sensitive.driver_license.number`、`sensitive.driver_license.state_code`、`sensitive.passport.number`、`sensitive.passport.country`、`sensitive.bank_card.brand`、`sensitive.bank_card.number`、`sensitive.bank_card.expiry_month`、`sensitive.bank_card.expiry_year`、`sensitive.bank_card.cvv`

### 错误响应

```json
{
  "error": {
    "code": "validation_error",
    "message": "count must be between 1 and 100"
  }
}
```

## 文档

- 首页工具：`http://127.0.0.1:8080/`
- OpenAPI JSON: `http://127.0.0.1:8080/openapi.json`
- Swagger UI: `http://127.0.0.1:8080/docs`

说明：

- `/docs` 使用 Swagger UI CDN 资源，打开页面时需要网络

## 测试

运行全部测试：

```bash
go test ./...
```
