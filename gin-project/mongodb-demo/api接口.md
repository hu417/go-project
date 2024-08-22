---
title: 个人项目 v1.0.0
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.17"

---

# 个人项目

> v1.0.0

Base URLs:

* <a href="http://prod-cn.your-api-server.com">正式环境: http://prod-cn.your-api-server.com</a>

# go-gin/gin-demo/mongo-demo

## POST 创建单个文档

POST /api/v1/user

> Body 请求参数

```json
{
  "name": "hu6",
  "age": 17,
  "address": "us"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## GET 查询单个文档

GET /api/v1/user

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|address|query|string| 否 |none|

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## PUT 更新单个文档

PUT /api/v1/user

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## DELETE 删除单个文档

DELETE /api/v1/user

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## POST 创建多个文档

POST /api/v1/users

> Body 请求参数

```json
[
  {
    "name": "hu6",
    "age": 25,
    "address": "us"
  },
  {
    "name": "hu7",
    "age": 26,
    "address": "us"
  },
  {
    "name": "hu8",
    "age": 37,
    "address": "us"
  },
  {
    "name": "hu5",
    "age": 29,
    "address": "us"
  }
]
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## GET 查询多个文档

GET /api/v1/users

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## PUT 更新多个文档

PUT /api/v1/users

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## DELETE 删除多个文档

DELETE /api/v1/users

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## PUT 先查后更新

PUT /api/v1/users/update

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## GET 查询文档数量

GET /api/v1/users/count

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

## GET 聚合查询

GET /api/v1/users/aggregate

> 返回示例

> 201 Response

```json
{
  "message": "string",
  "result": {
    "id": "string",
    "data": {
      "name": "string",
      "age": 0,
      "address": "string"
    }
  },
  "status": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|201|[Created](https://tools.ietf.org/html/rfc7231#section-6.3.2)|成功|Inline|

### 返回数据结构

状态码 **201**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» result|object|true|none||none|
|»» id|string|true|none||none|
|»» data|object|true|none||none|
|»»» name|string|true|none||none|
|»»» age|integer|true|none||none|
|»»» address|string|true|none||none|
|» status|string|true|none||none|

# 数据模型

