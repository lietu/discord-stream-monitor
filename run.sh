#!/usr/bin/env bash

export DSM_TOKEN="DISCORD_BOT_TOKEN_GOES_HERE"
export DSM_CLIENT_ID="TWITCH_CLIENT_ID_GOES_HERE"

rm -f dsm
go build dsm.go
./dsm
