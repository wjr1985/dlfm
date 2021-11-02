package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-ini/ini"
	"github.com/shkh/lastfm-go/lastfm"
)

func scrobbler() error {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Println(err)
		return err
	}

	token := cfg.Section("discord").Key("token").String()
	apiKey := cfg.Section("lastfm").Key("api_key").String()
	username := cfg.Section("lastfm").Key("username").String()
	title := cfg.Section("app").Key("title").String()
	endlessMode, err := strconv.ParseBool(cfg.Section("app").Key("endless_mode").String())
	configInterval, err := cfg.Section("lastfm").Key("check_interval").Int()

	if err != nil {
		log.Println(err)
		return err
	}

	api := lastfm.New(apiKey, "")

	log.Println("Settings loaded: config.ini")

	if endlessMode {
		log.Println("Endless mode! Ctrl+C to exit")
	}
	dg, err := discordgo.New(token)
	if err != nil {
		log.Println("Discord error: ", err)
		return err
	}
	log.Println("Authorized to Discord")
	if err := dg.Open(); err != nil {
		log.Println("Discord error: ", err)
		return err
	}
	defer dg.Close()
	log.Println("Connected to Discord")

	interval := time.Duration(configInterval*1000) * time.Millisecond
	ticker := time.NewTicker(interval)
	var prevTrack string
	for {
		select {
		case <-ticker.C:
			result, err := api.User.GetRecentTracks(lastfm.P{"limit": "1", "user": username})
			if err != nil {
				log.Println("LastFM error: ", err)
			} else {
				if len(result.Tracks) > 0 {
					currentTrack := result.Tracks[0]
					isNowPlaying, _ := strconv.ParseBool(currentTrack.NowPlaying)
					trackName := currentTrack.Artist.Name + " - " + currentTrack.Name
					if isNowPlaying {
						statusData := discordgo.UpdateStatusData{
							Game: &discordgo.Game{
								Name:    title,
								Type:    discordgo.GameTypeListening,
								Details: currentTrack.Name,
								State:   currentTrack.Artist.Name,
							},
							AFK:    false,
							Status: "online",
						}
						if err := dg.UpdateStatusComplex(statusData); err != nil {
							log.Println("Discord error: ", err)
							if !endlessMode {
								return err
							}
						}
						if prevTrack != trackName {
							log.Println("Now playing: " + trackName)
							prevTrack = trackName
						}
					} else if !isNowPlaying {
						log.Println("!")
						statusData := discordgo.UpdateStatusData{
							Game:   nil,
							Status: "offline",
						}
						if err := dg.UpdateStatusComplex(statusData); err != nil {
							log.Println("Discord error: ", err)
							if !endlessMode {
								return err
							}
						}
					}
				}
			}
		}
	}
}

func main() {
	scrobbler()
	log.Println("Press the Enter Key to terminate the console screen!")
	fmt.Scanln()
}
