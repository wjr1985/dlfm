package main

import (
	dgo "github.com/bwmarrin/discordgo"
	rgo "github.com/dikey0ficial/rich-go/v2/client" // my fork with some fixes
	"github.com/shkh/lastfm-go/lastfm"
	"strings"
)

func replaceTags(original string, name, artist, album, albumCoverURL string) string {
	var result string
	result = strings.Replace(original, "{{name}}", name, -1)
	result = strings.Replace(result, "{{artist}}", artist, -1)
	result = strings.Replace(result, "{{album}}", album, -1)
	result = strings.Replace(result, "{{album_image}}", albumCoverURL, -1)
	return result
}

type RT = lastfm.UserGetRecentTracks

type StatusUpdater interface {
	Login(string) error
	Logout() error
	Set(RT) error
	Clear() error
}

type AppStatusUpdater struct{}

func (AppStatusUpdater) Login(id string) error {
	return rgo.Login(id)
}

func (AppStatusUpdater) Logout() error {
	rgo.Logout()
	return nil
}

func (AppStatusUpdater) Set(t RT) error {
	ctrack := t.Tracks[0]
	ffirstl, fsecline := conf.App.FirstLine, conf.App.SecondLine
	fltext, fstext := conf.App.LargeText, conf.App.SmallText
	flimg := conf.App.LargeImage
	for _, v := range []*string{&fltext, &fstext, &ffirstl, &fsecline, &flimg} {
		*v = replaceTags(*v, ctrack.Name, ctrack.Artist.Name, ctrack.Album.Name, ctrack.Images[3].Url)
	}
	var bs = make([]*rgo.Button, 0)
	if conf.App.ShowButton {
		bs = []*rgo.Button{&rgo.Button{
			Label: "This track on last.fm",
			URL:   ctrack.Url,
		}}
	}
	return rgo.SetActivity(
		rgo.Activity{
			Details:    ffirstl,
			State:      fsecline,
			LargeImage: flimg,
			LargeText:  fltext,
			SmallImage: conf.App.SmallImage,
			SmallText:  fstext,
			Buttons:    bs,
		},
	)
}

func (AppStatusUpdater) Clear() error {
	return nil
}

type TokenModeStatusUpdater struct {
	Session *dgo.Session
}

func (tmsu *TokenModeStatusUpdater) Login(token string) error {
	dg, err := dgo.New(conf.Discord.Token)
	if err != nil {
		return err
	}
	err = dg.Open()
	if err != nil {
		return err
	}
	tmsu.Session = dg
	return nil
}

func (tmsu *TokenModeStatusUpdater) Logout() error {
	if tmsu.Session == nil {
		return ErrNilDGoSession
	}
	return tmsu.Session.Close()
}

func (tmsu *TokenModeStatusUpdater) Set(t RT) error {
	ctrack := t.Tracks[0]
	ftitle := conf.App.Title
	ffirstl, fsecline := conf.App.FirstLine, conf.App.SecondLine
	for _, v := range []*string{&ftitle, &ffirstl, &fsecline} {
		*v = replaceTags(*v, ctrack.Name, ctrack.Artist.Name, ctrack.Album.Name, ctrack.Images[3].Url)
	}
	if tmsu.Session == nil {
		return ErrNilDGoSession
	}
	return tmsu.Session.UpdateStatusComplex(
		dgo.UpdateStatusData{
			Game: &dgo.Game{
				Name:    ftitle,
				Type:    2,
				Details: ffirstl,
				State:   fsecline,
			},
		},
	)
}

func (tmsu *TokenModeStatusUpdater) Clear() error {
	if tmsu.Session == nil {
		return ErrNilDGoSession
	}
	return tmsu.Session.UpdateStatusComplex(dgo.UpdateStatusData{Game: nil})
}
