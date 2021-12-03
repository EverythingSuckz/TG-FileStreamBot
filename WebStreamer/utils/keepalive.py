import asyncio
import logging
import aiohttp
import traceback
from WebStreamer.vars import Var

async def ping_server():
    sleep_time = Var.PING_INTERVAL
    while True:
        await asyncio.sleep(sleep_time) 
        try:
            async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=10)) as session:
                async with session.get(Var.URL) as resp:
                    logging.info("Pinged server with response: {}".format(resp.status))
        except TimeoutError:
            logging.warning("Couldn't connect to the site URL..!")
        except Exception:
            traceback.print_exc()
