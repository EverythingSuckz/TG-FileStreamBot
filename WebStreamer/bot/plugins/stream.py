# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import logging
from pyrogram import filters
from WebStreamer.vars import Var
from urllib.parse import quote_plus
from WebStreamer.bot import StreamBot
from pyrogram.types import Message, InlineKeyboardMarkup, InlineKeyboardButton

def detect_type(m: Message):
    if m.document:
        return m.document
    elif m.video:
        return m.video
    elif m.audio:
        return m.audio
    else:
        return
    

@StreamBot.on_message(filters.private & (filters.document | filters.video | filters.audio), group=4)
async def media_receive_handler(_, m: Message):
    file = detect_type(m)
    file_name = ''
    if file:
        file_name = file.file_name
    log_msg = await m.forward(chat_id=Var.BIN_CHANNEL)
    stream_link = f"{Var.URL}{log_msg.message_id}"
    if file_name:
        stream_link += f'/{quote_plus(file_name)}'
    logging.info(f"Generated link: {stream_link} for {m.from_user.first_name}")
    await m.reply_text(
        text = "<code>{}</code>".format(stream_link),
        quote = True,
        reply_markup = InlineKeyboardMarkup([[InlineKeyboardButton('Open', url=stream_link)]])
    )