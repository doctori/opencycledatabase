# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html
import requests
import scrapy
import logging
from os.path import basename
from openBicycleDatabase.items import ImageItem
from scrapy.utils.serialize import ScrapyJSONEncoder
from scrapy.pipelines.images import ImagesPipeline
from scrapy.exceptions import DropItem

_encoder = ScrapyJSONEncoder()

logger = logging.getLogger(__name__)
backendHost = '127.0.0.1'
backendPort = '8080'


class ImagePipeline(ImagesPipeline):
    """docstring for ImagePipeline"""

    def get_media_requests(self, item, info):
        for imageUrl in item['Images']:
            yield scrapy.Request(imageUrl['URL'])

    def item_completed(self, results, item, info):
        image_items = [x for ok, x in results if ok]
        images = []
        image = ImageItem()
        for image_item in image_items:
            image['Name'] = basename(image_item['path'])
            image['Source'] = image_item['url']
            logger.info("uploading {0}".format(image['Name']))
            file = {
                'file': (
                    image['Name'],
                    open('/tmp/' + image_item['path'], 'rb'),
                    'image/jpeg',
                    {
                        'Source': image_item['url']
                    }
                )
            }
            #file = {
            #    'file': open('/tmp/' + image_item['path'], 'rb')
            #}
            data = {
                "IMGSource": image_item['url']
            }
            result = requests.post(
                'http://127.0.0.1:8080/images',
                files=file,
                data=data
            )
            if result.status_code != 200:
                raise DropItem("Could Not Upload the Images")
            print result.status_code
            print result.text
            logger.info(result.json()["ID"])
            images.append(result.json())

        item['Images'] = images
        return item


class ComponentPipeline(object):

    def __init__(self):
        pass

    def process_item(self, item, spider):
        if type(item).__name__ == 'ComponentItem':
            item = self.save_component(item)
            return item
        elif type(item).__name__ == 'BikeItem':
            # Save the component before saving the bike
            # (GORM seems to behave better that way)
            if 'Brand' in item:
                item['Brand'] = self.save_brand(item['Brand'])
            for i, component in enumerate(item['Components']):
                item['Components'][i] = self.save_component(component)

            resp = requests.post(
                url='http://{0}:{1}/bikes'.format(backendHost, backendPort),
                data=_encoder.encode(item))
            if resp.status_code != 200:
                raise DropItem("Could not Save the Bike")

            else:
                print('Created component. ID: {0}'.format(resp.json()["ID"]))
                return item
        else:
            raise DropItem("Item Unkown")

    def save_component(self, component):
        if 'Brand' in component:
            component['Brand'] = self.save_brand(component['Brand'])
        resp = requests.post(
            url='http://{0}:{1}/components'.format(backendHost, backendPort),
            data=_encoder.encode(component))
        if resp.status_code != 200:
            raise DropItem("Could not save the component")
        return resp.json()

    def save_brand(self, brand):
        resp = requests.post(
            url='http://{0}:{1}/brands'.format(backendHost, backendPort),
            data=_encoder.encode(brand))
        if resp.status_code != 200:
            raise DropItem("Could not save the component")
        return resp.json()
