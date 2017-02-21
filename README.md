# Discord Stream Monitor

This is a tool to automatically delete expired Twitch stream links on your 
Discord servers.

If you have a channel like `#are-you-streaming` where Twitch streamers post
links to their streams when they go live, it's common to have a rule to delete
the post once the stream ends. However, since it's manual effort, it often gets
forgotten. This tool will help automate that task.

When the bot notices a message in one of the configured channels that mentions
what looks like a valid stream URL, it tries to figure out the Twitch stream
in question, and then monitors the stream status every 2 minutes. If the stream
is seen to be offline, the message mentioning it will be deleted.


## Configuration and setup

You will need to register yourself a bot application to Discord, add that bot
to your server, and then configure the channel filtering. You will also need
to register an app on Twitch to get API access.

Additionally you will need Golang set up on your computer. Get the installer
from https://golang.org/dl/.
 

### Register a Discord bot

Open up the address https://discordapp.com/developers/applications/me

Click on the giant `New App` -button. Give your app a name, and then click
`Create App`.

Click on the `Create a Bot User` -button, and set up the bot account.

Copy click on the `click to reveal` next to the `Token:` -field, and copy the
token for later use.


### Register a Twitch app

Open up the address https://www.twitch.tv/settings/connections

Scroll to the bottom of the page, and click "Register your application".

Give it a unique name, use e.g. `http://localhost` for redirect URI, pick a 
category, accept terms, and click on `Register`. Get the client ID for later
use.


### Add the bot to your server

Replace the `<CLIENT_ID>` in this URL with the one for your bot app, and then 
open it in a browser: 
https://discordapp.com/oauth2/authorize?client_id=<CLIENT_ID>&scope=bot


### Configure Discord and Twitch API access in the bot

Edit the launch script (`run.bat` on Windows, `run.sh` on *nix). Update the
`DSM_TOKEN` with the Discord chat token, and `DSM_CLIENT_ID` with the Twitch
Client ID.


### Configure channel filters

Open up the `dsm.go` -file in your favorite editor, and locate the variable
`MonitorChannels`.

At this stage the bot should launch when you run it (`run.bat` on Windows, 
`run.sh` on *nix). Start it up, and say something in the channel(s) you want
the bot to monitor.

You should see messages like this in the log:

```
2017/02/21 22:44:09 #87867525001396224 <lietu> test
```

The `#87867525001396224` indicates the Channel ID, remove the `#` from it and
that's the Channel ID you should add to the `MonitorChannels` -list.

At this stage also open up Discord, and configure the roles for the bot user.
It should have the `Manage Messages` -permission to delete messages.

Restart the bot, and confirm that it now reacts only to stream links on your
configured channels.

Open up Discord and in the configured channel(s) say e.g. 
`Live now on https://twitch.tv/lietu`. Assuming the mentioned stream is 
offline, the message should get deleted about 2 minutes later.


## License

The [discordgo](https://github.com/bwmarrin/disgord/blob/master/LICENSE) code
is licensed under it's own license. This code is licensed under the MIT 
license. Basically that means: go ahead, use it.
