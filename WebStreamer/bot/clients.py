# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

from ..vars import Var
from pyrogram import Client
from . import multi_clients, work_loads, StreamBot
from WebStreamer.utils.config_parser import TokenParser


async def initialize_clients():
    multi_clients[0] = StreamBot
    work_loads[0] = 0
    if Var.MULTI_CLIENT:
        all_tokens = TokenParser().parse_from_env()
        for client_id, token in all_tokens.items():
            instance = Client(
                session_name=":memory:",
                api_id=Var.API_ID,
                api_hash=Var.API_HASH,
                bot_token=token,
                sleep_threshold=Var.SLEEP_THRESHOLD,
                no_updates=True,
            )
            try:
                multi_clients[client_id] = await instance.start()
            except Exception as e:
                print(f"Failed starting Client - {client_id}; Error: {e}")
                continue
            work_loads[client_id] = 0
            print(f"Started - Client {client_id}")
