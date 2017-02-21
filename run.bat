@echo off

set DSM_TOKEN=DISCORD_BOT_TOKEN_GOES_HERE
set DSM_CLIENT_ID=TWITCH_CLIENT_ID_GOES_HERE

del dsm.exe
go build dsm.go
dsm.exe
