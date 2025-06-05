# import json
# import datetime
# from billboard.models import Category, Sale
# from django.utils import timezone
# from rest_framework.test import APIClient
# from django.utils import timezone
# from django.contrib.auth.models import User
# from api.views import SaleListView
# import datetime
# import pytest
# from django.test import RequestFactory

# @pytest.fixture
# def example_user():
#     user = User.objects.create_user(
#         username='login',
#         password='Y7h3mi*ow9'
#     )
#     return user


# @pytest.fixture
# def example_category():
#     category = Category.objects.create(title="Новая категория")
#     category.save()
#     return category


# @pytest.fixture
# def example_sale(example_category):
#     sale = Sale.objects.create(
#         title="Новая скидка", 
#         start_date=timezone.now(), 
#         end_date=timezone.now() + datetime.timedelta(days=2), 
#         category_id=example_category.id
#     )
#     sale.save()
#     return sale

# @pytest.mark.django_db
# def test_sale_list(client, example_sale, example_user):
#     client.login(username="login", password="Y7h3mi*ow9")
#     data = {
#         "username": "login",
#         "password": "Y7h3mi*ow9",
#         "password2": "Y7h3mi*ow9",
#     }

#     sale_list_view = SaleListView()
#     print(sale_list_view.get(RequestFactory().get(request="/api/sale/")))
#     assert 1 == 2
#     response = client.post("/api/auth/token/", data, format="json")
#     access_token = json.loads(response.content)["access"]
#     response = client.get("/api/sale/", format="json", headers={"Authorization": f"Bearer {access_token}"})
#     sales_response = json.loads(response.content)
#     for i in range(len(sales_response)):
#         sales_response[i]["category_id"] = sales_response[i]["category"]["id"]
#         del sales_response[i]["category"]
#         sales_response[i] = Sale(**sales_response[i])
#     sales = list(Sale.objects.all())
#     assert sales == sales_response