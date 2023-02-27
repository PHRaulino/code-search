from fastapi.testclient import TestClient
from unittest.mock import patch

from main import app


client = TestClient(app)


@patch("app.routes.users.get_db")
def test_read_users(mock_get_db):
    # Define o comportamento do mock
    mock_db = mock_get_db.return_value
    mock_db.query.return_value.all.return_value = [
        models.User(
            num_seqe_usua=1,
            dat_hor_atui=datetime.now(),
            nom_orig_usua="Usuário 1",
            idef_user="user1",
            groups=[],
        ),
        models.User(
            num_seqe_usua=2,
            dat_hor_atui=datetime.now(),
            nom_orig_usua="Usuário 2",
            idef_user="user2",
            groups=[],
        ),
    ]

    # Faz a chamada à rota com o mock injetado
    response = client.get("/users")

    # Verifica se a resposta tem o status code esperado
    assert response.status_code == 200

    # Verifica se a lista de usuários retornada pela rota é a mesma que a lista definida no mock
    assert response.json() == [
        {
            "num_seqe_usua": 1,
            "dat_hor_atui": mock_db.query.return_value.all.return_value[0].dat_hor_atui.isoformat(),
            "nom_orig_usua": "Usuário 1",
            "idef_user": "user1",
            "groups": [],
        },
        {
            "num_seqe_usua": 2,
            "dat_hor_atui": mock_db.query.return_value.all.return_value[1].dat_hor_atui.isoformat(),
            "nom_orig_usua": "Usuário 2",
            "idef_user": "user2",
            "groups": [],
        },
    ]
