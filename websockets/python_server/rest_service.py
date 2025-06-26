import requests
from requests.auth import HTTPBasicAuth

def post():
    url = 'http://localhost:8080/alert'
    data = {'message': 'god'}
    # response = requests.post(url, data=data, auth=HTTPBasicAuth("", "mySecretKey-10101"))# Handling the response object
    response = requests.post(url, data=data, verify=False)# Handling the response object
    if response.status_code == 201:
        print('Post request successful!')
        print('Response Content:', response.json())
    else:
        print('Request failed with status code:', response.status_code)

if __name__ == "__main__":
    post()