def test_set_is_active(api_client):
    body = {"user_id": "u2", "is_active": False}
    response = api_client("POST", "/users/setIsActive", json=body)
    assert response.status_code == 200
    assert not response.json()["user"]["is_active"]
    assert response.json()["user"]["user_id"] == "u2"


def test_set_is_active_with_invalid_user_id(api_client):
    body = {"user_id": "invalid_user_id", "is_active": False}
    response = api_client("POST", "/users/setIsActive", json=body)
    assert response.status_code == 404
    assert response.json()["error"]["code"] == "NOT_FOUND"
    assert response.json()["error"]["message"] == "resource not found"


def test_get_reviews(api_client):
    response = api_client("GET", "/users/getReview", params={"user_id": "u2"})
    assert response.status_code == 200
    assert response.json()["user_id"] == "u2"
    assert len(response.json()["pull_requests"]) > 0
