from django.db import models

class Category(models.Model):
    category_id = models.AutoField(primary_key=True)
    category_name = models.CharField(max_length=255)

    class Meta:
        db_table = "categories"

    def __str__(self):
        return self.category_name

class SubCategory(models.Model):
    subcategory_id = models.AutoField(primary_key=True)
    category = models.ForeignKey(Category, on_delete=models.CASCADE)
    subcategory_name = models.CharField(max_length=255)

    class Meta:
        db_table = "subcategories"

    def __str__(self):
        return f"{self.category.category_name} - {self.subcategory_name}"

class Service(models.Model):
    service_id = models.AutoField(primary_key=True)
    category = models.ForeignKey(Category, on_delete=models.CASCADE)
    subcategory = models.ForeignKey(SubCategory, on_delete=models.CASCADE)
    description = models.TextField()

    class Meta:
        db_table = "service"

    def __str__(self):
        return f"{self.subcategory.subcategory_name} - {self.description[:50]}"
