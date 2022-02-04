# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import logging
from pyrogram import filters
from WebStreamer.vars import Var
from urllib.parse import quote_plus
from WebStreamer.bot import StreamBot
from WebStreamer.utils import get_hash, get_name
from pyrogram.types import Message, InlineKeyboardMarkup, InlineKeyboardButton




@StreamBot.on_message(
    filters.private
    & (
        filters.document
        | filters.video
        | filters.audio
        | filters.animation
        | filters.voice
        | filters.video_note
        | filters.photo
        | filters.sticker
    ),
    group=4,
)
async def media_receive_handler(_, m: Message):
    log_msg = await m.forward(chat_id=Var.BIN_CHANNEL)
    stream_link = f"{Var.URL}{log_msg.message_id}/{quote_plus(get_name(m))}?hash={get_hash(log_msg)}"
    short_link = f"{Var.URL}{get_hash(log_msg)}{log_msg.message_id}"
    logging.info(f"Generated link: {stream_link} for {m.from_user.first_name}")
    await m.reply_text(
        text="<code>{}</code>\n(<a href='{}'>shortened</a>)".format(
            stream_link, short_link
        ),
        quote=True,
        parse_mode="html",
        reply_markup=InlineKeyboardMarkup(
            [[InlineKeyboardButton("Open", url=stream_link)]]
        ),
    )
