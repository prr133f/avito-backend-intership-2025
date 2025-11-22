def test_set_is_active(api_client):
    body = {"user_id": "u2", "is_active": False}
    response = api_client("POST", "/user/setIsActive", json=body)
    assert response.status_code == 200, response.json()


def test_set_is_active_with_invalid_user_id(api_client):
    body = {"user_id": "uu13", "is_active": False}
    response = api_client("POST", "/user/setIsActive", json=body)
    assert response.status_code == 404, response.json()
    assert response.json()["error"]["code"] == "NOT_FOUND", response.json()


def test_get_reviews(api_client):
    response = api_client("GET", "/user/getReview", params={"user_id": "u2"})
    assert response.status_code == 200, response.json()
