import pytest
from unittest.mock import Mock

@pytest.fixture(scope="session")
def mock_db_session():
    return Mock()

@pytest.fixture(scope="session")
def empty_user_list():
    return []

@pytest.fixture
def db_session(empty_user_list, mock_db_session):
    mock_db_session.query.return_value.all.return_value = empty_user_list
    return mock_db_session
