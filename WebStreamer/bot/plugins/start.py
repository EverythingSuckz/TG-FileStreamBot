# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

from pyrogram import filters
from WebStreamer.bot import StreamBot

@StreamBot.on_message(filters.command(["start", "help"]))
async def start(bot, message):
    await message.reply(f"Hello, {message.from_user.mention(style="md")}, Send me a file to get stream link.")
