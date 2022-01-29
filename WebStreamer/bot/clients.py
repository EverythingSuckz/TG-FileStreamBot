# This file is a part of TG-FileStreamBot
# Coding : Jyothis Jayanth [@EverythingSuckz]

import logging
from ..vars import Var
from pyrogram import Client
from . import multi_clients, work_loads, StreamBot
from WebStreamer.utils.config_parser import TokenParser

if Var.MULTI_CLIENT:
    all_tokens = TokenParser().parse_from_env()
    multi_clients[0] = StreamBot
    work_loads[0] = 0
    for client_id, token in all_tokens.items():
        multi_clients.update(client_id, Client(
            session_name= ':memory:',
            api_id=Var.API_ID,
            api_hash=Var.API_HASH,
            bot_token=token,
            sleep_threshold=Var.SLEEP_THRESHOLD,
            workers=Var.WORKERS
        ).start())
        work_loads[client_id] = 0
        print(f"Started - Client {client_id}")
