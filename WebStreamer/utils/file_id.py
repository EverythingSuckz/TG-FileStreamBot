def get_unique_id(media_msg) -> str:
    for attr in (
        "audio",
        "document",
        "photo",
        "sticker",
        "animation",
        "video",
        "voice",
        "video_note",
    ):
        try:
            attrd = getattr(media_msg, attr)["file_unique_id"]
        except (AttributeError, TypeError):
            continue
    return attrd[:6]
