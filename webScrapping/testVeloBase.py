import scrapy
import time
from scrapy.selector import Selector
from selenium import webdriver
import traceback


class VeloBaseScrapper(scrapy.Spider):
    name = 'velobase'
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
                self.logger.info("="*60)
                self.logger.info("HXS  %s ", hxs.xpath('//tr[@class="GroupHeader1"]/td/a/text()').extract()[0])
                self.logger.info("="*60)
                self.logger.info("="*60)
                # self.logger.info("LINKS  %s ", hxs.xpath('//table[@class="content"]/tr[@class="content_normal" or @class="content_alternate"]/td[0]/a[@class="ttiptxt"]/@href).extract()[0]'))
                links.extend(hxs.xpath(
                    '//table[@class="content"]//tr[@class="content_normal" or @class="content_alternate"]/td/a[@class=" ttiptxt"]/@href'
                    ).extract())
                self.logger.info("LINKS  %s ", links)
                self.logger.info("="*60)

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
        componentCategory = response.xpath('//td[@id="ctl00_ContentPlaceHolder1_GenInfo"]/table/tr[1]/td/text()').extract()[-1]
        componentName = response.xpath('//td[@id="ctl00_ContentPlaceHolder1_GenInfo"]/table/tr[2]/td/text()').extract()[-1]
        componentBrand = response.xpath('//td[@id="ctl00_ContentPlaceHolder1_GenInfo"]/table/tr[3]/td/a/text()').extract()[-1]
        componentCountry = response.xpath('//td[@id="ctl00_ContentPlaceHolder1_GenInfo"]/table/tr[7]/td/text()').extract()[-1]
        componentDescription = response.xpath('//td[@id="ctl00_ContentPlaceHolder1_GenInfo"]/table/tr/td[contains(text(),"Country:")]/following-sibling::td').extract()

        self.logger.info("-" * 70)
        self.logger.info(" COMPONENT %s ", componentName)
        self.logger.info("-" * 70)
        yield {
            'category': componentCategory,
            'name': componentName,
            'brand': componentBrand,
            'country': componentCountry,
            'description': componentDescription
        }
