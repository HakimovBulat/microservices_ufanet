def test_index_view_redirect(client):
    response = client.get("/", follow=True)
    assert response.redirect_chain[0][0] == "/billboard"