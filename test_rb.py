def test_read_users_returns_empty_list(db_session):
    # Act
    response = client.get("/users")

    # Assert
    assert response.status_code == 200
    assert response.json() == []
