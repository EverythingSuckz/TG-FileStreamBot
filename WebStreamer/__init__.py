# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]


import time
from .vars import Var
from WebStreamer.bot.clients import StreamBot

print("\n")
print("------------------- Initializing Telegram Bot -------------------")

StreamBot.start()
bot_info = StreamBot.get_me()
__version__ = 2.0
StartTime = time.time()
