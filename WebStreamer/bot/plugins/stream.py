# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

from pyrogram import filters, errors
from pyrogram.enums.parse_mode import ParseMode
from pyrogram.types import Message, InlineKeyboardMarkup, InlineKeyboardButton
from WebStreamer.vars import Var
from WebStreamer.bot import StreamBot, logger
from WebStreamer.utils import get_hash, get_mimetype


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
    mimetype = get_mimetype(log_msg)
    stream_link = f"{Var.URL}{log_msg.id}?hash={file_hash}"
    logger.info("Generated link: %s for %s", stream_link, m.from_user.first_name)
    markup = [InlineKeyboardButton("Download", url=stream_link+"&d=true")]
    if set(mimetype.split("/")) & {"video","audio","pdf"}:
        markup.append(InlineKeyboardButton("Stream", url=stream_link))
    try:
        await m.reply_text(
            text=f"<code>{stream_link}</code>",
            quote=True,
            parse_mode=ParseMode.HTML,
            reply_markup=InlineKeyboardMarkup([markup]),
        )
    except errors.ButtonUrlInvalid:
        await m.reply_text(
            text=f"<code>{stream_link}</code>)",
            quote=True,
            parse_mode=ParseMode.HTML,
        )
