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

> [!NOTE]
> Checkout [python branch](https://github.com/EverythingSuckz/TG-FileStreamBot/tree/python) if you are interested in that.

<hr>

<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#how-to-make-your-own">How to make your own</a>
      <ul>
        <li><a href="#deploy-to-heroku">Deploy to Heroku</a></li>
      </ul>
      <ul>
        <li><a href="#download-from-releases">Download and run</a></li>
        <li><a href="#run-using-docker-compose">Run via Docker compose</a></li>
        <li><a href="#run-using-docker">Run via Docker</a></li>
        <li><a href="#build-from-source">Build and run</a>
          <ul>
            <li><a href="#ubuntu">Ubuntu</a></li>
            <li><a href="#windows">Windows</a></li>
          </ul>
        </li>
      </ul>
    </li>
    <li>
      <a href="#setting-up-things">Setting up Things</a>
      <ul>
        <li><a href="#required-vars">Required environment variables</a></li>
        <li><a href="#optional-vars">Optional environment variables</a></li>
        <li><a href="#use-multiple-bots-to-speed-up">Using multiple bots</a></li>
        <li><a href="#use-multiple-bots-to-speed-up">Using user session to auto add bots</a>
          <ul>
            <li><a href="#what-it-does">What it does?</a></li>
            <li><a href="#how-to-generate-a-session-string">How to generate a session string?</a></li>
          </ul>
        </li>
      </ul>
    </li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#contact-me">Contact me</a></li>
    <li><a href="#credits">Credits</a></li>
  </ol>
</details>



## How to make your own

### Deploy to Heroku

> [!NOTE]
> You'll have to [fork](https://github.com/EverythingSuckz/TG-FileStreamBot/fork) this repository to deploy to Heroku.

Press the below button to fast deploy to Heroku

[![Deploy To Heroku](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

[Click Here](https://devcenter.heroku.com/articles/config-vars#using-the-heroku-dashboard) to know how to add / edit [environment variables](#required-vars) in Heroku.

<hr>

### Download from releases
- Head over to [releases](https://github.com/EverythingSuckz/TG-FileStreamBot/releases) tab, from the *pre release* section, download the one for your platform and architecture.
- Extract the zip file to a folder.
- Create an a file named `fsb.env` and add all the variables there (see `fsb.sample.env` file for reference).
- Give the executable file permission to execute using the command `chmod +x fsb` (Not required for windows).
- Run the bot using `./fsb run` command. ( `./fsb.exe run` for windows)

<hr>

### Run using docker-compose

- Clone the repository
```sh
git clone https://github.com/EverythingSuckz/TG-FileStreamBot
cd TG-FileStreamBot
```

- Create an a file named `fsb.env` and add all the variables there (see `fsb.sample.env` file for reference).

```sh
nano fsb.env
```

- Build and run the docker-compose file

```sh
docker-compose up -d
```
 OR

```sh
docker compose up -d
```

<hr>

### Run using docker

```sh
docker run --env-file fsb.env ghcr.io/everythingsuckz/fsb:latest
```
Where `fsb.env` is the environment file containing all the variables.

<hr>

### Build from source

#### Ubuntu

> [!NOTE]
> Make sure to install go 1.21 or above.
> Refer https://stackoverflow.com/a/17566846/15807350

```sh
git clone https://github.com/EverythingSuckz/TG-FileStreamBot
cd TG-FileStreamBot
go build ./cmd/fsb/
chmod +x fsb
mv fsb.sample.env fsb.env
nano fsb.env
# (add your environment variables, see the next section for more info)
./fsb run
```

and to stop the program,
 do <kbd>CTRL</kbd>+<kbd>C</kbd>

#### Windows

> [!NOTE]
> Make sure to install go 1.21 or above.

```powershell
git clone https://github.com/EverythingSuckz/TG-FileStreamBot
cd TG-FileStreamBot
go build ./cmd/fsb/
Rename-Item -LiteralPath ".\fsb.sample.env" -NewName ".\fsb.env"
notepad fsb.env
# (add your environment variables, see the next section for more info)
.\fsb run
```

and to stop the program,
 do <kbd>CTRL</kbd>+<kbd>C</kbd>

## Setting up things

If you're locally hosting, create a file named `fsb.env` in the root directory and add all the variables there.
You may check the `fsb.sample.env`.
An example of `fsb.env` file:

```sh
API_ID=452525
API_HASH=esx576f8738x883f3sfzx83
BOT_TOKEN=55838383:yourbottokenhere
LOG_CHANNEL=-10045145224562
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

- `LOG_CHANNEL` :  This is the channel ID for the log channel where the bot will forward media messages and store these files to make the generated direct links work. To obtain a channel ID, create a new telegram channel (public or private), post something in the channel, forward the message to [@missrose_bot](https://telegram.dog/MissRose_bot) and **reply the forwarded message** with the /id command. Copy the forwarded channel ID and paste it into the this field.

### Optional Vars
In addition to the mandatory variables, you can also set the following optional variables:

- `PORT` : This sets the port that your webapp will listen to. The default value is 8080.

- `HOST` :  A Fully Qualified Domain Name if present or use your server IP. (eg. `https://example.com` or `http://14.1.154.2:8080`)

- `HASH_LENGTH` : Custom hash length for generated URLs. The hash length must be greater than 5 and less than or equal to 32. The default value is 6.

- `USE_SESSION_FILE` : Use session files for worker client(s). This speeds up the worker bot startups. (default: `false`)

- `USER_SESSION` : A pyrogram session string for a user bot. Used for auto adding the bots to `LOG_CHANNEL`. (default: `null`)

- `ALLOWED_USERS` : A list of user IDs separated by comma (`,`). If this is set, only the users in this list will be able to use the bot. (default: `null`)

<hr>

### Use Multiple Bots to speed up

> [!NOTE]
> **What it multi-client feature and what it does?** <br>
> This feature shares the Telegram API requests between worker bots to speed up download speed when many users are using the server and to avoid the flood limits that are set by Telegram. <br>

> [!NOTE]
> You can add up to 50 bots since 50 is the max amount of bot admins you can set in a Telegram Channel.

To enable multi-client, generate new bot tokens and add it as your `fsb.env` with the following key names. 

`MULTI_TOKEN1`: Add your first bot token here.

`MULTI_TOKEN2`: Add your second bot token here.

you may also add as many as bots you want. (max limit is 50)
`MULTI_TOKEN3`, `MULTI_TOKEN4`, etc.

> [!WARNING]
> Don't forget to add all these worker bots to the `LOG_CHANNEL` for the proper functioning

### Using user session to auto add bots

> [!WARNING]
> This might sometimes result in your account getting resticted or banned.
> **Only newly created accounts are prone to this.**

To use this feature, you need to generate a pyrogram session string for the user account and add it to the `USER_SESSION` variable in the `fsb.env` file.

#### What it does?

This feature is used to auto add the worker bots to the `LOG_CHANNEL` when they are started. This is useful when you have a lot of worker bots and you don't want to add them manually to the `LOG_CHANNEL`.

#### How to generate a session string?

The easiest way to generate a session string is by running

```sh
./fsb session --api-id <your api id> --api-hash <your api hash>
```

<img src="https://github.com/EverythingSuckz/TG-FileStreamBot/assets/65120517/b5bd2b88-0e1f-4dbc-ad9a-faa6d5a17320" height=300>

<br><br>

This will generate a session string for your user account using QR code authentication. Authentication via phone number is not supported yet and will be added in the future.

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
