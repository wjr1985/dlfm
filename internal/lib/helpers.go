package lib

import (
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
