package lib

import (
	"github.com/shkh/lastfm-go/lastfm"

	"log"
	"strconv"
	"time"
)

type App struct {
	conf Config

	updater    StatusUpdater
	lastfm     *lastfm.Api
	stdl, errl *log.Logger
}

func NewApp(configPath string, stdl, errl *log.Logger) (App, error) {
	a := App{}

	if stdl == nil || errl == nil {
		return a, ErrNilLoggerPtr
	}

	a.stdl = stdl
	a.errl = errl

	c, err := ParseConfig(configPath)

	if err != nil {
		return a, err
	}

	a.conf = c

	stdl.Println("Parsed config")

	{
		a.lastfm = lastfm.New(c.LastFM.APIKey, "")

		stdl.Println("Initialized last.fm api")
	}

	{
		var key string
		if c.Discord.UseAppMode {
			a.updater = NewAppStatusUpdater(&a)
			key = strconv.Itoa(c.Discord.AppID)
		} else {
			a.updater = NewTokenModeStatusUpdater(&a)
			key = c.Discord.Token
		}

		if err := a.updater.Login(key); err != nil {
			return a, err
		}

		stdl.Println("Connected to discord")
	}

	return a, nil
}

func (a App) Run() error {

	defer func() { a.updater.Logout() }()

	interval := time.Duration(a.conf.LastFM.CheckInterval) * time.Second
	mainTicker := time.NewTicker(interval)
	endless := a.conf.App.EndlessMode

	var (
		prevURL string
	)

	a.stdl.Println("Started!\n")

	for {
		select {
		case <-mainTicker.C:
			result, err := a.lastfm.User.GetRecentTracks(
				lastfm.P{
					"limit": "1",
					"user":  a.conf.LastFM.Username,
				},
			)

			if err != nil && !endless {
				return err
			}

			if len(result.Tracks) > 0 {

				ctrack := result.Tracks[0]
				isNowPlaying, err := strconv.ParseBool(ctrack.NowPlaying)

				if err != nil && ctrack.NowPlaying != "" {
					return ErrIncorrectNowPlaying
				}

				if isNowPlaying {
					if err := a.updater.Set(result); err != nil {

						if !endless {
							if err.Error() == "The pipe is being closed." {
								return ErrDiscordDisconnected
							}

							return err
						}

						a.errl.Println("Discord error: ", err)

					}

					if prevURL != ctrack.Url {
						a.stdl.Println("Now playing: " + ctrack.Artist.Name + " - " + ctrack.Name)
						prevURL = ctrack.Url
					}

				} else {
					if err := a.updater.Clear(); err != nil {

						if !endless {
							return err
						}

						a.errl.Println("Discord error: ", err)
					}
				}

			}
		}
	}
}
