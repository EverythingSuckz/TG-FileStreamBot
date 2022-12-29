[![StandWithUkraineBanner](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/banner2-direct.svg)](https://supportukrainenow.org/)
<h1 align="center">Telegram File Stream Bot</h1>
<p align="center">
  <a href="https://github.com/EverythingSuckz/TG-FileStreamBot">
    <img src="https://socialify.git.ci/EverythingSuckz/TG-FileStreamBot/image?description=1&font=Source%20Code%20Pro&forks=1&issues=1&logo=https://telegra.ph/file/01385a9f4cf0419682b87.png&pattern=Circuit%20Board&pulls=1&stargazers=1&theme=Dark" alt="Cover Image" width="650">
  </a>
  <p align="center">
    A Telegram bot to stream files to web
    <br />
    <a href="https://telegram.dog/TG_FileStreamBot"><strong>Demo Bot »</strong></a>
    <br />
    <a href="https://github.com/EverythingSuckz/TG-FileStreamBot/issues">Report a Bug</a>
    |
    <a href="https://github.com/EverythingSuckz/TG-FileStreamBot/issues">Request Feature</a>
  </p>
</p>

<hr>

<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-this-bot">About this Bot</a>
      <ul>
        <li><a href="#original-repository">Original Repository</a></li>
      </ul>
    </li>
    <li>
      <a href="#how-to-make-your-own">How to make your own</a>
      <ul>
        <li><a href="#host-it-on-vps-or-locally">Run it in a VPS / local</a></li>
        <li><a href="#deploy-using-docker">Deploy using Docker</a></li>
      </ul>
    </li>
    <li><a href="#setting-up-things">Setting up things</a></li>
    <ul>
      <li><a href="#mandatory-vars">Mandatory Vars</a></li>
      <li><a href="#optional-vars">Optional Vars</a></li>
    </ul>
    <li><a href="#how-to-use-the-bot">How to use the bot</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#contact-me">Contact me</a></li>
    <li><a href="#credits">Credits</a></li>
  </ol>
</details>

## About This Bot

<p align="center">
    <a herf="https://github.com/EverythingSuckz/TG-FileStreamBot">
        <img src="https://telegra.ph/file/a8bb3f6b334ad1200ddb4.png" height="100" width="100" alt="Telegram Logo">
    </a>
</p>
<p align='center'>
    This bot will give you stream links for Telegram files without the need of waiting till the download completes
</p>

### Original Repository

The main working part was taken from [Megatron](https://github.com/eyaadh/megadlbot_oss) and thanks to [eyaadh](https://github.com/eyaadh) for his awesome project.

## How to make your own

Either you could locally host or deploy on ~~[Heroku](https://heroku.com)~~ Free tier is dead.

### Host it on VPS or Locally

```sh
git clone https://github.com/EverythingSuckz/TG-FileStreamBot
cd TG-FileStreamBot
python3 -m venv ./venv
. ./venv/bin/activate
pip3 install -r requirements.txt
python3 -m WebStreamer
```

and to stop the whole bot,
 do <kbd>CTRL</kbd>+<kbd>C</kbd>

> **If you wanna run this bot 24/7 on the VPS, follow thesesteps.**
> ```sh
> sudo apt install tmux -y
> tmux
> python3 -m WebStreamer
> ```
> now you can close the VPS and the bot will run on it.

### Deploy using Docker
First clone the repository
```sh
git clone https://github.com/EverythingSuckz/TG-FileStreamBot
cd TG-FileStreamBot
```
then build the docker image
```sh
docker build . -t stream-bot
```
now create the `.env` file with your variables. and start your container:
```sh
docker run -d --restart unless-stopped --name fsb \
-v /PATH/TO/.env:/app/.env \
-p 8000:8000 \
stream-bot
```

your `PORT` variable has to be consistent with the container's exposed port since it's used for URL generation. so remember if you changed the `PORT` variable your docker run command changes too. (example: `PORT=9000` -> `-p 9000:9000`)

if you need to change the variables in `.env` file after your bot was already started, all you need to do is restart the container for the bot settings to get updated:
```sh
docker restart fsb
```

### Deploy using docker-compose
First install docker-compose. For debian based, run 
```sh
sudo apt install docker-compose -y
```
Afterwards, clone the repository
```sh
git clone https://github.com/EverythingSuckz/TG-FileStreamBot
cd TG-FileStreamBot
```
No need to create .env file, just edit the variables in the docker-compose.yml

Now run the compose file
```sh
sudo docker compose up -d
```

## Setting up things

If you're locally hosting, create a file named `.env` in the root directory and add all the variables there.
An example of `.env` file:

```sh
API_ID=452525
API_HASH=esx576f8738x883f3sfzx83
BOT_TOKEN=55838383:yourtbottokenhere
MULTI_TOKEN1=55838383:yourfirstmulticlientbottokenhere
MULTI_TOKEN2=55838383:yoursecondmulticlientbottokenhere
MULTI_TOKEN3=55838383:yourthirdmulticlientbottokenhere
BIN_CHANNEL=-100
PORT=8080
FQDN=yourserverip
HAS_SSL=False
```

### Mandatory Vars

`API_ID` : Goto [my.telegram.org](https://my.telegram.org) to obtain this.

`API_HASH` : Goto [my.telegram.org](https://my.telegram.org) to obtain this.

`BOT_TOKEN` : Get the bot token from [@BotFather](https://telegram.dog/BotFather)

`BIN_CHANNEL` : Create a new channel (private/public), post something in your channel. Forward that post to [@missrose_bot](https://telegram.dog/MissRose_bot) and **reply** `/id`. Now copy paste the forwarded channel ID in this field. 

### For making use of Multi-Client support

> **What it does?** <br>
> Shares the workload between other bots to avoid getting floodwaited and to make the server handle more requests.
`MULTI_TOKEN1`: Add your first bot token here.

`MULTI_TOKEN2`: Add your second bot token here.

you may also add as many as bots you want. (max limit is not tested yet)
`MULTI_TOKEN3`, `MULTI_TOKEN4`, etc.

> Don't forget to add all these bots to the `BIN_CHANNEL`

### Optional Vars

-`HASH_LENGTH` : Set custom hash length for generated urls
> **NOTE**: Hash length should be greater than 5 and less than 64.


- `SLEEP_THRESHOLD` : Set a sleep threshold for flood wait exceptions happening globally in this telegram bot instance, below which any request that raises a flood wait will be automatically invoked again after sleeping for the required amount of time. Flood wait exceptions requiring higher waiting times will be raised. Defaults to 60 seconds.


- `WORKERS` : Number of maximum concurrent workers for handling incoming updates.
> Defaults to `3`


- `PORT` : The port that you want your webapp to be listened to.
> Defaults to `8080`


- `WEB_SERVER_BIND_ADDRESS` : Your server bind address.
> Defaults to `0.0.0.0`

- `NO_PORT` : (can be either `True` or `False`) If you don't want your port to be displayed.
> You should point your `PORT` to `80` (http) or `443` (https) for the links to work.

- `FQDN` :  A Fully Qualified Domain Name if present.
> Defaults to `WEB_SERVER_BIND_ADDRESS`

- `HAS_SSL` : (can be either `True` or `False`) If you want the generated links in https format.

- `KEEP_ALIVE`: If you want to make the server ping itself every `PING_INTERVAL` seconds to avoid sleeping. Helpful in PaaS Free tiers. 
> Defaults to `False`

- `PING_INTERVAL` : The time in ms you want the servers to be pinged each time to avoid sleeping (If you're on some PaaS). 
> Defaults to `1200` or 20 minutes.

- `USE_SESSION_FILE` : Use session files for client(s) rather than storing the pyrogram sqlite db in the memory

## How to use the bot

> :warning: **Before using the  bot, don't forget to add all the bots (multi-client ones too) to the `BIN_CHANNEL` as an admin**
 
`/start` : To check if the bot is alive or not.

To get an instant stream link, just forward any media to the bot and boom, its fast af.

## faQ

- How long the links will remain valid or is there any expiration time for the links generated by the bot?
> The links will will be valid as longs as your bot is alive and you haven't deleted the log channel.

## Contributing

Feel free to contribute to this project if you have any further ideas

## Contact me

[![Telegram Channel](https://img.shields.io/static/v1?label=Join&message=Telegram%20Channel&color=blueviolet&style=for-the-badge&logo=telegram&logoColor=violet)](https://xn--r1a.click/WhySooSerious)
[![Telegram Group](https://img.shields.io/static/v1?label=Join&message=Telegram%20Group&color=blueviolet&style=for-the-badge&logo=telegram&logoColor=violet)](https://xn--r1a.click/WhyThisUsername)

You can contact either via my [Telegram Group](https://xn--r1a.click/WhyThisUsername) ~~or you can PM me on [@EverythingSuckz](https://xn--r1a.click/EverythingSuckz)~~


## Credits

- Me
- [eyaadh](https://github.com/eyaadh) for his awesome [Megatron Bot](https://github.com/eyaadh/megadlbot_oss).
- [BlackStone](https://github.com/eyMarv) for adding multi-client support.
- [Dan Tès](https://telegram.dog/haskell) for his [Pyrogram Library](https://github.com/pyrogram/pyrogram)
- [TheHamkerCat](https://github.com/TheHamkerCat)

## Copyright

Copyright (C) 2022 [EverythingSuckz](https://github.com/EverythingSuckz) under [GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.en.html).

TG-FileStreamBot is Free Software: You can use, study share and improve it at your
will. Specifically you can redistribute and/or modify it under the terms of the
[GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.en.html) as
published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version. 
