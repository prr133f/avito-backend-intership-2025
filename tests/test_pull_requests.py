def test_create_pull_request(api_client, service_available):
    payload = {
        "pull_request_id": "pr-1001",
        "pull_request_name": "Add search",
        "author_id": "u1",
    }
    resp = api_client("post", "/pullRequest/create", json=payload)
    assert resp.status_code == 201
    assert resp.json()["pr"] == "pr-1001"


def test_create_pull_request_with_existing_id(api_client, service_available):
    payload = {
        "pull_request_id": "pr-1001",
        "pull_request_name": "Add search",
        "author_id": "u1",
    }
    resp = api_client("post", "/pullRequest/create", json=payload)
    assert resp.status_code == 409, (
        f"expected 409 for existing PR, got {resp.status_code}"
    )


def test_create_pull_request_with_invalid_author_id(api_client, service_available):
    payload = {
        "pull_request_id": "pr-1002",
        "pull_request_name": "Add search",
        "author_id": "invalid id",
    }
    resp = api_client("post", "/pullRequest/create", json=payload)
    assert resp.status_code == 404, (
        f"expected 404 for invalid author ID, got {resp.status_code}"
    )


def test_merge_pr(api_client):
    resp = api_client("post", "/pullRequest/merge", json={"pull_request_id": "pr-1001"})
    assert resp.status_code == 200, f"expected 200 for merge PR, got {resp.status_code}"
    assert resp.json()["pr"]["status"] == "MERGED"


def test_merge_nonexistent_pr(api_client, service_available):
    resp = api_client(
        "post", "/pullRequest/merge", json={"pull_request_id": "pr-does-not-exist"}
    )
    assert resp.status_code == 404, (
        f"expected 404 for nonexistent PR, got {resp.status_code}"
    )
    assert resp.json()["error"]["code"] == "NOT_FOUND"
    assert resp.json()["error"]["message"] == "resource not found"


def test_reassign_pr(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr1", "old_reviewer_id": "u2"},
    )
    assert resp.status_code == 200, (
        f"expected 200 for reassign PR, got {resp.status_code}"
    )


def test_reassign_nonexistent_user_pr(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr-1001", "old_reviewer_id": "unknown user"},
    )
    assert resp.status_code == 404, (
        f"expected 404 for nonexistent PR reassign, got {resp.status_code}"
    )
    assert resp.json()["error"]["code"] == "NOT_FOUND"
    assert resp.json()["error"]["message"] == "resource not found"


def test_reassign_nonexistent_pr(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "unknown pr", "old_reviewer_id": "u2"},
    )
    assert resp.status_code == 404, (
        f"expected 404 for nonexistent PR reassign, got {resp.status_code}"
    )
    assert resp.json()["error"]["code"] == "NOT_FOUND"
    assert resp.json()["error"]["message"] == "resource not found"


def test_reassign_after_merge(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr1", "old_reviewer_id": "u2"},
    )
    assert resp.status_code == 409, (
        f"expected 409 for reassign after merge, got {resp.status_code}"
    )
    assert resp.json()["error"]["code"] == "PG_MERGED"
    assert resp.json()["error"]["message"] == "cannot reassign on merged PR"


def test_reassign_user_not_reviewer(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr1", "old_reviewer_id": "u3"},
    )
    assert resp.status_code == 409, (
        f"expected 409 for reassign user not reviewer, got {resp.status_code}"
    )
    assert resp.json()["error"]["code"] == "NOT_ASSIGNED"
    assert resp.json()["error"]["message"] == "reviewer is not assigned to this PR"


def test_reassign_no_candidates(api_client):
    resp = api_client(
        "post",
        "/pullRequest/reassign",
        json={"pull_request_id": "pr1", "old_reviewer_id": "u1"},
    )
    assert resp.status_code == 409, (
        f"expected 409 for reassign no candidates, got {resp.status_code}"
    )
    assert resp.json()["error"]["code"] == "NO_CANDIDATE"
    assert resp.json()["error"]["message"] == "no active replacement candidate in team"
