# -*- coding: utf-8 -*-
import scrapy
import re
import json
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
                },
            callback=self.parse_brands)]
    
    def parse_brands(self, response):
        #i = OpenbicycledatabaseItem()
        #i['domain_id'] = response.xpath('//input[@id="sid"]/@value').extract()
        #i['name'] = response.xpath('//div[@id="name"]').extract()
        #i['description'] = response.xpath('//div[@id="description"]').extract()
        brands = response.xpath('//section[@id="brand_tab"]/script/text()').extract()
        self.logger.info("="*60)
        self.logger.info("URL  %s ", response.url)
        #self.logger.info(brands)
        for brand in brands:
            for link in self.extract_links(brand.strip()):
                full_url = response.urljoin(link)
                yield scrapy.Request(full_url, callback=self.parse_articles)
        self.logger.info("="*60)        
        pass

    def parse_articles(self,response):
        articles = response.xpath('//li[@class="description"]/a/@href').extract()
        for article in articles:
            full_url = response.urljoin(article)
            yield scrapy.Request(full_url,callback=self.parse_details)


    def parse_details(self,response):
        # fetch the JS variable in order to have more details
        universal_variables = response.xpath('/html/head/script/text()[contains(.,"universal_variable")]').extract()
        if len(universal_variables) > 0:
            universal_variables = self.get_universal_vars(universal_variables[0].strip())
            self.logger.info(universal_variables.get("name"))
            self.logger.info(universal_variables.get("manufacturer"))
            self.logger.info(universal_variables.get("category"))
            self.logger.info(universal_variables.get("subcategory"))
        description = response.xpath('//div[@id="crcPDPComponentDescription"]/p/text()').extract()
        description.pop()
        description.pop()
        technical_details = response.xpath('//div[@id="crcPDPComponentDescription"]/ul/li/text()').extract()
        self.logger.warn("============= {0} ============".format(description))
        for technical_detail in technical_details:
            self.logger.info(technical_detail)

    def get_universal_vars(self,string):
        regex = re.compile(".*('product'\:{.*}),.*",re.MULTILINE | re.DOTALL)
        matcher = regex.match(string)
        if matcher:
            result = matcher.group(1).replace("'product':","").replace("'",'"')
            if result:
                return json.loads(result)
        else:
            self.logger.info("LOOOOSER : [{0}]".format(string))

    def extract_links(self,string):
        # String shouild looks like something like : 
        # buildLinksContainer('{WTB=/wtb, Wahoo=/wahoo, Walz=/walz, WeThePeople=/wethepeople, Weldtite=/weldtite, Wellgo=/wellgo, Wilier=/wilier, Wippermann=/wippermann}', 'brands_w', 'W', 'brand');

        matcher = re.match("buildLinksContainer\('{(.*)}',.*\);",string)
        if matcher:
            splited_matches = matcher.group(1).split(',')
            for splited_match in splited_matches:
                self.logger.info(splited_match.strip())
                key_values=splited_match.split('=')
                if len(key_values) == 2:
                    link = key_values[1]
                    yield link
        else: 
            self.logger.info("LOOOOSER : [{0}])".format(string))
