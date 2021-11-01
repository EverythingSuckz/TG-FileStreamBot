# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

from WebStreamer.bot import StreamBot
from WebStreamer.vars import Var
from pyrogram import filters, emoji
from pyrogram.types import InlineKeyboardMarkup, InlineKeyboardButton

@StreamBot.on_message(filters.command(['start', 'help']))
async def start(b, m):
    await m.reply('Hi, Send me a file to get an instant stream link.',
                  reply_markup=InlineKeyboardMarkup(
                      [
                          [
                              InlineKeyboardButton(
                                  f'{emoji.STAR} Channel {emoji.STAR}',
                                  url='https://t.me/Cinema_Great'
                              )
                          ]
                      ]
                  ))
