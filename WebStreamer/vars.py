# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

from os import getenv, environ
from dotenv import load_dotenv

load_dotenv()

class Var(object):
    ENV = bool(getenv('ENV', False))
    API_ID = int(getenv('API_ID'))
    API_HASH = str(getenv('API_HASH'))
    BOT_TOKEN = str(getenv('BOT_TOKEN'))
    SLEEP_THRESHOLD = int(getenv('SLEEP_THRESHOLD', '300'))
    WORKERS = int(getenv('WORKERS', '3'))
    BIN_CHANNEL = int(getenv('BIN_CHANNEL', None))
    FQDN = str(getenv('FQDN', 'localhost'))
    PORT = int(getenv('PORT', 8080))
    BIND_ADRESS = str(getenv('BIND_ADRESS', '0.0.0.0'))
    CACHE_DIR = str(getenv('CACHE_DIR', 'WebStreamer/bot/cache'))
    OWNER_ID = int(getenv('OWNER_ID'))
    if 'DYNO' in environ:
        ON_HEROKU = True
    else:
        ON_HEROKU = False