from django.shortcuts import render
from .models import Category, Sale
from django.db.models import Q
from django.contrib.auth.decorators import login_required
from confluent_kafka import Producer
import json
from django.forms.models import model_to_dict
import datetime
import uuid

conf = {'bootstrap.servers': 'localhost:29092'}

producer = Producer(conf)

TOPIC_SALE = 'wal_listener.public_billboard_sale'
TOPIC_CATEGORY = 'wal_listener.public_billboard_category'

def model_to_json(model):
    data_model = model_to_dict(model)
    print(str(uuid.uuid1()))
    json_data = {
        "id": uuid.uuid1(),
        "data": data_model,
        "schema": "public",
        "action": "SELECT",
        "dataOld": {},
        "commitTime": datetime.datetime.now()
    }
    json_data["data"]["photo"] = model.photo.url

    if "end_date" in data_model.keys():
        json_data["table"] = "billboard_sale"
    else:
        json_data["table"] = "billboard_category"

    return json.dumps(json_data, default=str)


def post_to_topic(model, topic):
    producer.produce(topic, value=model_to_json(model), callback=delivery_report)
    producer.poll(0)
    producer.flush()


def delivery_report(err, msg):
    if err is not None:
        print(f'Ошибка доставки сообщения: {err}')
    else:
        print(f'Сообщение доставлено в {msg.topic()} [{msg.partition()}]')


@login_required
def index(request):
    if request.method == "POST":
        sale_search = request.POST.get("sale_search", "").lower()

        query = Q(title__icontains=sale_search)
        query.add(Q(subtitle__icontains=sale_search), Q.OR)
        query.add(Q(about_partner__icontains=sale_search), Q.OR)
        query.add(Q(description__icontains=sale_search), Q.OR)

        sales = Sale.objects.filter(query)
        context = {
            "sales": sales,
            "categories": [],
        }
    else:
        categories = Category.objects.all()
        sales = Sale.objects.all()
        context = {
            "categories": categories,
            "sales": sales,
        }
    return render(request, "billboard/index.html", context=context)


@login_required
def category_sales(request, pk):
    category = Category.objects.get(pk=pk)
    context = {
        "category_title": category.title,
        "sales": category.sale_set.all()
    }
    post_to_topic(category, TOPIC_CATEGORY)
    return render(request, "billboard/category_sales.html", context=context)


@login_required
def sale_info(request, sale_pk):
    sale = Sale.objects.get(pk=sale_pk)
    context = {
        "sale": sale
    }
    post_to_topic(sale, TOPIC_SALE)
    return render(request, "billboard/sale_info.html", context=context)