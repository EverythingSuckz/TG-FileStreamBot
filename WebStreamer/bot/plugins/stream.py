# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import logging
from pyrogram import filters
from WebStreamer.vars import Var
from urllib.parse import quote_plus
from WebStreamer.bot import StreamBot
from pyrogram.types import Message, InlineKeyboardMarkup, InlineKeyboardButton
from WebStreamer.utils.file_id import get_unique_id


def detect_type(m: Message):
    if m.document:
        return m.document
    elif m.video:
        return m.video
    elif m.audio:
        return m.audio
    else:
        return


@StreamBot.on_message(
    filters.private
    & (
        filters.document
        | filters.video
        | filters.audio
        | filters.animation
        | filters.voice
        | filters.video_note
    ),
    group=4,
)
async def media_receive_handler(_, m: Message):
    file = detect_type(m)
    file_name = ""
    if file:
        file_name = file.file_name
    log_msg = await m.forward(chat_id=Var.BIN_CHANNEL)
    stream_link = f"{Var.URL}{log_msg.message_id}/{quote_plus(file_name)}?hash={get_unique_id(log_msg)}"
    short_link = f"{Var.URL}{get_unique_id(log_msg)}{log_msg.message_id}"
    logging.info(f"Generated link: {stream_link} for {m.from_user.first_name}")
    await m.reply_text(
        text="<code>{}</code> (<a href='{}'>shortened</a>)".format(
            stream_link, short_link
        ),
        quote=True,
        parse_mode="html",
        reply_markup=InlineKeyboardMarkup(
            [[InlineKeyboardButton("Open", url=stream_link)]]
        ),
    )
