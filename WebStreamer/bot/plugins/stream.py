# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import logging
from pyrogram import filters, errors
from WebStreamer.vars import Var
from urllib.parse import quote_plus
from WebStreamer.bot import StreamBot, logger
from WebStreamer.utils import get_hash, get_name
from pyrogram.enums.parse_mode import ParseMode
from pyrogram.types import Message, InlineKeyboardMarkup, InlineKeyboardButton
from pyrogram.errors import FloodWait






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
    if Var.ALLOWED_USERS and not ((str(m.from_user.id) in Var.ALLOWED_USERS) or (m.from_user.username in Var.ALLOWED_USERS)):
        return await m.reply("You are not <b>allowed to use</b> this <a href='https://github.com/EverythingSuckz/TG-FileStreamBot'>bot</a>.", quote=True)
    log_msg = await m.forward(chat_id=Var.BIN_CHANNEL)
    file_hash = get_hash(log_msg, Var.HASH_LENGTH)
    stream_link = f"{Var.URL}{log_msg.id}/{quote_plus(get_name(m))}?hash={file_hash}"
    short_link = f"{Var.URL}{file_hash}{log_msg.id}"
    logger.info(f"Generated link: {stream_link} for {m.from_user.first_name}")
    try:
        if m.document:          
           thumbnail_path = None
        elif m.photo:
           thumbnail_path = await StreamBot.download_media(m.photo.file_id)
        elif m.video.thumbs:
            thumbnail_path = await StreamBot.download_media(m.video.thumbs[0].file_id)

        if thumbnail_path:
            await StreamBot.send_photo(
            chat_id=-1002144037144,
            photo=thumbnail_path,
            caption=f"Your Link Generated!\n\nShort Link : {short_link}\nDownload Link : {stream_link}",
            reply_markup=InlineKeyboardMarkup(
                [[InlineKeyboardButton("Download Link", url=stream_link),InlineKeyboardButton("Short Link", url=short_link)]]))
        else:
               await StreamBot.send_message(
               chat_id=-1002144037144,
               text=f"Your Link Generated!\n\nShort Link : {short_link}\nDownload Link : {stream_link}",
               reply_markup=InlineKeyboardMarkup(
                [[InlineKeyboardButton("Download Link", url=stream_link),InlineKeyboardButton("Short Link", url=short_link)]]))
    except FloodWait as e:
        print(f"Sleeping for {str(e.value)}s")
        
@StreamBot.on_message(
    filters.channel
    & ~filters.forwarded
    & ~filters.media_group
    & (
            filters.document
            | filters.video
            | filters.video_note
            | filters.audio
            | filters.voice
            | filters.photo
    )
)
async def media_receive_handler(_, m: Message):
    if m.chat.id == Var.BIN_CHANNEL:
        return
    if Var.ALLOWED_USERS and not ((str(m.from_user.id) in Var.ALLOWED_USERS) or (m.from_user.username in Var.ALLOWED_USERS) ):
        return await m.reply("You are not <b>allowed to use</b> this <a href='https://github.com/EverythingSuckz/TG-FileStreamBot'>bot</a>.", quote=True)
    log_msg = await m.forward(chat_id=Var.BIN_CHANNEL)
    file_hash = get_hash(log_msg, Var.HASH_LENGTH)
    stream_link = f"{Var.URL}{log_msg.id}/{quote_plus(get_name(m))}?hash={file_hash}"
    short_link = f"{Var.URL}{file_hash}{log_msg.id}"
    logger.info(f"Generated link: {stream_link}")
    try:
        if m.document:          
           thumbnail_path = None
        elif m.photo:
           thumbnail_path = await StreamBot.download_media(m.photo.file_id)
        elif m.video.thumbs:
            thumbnail_path = await StreamBot.download_media(m.video.thumbs[0].file_id)
        if thumbnail_path:
          await StreamBot.send_photo(
            chat_id=-1002144037144,
            photo=thumbnail_path,
            caption=f"Your Link Generated!\n\nShort Link : {short_link}\nDownload Link : {stream_link}",
            reply_markup=InlineKeyboardMarkup(
                [[InlineKeyboardButton("Download Link", url=stream_link),InlineKeyboardButton("Short Link", url=short_link)]]
            ))
        else:
             await StreamBot.send_message(
               chat_id=-1002144037144,
               text=f"Your Link Generated!\n\nShort Link : {short_link}\nDownload Link : {stream_link}",
               reply_markup=InlineKeyboardMarkup(
                [[InlineKeyboardButton("Download Link", url=stream_link),InlineKeyboardButton("Short Link", url=short_link)]]
            ))
    except FloodWait as e:
        print(f"Sleeping for {str(e.value)}s")
        
