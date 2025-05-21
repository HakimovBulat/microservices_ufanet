from django.contrib import admin
from .models import Category, Sale


@admin.register(Category)
class CategoryAdmin(admin.ModelAdmin):
    list_display = ["title", "photo"]
    

@admin.register(Sale)
class SaleAdmin(admin.ModelAdmin):
    list_display = ["title", "subtitle", "start_date", "end_date", "category"]