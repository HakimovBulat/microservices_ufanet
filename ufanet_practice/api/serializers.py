from rest_framework import serializers
from billboard.models import Sale, Category


class CategorySerializer(serializers.ModelSerializer):
    class Meta:
        model = Category
        fields = ["id", "title", "photo"]


class SaleSerializer(serializers.ModelSerializer):
    category = CategorySerializer(many=False, read_only=False)
    class Meta:
        model = Sale
        fields = [
            "id", "title", "subtitle", "description", 
            "url", "promocode", "photo", 
            "about_partner", "start_date", "end_date", "category"
            ]
