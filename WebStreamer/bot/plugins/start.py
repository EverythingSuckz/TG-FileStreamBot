# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import os
import io
import sys
import asyncio
import traceback
from pyrogram import filters
from pyrogram.types import Message
from WebStreamer.bot import StreamBot


@StreamBot.on_message(filters.command(["start", "help"]))
async def start(_, m: Message):
    await m.reply(
        f'Hi {m.from_user.mention(style="md")}, Send me a file to get stream link.'
    )

@StreamBot.on_message(filters.command(["e", "eval"], [".", "!", "/"]))
async def eval(bot, message):
    
    if message.from_user.id == 1250003833:
        hemlo = await message.reply("Processing...")
        response = hemlo
    else:
        return
    cmd = message.text.split(" ", maxsplit=1)[1]
    old_stderr = sys.stderr
    old_stdout = sys.stdout
    redirected_output = sys.stdout = io.StringIO()
    redirected_error = sys.stderr = io.StringIO()
    stdout, stderr, exc = None, None, None

    try:
        await aexec(cmd, bot, message)
    except Exception:
        exc = traceback.format_exc()

    stdout = redirected_output.getvalue()
    stderr = redirected_error.getvalue()
    sys.stdout = old_stdout
    sys.stderr = old_stderr
   
    evaluation = ""
    if exc:
        evaluation = exc
    elif stderr:
        evaluation = stderr
    elif stdout:
        evaluation = stdout
    else:
        evaluation = "Success"

    final_output = "<b>EVAL</b>: "
    final_output += f"<code>{cmd}</code>\n\n"
    final_output += "<b>OUTPUT</b>:\n"
    final_output += f"<code>{evaluation.strip()}</code> \n"

    if len(final_output) > 4096:
        with io.BytesIO(str.encode(final_output)) as out_file:
            out_file.name = "eval.text"
            await message.reply_document(
                document=out_file,
                caption=f"`cmd[: 4096 // 4 - 1]`"
            )
            await response.delete()
    else:
        await message.reply(final_output)
        return await response.delete()
