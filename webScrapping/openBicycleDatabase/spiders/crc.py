# -*- coding: utf-8 -*-
import scrapy
from scrapy.linkextractors.lxmlhtml import LxmlLinkExtractor
from scrapy.spiders import CrawlSpider, Rule

#from openBicycleDatabase.items import bikeItem


class CrcSpider(CrawlSpider):
    name = 'crc'
    allowed_domains = ['www.chainreactioncycles.com']
#    start_urls = ['http://www.chainreactioncycles.com/sitemap']

    rules = (
        Rule(LxmlLinkExtractor(allow_domains='chainreactioncycles.com',restrict_xpaths='//*[contains(@id,"brand")]'), callback='parse_item', follow=True),
        # Rule(LinkExtractor(allow_domains='chainreactioncycles.com',restrict_css='.link_container'), callback='parse_item'),
    )
    def start_requests(self):
        return [scrapy.Request("http://www.chainreactioncycles.com/sitemap",
                cookies = {
                    'languageCode': 'en',
                    'countryCode': 'GB',
                    'currencyCode': 'GBP'},
                meta = {
                'dont_redirect': True,
                'handle_httpstatus_list': [301, 302]
                    })]
    def parse_item(self, response):
        #i = OpenbicycledatabaseItem()
        #i['domain_id'] = response.xpath('//input[@id="sid"]/@value').extract()
        #i['name'] = response.xpath('//div[@id="name"]').extract()
        #i['description'] = response.xpath('//div[@id="description"]').extract()
        products = response.xpath('//div[@id="grid-view"]/div[@class="grid_view_row"]//ul')
        self.logger.info("="*60)
        self.logger.info("URL  %s ", response.url)

        for product in products:
            description = product.xpath('li[@class="description"]/a/text()').extract()
            link = product.xpath('li[@class="description"]/a/@href').extract()
            self.logger.info('DESCRIPTION :: {0} ::'.format(description))
            self.logger.info('>>>>NEXT URL :: {0} :: <<<<<<<'.format(link))
        self.logger.info("="*60)        
        pass
