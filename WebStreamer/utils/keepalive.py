import logging
import aiohttp
from WebStreamer.vars import Var

async def ping_server():
    try:
        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=10)) as session:
            async with session.get(Var.URL) as resp:
                logging.info("Pinged server with response: {}".format(resp.status))
    except TimeoutError:
        logging.warning("Couldn't connect to the site URL..!")