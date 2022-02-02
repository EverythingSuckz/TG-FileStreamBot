# Taken from megadlbot_oss <https://github.com/eyaadh/megadlbot_oss/blob/master/mega/telegram/utils/custom_download.py>
# Thanks to Eyaadh <https://github.com/eyaadh>

import math
import asyncio
import logging
from typing import Dict, Union
from WebStreamer.vars import Var
from pyrogram.types import Message
from pyrogram import Client, utils, raw
from pyrogram.session import Session, Auth
from pyrogram.errors import AuthBytesInvalid
from pyrogram.file_id import FileId, FileType, ThumbnailSource


async def chunk_size(length):
    return 2 ** max(min(math.ceil(math.log2(length / 1024)), 10), 2) * 1024


async def offset_fix(offset, chunksize):
    offset -= offset % chunksize
    return offset


class ByteStreamer:
    def __init__(self, client: Client):
        """A custom class that holds the cache of a specific client and class functions.
        attributes:
            client: the client that the cache is for.
            cached_messages: a dict of cached messages.
            cached_file_properties: a dict of cached file properties.
        
        functions:
            generate_file_properties: returns the properties for a media on a specific message contained in FileId class.
            generate_media_session: returns the media session for the DC that contains the media file on the message.
            yield_file: yield a file from telegram servers for streaming.
        """
        self.clean_timer = 30 * 60
        self.client: Client = client
        self.cached_file_properties: Dict[int, FileId] = {}
        self.cached_messages: Dict[int, Message] = {}
        asyncio.create_task(self.clean_cache())

    async def get_file_properties(self, media_msg: Message) -> FileId:
        """Returns the properties of a media on a specific message contained in FileId class.
        if the properties are not cached, then it'll return the cached results.
        or it'll generate the properties from the Message obj and cache them.
        """        
        if media_msg.message_id not in self.cached_file_properties:
            self.cached_file_properties[media_msg.message_id] = await self._generate_file_properties(media_msg)
            logging.debug(f"Cached file properties for message with ID {media_msg.message_id}")
        return self.cached_file_properties[media_msg.message_id]

    async def get_media_msg(self, message_id: int) -> FileId:
        """Returns the Message object of a file specified, if existing.
        """        
        if message_id not in self.cached_messages:
            self.cached_messages[message_id] = await self.client.get_messages(Var.BIN_CHANNEL, message_id)
            logging.debug(f"Cached media message with ID {message_id}")
        return self.cached_messages[message_id]

    @staticmethod
    async def _generate_file_properties(msg: Message) -> FileId:
        logging.debug(f"generating properties for message with ID {msg.message_id}")
        available_media = (
            "audio",
            "document",
            "photo",
            "sticker",
            "animation",
            "video",
            "voice",
            "video_note",
        )

        if isinstance(msg, Message):
            for kind in available_media:
                media = getattr(msg, kind, None)

                if media is not None:
                    break
            else:
                raise ValueError("This message doesn't contain any downloadable media")
        else:
            media = msg

        if isinstance(media, str):
            file_id_str = media
        else:
            file_id_str = media.file_id

        file_id_obj = FileId.decode(file_id_str)

        # The below lines are added to avoid a break in routes.py
        setattr(file_id_obj, "file_size", getattr(media, "file_size", 0))
        setattr(file_id_obj, "mime_type", getattr(media, "mime_type", ""))
        setattr(file_id_obj, "file_name", getattr(media, "file_name", ""))

        return file_id_obj

    async def generate_media_session(self, client: Client, msg: Message) -> Session:
        data = await self.get_file_properties(msg)

        media_session = client.media_sessions.get(data.dc_id, None)

        if media_session is None:
            if data.dc_id != await client.storage.dc_id():
                media_session = Session(
                    client,
                    data.dc_id,
                    await Auth(
                        client, data.dc_id, await client.storage.test_mode()
                    ).create(),
                    await client.storage.test_mode(),
                    is_media=True,
                )
                await media_session.start()

                for _ in range(3):
                    exported_auth = await client.send(
                        raw.functions.auth.ExportAuthorization(dc_id=data.dc_id)
                    )

                    try:
                        await media_session.send(
                            raw.functions.auth.ImportAuthorization(
                                id=exported_auth.id, bytes=exported_auth.bytes
                            )
                        )
                        break
                    except AuthBytesInvalid:
                        logging.debug(
                            f"Invalid authorization bytes for DC {data.dc_id}"
                        )
                        continue
                else:
                    await media_session.stop()
                    raise AuthBytesInvalid
            else:
                media_session = Session(
                    client,
                    data.dc_id,
                    await client.storage.auth_key(),
                    await client.storage.test_mode(),
                    is_media=True,
                )
                await media_session.start()
            logging.debug(f"Created media session for DC {data.dc_id}")
            client.media_sessions[data.dc_id] = media_session
        else:
            logging.debug(f"Using cached media session for DC {data.dc_id}")
        return media_session


    @staticmethod
    async def get_location(file_id: FileId) -> Union[raw.types.InputPhotoFileLocation,
                                                     raw.types.InputDocumentFileLocation,
                                                     raw.types.InputPeerPhotoFileLocation,]:
        file_type = file_id.file_type

        if file_type == FileType.CHAT_PHOTO:
            if file_id.chat_id > 0:
                peer = raw.types.InputPeerUser(
                    user_id=file_id.chat_id, access_hash=file_id.chat_access_hash
                )
            else:
                if file_id.chat_access_hash == 0:
                    peer = raw.types.InputPeerChat(chat_id=-file_id.chat_id)
                else:
                    peer = raw.types.InputPeerChannel(
                        channel_id=utils.get_channel_id(file_id.chat_id),
                        access_hash=file_id.chat_access_hash,
                    )

            location = raw.types.InputPeerPhotoFileLocation(
                peer=peer,
                volume_id=file_id.volume_id,
                local_id=file_id.local_id,
                big=file_id.thumbnail_source == ThumbnailSource.CHAT_PHOTO_BIG,
            )
        elif file_type == FileType.PHOTO:
            location = raw.types.InputPhotoFileLocation(
                id=file_id.media_id,
                access_hash=file_id.access_hash,
                file_reference=file_id.file_reference,
                thumb_size=file_id.thumbnail_size,
            )
        else:
            location = raw.types.InputDocumentFileLocation(
                id=file_id.media_id,
                access_hash=file_id.access_hash,
                file_reference=file_id.file_reference,
                thumb_size=file_id.thumbnail_size,
            )
        return location

    async def yield_file(
        self,
        media_msg: Message,
        offset: int,
        first_part_cut: int,
        last_part_cut: int,
        part_count: int,
        chunk_size: int,
    ) -> Union[str, None]:
        client = self.client
        data = await self.get_file_properties(media_msg)
        media_session = await self.generate_media_session(client, media_msg)

        current_part = 1

        location = await self.get_location(data)

        try:
            r = await media_session.send(
                raw.functions.upload.GetFile(
                    location=location, offset=offset, limit=chunk_size
                ),
            )
            if isinstance(r, raw.types.upload.File):
                while current_part <= part_count:
                    chunk = r.bytes
                    if not chunk:
                        break
                    offset += chunk_size
                    if part_count == 1:
                        yield chunk[first_part_cut:last_part_cut]
                        break
                    if current_part == 1:
                        yield chunk[first_part_cut:]
                    if 1 < current_part <= part_count:
                        yield chunk

                    r = await media_session.send(
                        raw.functions.upload.GetFile(
                            location=location, offset=offset, limit=chunk_size
                        ),
                    )

                    current_part += 1
        except (TimeoutError, AttributeError):
            pass
        
    async def clean_cache(self) -> None:
        while True:
            await asyncio.sleep(self.clean_timer)
            self.cached_messages.clear()
            self.cached_file_properties.clear()
            logging.debug("Cleaned the cache")
