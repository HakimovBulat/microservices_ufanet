from django.utils import timezone
from django.contrib.auth.models import User
from .models import Sale, Category
import datetime
import pytest


@pytest.fixture
def example_user():
    user = User.objects.create_user(
        username='login',
        password='Y7h3mi*ow9'
    )
    return user


@pytest.fixture
def example_category():
    category = Category.objects.create(title="Тестовая категория")
    category.save()
    return category


@pytest.fixture
def example_sale(example_category):
    sale = Sale.objects.create(
        title="Тестовая скидка", 
        start_date=timezone.now(), 
        end_date=timezone.now() + datetime.timedelta(days=2), 
        category_id=example_category.id
    )
    sale.save()
    return sale


@pytest.mark.django_db
def test_category_to_string(example_category):
    assert example_category.title == str(example_category)


@pytest.mark.django_db
def test_sale_to_string(example_sale):
    assert example_sale.title == str(example_sale)


@pytest.mark.django_db
def test_sale_info_context_sale(client, example_sale, example_user):
    client.login(username="login", password="Y7h3mi*ow9")
    response = client.get(f"/billboard/sale/{example_sale.id}")
    assert response.context["sale"] == example_sale


@pytest.mark.django_db
def test_index_search_post(client, example_user, example_sale):
    client.login(username="login", password="Y7h3mi*ow9")
    data = {
        "sale_search": "скидка"
    }
    response = client.post("/billboard/", data)
    sales = Sale.objects.filter(title__icontains="Тестовая")
    sales = list(sales)
    sales_search = list(response.context["sales"])

    assert sales == sales_search


@pytest.mark.django_db
def test_index_get_categories(client, example_category, example_user):
    client.login(username="login", password="Y7h3mi*ow9")
    response = client.get("/billboard/")
    categories = Category.objects.all()
    categories_response = response.context["categories"]
    for i in range(len(categories)):
        assert categories[i] == categories_response[i]


@pytest.mark.django_db
def test_index_get_sales(client, example_sale, example_user):
    client.login(username="login", password="Y7h3mi*ow9")
    response = client.get("/billboard/")
    sales = Sale.objects.all()
    sales_response = response.context["sales"]
    for i in range(len(sales)):
        assert sales[i] == sales_response[i]


@pytest.mark.django_db
def test_category_sales_get_category_title(client, example_category, example_user):
    client.login(username="login", password="Y7h3mi*ow9")
    response = client.get(f"/billboard/{example_category.id}")
    assert example_category.title == response.context["category_title"]


@pytest.mark.django_db
def test_category_sales_get_sales(client, example_sale, example_category, example_user):
    client.login(username="login", password="Y7h3mi*ow9")
    response = client.get(f"/billboard/{example_category.id}")
    sales = Sale.objects.filter(category_id=example_category.id)
    sales = list(sales)
    response_sales = list(response.context["sales"])
    for i in range(len(sales)):
        sale_dict = sales[i].__dict__
        response_sale_dict = response_sales[i].__dict__
        for key in sale_dict:
            if key != "_state":
                assert sale_dict.get(key) == response_sale_dict.get(key)