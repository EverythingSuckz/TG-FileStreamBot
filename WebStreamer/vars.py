# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

from os import environ
from dotenv import load_dotenv

load_dotenv()

class Var(object):
    MULTI_CLIENT = False 
    API_ID = int(environ.get("API_ID")) # Get it from my.telegram.org.
    API_HASH = str(environ.get("API_HASH")) # Get it from my.telegram.org.
    BOT_TOKEN = str(environ.get("BOT_TOKEN")) # Get it from @BotFather.
    SLEEP_THRESHOLD = int(environ.get("SLEEP_THRESHOLD", "60"))  # Default value is 1 minute.
    WORKERS = int(environ.get("WORKERS", "6"))  # Default value is 6, 6 commands at once.
    BIN_CHANNEL = int(environ.get("BIN_CHANNEL", None))  # A channel where bot will share files for stream.
    PORT = int(environ.get("PORT", 8080)) 
    BIND_ADDRESS = str(environ.get("WEB_SERVER_BIND_ADDRESS", "0.0.0.0"))
    PING_INTERVAL = int(environ.get("PING_INTERVAL", "1200"))  # Default value is 20 minutes.
    HAS_SSL = environ.get("HAS_SSL", False)
    HAS_SSL = True if str(HAS_SSL).lower() == "true" else False
    NO_PORT = environ.get("NO_PORT", False)
    NO_PORT = True if str(NO_PORT).lower() == "true" else False
    CUSTOM_CAPTION = environ.get("CUSTOM_CAPTION", "Www.Hagadmansa.Com") # Set a custom caption in the starting of file name.
    REDIRECT_TO = environ.get("REDIRECT_TO", "https://hagadmansa.com") # Redirect to your own website or channel.
    if "DYNO" in environ:
        ON_HEROKU = True
        APP_NAME = str(environ.get("APP_NAME"))
    else:
        ON_HEROKU = False
    FQDN = (
        str(environ.get("FQDN", BIND_ADDRESS))
        if not ON_HEROKU or environ.get("FQDN")
        else APP_NAME + ".herokuapp.com"
    )
    if ON_HEROKU:
        URL = f"https://{FQDN}/"
    else:
        URL = "http{}://{}{}/".format("s" if HAS_SSL else "", FQDN, "" if NO_PORT else ":" + str(PORT))
