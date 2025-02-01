package main

import (
	"dirbackup/internal/backup_client"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.Infoln("Directory backup client is started")

	apiKey := os.Getenv("API_KEY")
	if len(apiKey) == 0 {
		log.Fatalln("OS ENVironment variable for API_KEY is empty!")
	}

	srvUrl := os.Getenv("SERVER")
	if len(srvUrl) == 0 {
		log.Fatalln("OS ENVironment variable for SERVER is empty!")
	}

	backupDir := os.Getenv("BACKUP_DIR")
	if len(backupDir) == 0 {
		log.Fatalln("OS ENVironment variable for BACKUP_DIR is empty!")
	} else if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		log.Fatalln("BACKUP_DIR does not exist!")
	}

	client := backup_client.NewBackupClient(backupDir, apiKey, srvUrl)
	client.CreateAndSend()
}
