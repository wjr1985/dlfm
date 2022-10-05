package lib

import (
	"github.com/BurntSushi/toml"

	"io"
	"io/ioutil"
	"os"
	"strings"
)

// ===================== Structs =====================

type Config struct {
	LastFM  LastFMConfig  `toml:"lastfm"`
	Discord DiscordConfig `toml:"discord"`
	App     AppConfig     `toml:"app"`
}

type LastFMConfig struct {
	APIKey        string `toml:"api_key"`
	Username      string `toml:"username"`
	CheckInterval int    `toml:"check_interval"`
}

type DiscordConfig struct {
	UseAppMode bool   `toml:"use_app_mode"`
	Token      string `toml:"token"`
	AppID      int    `toml:"app_id"`
}

type AppConfig struct {
	Title       string `toml:"title"`
	FirstLine   string `toml:"first_line"`
	SecondLine  string `toml:"second_line"`
	LargeImage  string `toml:"large_image"`
	LargeText   string `toml:"large_text"`
	SmallImage  string `toml:"small_image"`
	SmallText   string `toml:"small_text"`
	ShowButton  bool   `toml:"show_button"`
	EndlessMode bool   `toml:"endless_mode"`
}

// ===================== Structs =====================

// ==================== Variables ====================

var defaultConfig = Config{
	LastFM: LastFMConfig{
		CheckInterval: 5,
	},
	Discord: DiscordConfig{
		UseAppMode: true,
	},
	App: AppConfig{
		EndlessMode: true,

		Title:      "last.fm",
		FirstLine:  "{{name}}",
		SecondLine: "{{artist}}",
		LargeImage: "{{album_image}}",
		LargeText:  "{{album}}",
		SmallImage: "{{lastfm}}",
		SmallText:  "github.com/dikey0ficial/dlfm",
	},
}

var servicesTags = map[string]string{

	// !!WARNING!!: I did't checked licenses of these icons

	"lastfm":     "http://icons.iconarchive.com/icons/danleech/simple/512/LastFM-icon.png",
	"deezer":     "https://www.macupdate.com/images/icons512/60905.png",
	"youtube":    "https://seeklogo.com/images/Y/youtube-music-logo-50422973B2-seeklogo.com.png",
	"apple":      "http://ixd.prattsi.org/wp-content/uploads/2017/01/apple_music_logo_by_mattroxzworld-d982zrj.png",
	"vk":         "https://seeklogo.com/images/V/vk-icon-logo-10188561D5-seeklogo.com.png",
	"yandex":     "https://download.cdn.yandex.net/from/yandex.ru/support/ru/music/files/icon_main.png",
	"soundcloud": "https://icons.iconarchive.com/icons/sicons/basic-round-social/512/soundcloud-icon.png",
}

// ==================== Variables ====================

func ParseConfig(path string) (Config, error) {
	var conf Config = defaultConfig

	f, err := os.Open(path)
	if err != nil {
		return conf, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil && err != io.EOF {
		return conf, err
	}
	if err := toml.Unmarshal(data, &conf); err != nil {
		return conf, err
	}

	var (
		strs = []*string{
			&conf.LastFM.APIKey, &conf.LastFM.Username,
			&conf.Discord.Token, &conf.App.Title,
			&conf.App.LargeImage, &conf.App.LargeText,
			&conf.App.SmallImage, &conf.App.SmallText,
		}
	)
	for _, v := range strs {
		*v = strings.TrimSpace(*v)
	}

	if conf.LastFM.APIKey == "" {

		return conf, ErrEmptyField("lastfm.api_key")
	} else if conf.LastFM.Username == "" {
		return conf, ErrEmptyField("lastfm.username")
	} else if conf.Discord.Token == "" && !conf.Discord.UseAppMode {
		return conf, ErrEmptyField("discord.token")
	} else if conf.Discord.AppID == 0 && conf.Discord.UseAppMode {
		return conf, ErrEmptyField("discord.app_id")
	} else if conf.LastFM.CheckInterval < 0 {
		// why < 0 instead of 1? because zero is Grindcore mode!!!!!!
		return conf, ErrEmptyField("Error: invalid check_interval (< 0)")
	}

	for name, url := range servicesTags {
		if conf.App.SmallImage == "{{"+name+"}}" {
			conf.App.SmallImage = url
			break
		}
	}

	return conf, nil
}
