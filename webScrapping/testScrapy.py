import scrapy

class BikePediaScapper(scrapy.Spider):
	name = 'bikepedia'
	start_urls = ['http://www.bikepedia.com/QuickBike/Bikes.aspx?Year=2015']

	def parse(self, response):
		for href in response.css('.qbBrandLI').xpath('a/@href'):
			full_url = response.urljoin(href.extract())
			yield scrapy.Request(full_url,callback=self.parse_brand)

			
	def parse_brand(self, response):
		for bike in response.selector.css('.qbBrandLI'):
			bikeName = bike.xpath('a/span/text()').extract()
			bikeBrand = bike.xpath('a/@href').re('.*brand=(\w+).*')
			yield {
				'bike': bikeName,
				'brand': bikeBrand

			}