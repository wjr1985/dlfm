package main

import (
	dgo "github.com/bwmarrin/discordgo"
	rgo "github.com/dikey0ficial/rich-go/v2/client" // my fork with some fixes
	"github.com/shkh/lastfm-go/lastfm"
	"strings"
)

type RT = lastfm.UserGetRecentTracks

type StatusUpdater interface {
	Login(string) error
	Set(RT) error
	Clear() error
}

type AppStatusUpdater struct{}

func (AppStatusUpdater) Login(id string) error {
	return rgo.Login(id)
}

func (AppStatusUpdater) Set(t RT) error {
	ctrack := t.Tracks[0]
	fltext, fstext := conf.App.LargeText,
		conf.App.SmallText
	for _, v := range []*string{&fltext, &fstext} {
		*v = strings.Replace(*v, "{{name}}", ctrack.Name, -1)
		*v = strings.Replace(*v, "{{artist}}", ctrack.Artist.Name, -1)
		*v = strings.Replace(*v, "{{album}}", ctrack.Album.Name, -1)
	}
	flimg := conf.App.LargeImage
	flimg = strings.Replace(flimg, "{{album_image}}", ctrack.Images[3].Url, -1)
	var bs = make([]*rgo.Button, 0)
	if conf.App.ShowButton {
		bs = []*rgo.Button{&rgo.Button{
			Label: "This track on last.fm",
			URL:   ctrack.Url,
		}}
	}
	return rgo.SetActivity(
		rgo.Activity{
			Details:    ctrack.Name,
			State:      ctrack.Artist.Name,
			LargeImage: flimg,
			LargeText:  fltext,
			SmallImage: conf.App.SmallImage,
			SmallText:  fstext,
			Buttons:    bs,
		},
	)
}

func (AppStatusUpdater) Clear() error {
	// rgo.Logout()
	return nil
}

type TokenModeStatusUpdater struct {
	Session dgo.Session
}

func (tmsu TokenModeStatusUpdater) Login(token string) error {
	dg, err := dgo.New(conf.Discord.Token)
	if err != nil {
		return err
	}
	tmsu.Session = *dg
	err = tmsu.Session.Open()
	return err
}

func (tmsu TokenModeStatusUpdater) Set(t RT) error {
	ctrack := t.Tracks[0]
	ftitle := conf.App.Title
	ftitle = strings.Replace(ftitle, "{{name}}", ctrack.Name, -1)
	ftitle = strings.Replace(ftitle, "{{artist}}", ctrack.Artist.Name, -1)
	ftitle = strings.Replace(ftitle, "{{album}}", ctrack.Album.Name, -1)
	return tmsu.Session.UpdateStatusComplex(
		dgo.UpdateStatusData{
			Game: &dgo.Game{
				Name:    ftitle,
				Type:    2,
				Details: ctrack.Name,
				State:   ctrack.Artist.Name,
			},
		},
	)
}

func (tmsu TokenModeStatusUpdater) Clear() error {
	return tmsu.Session.UpdateStatusComplex(dgo.UpdateStatusData{Game: nil})
}
