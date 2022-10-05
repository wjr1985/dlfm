package lib

import (
	dgo "github.com/bwmarrin/discordgo"
	rgo "github.com/dikey0ficial/rich-go/v2/client" // my fork with some fixes
	"github.com/shkh/lastfm-go/lastfm"
)

// ==================== Typedefs ====================

type RT = lastfm.UserGetRecentTracks

type StatusUpdater interface {
	Login(string) error
	Logout() error
	Set(RT) error
	Clear() error
}

// ==================== Typedefs ====================

type AppStatusUpdater struct {
	app *App
}

// ==================== Methods ====================

func NewAppStatusUpdater(app *App) *AppStatusUpdater {
	return &AppStatusUpdater{
		app: app,
	}
}

func (*AppStatusUpdater) Login(id string) error {
	return rgo.Login(id)
}

func (*AppStatusUpdater) Logout() error {
	rgo.Logout()
	return nil
}

func (a *AppStatusUpdater) Set(t RT) error {
	conf := a.app.conf

	var act = rgo.Activity{
		Details:    conf.App.FirstLine,
		State:      conf.App.SecondLine,
		LargeImage: conf.App.LargeImage,
		LargeText:  conf.App.LargeText,
		SmallImage: conf.App.SmallImage,
		SmallText:  conf.App.SmallText,
		Buttons:    []*rgo.Button{},
	}

	ctrack := t.Tracks[0]

	for _, v := range []*string{
		&act.Details,
		&act.State,
		&act.LargeImage,
		&act.LargeText,
		&act.SmallText,
	} {
		*v = replaceTags(
			*v,
			ctrack.Name,
			ctrack.Artist.Name,
			ctrack.Album.Name,
			ctrack.Images[3].Url,
		)
	}

	if conf.App.ShowButton {
		act.Buttons = append(act.Buttons, &rgo.Button{
			Label: "This track on last.fm",
			URL:   ctrack.Url,
		})
	}

	return rgo.SetActivity(act)
}

func (*AppStatusUpdater) Clear() error {
	return nil
}

// ==================== Methods ====================

type TokenModeStatusUpdater struct {
	app *App

	Session *dgo.Session
}

// ==================== Methods ====================

func NewTokenModeStatusUpdater(app *App) *TokenModeStatusUpdater {
	return &TokenModeStatusUpdater{
		app: app,
	}
}

func (tmsu *TokenModeStatusUpdater) Login(token string) error {
	dg, err := dgo.New(tmsu.app.conf.Discord.Token)

	if err != nil {
		return err
	}

	if err := dg.Open(); err != nil {
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
	conf := tmsu.app.conf

	ctrack := t.Tracks[0]

	act := dgo.Game{
		Name:    conf.App.Title,
		Type:    2,
		Details: conf.App.FirstLine,
		State:   conf.App.SecondLine,
	}

	for _, v := range []*string{
		&act.Name,
		&act.Details,
		&act.State,
	} {
		*v = replaceTags(
			*v,
			ctrack.Name,
			ctrack.Artist.Name,
			ctrack.Album.Name,
			ctrack.Images[3].Url,
		)
	}

	if tmsu.Session == nil {
		return ErrNilDGoSession
	}

	return tmsu.Session.UpdateStatusComplex(
		dgo.UpdateStatusData{Game: &act},
	)
}

func (tmsu *TokenModeStatusUpdater) Clear() error {

	if tmsu.Session == nil {
		return ErrNilDGoSession
	}

	return tmsu.Session.UpdateStatusComplex(dgo.UpdateStatusData{Game: nil})
}

// ==================== Methods ====================
