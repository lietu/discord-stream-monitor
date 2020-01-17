package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"regexp"
	"net/http"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"time"
)

// List of Discord channels the monitor is active on
var MonitorChannels = []string{
	"283689700802166785",
}

// Regular expressions to match streams
// Valid twitch.tv/username -links, as well as username.tv/live -links work
var twitchLinkRe = regexp.MustCompile("(?:https?://)?(?:www\\.)?twitch\\.tv/([a-zA-Z0-9_]+)[/]?(?: |$)")
var liveLinkRe = regexp.MustCompile("(?:https?://)?(?:www\\.)?([a-zA-Z0-9]+)\\.tv/live[/]?(?: |$)")

// ----- Data structs for Twitch API ----- //

type User struct {
	Id          string `json:"_id"`
	Bio         string `json:"bio"`
	CreatedAt   string `json:"created_at"`
	DisplayName string `json:"display_name"`
	Logo        string `json:"logo"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	UpdatedAt   string `json:"updated_at"`
}

type Follow struct {
	CreatedAt     string `json:"created_at"`
	Notifications bool `json:"notifications"`
	User          *User `json:"user"`
}

type Preview struct {
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Large    string `json:"large"`
	Template string `json:"template"`
}

type Channel struct {
	Mature                       bool `json:"mature"`
	Status                       string `json:"status"`
	BroadcasterLanguage          string `json:"broadcaster_language"`
	DisplayName                  string `json:"display_name"`
	Game                         string `json:"game"`
	Language                     string `json:"language"`
	Id                           int `json:"_id"`
	Name                         string `json:"name"`
	CreatedAt                    string `json:"created_at"`
	UpdatedAt                    string `json:"updated_at"`
	Partner                      bool `json:"partner"`
	Logo                         string `json:"logo"`
	VideoBanner                  string `json:"video_banner"`
	ProfileBanner                string `json:"profile_banner"`
	ProfileBannerBackgroundColor string `json:"profile_banner_background_color"`
	Url                          string `json:"url"`
	Views                        int `json:"views"`
	Followers                    int `json:"followers"`
}

type Stream struct {
	Id          int `json:"_id"`
	Game        string `json:"game"`
	Viewers     int `json:"viewers"`
	VideoHeight int`json:"video_height"`
	AverageFPS  int `json:"average_fps"`
	Delay       float32 `json:"delay"`
	CreatedAt   string `json:"created_at"`
	is_playlist bool `json:"is_playlist"`
	Preview     *Preview `json:"preview"`
	Channel     *Channel `json:"channel"`
}

type UsersResponse struct {
	Total int `json:"_total"`
	Users []*User `json:"users"`
}

type FollowsResponse struct {
	Cursor  string `json:"_cursor"`
	Total   int `json:"_total"`
	Follows []*Follow `json:"follows"`
}

type StreamResponse struct {
	Stream *Stream `json:"stream"`
}

// Session is declared in the global space so it can be easily used
// throughout this program.
// In this use case, there is no error that would be returned.
var Session, _ = discordgo.New()
var Token = ""
var ClientID = ""

// Read in all configuration options from both environment variables and
// command line arguments.
func init() {
	// Discord Authentication Token
	Token = os.Getenv("DSM_TOKEN")
	ClientID = os.Getenv("DSM_CLIENT_ID")

	if Token == "" {
		flag.StringVar(&Token, "t", "", "Discord Authentication Token")
	}

	if ClientID == "" {
		flag.StringVar(&ClientID, "c", "", "Twitch Client ID")
	}
}

// Print a Discord message in the log
func printMsg(m *discordgo.MessageCreate) {
	log.Printf("#%s <%s> %s", m.ChannelID, m.Author.Username, m.Content)
}

// Try and determine the advertised stream in the message (if any)
func getStreamerAdvertised(content string) string {
	if matches := twitchLinkRe.FindStringSubmatch(content); len(matches) > 0 {
		return matches[1]
	}

	if matches := liveLinkRe.FindStringSubmatch(content); len(matches) > 0 {
		return matches[1]
	}

	return ""
}

// Figure out the stream state for the given Twitch channel ID
func getStreamState(channelId string) string {
	state := "unknown"

	path := fmt.Sprintf("streams/%s", channelId)
	res, err := kraken(path)
	if err != nil {
		return state
	}
	defer res.Body.Close()

	ksr := StreamResponse{}
	err = json.NewDecoder(res.Body).Decode(&ksr)
	if err != nil {
		// Unknown response, something is wrong
		return state
	}

	if ksr.Stream == nil {
		state = "offline"
	} else {
		state = "live"
	}

	log.Printf("Twitch stream %s is %s", channelId, state)

	return state
}

// Figure out the Twitch channel ID from the Twitch username
func getTwitchChannelId(username string) string {
	path := fmt.Sprintf("users/?login=%s", username)
	res, err := kraken(path)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	kur := UsersResponse{}
	err = json.NewDecoder(res.Body).Decode(&kur)
	if err != nil {
		return ""
	}

	if len(kur.Users) == 0 {
		return ""
	}

	return kur.Users[0].Id
}

// Make a request to the Twitch kraken API
func kraken(path string) (*http.Response, error) {
	c := &http.Client{}

	url := fmt.Sprintf("https://api.twitch.tv/kraken/%s", path)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("Client-ID", ClientID)

	return c.Do(req)
}

// Check if the given channel ID is ok for monitoring
func channelOk(channelID string) bool {
	for _, c := range MonitorChannels {
		if c == channelID {
			return true
		}
	}

	return false
}

// Monitor the stream, once the stream is detected to go offline the original
// message mentioning it will be deleted.
func monitor(messageID string, channelID string, streamer string) {
	//Session.ChannelMessageSend(channelID, fmt.Sprintf("Twitch stream for %s mentioned, notifying when it goes offline.", streamer))

	// Determing the Twitch channel ID
	tries := 0
	twitchId := ""
	for twitchId == "" {
		twitchId = getTwitchChannelId(streamer)

		if twitchId == "" {
			tries += 1

			if tries > 10 {
				log.Printf("Failed to determine channel ID for Twitch stream %s, giving up.", streamer)
			}

			time.Sleep(time.Second * 5)
		}
	}

	log.Printf("Resolved channel ID for %s to %s", streamer, twitchId)

	for {
		time.Sleep(time.Minute * 2) // Time between checks

		state := getStreamState(twitchId)
		if state == "offline" {
			//Session.ChannelMessageSend(channelID, fmt.Sprintf("Stream for %s has gone offline.", streamer))

			err := Session.ChannelMessageDelete(channelID, messageID)
			if err != nil {
				// Generally caused by lack of permissions
				log.Printf("Got error when trying to delete message: %s", err)
			}
			return
		} else {
			log.Printf("Channel %s / %s is still live", streamer, twitchId)
		}
	}
}

// Handle a message in Discord that the bot received
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !channelOk(m.ChannelID) {
		return
	}

	streamer := getStreamerAdvertised(m.Content)
	if streamer != "" {
		log.Printf("Stream from %s mentioned, monitoring..", streamer)

		go monitor(m.ID, m.ChannelID, streamer)
	}
}

func main() {
	var err error

	log.Print("Starting up Discord Stream Monitor.")
	log.Print("")

	flag.Parse()

	if Token == "" {
		log.Println("You must provide a Discord authentication token.")
		return
	}

	if ClientID == "" {
		log.Println("You must provide a Twitch Client ID.")
		return
	}

	Session.Token = fmt.Sprintf("Bot %s", Token)

	Session.State.User, err = Session.User("@me")
	if err != nil {
		log.Printf("error fetching user information, %s\n", err)
	}

	Session.AddHandler(handleMessage)

	err = Session.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
	}

	// Wait for a CTRL-C
	log.Print("Now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	Session.Close()

	// Exit Normally.
}
