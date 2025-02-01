package backup_server

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

func (this *BackupServer) executeCleanup() {
	// delay to let the caller close/flush the file
	time.Sleep(time.Second * 5)

	log.Infoln("Executing backup cleanup.")

	files, err := filepath.Glob(filepath.Join(this.directory, "backup_*.tgz"))
	if err != nil {
		log.Errorln("Cannot GLOB: ", this.directory, err)
		return
	}

	treshold := time.Now().Add(-7 * time.Hour * 24)

	for _, file := range files {
		fstat, err := os.Stat(file)
		if err != nil {
			log.Errorln("Cannot stat file: ", file, err)
			continue
		}

		if fstat.ModTime().Before(treshold) {
			log.Warningln("Removing an old file: ", filepath.Base(file), fstat.ModTime().Format(time.DateTime))
			err = os.Remove(file)
			if err != nil {
				log.Errorln("Cannot remove a file: ", filepath.Base(file), err)
				continue
			}
		}
	}

	log.Infoln("Backup cleanup done!")
}
