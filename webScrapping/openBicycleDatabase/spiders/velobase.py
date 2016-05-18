# -*- coding: utf-8 -*-
import scrapy
import time
from scrapy.selector import Selector
from selenium import webdriver
from openBicycleDatabase.items import ComponentItem, ComponentType, BrandItem, ImageItem
import traceback


class VelobaseSpider(scrapy.Spider):
    name = 'velobase'
    allowed_domains = ["velobase.com"]
    start_urls = ['http://velobase.com/ListComponents.aspx?ClearFilter=true']

    def __init__(self):
        self.driver = webdriver.Firefox()

    def __del__(self):
        self.driver.stop()

    def parse(self, response):
        hxs = Selector(response)
        self.driver.get(response.url)

        links = []
        while True:
            try:
                self.logger.info("=" * 60)
                self.logger.info(
                    "HXS  %s ", hxs.xpath('//tr[@class="GroupHeader1"]/td/a/text()').extract()[0])
                self.logger.info("=" * 60)

                links.extend(hxs.xpath(
                    '//table[@class="content"]//tr[@class="content_normal" or @class="content_alternate"]/td/a[@class=" ttiptxt"]/@href'
                ).extract())
                
                for link in links:
                    full_url = response.urljoin(link)
                    yield scrapy.Request(full_url, callback=self.parse_details)

                nextPage = self.driver.find_element_by_link_text('Next')
                nextPage.click()
                time.sleep(3)
                hxs = Selector(text=self.driver.page_source)

            except:
                traceback.print_exc()
                break

    def parse_details(self, response):
        component = ComponentItem()
        comptype = ComponentType()
        brand = BrandItem()
        image = ImageItem()
        comptype['Name'] = response.xpath(
            '//td[@id="ctl00_ContentPlaceHolder1_GenInfo"]/table/tr[1]/td/text()'
        ).extract()[-1]

        component['Name'] = response.xpath(
            '//td[@id="ctl00_ContentPlaceHolder1_GenInfo"]/table/tr[2]/td/text()'
        ).extract()[-1].strip()

        brand['Name'] = response.xpath(
            '//td[@id="ctl00_ContentPlaceHolder1_GenInfo"]/table/tr[3]/td/a/text()'
        ).extract()[-1].strip()

        component['Country'] = response.xpath(
            '//td[@id="ctl00_ContentPlaceHolder1_GenInfo"]/table/tr/td[contains(text(),"Country:")]/following-sibling::td/text()'
            ).extract()[-1].strip()

        component['Description'] = response.xpath(
            '//td[@id="ctl00_ContentPlaceHolder1_AddInfoCell"]/p/text()'
        ).extract()
        if len(component['Description']) != 0:
            component['Description'] = component['Description'][-1]
        else:
            component['Description'] = ''
        # Get Those Images !!!
        component['Images'] = []
        imageThumbs = response.xpath(
            '//span[@class="PhotoThumbContainer"]/a/img/@src'
        ).extract()
        # Convert the thumbs URL into Real Image URL (remove Thumbs etc ...)
        for imageThumb in imageThumbs:
            # Join to have the Full URL to be processed later 
            # (push to /images and ID retrieved)
            image['URL'] = response.urljoin(
                imageThumb.replace('Thumbs/tn_', '')
            )
            component['Images'].append(image.copy())

        component['Year'] = "0"
        component['Brand'] = dict(brand)
        component['Type'] = dict(comptype)

        self.logger.info("-" * 70)
        self.logger.info(" COMPONENT %s ", component['Name'])
        self.logger.info("-" * 70)
        yield component
