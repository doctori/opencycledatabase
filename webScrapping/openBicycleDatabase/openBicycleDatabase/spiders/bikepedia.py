# -*- coding: utf-8 -*-
import scrapy


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
        bikeName = response.xpath('//span[@id="ctl00_MainContent_TitleOfBike_modelLabel2"]/text()').extract()
        bikeBrand = response.xpath('//span[@id="ctl00_MainContent_TitleOfBike_brandLabel2"]/text()').extract()
        bikeYear = response.xpath('//span[@id="ctl00_MainContent_TitleOfBike_yearLabel2"]/text()').extract()

        bikeFrontDerailleur = self.get_bike_element(
            response=response, elementToFind='Front Derailleur')
        bikeRearDerailleur = self.get_bike_element(
            response, 'Rear Derailleur')
        bikeCrankSet = self.get_bike_element(
            response, 'Crackset')
        bikePedals = self.get_bike_element(
            response, 'Pedals')
        bikeBottomBracket = self.get_bike_element(
            response, 'Bottom Bracket')
        bikeRearCogs = self.get_bike_element(
            response, 'Rear Cogs')
        bikeChain = self.get_bike_element(
            response, 'Chain')
        bikeSeatPost = self.get_bike_element(
            response, 'Seatpost')
        bikeSaddle = self.get_bike_element(
            response, 'Saddle')
        bikeHandlebar = self.get_bike_element(
            response, 'Handlebar')
        bikeHandlebarExt = self.get_bike_element(
            response, 'Handlebar Extensions')
        bikeHandlebarStem = self.get_bike_element(
            response, 'Handlebar Stem')
        bikeHeadSet = self.get_bike_element(
            response, 'Headset')
        bikeFrameType = self.get_frame_element(
            response, 'Frame Construction')
        bikeFrameTubing = self.get_frame_element(
            response, 'Frame Tubing Material')
        bikeForkModel = self.get_frame_element(
            response, 'Fork Brand & Model')
        bikeForkMaterial = self.get_frame_element(
            response, 'Fork Material')
        # TODO : Wheels Rear & Front

        yield {
            'bikeName': bikeName,
            'bikeBrand': bikeBrand,
            'bikeYear': bikeYear,
            'components': {
                'frontDerailleur': bikeFrontDerailleur,
                'rearDerailleur': bikeRearDerailleur,
                'crankSet': bikeCrankSet,
                'bottomBracket': bikeBottomBracket,
                'rearCogs': bikeRearCogs,
                'chain': bikeChain,
                'pedals': bikePedals,
                'seatpost': bikeSeatPost,
                'saddle': bikeSaddle,
                'handlebar': bikeHandlebar,
                'handlebarExt': bikeHandlebarExt,
                'stem': bikeHandlebarStem,
            },
            'frame': {
                'type': bikeFrameType,
                'tubing': bikeFrameTubing,
                'headset': bikeHeadSet
            },
            'fork': {
                'model': bikeForkModel,
                'material': bikeForkMaterial
            }
        }

    def get_bike_element(self, response, elementToFind):
        element =  response.xpath('//table[@id="ctl00_MainContent_CBSDetailsView3"]/tr/td[@class="FieldHeader" and contains(text(),"{0}")]/following-sibling::td/text()'.format(elementToFind)).extract()
        if isinstance(element, list) & len(element) > 0:
            return element[0]
        else:
            return element

    def get_frame_element(self, response, elementToFind):
    	element = response.xpath('//table[@id="ctl00_MainContent_CBSDetailsView2"]/tr/td[@class="FieldHeader" and contains(text(),"{0}")]/following-sibling::td/text()'.format(elementToFind)).extract()
        if isinstance(element, list) & len(element) > 0:
            return element[0]
        else:
            return element