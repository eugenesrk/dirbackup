package main

import (
	"dirbackup/internal/backup_server"
	"dirbackup/internal/keygen"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func main() {
	log.Infoln("Directory Backup Server is starting up!")

	apiKey := os.Getenv("API_KEY")
	if len(apiKey) < 5 {
		log.Warningln("OS ENVironment variable for API_KEY is empty!")
		log.Fatalln("Recommended key: ", keygen.GenerateApiKey())
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = "0.0.0.0:14800"
		log.Infoln("Listetning on default address of: ", listenAddr)
	}

	backupDir := os.Getenv("BACKUP_DIR")
	if len(backupDir) == 0 {
		log.Infoln("BACKUP_DIR ENVironment variable is empty, using default directory.")
		exePath, err := os.Executable()
		if err != nil {
			log.Fatalln("Cannot get EXE path: ", err)
		}
		exePath = filepath.Dir(exePath)
		backupDir = filepath.Join(exePath, "backups")

		if _, err := os.Stat(backupDir); err != nil {
			log.Infoln("Default Directory does not exist, attempting creation...")
			err = os.Mkdir(backupDir, 0700)
			if err != nil {
				log.Fatalln("Cannot create default backup directory: ", backupDir, err)
			}
		}

		log.Infoln("Running with backup directory: ", backupDir)

	} else {
		if _, err := os.Stat(backupDir); err != nil {
			log.Fatalln("BACKUP_DIR ENVironment variable contains a directory that does not exist: ", backupDir, err)
		}
	}

	svr := backup_server.NewBackupServer(apiKey, backupDir, listenAddr)
	svr.Start()
}
