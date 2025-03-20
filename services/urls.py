from django.urls import path
from .views import CategoryListCreateView, SubcategoryListCreateView, ServiceListCreateView

urlpatterns = [
    path('categories/', CategoryListCreateView.as_view(), name='category-list'),
    path('subcategories/', SubcategoryListCreateView.as_view(), name='subcategory-list'),
    path('services/', ServiceListCreateView.as_view(), name='service-list'),
]
