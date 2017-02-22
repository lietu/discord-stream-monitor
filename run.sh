#!/usr/bin/env bash

# Ensure the environment is properly set up if running in supervisord
if [[ -f "/etc/profile" ]]; then
    source "/etc/profile"
fi

if [[ -f "$HOME/.profile" ]]; then
    source "$HOME/.profile"
fi

# Configure your tokens
export DSM_TOKEN="DISCORD_BOT_TOKEN_GOES_HERE"
export DSM_CLIENT_ID="TWITCH_CLIENT_ID_GOES_HERE"

# Update to latest version and build (feel free to comment out)
go get -u github.com/lietu/discord-stream-monitor
rm -f dsm
go build dsm.go

# Run it and let it take over this process, that way events from supervisord
# etc. are properly relayed to dsm, not this script.
exec ./dsm
