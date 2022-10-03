package main

import (
	"github.com/BurntSushi/toml"
	"github.com/shkh/lastfm-go/lastfm"

	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func init() {
	initConfig()
}

func end(err error) {
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Press enter to close window...")
	fmt.Scanln()
}

func main() {
	api := lastfm.New(conf.LastFM.APIKey, "")

	log.Println("Connected to last.fm")

	if conf.App.EndlessMode {
		log.Println("Endless mode! Press `Ctrl+C` to exit")
	}

	var (
		indentificator string
		updater        StatusUpdater
	)

	if conf.Discord.UseAppMode {
		updater = AppStatusUpdater{}
		indentificator = strconv.Itoa(conf.Discord.AppID)
	} else {
		updater = &TokenModeStatusUpdater{}
		indentificator = conf.Discord.Token
	}

	if err := updater.Login(indentificator); err != nil {
		end(err)
		os.Exit(1)
	}

	defer func() { updater.Logout() }()
	log.Println("Authorized and connected to discord")

	interval := time.Duration(conf.LastFM.CheckInterval) * time.Second
	ticker := time.NewTicker(interval)

	var (
		prevTrack string
	)
BIGFOR:
	for {
		select {
		case <-ticker.C:
			result, err := api.User.GetRecentTracks(lastfm.P{"limit": "1",
				"user": conf.LastFM.Username})
			if err != nil {
				log.Println("last.fm error: ", err)
				if !conf.App.EndlessMode {
					end(nil)
					return
				}
			} else {
				if len(result.Tracks) > 0 {
					ctrack := result.Tracks[0]
					isNowPlaying, _ := strconv.ParseBool(ctrack.NowPlaying)
					trackName := ctrack.Artist.Name + " - " + ctrack.Name
					if isNowPlaying {
						if err := updater.Set(result); err != nil {
							if err.Error() == "The pipe is being closed." {
								log.Println("Error: discord disconnected =(")
								break BIGFOR
							}
							log.Println("Discord error: ", err)
							if !conf.App.EndlessMode {
								end(nil)
								break
							}
						}
						if prevTrack != trackName {
							log.Println("Now playing: " + trackName)
							prevTrack = trackName
						}
					} else if !isNowPlaying {
						if err := updater.Clear(); err != nil {
							log.Println("Discord error: ", err)
							if !conf.App.EndlessMode {
								end(nil)
								break
							}
						}
					}
				}
			}
		}
	}
	end(nil)
}
