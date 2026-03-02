import pytest
import httpx
import asyncio

# 配置后端服务地址（根据你的 Go Server 启动端口调整）
BASE_URL = "http://localhost:8080"

@pytest.fixture(scope="session")
def event_loop():
    """创建一个 session 级别的事件循环，用于异步测试"""
    loop = asyncio.get_event_loop_policy().new_event_loop()
    yield loop
    loop.close()

@pytest.fixture(scope="session")
async def client():
    """创建一个全局共享的异步 HTTP 客户端"""
    async with httpx.AsyncClient(base_url=BASE_URL, timeout=10.0) as c:
        yield c

@pytest.fixture
async def auth_client(client):
    """
    一个带自动登录逻辑的 fixture。
    测试需要权限的接口（如退出）时直接引用此 fixture。
    """
    # 模拟一个预设的测试账号
    test_user = {"username": "test_env_user", "password": "secure_password_123"}
    
    # 自动注册并登录获取 Token
    await client.post("/api/v1/register", json=test_user)
    login_res = await client.post("/api/v1/login", json=test_user)
    token = login_res.json().get("data", {}).get("token")
    
    # 在 Header 中注入 Token
    client.headers.update({"Authorization": f"Bearer {token}"})
    yield client
    # 清理：移除 Header
    client.headers.pop("Authorization", None)
