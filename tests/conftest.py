import os

import pytest
import requests

BASE_URL = "http://localhost:8000"


@pytest.fixture(scope="module")
def service_available():
    """Проверка доступности сервиса перед запуском тестов."""
    try:
        response = requests.get(
            f"{BASE_URL}/health",
            timeout=5,
        )
        assert response.status_code == 200, "Сервис недоступен"
    except requests.ConnectionError:
        pytest.fail("Не удалось подключиться к сервису")
    yield


@pytest.fixture
def api_client():
    """Фикстура для выполнения HTTP-запросов."""

    def _request(method, endpoint, **kwargs):
        url = f"{BASE_URL}{endpoint}"
        return requests.request(method, url, **kwargs)

    return _request
