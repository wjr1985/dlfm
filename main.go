package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/shkh/lastfm-go/lastfm"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var conf = struct {
	LastFm struct {
		APIKey        string `toml:"api_key"`
		Username      string `toml:"username"`
		CheckInterval int    `toml:"check_interval"`
	} `toml:"lastfm"`
	Discord struct {
		UseAppMode bool   `toml:"use_app_mode"`
		Token      string `toml:"token"`
		AppID      int    `toml:"app_id"`
	} `toml:"discord"`
	App struct {
		Title       string `toml:"title"`
		FirstLine   string `toml:"first_line"`
		SecondLine  string `toml:"second_line"`
		LargeImage  string `toml:"large_image"`
		LargeText   string `toml:"large_text"`
		SmallImage  string `toml:"small_image"`
		SmallText   string `toml:"small_text"`
		ShowButton  bool   `toml:"show_button"`
		EndlessMode bool   `toml:"endless_mode"`
	} `toml:"app"`
}{}

func init() {
	f, err := os.Open("config.toml")
	if err != nil {
		log.Println("Error loading config:", err)
		end(nil)
		os.Exit(1)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil && err != io.EOF {
		log.Println("Error loading config:", err)
		end(nil)
		os.Exit(1)
	}
	if err := toml.Unmarshal(data, &conf); err != nil {
		log.Println("Error loading config:", err)
		end(nil)
		os.Exit(1)
	}
	log.Println("Config loaded")
	var (
		noErr bool
		strs  = []*string{
			&conf.LastFm.APIKey, &conf.LastFm.Username,
			&conf.Discord.Token, &conf.App.Title,
			&conf.App.LargeImage, &conf.App.LargeText,
			&conf.App.SmallImage, &conf.App.SmallText,
		}
	)
	for _, v := range strs {
		*v = strings.TrimSpace(*v)
	}
	if conf.LastFm.APIKey == "" {
		log.Println("Error: lastfm.api_key can't be empty")
	} else if conf.LastFm.Username == "" {
		log.Println("Error: lastfm.username can't be empty")
	} else if conf.Discord.Token == "" && !conf.Discord.UseAppMode {
		log.Println("Error: discord.token can't be empty in token mode")
	} else if conf.Discord.AppID == 0 && conf.Discord.UseAppMode {
		log.Println("Error: discord.app_id can't be empty in app mode")
	} else if conf.LastFm.CheckInterval < 1 {
		if conf.LastFm.CheckInterval == 0 {
			log.Println("Grindcore mode on!")
			noErr = true
		} else {
			log.Println("Error: invalid check_interval (< 0)")
		}
	} else {
		noErr = true
	}
	if !noErr {
		end(nil)
		os.Exit(1)
	}
	if conf.App.Title == "" {
		conf.App.Title = "last.fm"
	}

	if conf.App.FirstLine == "" {
		conf.App.FirstLine = "{{name}}"
	}
	if conf.App.SecondLine == "" {
		conf.App.SecondLine = "{{artist}}"
	}
	if conf.App.LargeImage == "" {
		conf.App.LargeImage = "{{album_image}}"
	}
	if conf.App.LargeText == "" {
		conf.App.LargeText = "{{album}}"
	}
	if conf.App.SmallImage == "" {
		conf.App.SmallImage = "{{lastfm}}"
	}
	if conf.App.SmallText == "" {
		conf.App.SmallText = "github.com/dikey0ficial/dlfm"
	}
	for name, url := range map[string]string{
		// got this URLs from search, they aren't mine
		// and can be with any license (i haven't checked)
		"lastfm":     "http://icons.iconarchive.com/icons/danleech/simple/512/lastfm-icon.png",
		"deezer":     "https://www.macupdate.com/images/icons512/60905.png",
		"youtube":    "https://seeklogo.com/images/Y/youtube-music-logo-50422973B2-seeklogo.com.png",
		"apple":      "http://ixd.prattsi.org/wp-content/uploads/2017/01/apple_music_logo_by_mattroxzworld-d982zrj.png",
		"vk":         "https://seeklogo.com/images/V/vk-icon-logo-10188561D5-seeklogo.com.png",
		"yandex":     "https://download.cdn.yandex.net/from/yandex.ru/support/ru/music/files/icon_main.png",
		"soundcloud": "https://icons.iconarchive.com/icons/sicons/basic-round-social/512/soundcloud-icon.png",
	} {
		if conf.App.SmallImage == "{{"+name+"}}" {
			conf.App.SmallImage = url
			break
		}
	}
}

func end(err error) {
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Press enter to close window...")
	fmt.Scanln()
}

func main() {
	api := lastfm.New(conf.LastFm.APIKey, "")

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

	interval := time.Duration(conf.LastFm.CheckInterval) * time.Second
	ticker := time.NewTicker(interval)

	var (
		prevTrack string
	)
BIGFOR:
	for {
		select {
		case <-ticker.C:
			result, err := api.User.GetRecentTracks(lastfm.P{"limit": "1",
				"user": conf.LastFm.Username})
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
