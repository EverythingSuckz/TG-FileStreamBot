from pyrogram import Client
from pyrogram.types import Message
from typing import Any, Optional, Tuple
from pyrogram.raw.types.messages import Messages
from pyrogram.file_id import FileId, FileUniqueId
from WebStreamer.server.exceptions import FIleNotFound


async def parse_file_id(message: "Message") -> Optional[FileId]:
    media = get_media_from_message(message)
    if media:
        return FileId.decode(media.file_id)

async def parse_file_unique_id(message: "Messages") -> Optional[FileUniqueId]:
    media = get_media_from_message(message)
    if media:
        return FileUniqueId.decode(media.file_unique_id)

async def get_file_ids(client: Client, chat_id: int, message_id: int) -> Tuple[Optional[FileId], Optional[FileUniqueId]]:
    message = await client.get_messages(chat_id, message_id)
    if message.empty:
        raise FIleNotFound
    media = get_media_from_message(message)
    file_unique_id = await parse_file_unique_id(message)
    file_id = await parse_file_id(message)
    setattr(file_id, "file_size", getattr(media, "file_size", 0))
    setattr(file_id, "mime_type", getattr(media, "mime_type", ""))
    setattr(file_id, "file_name", getattr(media, "file_name", ""))
    return (file_id, file_unique_id)

def get_media_from_message(message: "Message") -> Any:
    media_types = (
        "audio",
        "document",
        "photo",
        "sticker",
        "animation",
        "video",
        "voice",
        "video_note",
    )
    for attr in media_types:
        media = getattr(message, attr, None)
        if media:
            return media


def get_hash(media_msg: Message) -> str:
    media = get_media_from_message(media_msg)
    return getattr(media, "file_unique_id", "")[:6]

def get_name(media_msg: Message) -> str:
    media = get_media_from_message(media_msg)
    return getattr(media, "file_name", "")