# Taken from megadlbot_oss <https://github.com/eyaadh/megadlbot_oss/blob/master/mega/webserver/__init__.py>
# Thanks to Eyaadh <https://github.com/eyaadh>
# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import logging
from aiohttp import web
from .stream_routes import routes

logger = logging.getLogger("server")

def web_server():
    logger.info("Initializing..")
    web_app = web.Application(client_max_size=30000000)
    web_app.add_routes(routes)
    logger.info("Added routes")
    return web_app
