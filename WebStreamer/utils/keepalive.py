import logging
import requests
from ..vars import Var
def ping_server():
    k = requests.get(f'https://ping-pong-sn.herokuapp.com/pingback?link={Var.URL}').json()
    if not k.get('error'):
        logging.info('KeepAliveService: Pinged {} with status {}'.format(Var.FQDN, k.get('Status')))
    else:
        logging.error('Couldn\'t Ping the Server!')
