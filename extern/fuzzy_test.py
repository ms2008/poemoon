#!/usr/bin/env python
# -*- coding:utf-8 -*-

import urllib
import urllib2
import random
import re
import time
import webbrowser

URL = "http://192.168.161.2/"


def post(url, data=None):
    if data:
        data = urllib.urlencode(data)
    req = urllib2.Request(url, data)
    req.add_header('User-Agent', 'Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.116 Safari/537.36')
    req.add_header('Referer', URL)
    response = urllib2.urlopen(req)
    return response

def check(response):
    lineBeChecked = response.readlines()[5]
    if lineBeChecked[1] == 's':
        return True
    elif lineBeChecked[1] == 'S':
        return False
    else:
        print("Unexpected response!")
        return False

def login(username, upass):
    data = {}
    data["DDDDD"] = username
    data["upass"] = upass
    data["0MKKey"] = ''
    return post(URL, data)

def logout():
    post(URL+"F.htm")

def flux(lineBeChecked):
    indexStart = lineBeChecked.find("time") + 6
    indexEnd = lineBeChecked[indexStart:].find(" ") + indexStart
    #print("indexStart: ", indexStart, "indexEnd: ", indexEnd, "calced: ", lineBeChecked[indexStart:indexEnd])
    return int(lineBeChecked[indexStart:indexEnd])

def fee(lineBeChecked):
    feeArray = re.findall(r";fee='(\d+) *';", lineBeChecked, re.S)
    amount = (float(feeArray[0])-float(feeArray[0])%100)/10000
    return amount


def main():
    f = open('ids.txt','r')
    reader = f.readlines()
    random.seed(time.time())

    while True:
        line = random.choice(reader)
        line = line.strip('\r\n')
        i = line.split('\t')
        username = i[0]
        password = i[-1]
        #print(username, password)

        try:
            if check(login(username, password)):
                #print("Connected!")
                balanceAmount = post(URL).readlines()[6]
                usedFlux = flux(balanceAmount)
                remainedFee = fee(balanceAmount)

                print username, password, "Used Time:", usedFlux, "Balance:", remainedFee

                if usedFlux >= 2400 and remainedFee == 0:
                    #print username, "time exceed!"
                    pass
                else:
                    webbrowser.open(URL)
                    break

            else:
                print(username, "login failed!")
                #pass
        except:
            print(username, "Network Error!")
            continue

    f.close()

if __name__ == '__main__':
    main()
    print("看到这个，就是想证明下我不是个恶意程序，5s 之后就看不到我啦 :-)")
    time.sleep(5)
