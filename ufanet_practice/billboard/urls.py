from django.urls import path
from . import views


urlpatterns = [
    path("", views.index, name="index"),
    path("<int:pk>", views.category_sales, name="category_sales"),
    path("sale/<int:sale_pk>", views.sale_info, name="sale_info"),
]