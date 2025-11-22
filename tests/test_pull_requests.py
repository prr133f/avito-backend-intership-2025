def test_create_pull_request(api_client, service_available):
    payload = {
        "pull_request_id": "pr-1001",
        "pull_request_name": "Add search",
        "author_id": "u1",
    }
    resp = api_client("post", "/pullRequest/create", json=payload)
    assert resp.status_code == 201, resp.json()


def test_create_pull_request_with_existing_id(api_client, service_available):
    payload = {
        "pull_request_id": "pr-1001",
        "pull_request_name": "Add search",
        "author_id": "u1",
    }
    resp = api_client("post", "/pullRequest/create", json=payload)
    assert resp.status_code == 409, resp.json()
    assert resp.json()["error"]["code"] == "PR_EXISTS", resp.json()


def test_create_pull_request_with_invalid_author_id(api_client, service_available):
    payload = {
        "pull_request_id": "pr-1002",
        "pull_request_name": "Add search",
        "author_id": "u100",
    }
    resp = api_client("post", "/pullRequest/create", json=payload)
    assert resp.status_code == 404, resp.json()


def test_reassign_pr(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr-1001", "old_reviewer_id": "u2"},
    )
    assert resp.status_code == 200, resp.json()


def test_reassign_nonexistent_user_pr(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr-1001", "old_reviewer_id": "u5"},
    )
    assert resp.status_code == 404, resp.json()
    assert resp.json()["error"]["code"] == "NOT_FOUND", resp.json()


def test_reassign_nonexistent_pr(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr-1", "old_reviewer_id": "u2"},
    )
    assert resp.status_code == 404, resp.json()
    assert resp.json()["error"]["code"] == "NOT_FOUND", resp.json()


def test_reassign_user_not_reviewer(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr-1001", "old_reviewer_id": "u4"},
    )
    assert resp.status_code == 409, resp.json()
    assert resp.json()["error"]["code"] == "NOT_ASSIGNED", resp.json()


# def test_reassign_no_candidates(api_client):
#     resp = api_client(
#         "post",
#         "/pullRequest/reassign",
#         json={"pull_request_id": "pr-1001", "old_reviewer_id": "u1"},
#     )
#     assert resp.status_code == 409, (
#         f"expected 409 for reassign no candidates, got {resp.status_code}"
#     )
#     assert resp.json()["error"]["code"] == "NO_CANDIDATE"
#     assert resp.json()["error"]["message"] == "no active replacement candidate in team"


def test_merge_pr(api_client):
    resp = api_client("post", "/pullRequest/merge", json={"pull_request_id": "pr-1001"})
    assert resp.status_code == 200, f"expected 200 for merge PR, got {resp.status_code}"
    print(resp.json())
    assert resp.json()["pr"]["status"] == "MERGED", resp.json()


def test_merge_nonexistent_pr(api_client, service_available):
    resp = api_client("post", "/pullRequest/merge", json={"pull_request_id": "pr-1"})
    assert resp.status_code == 404, resp.json()
    assert resp.json()["error"]["code"] == "NOT_FOUND", resp.json()


def test_reassign_after_merge(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr-1001", "old_reviewer_id": "u2"},
    )
    assert resp.status_code == 409, resp.json()
    assert resp.json()["error"]["code"] == "PR_MERGED", resp.json()
