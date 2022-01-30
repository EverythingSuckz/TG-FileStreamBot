def get_unique_id(media_msg) -> str:
    for attr in ("document", "video", "audio", "animation"):
        try:
            attrd = getattr(media_msg, attr)["file_unique_id"]
        except (AttributeError, TypeError):
            continue
    return attrd[:6]
