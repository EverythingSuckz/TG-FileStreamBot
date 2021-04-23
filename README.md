<h1 align="center">Telegram File Stream Bot</h1>
<p align="center">
  <a href="https://github.com/AbirHasan2005/Pyro-FileStreamBot">
    <img src="https://socialify.git.ci/AbirHasan2005/Pyro-FileStreamBot/image?description=1&font=Source%20Code%20Pro&forks=1&issues=1&logo=https%3A%2F%2Fwww.flaticon.com%2Fpremium-icon%2Ficons%2Fsvg%2F2321%2F2321089.svg&pattern=Circuit%20Board&pulls=1&stargazers=1&theme=Dark" alt="Cover Image" width="650">
  </a>
  <p align="center">
    A Telegram bot to stream files to web
    <br />
    <a href="https://telegram.dog/AH_File2Link_Bot"><strong>Demo Bot »</strong></a>
    <br />
    <a href="https://github.com/AbirHasan2005/Pyro-FileStreamBot/issues">Report a Bug</a>
    |
    <a href="https://github.com/AbirHasan2005/Pyro-FileStreamBot/issues">Request Feature</a>
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
        <li><a href="#deploy-on-heroku">Deploy using Heroku</a></li>
        <li><a href="#host-it-on-vps-or-locally">Run it in a VPS / local</a></li>
      </ul>
    </li>
    <li><a href="#setting-up-things">Setting up things</a></li>
    <ul>
      <li><a href="#mandatory-vars">Mandatory Vars</a></li>
      <li><a href="#optional-vars">Optional Vars</a></li>
    </ul>
    <li><a href="#how-to-use-the-bot">How to use the bot</a></li>
    <li><a href="#special-features">Special Features</a></li>
    <li><a href="#demo-bot>Demo Bot</a></li>
    <li><a href="#credits">Credits</a></li>
  </ol>
</details>

## About This Bot

<p align="center">
    <a href="https://github.com/AbirHasan2005/Pyro-FileStreamBot">
        <img src="https://www.flaticon.com/premium-icon/icons/svg/2626/2626281.svg" height="100" width="100" alt="Telegram Logo">
    </a>
</p>
<p align='center'>
    This bot will give you stream links for Telegram files without the need of waiting till the download completes
</p>

### Original Repository

The main working part was taken from [Megatron](https://github.com/eyaadh/megadlbot_oss) and thanks to [eyaadh](https://github.com/eyaadh) for his awesome project.

## How to make your own

Either you could locally host or deploy on [Heroku](https://heroku.com)

### Deploy on Heroku

Press the below button to Fast deploy on Heroku

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

then goto the <a href="#mandatory-vars">variables tab</a> for more info on setting up environmental variables.

### Host it on VPS or Locally

```sh
git clone https://github.com/AbirHasan2005/Pyro-FileStreamBot
cd TG-FileStreamBot
virtualenv -p /usr/bin/python3 venv
. ./venv/bin/activate
pip install -r requirements.txt
python3 -m WebStreamer
```

and to stop the whole bot,
 do <kbd>CTRL</kbd>+<kbd>C</kbd>

## Setting up things

If you're on Heroku, just add these in the Environmental Variables
or if you're Locally hosting, create a file named `.env` in the root directory and add all the variables there.
An example of `.env` file:

```sh
API_ID=452525
API_HASH=esx576f8738x883f3sfzx83
BOT_TOKEN=55838383:yourtbottokenhere
BIN_CHANNEL=-100
PORT=8080
FQDN=your_server_ip
OWNER_ID=your_user_id
DATABASE_URL=mongodb_uri
```

### Mandatory Vars

`API_ID` : Goto [my.telegram.org](https://my.telegram.org) to obtain this.

`API_HASH` : Goto [my.telegram.org](https://my.telegram.org) to obtain this.

`BOT_TOKEN` : Get the bot token from [@BotFather](https://telegram.dog/BotFather)

`BIN_CHANNEL` : Create a new channel (private/public), add [@missrose_bot](https://telegram.dog/MissRose_bot) as admin to the channel and type /id. Now copy paste the ID into this field.

`OWNER_ID` : Your Telegram User ID

`DATABASE_URL` : MongoDB URI for saving User IDs when they first Start the Bot. We will use that for Broadcasting to them. I will try to add more features related with Database. If you need help to get the URI you can ask in [Support Group](https://t.me/linux_repo).

### Optional Vars

`UPDATES_CHANNEL` : Put a Public Channel Username, so every user have to Join that channel to use the bot. Must add bot to channel as Admin to work properly.

`BANNED_CHANNELS` : Put IDs of Banned Channels where bot will not work. You can add multiple IDs & separate with <kbd>Space</kbd>.

`SLEEP_THRESHOLD` : Set a sleep threshold for flood wait exceptions happening globally in this telegram bot instance, below which any request that raises a flood wait will be automatically invoked again after sleeping for the required amount of time. Flood wait exceptions requiring higher waiting times will be raised. Defaults to 60 seconds.

`WORKERS` : Number of maximum concurrent workers for handling incoming updates. Defaults to `3`

`PORT` : The port that you want your webapp to be listened to. Defaults to `8080`

`WEB_SERVER_BIND_ADDRESS` : Your server bind adress. Defauls to `0.0.0.0`

`NO_PORT` : If you don't want your port to be displayed. You should point your `PORT` to `80` (http) or `443` (https) for the links to work. Ignore this if you're on Heroku.

`FQDN` :  A Fully Qualified Domain Name if present. Defaults to `WEB_SERVER_BIND_ADDRESS`

## How to use the bot

:warning: **Before using the  bot, don't forget to add the bot to the `BIN_CHANNEL` as an Admin**
 
`/start` : To check if the bot is alive or not.

To get an instant stream link, just forward any media to the bot and boom, its fast af.

## Demo Bot:
<a href="https://t.me/AH_File2Link_Bot"><img src="https://img.shields.io/badge/Demo-Telegram%20Bot-blue.svg?logo=telegram"></a>

[![Heroku Awake](https://github.com/AbirHasan2005/Pyro-FileStreamBot/actions/workflows/Alive.yml/badge.svg?event=status)](https://github.com/AbirHasan2005/Pyro-FileStreamBot/actions/workflows/Alive.yml)

## Special Features:
- Added Mongo Database Support!
- Added Broadcast Feature!
- Added Both Channel & Private Chat Support!
- Made Better Bin Channel Logging!

### Channel Support:
Bot also Supported with Channels. Just add bot Channel as Admin. If any new file comes in Channel it will edit it with **Get Download Link** Button.

## Credits

- [@EverythingSuckz](https://github.com/EverythingSuckz) & [@AbirHasan2005](https://github.com/AbirHasan2005)
- [eyaadh](https://github.com/eyaadh) for his awesome [Megatron Bot](https://github.com/eyaadh/megadlbot_oss)
- [Dan Tès](https://telegram.dog/haskell) for his [Pyrogram Library](https://github.com/pyrogram/pyrogram)
