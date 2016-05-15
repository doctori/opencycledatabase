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


class OpenbicycledatabasePipeline(object):

    def __init__(self):
        self.file = open('items.jl', 'wb')

    def process_item(self, item, spider):
        line = _encoder.encode(item) + "\n"
        self.file.write(line)
        return item


class ImagePipeline(ImagesPipeline):
    """docstring for ImagePipeline"""

    def get_media_requests(self, item, info):
        for imageUrl in item['Images']:
            yield scrapy.Request(imageUrl['URL'])

    def item_completed(self, results, item, info):
        image_paths = [x['path'] for ok, x in results if ok]
        images = []
        image = ImageItem()
        for image_path in image_paths:
            image['Name'] = basename(image_path)
            logger.info("uploading {0}".format(image_path))
            file = {
                'file': (
                    # image['Name'],
                    # Should retrieve the configuration Key
                    open('/tmp/' + image_path, 'rb')
                )
            }
            result = requests.post(
                'http://127.0.0.1:8080/images',
                files=file)
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
            resp = requests.post(
                url='http://127.0.0.1:8080/components',
                data=_encoder.encode(item))
            if resp.status_code != 200:
                print("ERROR {0}".format(resp))
                return item
            else:
                print('Created component. ID: {0}'.format(resp.json()["ID"]))
                return item
        else:
            resp = requests.post(
                url='http://127.0.0.1:8080/bikes',
                data=_encoder.encode(item))
            if resp.status_code != 200:
                print("ERROR {0}".format(resp))
                return item
            else:
                print('Created component. ID: {0}'.format(resp.json()["ID"]))
                return item
