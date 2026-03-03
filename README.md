### 来自人类的提示
这是一个简单示例的用户注册登录系统，除了基本的md约束描述和需求描述说明文档，整个项目的代码都是cursor agent自动生成的！并且前述的说明文档是由Gemini生成的，依然不是人干的！！！

#### 起初我只是对Gemini说
我想要做全新的一个用户系统，基于核心接口是用户注册、登录和退出，要求高并发。如何使用cursor agent，通过编写完整的rules来约束使用的技术栈、编码规范（错误规范、编码风格、db规范等）。并依据我提供的需求描述文档，自动帮我搭建项目框架，实现接口功能，并实现docs/api.md接口文档。还需要帮我实现pytest接口黑盒测试文件。你帮我提供一份完整的rules mdc文件吧，剩下的我拿去让cursor agent帮我实现代码。

#### 然后就有了这个项目的全部内容

# User System（用户系统）

生产级用户中心模块，支持高并发下的**注册、登录、退出**，采用 DDD 四层架构，密码 bcrypt 哈希，Session 存 Redis，MySQL 连接池与索引满足高并发与安全规范。

---

## 技术栈

| 类别     | 技术 |
|----------|------|
| 语言     | Go 1.22+ |
| Web 框架 | [Gin](https://gin-gonic.com) |
| ORM      | [GORM](https://gorm.io) |
| 数据库   | MySQL 8.0+ |
| 缓存     | Redis（Session/Token） |
| 密码哈希 | bcrypt（golang.org/x/crypto） |
| 测试     | Python 3.10+、[Pytest](https://docs.pytest.org)、[httpx](https://www.python-httpx.org/) |

---

## 项目结构（DDD 四层）

```
user-system/
├── cmd/
│   └── main.go                 # 程序入口，依赖注入与启动
├── internal/
│   ├── interfaces/http/        # 接口层：HTTP 请求/响应、参数校验
│   │   ├── response.go         # 统一 JSON 响应
│   │   ├── user_handler.go     # 注册/登录/退出 Handler
│   │   └── router.go           # 路由注册
│   ├── application/            # 应用层：用例编排，无业务逻辑
│   │   ├── dto/                # 请求/响应 DTO
│   │   └── usecase/            # 注册、登录、退出用例
│   ├── domain/                 # 领域层：实体与 Repository 接口（不依赖外层）
│   │   ├── entity/             # User 聚合根、工厂方法、密码校验
│   │   └── repository/         # UserRepository、SessionRepository 接口
│   └── infrastructure/        # 基础设施层：持久化与缓存实现
│       ├── persistence/        # GORM/MySQL、连接池、UserRepository 实现
│       └── cache/              # Redis、SessionRepository 实现（key: user:session:{token}）
├── pkg/
│   └── token/                  # 随机 Token 生成（Session 用）
├── docs/
│   ├── api.md                  # 接口文档（含 cURL 示例）
│   ├── db.sql                  # 建表语句（含索引与 comment）
│   └── agent-logs.md           # Agent 对话日志（可选）
├── tests/                      # Pytest 黑盒测试
│   ├── conftest.py             # client、api_client、auth_client
│   ├── test_user_system.py     # 正向/逆向用例
│   └── requirements.txt        # pytest, pytest-asyncio, httpx, pydantic
├── go.mod
├── .gitignore
└── README.md
```

**依赖方向**：Interfaces → Application → Domain ← Infrastructure（Domain 定义接口，Infrastructure 实现，依赖倒置）。

---

## 功能概览

| 能力       | 说明 |
|------------|------|
| 用户注册   | 用户名唯一、密码至少 8 位、邮箱可选；bcrypt 哈希存储；返回用户信息（不含密码） |
| 用户登录   | 校验用户名与密码；Token 写入 Redis（2 小时过期）；DB 更新 `is_login`、`login_at` |
| 用户退出   | 从 Redis 删除 Token；DB 更新登录状态为未登录 |
| 高并发     | MySQL 连接池 MaxOpenConns=100、MaxIdleConns=20；username 唯一索引 |
| 安全       | 明文密码不落库、不打印；统一 JSON `{code, data, msg}`，业务码与 HTTP 状态码分离 |

---

## 快速开始

### 环境要求

- Go 1.22+
- MySQL 8.0+（建议 9.0）
- Redis
- Python 3.10+（仅运行/编写测试时需要）

#### 持久化组件
这一小段是人工(我)补充的！
```
docker run --name mysql -e MYSQL_ROOT_PASSWORD=root -p 3306:3306 -d mysql:9
docker run --name redis -p 6379:6379 -d redis:7.4-alpine
```

### 1. 数据库与缓存

```bash
# 创建库并执行建表（按需取消注释 db 创建语句）
mysql -u root -p < docs/db.sql
```

Redis 无需建表，确保服务可用即可。

### 2. 配置

通过环境变量覆盖默认值（可选）：

| 变量            | 默认值 | 说明 |
|-----------------|--------|------|
| `MYSQL_DSN`     | `root:root@tcp(127.0.0.1:3306)/user_system?charset=utf8mb4&parseTime=True` | MySQL 连接串 |
| `REDIS_ADDR`    | `127.0.0.1:6379` | Redis 地址 |
| `REDIS_PASSWORD`| 空    | Redis 密码 |

生产环境请使用 `.env` 或配置中心，勿将敏感信息提交仓库（参见 `.cursor/rules/050-git-ignore.mdc`）。

### 3. 运行服务

```bash
cd /root/go/src/user-system
go mod tidy
go run ./cmd/main.go
```

默认监听 **8080**。如需改端口，可修改 `cmd/main.go` 中的 `r.Run(":8080")` 或通过环境变量扩展。

### 4. 验证接口

```bash
# 注册
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123","email":"alice@example.com"}'

# 登录（将返回 token）
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}'

# 退出（替换 YOUR_TOKEN）
curl -X POST http://localhost:8080/api/v1/logout \
  -H "Authorization: Bearer YOUR_TOKEN"
```

更多请求/响应说明与错误码见 **docs/api.md**。

---

## 测试

### Go 单元测试

```bash
go test ./internal/...
```

### Pytest 黑盒测试（需先启动服务）

测试默认请求 `http://localhost:8080`，请先启动 Go 服务。

```bash
# 推荐：使用项目根目录或 tests 目录下的虚拟环境
python3 -m venv .venv
source .venv/bin/activate   # Windows: .venv\Scripts\activate
pip install -r tests/requirements.txt

# 运行全部用例
python3 -m pytest tests/ -v

# 仅运行用户系统用例
python3 -m pytest tests/test_user_system.py -v
```

**覆盖场景**：

- 正向：注册 → 登录 → 携带 Token 退出
- 逆向：重复注册（409）、错误密码登录（401）、未携带 Token 退出（401）
- Fixture：`client`、`api_client`、`auth_client`（已带 Token，可直接测退出）

---

## 接口一览

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/register` | 注册 |
| POST | `/api/v1/login`     | 登录，返回 Token |
| POST | `/api/v1/logout`    | 退出，Header: `Authorization: Bearer <token>` |

统一响应格式：`{"code": 0, "data": {...}, "msg": "success"}`。错误码见 **docs/api.md**。

---

## 数据库与连接池

- **表**：`users`（id, username, password_hash, email, is_login, login_at, created_at, updated_at），详见 **docs/db.sql**。
- **索引**：`username` 唯一索引 `uk_username`。
- **连接池**：在 `internal/infrastructure/persistence/mysql.go` 中配置 `MaxOpenConns=100`、`MaxIdleConns=20`。

Session 存 Redis，Key：`user:session:{token}`，过期时间 2 小时。

---

## 规范与约定

- **DDD**：严格四层（Interfaces / Application / Domain / Infrastructure），Domain 不依赖框架与 DB，详见 `.cursor/rules/110-ddd-architecture.mdc`。
- **错误与安全**：不忽略错误；不提交明文密码与密钥；敏感配置使用环境变量或 `.env`，见 `.cursor/rules/050-git-ignore.mdc`。
- **接口与测试**：新增接口需同步更新 **docs/api.md**，并为每个接口编写至少一个 Pytest 用例（见 `.cursor/rules/000-project-main.mdc`、`300-pytest-spec.mdc`）。

---

## 文档索引

| 文档 | 说明 |
|------|------|
| [docs/api.md](docs/api.md) | 接口说明、参数、响应示例与 cURL |
| [docs/db.sql](docs/db.sql) | MySQL 建表语句（含索引与 comment） |
| [docs/agent-logs.md](docs/agent-logs.md) | Agent 对话与任务记录（可选） |

---

## License

按项目所在仓库约定执行。
