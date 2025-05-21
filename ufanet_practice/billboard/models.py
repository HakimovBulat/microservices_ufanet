from django.db import models


class Category(models.Model):
    title = models.CharField("Заголовок", max_length=30)
    photo = models.ImageField("Фото", upload_to ='static/media/', default="static/media/category-new_MROs2uS.png")

    def __str__(self):
        return self.title

    class Meta:
        verbose_name = "Категория"
        verbose_name_plural = "Категории"
        ordering = ["id"]


class Sale(models.Model):
    title = models.CharField("Заголовок", max_length=50)
    subtitle = models.CharField("Подзаголовок", max_length=50, null=True, blank=True)
    description = models.TextField("Описание")

    photo = models.ImageField("Фото", upload_to ='static/media/', default="static/media/category-new_MROs2uS.png")
    promocode = models.CharField("Промокод", max_length=20, null=True, blank=True)
    url = models.URLField("Ссылка", blank=True, null=True)
    about_partner = models.TextField("О партнере")
    start_date = models.DateField("Начало акции")
    end_date = models.DateField("Окончание акции")
    category = models.ForeignKey(Category, on_delete=models.CASCADE, verbose_name="Категория")
    
    def __str__(self):
        return self.title

    class Meta:
        verbose_name = "Скидка"
        verbose_name_plural = "Скидки"
        ordering = ["id"]