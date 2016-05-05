# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html
import requests
from scrapy.utils.serialize import ScrapyJSONEncoder
_encoder = ScrapyJSONEncoder()


class OpenbicycledatabasePipeline(object):

    def __init__(self):
        self.file = open('items.jl', 'wb')

    def process_item(self, item, spider):
        line = _encoder.encode(item) + "\n"
        self.file.write(line)
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
