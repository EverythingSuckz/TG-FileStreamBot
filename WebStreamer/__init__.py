# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]


import time
from WebStreamer.bot import StreamBot

print('\n')
print('------------------- Initalizing Telegram Bot -------------------')

StreamBot.start()
bot_info = StreamBot.get_me()
__version__ = 1.03
StartTime = time.time()


