# the-ghost-of-you-spotify

[![Go Reference](https://img.shields.io/badge/go-reference-blue?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/elliotwutingfeng/the-ghost-of-you-spotify)
[![Go Report Card](https://goreportcard.com/badge/github.com/elliotwutingfeng/the-ghost-of-you-spotify?style=for-the-badge)](https://goreportcard.com/report/github.com/elliotwutingfeng/the-ghost-of-you-spotify)
[![License](https://img.shields.io/badge/LICENSE-BSD--3--CLAUSE-GREEN?style=for-the-badge)](LICENSE)

## What is this?

A paranormal fix for Spotify's "Liked Songs" sync bug.

There is an ongoing problem with the [Spotify](https://spotify.com) music player app where the "Liked Songs" playlist fails to sync consistently across devices. Reinstalling the Spotify app usually does not help.

A workaround is to add a new track to the "Liked Songs" playlist, in order to force a cross-device update.

This program performs the same workaround; it adds a new random "ghost" track from the Spotify catalogue to your "Liked Songs" playlist, removes it, and then vanishes, like a ghost.

Schedule it with a job scheduler like cron and let [the ghost of you](https://en.wikipedia.org/wiki/The_Ghost_of_You) haunt your playlist just long enough to keep it in sync.

> [!WARNING]
> This program is designed to modify the content of your Spotify account.
>
> The credentials in the [.env](./.env) file must be kept safe as they can be used to access your Spotify account.

## Requirements

Tested on Linux x64

-   Go 1.24
-   [Spotify Developer Account](https://developer.spotify.com) (can be the same as your Spotify account)

## Setup

### Setup Environment File

Run the following

```bash
cp --update=none .env.txt .env
```

Alternatively, make a copy of [.env.txt](./.env.txt) and name it `.env`.

### Create A Spotify App

1. Sign in to your [Spotify Developer](https://developer.spotify.com) account.
1. [Create](https://developer.spotify.com/dashboard/create) a spotify app.
1. Fill in the following details

    - **App name**: `The Ghost Of You Spotify`
    - **App description**: `Keep your liked songs synced up across your devices.`
    - **Redirect URIs**: `http://127.0.0.1:3000/callback`
    - **APIs used**: Select `Web API`.

1. Agree with the Terms of Service and click "Save".
1. In the Spotify dashboard, select your newly created app and copy the `Client ID` and `Client secret` to the [.env](./.env) file in this project's root directory under the fields `CLIENT_ID` and `CLIENT_SECRET`.
1. If your Spotify market country code is not "SG" (Singapore), replace the value of `MARKET` in [.env](./.env) file with the correct country code. A list of market country codes is provided by [Spotify API](https://developer.spotify.com/documentation/web-api/reference/get-available-markets) (sign-in required).

## Usage

> [!IMPORTANT]
> When running this program for the first time, your web browser will open up a Spotify sign-in webpage for you to sign in and grant your newly created app access permissions to your Spotify account's library.
>
> Sign-in should not be needed for subsequent runs as this program will use a refresh token stored in the [.env](./.env) file for future requests to the Spotify API.

Run the following. If prompted via web browser, sign-in to your Spotify account and grant app access.

```bash
go run main.go
```

If successful, the program should terminate with output similar to the following

```text
🔍 Looking for a suitable ghost track...
🎯 Found track   | ID: 0000000000000000000000
📝 Added track   | ID: 0000000000000000000000
❌ Removed track | ID: 0000000000000000000000
👻 Boo! Your "Liked Songs" playlist should be synced up now across all devices.
```

## Run In Job scheduler

You can use a job scheduler like cron to run this program at regular intervals to keep your Spotify "Liked Songs" playlist in sync across devices.

Build the program.

```bash
# This creates an executable program file at ./dist/theghostofyouspotify
make build
```

**Example:** Set up the cron job to run every 15 minutes. Ensure that your [.env](./.env) file is in place with the correct values before installing the new cron job.

```bash
*/15 * * * * /path/to/dist/theghostofyouspotify
```

## Disclaimer

-   This program is licensed under [BSD-3-Clause](LICENSE). There is no warranty.
-   This is not an official Spotify product or service.
