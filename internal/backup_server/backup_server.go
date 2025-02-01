package backup_server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type BackupServer struct {
	apiKey    string
	directory string
	listen    string
}

func NewBackupServer(apiKey, directory, listen string) *BackupServer {
	bserv := BackupServer{
		apiKey:    apiKey,
		directory: directory,
		listen:    listen,
	}

	return &bserv
}

func (this *BackupServer) indexPage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusForbidden)
	writer.Write([]byte(`<!DOCTYPE html>
<html>
<head>
<title>ACCESS DENIED!</title>
</head>
<body>
<h1>403, Freeze!</h1>
</body>
</html>`))
}

func (this *BackupServer) Start() {
	http.HandleFunc("/", this.indexPage)
	http.HandleFunc("/upload-backup", this.processUpload)

	err := http.ListenAndServe(this.listen, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
