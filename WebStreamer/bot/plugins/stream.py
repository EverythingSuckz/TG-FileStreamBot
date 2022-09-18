# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import logging
from pyrogram import filters
from WebStreamer.vars import Var
from WebStreamer.bot import StreamBot
from WebStreamer.utils.file_properties import getNew, fileId

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
    
    if message.photo:
        await message.reply(
            text="Don't send me photos, send them as document.",
            quote=True
        )
        
    try:
        await message.reply(
            text=f"{Var.URL}{getNew(fileId(message))[0]}",
            quote=True
        )
        await bot.copy_message(
            chat_id=Var.BIN_CHANNEL,
            from_chat_id=message.chat.id,
            message_id=message.id,
            caption='Hello'
        )
    except Exception as e:
        await message.reply(
            text=e,
            quote=True
        )
