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
        <li><a href="#deploy-to-koyeb">Deploy to Koyeb</a></li>
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

### Deploy to Koyeb

> [!IMPORTANT]
> You'll have to expand the "Environment variables and files" section and update the env variables before hitting the deploy button.

> [!NOTE]
> This deploys the **latest docker release and NOT the latest commit**. Since it uses prebuilt docker container, the deploy speed will be significantly faster.

[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?type=docker&name=file-stream-bot&image=ghcr.io/everythingsuckz/fsb:latest&env%5BAPI_HASH%5D=&env%5BAPI_ID%5D=&env%5BAPI_HASH%5D=&env%5BAPI_ID%5D=&env%5BBOT_TOKEN%5D=&env%5BHOST%5D=https%3A%2F%2F%7B%7B+KOYEB_PUBLIC_DOMAIN+%7D%7D&env%5BLOG_CHANNEL%5D=&env%5BPORT%5D=8038&ports=8038%3Bhttp%3B%2F&hc_protocol%5B8038%5D=tcp&hc_grace_period%5B8038%5D=5&hc_interval%5B8038%5D=30&hc_restart_limit%5B8038%5D=3&hc_timeout%5B8038%5D=5&hc_path%5B8038%5D=%2F&hc_method%5B8038%5D=get)

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
Before running the bot, set these mandatory variables:

- `API_ID` : Telegram API ID from https://my.telegram.org.
- `API_HASH` : Telegram API hash from https://my.telegram.org.
- `BOT_TOKEN` : Bot token from [@BotFather](https://telegram.dog/BotFather).
- `LOG_CHANNEL` : Telegram channel ID where bot messages/files are stored for streaming links.
  - How to get it: Create a channel, send a message, forward it to [@missrose_bot](https://telegram.dog/MissRose_bot), then reply with `/id`.

### Optional Vars
In addition to required variables, these optional ones are available:

- `DEV` (default: `false`)
  - Enables development logging behavior.

- `PORT` (default: `8080`)
  - HTTP port used by the stream server.

- `HOST` (default: auto-generated)
  - Base URL used in generated links (for example `https://example.com` or `http://192.168.1.10:8080`).
  - If omitted, bot auto-builds it from detected IP + `PORT`.

- `HASH_LENGTH` (default: `6`, min `5`, max `32`)
  - Length of URL hash used in stream links.

- `USE_SESSION_FILE` (default: `true`)
  - Reuse saved worker sessions to speed startup and reduce login overhead.

- `USER_SESSION` (default: empty)
  - Optional user session string for userbot features (for example auto-adding bots to `LOG_CHANNEL`).

- `USE_PUBLIC_IP` (default: `false`)
  - If enabled, bot tries to discover public IP for host generation.

- `ALLOWED_USERS` (default: empty)
  - Comma-separated Telegram user IDs that are allowed to use the bot (access allowlist).


#### Stream Performance Configuration

These optional variables allow you to tune the streaming performance. Most users won't need to change these defaults.

- `STREAM_CONCURRENCY` (default: `4`)
  - How many block downloads a single stream request runs in parallel.
  - Effective first-batch request fanout per stream request = `STREAM_CONCURRENCY`.
  - Higher values can improve throughput and startup latency, but increase Telegram API pressure.

- `STREAM_BUFFER_COUNT` (default: `8`)
  - Capacity of the in-memory block queue between downloader and HTTP writer.
  - It controls **how far ahead** prefetch can run before the reader catches up.
  - It does **not** increase parallel download count by itself.
  - Approximate extra memory per request: `STREAM_BUFFER_COUNT * blockSize` (blockSize is dynamic: 64KB/256KB/512KB/1MB).

- `STREAM_TIMEOUT_SEC` (default: `30`)
  - Per-block timeout for Telegram `UploadGetFile` request.

- `STREAM_MAX_RETRIES` (default: `3`)
  - Retry attempts per block before failing the stream.

##### Simple explanation (for non-technical users)

- `STREAM_CONCURRENCY` = **how many Telegram downloads happen at the same time** for one stream.
  - Bigger value = video may start faster and download faster.
  - But bigger value also means more chance of Telegram floodwait/rate-limit.

- `STREAM_BUFFER_COUNT` = **how many downloaded chunks are kept ready in memory**.
  - Bigger value = smoother playback on unstable networks.
  - It does **not** create more Telegram requests by itself.

- `STREAM_TIMEOUT_SEC` = **how long to wait before saying a chunk request is too slow**.

- `STREAM_MAX_RETRIES` = **how many times to retry a failed chunk**.

Quick rule:

- Telegram request pressure is mainly controlled by `STREAM_CONCURRENCY`.
- Approximate in-flight Telegram calls at peak = `active_streams × STREAM_CONCURRENCY`.

Safe starting presets:

- Small server / home host:
  - `STREAM_CONCURRENCY=4`
  - `STREAM_BUFFER_COUNT=8`
  - `STREAM_TIMEOUT_SEC=30`
  - `STREAM_MAX_RETRIES=3`

- Medium VPS:
  - `STREAM_CONCURRENCY=6`
  - `STREAM_BUFFER_COUNT=12`
  - `STREAM_TIMEOUT_SEC=45`
  - `STREAM_MAX_RETRIES=3`

- High bandwidth + multiple worker bots:
  - `STREAM_CONCURRENCY=8`
  - `STREAM_BUFFER_COUNT=16`
  - `STREAM_TIMEOUT_SEC=60`
  - `STREAM_MAX_RETRIES=5`

> [!NOTE]
> Increasing `STREAM_CONCURRENCY` increases concurrent Telegram requests from one client connection. If you push this too high, floodwait/rate-limit risk increases.
> Start with `4`, test `6` or `8`, and only go higher if logs show stable behavior.

**Example configuration for high-performance servers:**
```sh
STREAM_CONCURRENCY=8
STREAM_BUFFER_COUNT=16
STREAM_TIMEOUT_SEC=60
STREAM_MAX_RETRIES=5
```
This configuration would allow up to 8 blocks to be downloaded in parallel, with a buffer of 16 blocks, and a longer timeout for slow connections.

### `MULTI_TOKEN` variables

You can add worker bots to distribute stream requests across different bot tokens:

- `MULTI_TOKEN1`
- `MULTI_TOKEN2`
- `MULTI_TOKEN3`
- ...

Each active HTTP stream request uses one worker client in round-robin mode. Using multiple workers reduces floodwait risk versus sending all traffic through a single token.

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
- [@karu](https://github.com/krau) for adding image support

## Copyright

Copyright (C) 2023 [EverythingSuckz](https://github.com/EverythingSuckz) under [GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.en.html).

TG-FileStreamBot is Free Software: You can use, study share and improve it at your
will. Specifically you can redistribute and/or modify it under the terms of the
[GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.en.html) as
published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version. Also keep in mind that all the forks of this repository MUST BE OPEN-SOURCE and MUST BE UNDER THE SAME LICENSE.
