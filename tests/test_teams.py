def test_add_team(api_client):
    body = {
        "team_name": "Test Team",
        "members": [
            {"user_id": "u1", "username": "Alice", "is_active": True},
            {"user_id": "u2", "username": "Bob", "is_active": True},
            {"user_id": "u3", "username": "Charlie", "is_active": True},
            {"user_id": "u4", "username": "David", "is_active": True},
        ],
    }
    response = api_client("POST", "/team/add", json=body)
    assert response.status_code == 201, response.json()


def test_add_existed_team(api_client):
    body = {
        "team_name": "Test Team",
        "members": [
            {"user_id": "u1", "username": "Alice", "is_active": True},
            {"user_id": "u2", "username": "Bob", "is_active": True},
        ],
    }
    response = api_client("POST", "/team/add", json=body)
    assert response.status_code == 400, response.json()
    assert response.json()["error"]["code"] == "TEAM_EXISTS", response.json()


def test_get_team(api_client):
    response = api_client("GET", "/team/get", params={"team_name": "Test Team"})
    assert response.status_code == 200, response.json()


def test_get_unknown_team(api_client):
    response = api_client("GET", "/team/get", params={"team_name": "Unknown Team"})
    assert response.status_code == 404, response.json()
    assert response.json()["error"]["code"] == "NOT_FOUND", response.json()
