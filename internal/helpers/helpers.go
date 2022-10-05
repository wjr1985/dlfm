package helpers

import (
	gap "github.com/muesli/go-app-paths"

	"fmt"
	"os"
	"strings"
)

func GetConfigPath() (string, error) {
	var confPath string

	pathScope := gap.NewScope(gap.User, "dlfm")
	dirs, _ := pathScope.ConfigDirs()

	// workdir
	dirs = append(dirs, ".")

	for _, d := range dirs {
		if _, err := os.Stat(d + string(os.PathSeparator) + "config.toml"); err == nil {
			confPath = d + string(os.PathSeparator) + "config.toml"
		}
	}

	if confPath == "" {
		return "", fmt.Errorf("Can't found config.toml in these directories: %s\n", strings.Join(dirs, ", "))
	}

	return confPath, nil
}
