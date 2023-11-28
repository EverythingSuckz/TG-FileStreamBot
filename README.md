<h1 align="center">Telegram File Stream Bot</h1>
<p align="center">
  </a>
  <p align="center">
    <a herf="https://github.com/EverythingSuckz/TG-FileStreamBot">
        <img src="https://telegra.ph/file/a8bb3f6b334ad1200ddb4.png" height="100" width="100" alt="File Stream Bot Logo">
    </a>
</p>
  <p align="center">
    A Telegram bot to <b>generate direct link</b> for your Telegram files.
    <br />
  </p>
</p>

<hr>

> [!WARNING]
> This branch is still under development and it is advised to use the [python branch](https://github.com/EverythingSuckz/TG-FileStreamBot/tree/main).

<hr>

<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#how-to-make-your-own">How to make your own</a>
      <ul>
        <li><a href="#download-from-releases">Download and run</a></li>
      </ul>
      <ul>
        <li><a href="#build-from-source">Run it in a VPS / local</a></li>
      </ul>
    </li>
    <li>
      <a href="#setting-up-things">Setting up Things</a>
      <ul>
        <li><a href="#required-vars">Required environment variables</a></li>
      </ul>
      <ul>
        <li><a href="#optional-vars">Optional environment variables</a></li>
      </ul>
      <ul>
        <li><a href="#use-multiple-bots-to-speed-up">Using multiple bots</a></li>
      </ul>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#contact-me">Contact me</a></li>
    <li><a href="#credits">Credits</a></li>
  </ol>
</details>

## How to make your own

### Download from releases
(available soon)

### Build from source

```sh
git clone https://github.com/EverythingSuckz/TG-FileStreamBot -b golang
cd TG-FileStreamBot
go get -u
go build .
chmod +x fsb
./fsb
```

and to stop the program,
 do <kbd>CTRL</kbd>+<kbd>C</kbd>

## Setting up things

If you're locally hosting, create a file named `.env` in the root directory and add all the variables there.
An example of `.env` file:

```sh
API_ID=452525
API_HASH=esx576f8738x883f3sfzx83
BOT_TOKEN=55838383:yourbottokenhere
BIN_CHANNEL=-10045145224562
PORT=8080
HOST=http://yourserverip
# (if you want to set up multiple bots)
MULTI_TOKEN1=55838373:yourworkerbottokenhere
MULTI_TOKEN2=55838355:yourworkerbottokenhere
```

### Required Vars
Before running the bot, you will need to set up the following mandatory variables:

- `API_ID` : This is the API ID for your Telegram account, which can be obtained from my.telegram.org.

- `API_HASH` : This is the API hash for your Telegram account, which can also be obtained from my.telegram.org.

- `BOT_TOKEN` : This is the bot token for the Telegram Media Streamer Bot, which can be obtained from [@BotFather](https://telegram.dog/BotFather).

- `BIN_CHANNEL` :  This is the channel ID for the log channel where the bot will forward media messages and store these files to make the generated direct links work. To obtain a channel ID, create a new telegram channel (public or private), post something in the channel, forward the message to [@missrose_bot](https://telegram.dog/MissRose_bot) and **reply the forwarded message** with the /id command. Copy the forwarded channel ID and paste it into the this field.

### Optional Vars
In addition to the mandatory variables, you can also set the following optional variables:

- `PORT` : This sets the port that your webapp will listen to. The default value is 8080.

- `HOST` :  A Fully Qualified Domain Name if present or use your server IP. (eg. `https://example.com` or `http://14.1.154.2:8080`)

- `HASH_LENGTH` : Custom hash length for generated URLs. The hash length must be greater than 5 and less than or equal to 32. The default value is 6.

- `USE_SESSION_FILE` : Use session files for worker client(s). This speeds up the worker bot startups. (default: `false`)

### Use Multiple Bots to speed up

> **Note**
> What it multi-client feature and what it does? <br>
> This feature shares the Telegram API requests between worker bots to speed up download speed when many users are using the server and to avoid the flood limits that are set by Telegram. <br>

> **Note**
> You can add up to 50 bots since 50 is the max amount of bot admins you can set in a Telegram Channel.

To enable multi-client, generate new bot tokens and add it as your `.env` with the following key names. 

`MULTI_TOKEN1`: Add your first bot token here.

`MULTI_TOKEN2`: Add your second bot token here.

you may also add as many as bots you want. (max limit is 50)
`MULTI_TOKEN3`, `MULTI_TOKEN4`, etc.

> **Warning**
> Don't forget to add all these worker bots to the `BIN_CHANNEL` for the proper functioning

## Contributing

Feel free to contribute to this project if you have any further ideas

## Contact me

[![Telegram Channel](https://img.shields.io/static/v1?label=Join&message=Telegram%20Channel&color=blueviolet&style=for-the-badge&logo=telegram&logoColor=violet)](https://xn--r1a.click/wrench_labs)
[![Telegram Group](https://img.shields.io/static/v1?label=Join&message=Telegram%20Group&color=blueviolet&style=for-the-badge&logo=telegram&logoColor=violet)](https://xn--r1a.click/AlteredVoid)

You can contact either via my [Telegram Group](https://xn--r1a.click/AlteredVoid) or you can message me on [@EverythingSuckz](https://xn--r1a.click/EverythingSuckz)


## Credits

- [@celestix](https://github.com/celestix) for [gotgproto](https://github.com/celestix/gotgproto)
- [@divyam234](https://github.com/divyam234/teldrive) for his [Teldrive](https://github.com/divyam234/teldrive) Project

## Copyright

Copyright (C) 2023 [EverythingSuckz](https://github.com/EverythingSuckz) under [GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.en.html).

TG-FileStreamBot is Free Software: You can use, study share and improve it at your
will. Specifically you can redistribute and/or modify it under the terms of the
[GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.en.html) as
published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version. Also keep in mind that all the forks of this repository MUST BE OPEN-SOURCE and MUST BE UNDER THE SAME LICENSE.
