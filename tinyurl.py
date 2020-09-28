import requests

base_url = 'https://tinyurl.com/api-create.php?url='

def shorten(link='https://chand1012.dev/'):
    response = requests.get(base_url + link)
    return response.text