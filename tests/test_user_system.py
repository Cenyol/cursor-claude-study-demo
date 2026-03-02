"""
用户系统黑盒测试：注册、登录、退出。
正向路径：注册 -> 登录 -> 退出；
逆向路径：重复注册 409、错误密码 401、未登录退出 401。
"""
import pytest


@pytest.mark.asyncio
async def test_register_login_logout_success(client):
    """正向：注册 -> 登录 -> 携带 Token 退出"""
    username = "new_user_001"
    password = "password123"

    reg_res = await client.post(
        "/api/v1/register",
        json={"username": username, "password": password},
    )
    assert reg_res.status_code == 200
    body = reg_res.json()
    assert body.get("msg") == "success"
    assert "data" in body and "user" in body["data"]

    login_res = await client.post(
        "/api/v1/login",
        json={"username": username, "password": password},
    )
    assert login_res.status_code == 200
    token = login_res.json().get("data", {}).get("token")
    assert token is not None

    logout_res = await client.post(
        "/api/v1/logout",
        headers={"Authorization": f"Bearer {token}"},
    )
    assert logout_res.status_code == 200
    assert logout_res.json().get("msg") == "success"


@pytest.mark.asyncio
async def test_register_duplicate_username_409(client):
    """逆向：重复注册同一用户名，预期 409"""
    username = "dup_user_002"
    password = "password123"
    await client.post("/api/v1/register", json={"username": username, "password": password})
    dup_res = await client.post("/api/v1/register", json={"username": username, "password": password})
    assert dup_res.status_code == 409
    assert "already exists" in dup_res.json().get("msg", "").lower() or "冲突" in dup_res.json().get("msg", "")


@pytest.mark.asyncio
async def test_login_wrong_password_401(client):
    """逆向：错误密码登录，预期 401"""
    username = "wrong_pwd_user"
    password = "password123"
    await client.post("/api/v1/register", json={"username": username, "password": password})
    login_res = await client.post(
        "/api/v1/login",
        json={"username": username, "password": "wrongpassword"},
    )
    assert login_res.status_code == 401


@pytest.mark.asyncio
async def test_logout_without_token_401(client):
    """逆向：未携带 Token 调用退出，预期 401"""
    logout_res = await client.post("/api/v1/logout")
    assert logout_res.status_code == 401


@pytest.mark.asyncio
async def test_auth_client_fixture(auth_client):
    """使用 auth_client fixture：已带 Token，可直接调用退出"""
    logout_res = await auth_client.post("/api/v1/logout")
    assert logout_res.status_code == 200
