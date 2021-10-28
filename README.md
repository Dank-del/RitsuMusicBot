# Ritsu

## Telegram bot to flex your music taste, written in pure go

Currently a work in progress.

## Features

- Track scrobbles through `Last.FM` API
- Lyric search

### Current file structure

```
├───config
│       helpers.go
│       methods.go
│       types.go
│       vars.go
│
├───database
│       client.go
│       users.go
│
├───handlers
│       constants.go
│       error.go
│       getstatus.go
│       help.go
│       history.go
│       load.go
│       lyrics.go
│       me.go
│       misc.go
│       odesli.go
│       registration.go
│       start.go
│       status.go
│       topArtists.go
│
├───last.fm
│       constants.go
│       helpers.go
│       types.go
│
├───logging
│       constants.go
│       methods.go
│
├───lyrics
│       constants.go
│       helpers.go
│       types.go
│       vars.go
│
├───odesli
│       constants.go
│       methods.go
│       types.go
│
├───tests
│       getLastFMTrack_test.go
│       getLyrics_test.go
│       getRecentTracks_test.go
│       gettopartists_test.go
│       getUserTopTracks_test.go
│       getUser_test.go
│       odesli_test.go
│
└───utilities
        admin_check.go

```
