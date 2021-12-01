# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import sys
import glob
import asyncio
import logging
import importlib
import importlib
from aiohttp import web
from pathlib import Path
from pathlib import Path
from pyrogram import idle
from WebStreamer import bot_info
from WebStreamer.vars import Var
from WebStreamer.server import web_server
from WebStreamer.utils.keepalive import ping_server


logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logging.getLogger("pyrogram").setLevel(logging.ERROR)
logging.getLogger("aiohttp").setLevel(logging.ERROR)
logging.getLogger("aiohttp.web").setLevel(logging.ERROR)

loop = asyncio.get_event_loop()

_path = "WebStreamer/bot/plugins/*.py"
files = glob.glob(_path)

async def start_services():
    print('----------------------------- DONE -----------------------------')
    print('\n')
    print('--------------------------- Importing ---------------------------')
    for name in files:
        with open(name) as a:
            path_ = Path(a.name)
            plugin_name = path_.stem.replace(".py", "")
            plugins_dir = Path(f"WebStreamer/bot/plugins/{plugin_name}.py")
            import_path = ".plugins.{}".format(plugin_name)
            spec = importlib.util.spec_from_file_location(import_path, plugins_dir)
            load = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(load)
            sys.modules["WebStreamer.bot.plugins." + plugin_name] = load
            print("Imported => " + plugin_name)
    if Var.ON_HEROKU:
        print('------------------ Starting Keep Alive Service ------------------')
        print('\n')
        await asyncio.create_task(ping_server())
    print('-------------------- Initalizing Web Server --------------------')
    app = web.AppRunner(await web_server())
    await app.setup()
    bind_address = "0.0.0.0" if Var.ON_HEROKU else Var.BIND_ADDRESS
    await web.TCPSite(app, bind_address, Var.PORT).start()
    print('----------------------------- DONE -----------------------------')
    print('\n')
    print('----------------------- Service Started -----------------------')
    print('                        bot =>> {}'.format(bot_info.first_name))
    if bot_info.dc_id:
        print('                        DC ID =>> {}'.format(str(bot_info.dc_id)))
    print('                        server ip =>> {}'.format(bind_address, Var.PORT))
    if Var.ON_HEROKU:
        print('                        app running on =>> {}'.format(Var.FQDN))
    print('---------------------------------------------------------------')
    await idle()

if __name__ == '__main__':
    try:
        loop.run_until_complete(start_services())
    except KeyboardInterrupt:
        logging.info('----------------------- Service Stopped -----------------------')