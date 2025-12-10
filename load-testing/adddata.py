import requests

users = []
for i in range(200):
    users.append({"user_id": f"u{i}", "username": f"user{i}", "is_active": True})

for i in range(20):
    response = requests.post(
        "http://localhost:8000/team/add",
        json={
            "team_name": f"test_team_{i}",
            "members": users[i * 10 : (i + 1) * 10],
        },
    )
    print(response.status_code)
