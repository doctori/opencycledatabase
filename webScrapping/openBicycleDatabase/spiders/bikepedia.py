# -*- coding: utf-8 -*-
import scrapy
from openBicycleDatabase.items import BikeItem, BrandItem, ComponentItem, ComponentType


class BikepediaSpider(scrapy.Spider):
    name = "bikepedia"
    allowed_domains = ["www.bikepedia.com"]
    start_urls = (
        'http://www.bikepedia.com/QuickBike',
    )

    def parse(self, response):
        for href in response.css('#ctl00_MainContent_yearsDL').xpath('tr//a/@href'):
            full_url = response.urljoin(href.extract())
            yield scrapy.Request(full_url, callback=self.parse_year)

    def parse_year(self, response):
        for href in response.css('.qbBrandLI').xpath('a/@href'):
            full_url = response.urljoin(href.extract())
            yield scrapy.Request(full_url, callback=self.parse_brand)

    def parse_brand(self, response):
        for bike in response.selector.css('.qbBrandLI').xpath('a/@href'):
            full_url = response.urljoin(bike.extract())
            yield scrapy.Request(full_url, callback=self.parse_bike)

    def parse_bike(self, response):
        bike = BikeItem()
        bikeBrand = BrandItem()
        bikeComponents = []
        bikeComponent = ComponentItem()
        bikeComponentType = ComponentType()
        bikeComponentTypes = [
            'Front Derailleur',
            'Rear Derailleur',
            'Crackset',
            'Pedals',
            'Bottom Bracket',
            'Rear Cogs',
            'Chain',
            'Seatpost',
            'Saddle',
            'Handlebar',
            'Handlebar Extensions',
            'Handlebar Stem',
            'Headset',

        ]

        bike['Name'] = response.xpath(
            '//span[@id="ctl00_MainContent_TitleOfBike_modelLabel2"]/text()'
        ).extract()[0]

        bikeBrand['Name'] = response.xpath(
            '//span[@id="ctl00_MainContent_TitleOfBike_brandLabel2"]/text()'
        ).extract()[0]
        bike['Brand'] = bikeBrand
        bike['Year'] = response.xpath(
            '//span[@id="ctl00_MainContent_TitleOfBike_yearLabel2"]/text()'
        ).extract()[0]

        for compType in bikeComponentTypes:
            bikeComponentType['Name'] = compType
            bikeComponent['Type'] = bikeComponentType.copy()
            bikeComponent['Name'] = self.get_bike_element(
                response, bikeComponentType.get('Name'))
            # We Add that component to the list
            #  of installed component on that bike
            bikeComponents.append(bikeComponent.copy())

        # Frame
        bikeComponentType['Name'] = 'Frame'
        bikeComponent['Type'] = bikeComponentType.copy()
        bikeComponent['Name'] = self.get_frame_element(
            response, 'Frame Construction')
        bikeComponent['Description'] = self.get_frame_element(
            response, 'Frame Tubing Material')
        bikeComponents.append(bikeComponent.copy())

        # Fork
        bikeComponentType['Name'] = 'Front Fork'
        bikeComponent['Type'] = bikeComponentType.copy()
        bikeComponent['Name'] = self.get_frame_element(
            response, 'Fork Brand & Model')
        bikeComponent['Description'] = self.get_frame_element(
            response, 'Fork Material')
        bikeComponents.append(bikeComponent.copy())
        # TODO : Wheels Rear & Front
        bike['Components'] = bikeComponents
        yield bike

    def get_bike_element(self, response, elementToFind):
        element =  response.xpath('//table[@id="ctl00_MainContent_CBSDetailsView3"]/tr/td[@class="FieldHeader" and contains(text(),"{0}")]/following-sibling::td/text()'.format(elementToFind)).extract()
        if isinstance(element, list) & len(element) > 0:
            return element[0]
        else:
            return ""

    def get_frame_element(self, response, elementToFind):
        element = response.xpath('//table[@id="ctl00_MainContent_CBSDetailsView2"]/tr/td[@class="FieldHeader" and contains(text(),"{0}")]/following-sibling::td/text()'.format(elementToFind)).extract()
        if isinstance(element, list) & len(element) > 0:
            return element[0]
        else:
            return ""