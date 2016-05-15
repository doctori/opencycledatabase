# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# http://doc.scrapy.org/en/latest/topics/items.html

from scrapy.item import Item, Field


class BikeItem(Item):
    # define the fields for your item here like:
    # name = scrapy.Field()
    Name = Field()
    Year = Field()
    # Will be an array on Compoennt Item
    Components = Field()
    Description = Field()
    Brand = Field()


class BrandItem(Item):
    Name = Field()
    Description = Field()


class ComponentType(Item):
    Name = Field()
    Description = Field()


class ImageItem(Item):
    Name = Field()
    ID   = Field()
    URL  = Field()


class ComponentItem(Item):
    Name = Field()
    Brand = Field()
    Country = Field()
    Type = Field()
    Description = Field()
    Year = Field()
    Standard = Field()
    Images = Field()
