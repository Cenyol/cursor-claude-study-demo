# 用户系统 API 文档

## 1. 用户注册 (Register)

**描述**：注册新用户，用户名唯一，密码至少 8 位，邮箱可选。

**Method & URL**  
`POST /api/v1/register`

**请求参数**

| 字段     | 类型   | 必填 | 说明            |
|----------|--------|------|-----------------|
| username | string | 是   | 用户名，唯一    |
| password | string | 是   | 密码，至少 8 位 |
| email    | string | 否   | 邮箱            |

**请求体示例**

```json
{
  "username": "alice",
  "password": "password123",
  "email": "alice@example.com"
}
```

**成功响应 (200)**

```json
{
  "code": 0,
  "data": {
    "user": {
      "id": 1,
      "username": "alice",
      "email": "alice@example.com",
      "created_at": "2025-03-02T10:00:00Z"
    }
  },
  "msg": "success"
}
```

**错误码**：400 参数无效/密码不足 8 位；409 用户名已存在；500 服务端错误。

**cURL 示例**

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123","email":"alice@example.com"}'
```

---

## 2. 用户登录 (Login)

**描述**：用户名密码登录，成功后返回 Token，Session 存 Redis（2 小时）。

**Method & URL**  
`POST /api/v1/login`

**请求参数**：username (必填)、password (必填)。

**成功响应 (200)**

```json
{
  "code": 0,
  "data": { "token": "a1b2c3d4..." },
  "msg": "success"
}
```

**错误码**：400 参数无效；401 用户名或密码错误；500 服务端错误。

**cURL 示例**

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}'
```

---

## 3. 用户退出 (Logout)

**描述**：使当前 Token 失效，需携带 `Authorization: Bearer <token>`。

**Method & URL**  
`POST /api/v1/logout`

**成功响应 (200)**

```json
{
  "code": 0,
  "data": null,
  "msg": "success"
}
```

**错误码**：401 未提供或无效 Token；500 服务端错误。

**cURL 示例**（将 `YOUR_TOKEN` 替换为登录返回的 token）

```bash
curl -X POST http://localhost:8080/api/v1/logout \
  -H "Authorization: Bearer YOUR_TOKEN"
```
