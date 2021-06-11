from WebStreamer.vars import Var
from WebStreamer.bot import StreamBot
from WebStreamer.utils.custom_dl import TGCustomYield
import urllib.parse
import secrets
import mimetypes
import aiofiles
import logging

async def fetch_properties(message_id):
    media_msg = await StreamBot.get_messages(Var.BIN_CHANNEL, message_id)
    file_properties = await TGCustomYield().generate_file_properties(media_msg)
    file_name = file_properties.file_name if file_properties.file_name \
        else f"{secrets.token_hex(2)}.jpeg"
    mime_type = file_properties.mime_type if file_properties.mime_type \
        else f"{mimetypes.guess_type(file_name)}"
    return file_name, mime_type


async def render_page(message_id):
    file_name, mime_type = await fetch_properties(message_id)
    src = urllib.parse.urljoin(Var.URL, str(message_id))
    async with aiofiles.open('WebStreamer/template/req.html') as r:
        audio_formats = ['audio/mpeg', 'audio/mp4', 'audio/x-mpegurl', 'audio/vnd.wav']
        video_formats = ['video/mp4', 'video/avi', 'video/ogg', 'video/h264', 'video/h265', 'video/x-matroska']
        if mime_type.lower() in video_formats:
            heading = 'Watch {}'.format(file_name)
        elif mime_type.lower() in audio_formats:
            heading = 'Listen {}'.format(file_name)
        else:
            return None
        tag = mime_type.split('/')[0].strip()
        html = (await r.read()).replace('tag', tag) % (heading, file_name, src)
        logging.info(html)
        return html