from rest_framework import generics
from .models import Category, SubCategory, Service
from .serializers import CategorySerializer, SubcategorySerializer, ServiceSerializer

class CategoryListCreateView(generics.ListCreateAPIView):
    queryset = Category.objects.all()
    serializer_class = CategorySerializer

class SubcategoryListCreateView(generics.ListCreateAPIView):
    queryset = SubCategory.objects.all()
    serializer_class = SubcategorySerializer

class ServiceListCreateView(generics.ListCreateAPIView):
    queryset = Service.objects.all()
    serializer_class = ServiceSerializer
