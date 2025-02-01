package backup_client

import (
	"os"
	"strings"
)

var exclusionDirs = []string{
	"node_modules",
	".cache",
}

var exclusionExtensions = []string{
	".bin",
	".db3",
	".xz",
	".gz",
	".bak",
}

var maxFileSize = int64(1 * 1024 * 1024 * 1024)

func init() {
	exclDirVar := os.Getenv("EXCLUDE_DIRS")
	if len(exclDirVar) > 0 {
		exclusionDirs = strings.Split(exclDirVar, ":")
	}

	exclExtensionVar := os.Getenv("EXCLUDE_EXTENSIONS")
	if len(exclExtensionVar) > 0 {
		exclusionExtensions = strings.Split(exclExtensionVar, ":")
	}
}
