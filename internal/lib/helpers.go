package lib

import (
	"strings"
)

func multiStringReplace(origin string, replaces map[string]string) string {
	var result string = origin
	for temp, rep := range replaces {
		result = strings.Replace(result, temp, rep, -1)
	}
	return result
}

func replaceTags(original string, name, artist, album, albumCover string) string {
	var result string
	result = multiStringReplace(
		original,
		map[string]string{
			"{{name}}":        name,
			"{{artist}}":      artist,
			"{{album}}":       album,
			"{{album_image}}": albumCover,

			// small image tags
			"{{lastfm}}": "http://icons.iconarchive.com/icons/danleech/simple/512/lastfm-icon.png",
			"{{deezer}}": "https://www.macupdate.com/images/icons512/60905.png",
			"{{youtube}}": "https://seeklogo.com/images/Y/youtube-music-logo-50422973B2-seeklogo.com.png",
			"{{apple}}": "http://ixd.prattsi.org/wp-content/uploads/2017/01/apple_music_logo_by_mattroxzworld-d982zrj.png",
			"{{vk}}": "https://seeklogo.com/images/V/vk-icon-logo-10188561D5-seeklogo.com.png",
			"{{yandex}}": "https://download.cdn.yandex.net/from/yandex.ru/support/ru/music/files/icon_main.png",
			"{{soundcloud}}": "https://icons.iconarchive.com/icons/sicons/basic-round-social/512/soundcloud-icon.png",
		},
	)
	return result
}
