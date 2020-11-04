# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html
import json
import requests


class OpenbicycledatabasePipeline(object):

    def __init__(self):
        self.file = open('items.jl', 'wb')

    def process_item(self, item, spider):
        line = json.dumps(dict(item)) + "\n"
        self.file.write(line)
        return item


class ComponentPipeline(object):

    def __init__(self):
        pass

    def process_item(self, item, spider):
        resp = requests.post(
            url='http://127.0.0.1:8080/components',
            data=json.dumps(dict(item)))
        if resp.status_code != 200:
            print("ERROR {0}".format(resp))
            return False
        else:
            print('Created component. ID: {0}'.format(resp.json()["ID"]))
            return True
