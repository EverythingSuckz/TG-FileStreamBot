# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import logging
from pyrogram import filters
from WebStreamer.vars import Var
from WebStreamer.bot import StreamBot
from WebStreamer.utils.file_properties import getFileid

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
async def getStreamlink(bot, message):
    try:
        await message.reply(
            text=f"{Var.URL}/{getFileid(message)}",
            quote=True
        )
    except Exception as e:
        await message.reply(e)
