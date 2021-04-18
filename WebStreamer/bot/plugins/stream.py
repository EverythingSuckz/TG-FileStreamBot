# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

from WebStreamer.bot import StreamBot
from WebStreamer.vars import Var
from pyrogram import filters, Client, emoji
from pyrogram.types import Message, InlineKeyboardMarkup, InlineKeyboardButton


@StreamBot.on_message(filters.private & (filters.document | filters.video | filters.audio), group=4)
async def media_receive_handler(c: Client, m: Message):
    log_msg = await m.copy(chat_id=Var.BIN_CHANNEL)
    stream_link = "https://{}/{}".format(Var.FQDN, log_msg.message_id) if Var.ON_HEROKU or Var.NO_PORT else \
        "http://{}:{}/{}".format(Var.FQDN,
                                Var.PORT,
                                log_msg.message_id)
    await m.reply_text(
        text="`{}`".format(stream_link),
        quote=True
    )