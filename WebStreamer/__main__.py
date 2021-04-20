# This file is a part of TG-FileStreamBot
# Coding: @EverythingSuckz & @AbirHasan2005

import os
import sys
import glob
import asyncio
import logging
import importlib
from pathlib import Path
from pyrogram import idle
from .bot import StreamBot
from .vars import Var
from aiohttp import web
from .server import web_server

ppath = "WebStreamer/bot/plugins/*.py"
files = glob.glob(ppath)

loop = asyncio.get_event_loop()


async def start_services():
    print('\n')
    print('------------------- Initalizing Telegram Bot -------------------')
    await StreamBot.start()
    print('\n')
    print('---------------------- DONE ----------------------')
    print('\n')
    print('------------------- Importing -------------------')
    for name in files:
        with open(name) as a:
            patt = Path(a.name)
            plugin_name = patt.stem.replace(".py", "")
            plugins_dir = Path(f"WebStreamer/bot/plugins/{plugin_name}.py")
            import_path = ".plugins.{}".format(plugin_name)
            spec = importlib.util.spec_from_file_location(import_path, plugins_dir)
            load = importlib.util.module_from_spec(spec)
            spec.loader.exec_module(load)
            sys.modules["WebStreamer.bot.plugins." + plugin_name] = load
            print("Imported => " + plugin_name)
    print('\n')
    print('------------------- Initalizing Web Server -------------------')
    app = web.AppRunner(await web_server())
    await app.setup()
    bind_address = "0.0.0.0" if Var.ON_HEROKU else Var.FQDN
    await web.TCPSite(app, bind_address, Var.PORT).start()
    print('\n')
    print('----------------------- Service Started -----------------------')
    print('                        bot =>> {}'.format((await StreamBot.get_me()).first_name))
    print('                        server ip =>> {}:{}'.format(bind_address, Var.PORT))
    if Var.ON_HEROKU:
        print('                        app runnng on =>> {}'.format(Var.FQDN))
    print('---------------------------------------------------------------')
    await idle()

if __name__ == '__main__':
    try:
        loop.run_until_complete(start_services())
    except KeyboardInterrupt:
        print('----------------------- Service Stopped -----------------------')