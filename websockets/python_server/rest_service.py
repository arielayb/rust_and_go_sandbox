import requests
from requests.auth import HTTPBasicAuth
import argparse

'''
def post(userid, global):
    url = 'http://127.0.0.1:8080/alert'
    # response = requests.post(url, data=data, auth=HTTPBasicAuth("", "mySecretKey-10101"))# Handling the response object
    for i in range(4):
        data = {"user_id": userid, "method": "USER_INFO", "msg": f'hello {i}', "global": f'{global}'}
        response = requests.post(url, json=data, verify=False)# Handling the response object
        if response.status_code == 201:
            print('Post request successful!')
            print('Response Content:', response.json())
        else:
            print('Request failed with status code:', response.status_code)
'''

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog='python REST POST Test Script',
        description='make REST POST calls',
        epilog='Text at the bottom of help')
    
    parser.add_argument("-p", "--post", help = "Show Output")
    parser.add_argument("-g", "--gmsg", help = "global message to all")
    args = parser.parse_args()

    url = 'http://127.0.0.1:8080/alert'
    # response = requests.post(url, data=data, auth=HTTPBasicAuth("", "mySecretKey-10101"))# Handling the response object
    for i in range(4):
        print("the arg for global: ", args.gmsg)
        data = {"user_id": f'{args.post}', "method": "USER_INFO", "msg": f'hello {i}', "global": args.gmsg}
        response = requests.post(url, json=data, verify=False)# Handling the response object
        if response.status_code == 200:
            print('Post request successful!')
            print("Response Content:", f'{response.json}')
        else:
            print('Request failed with status code:', response.status_code)
            print('Request failed with status code:', response.json)

