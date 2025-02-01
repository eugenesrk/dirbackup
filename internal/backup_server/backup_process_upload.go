package backup_server

import (
	"dirbackup/internal/tools"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (this *BackupServer) processUpload(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "PUT" {
		tools.ShowError(writer, "Invalid method!")
		return
	}

	err := request.ParseMultipartForm(1024 * 1024 * 500) // 500 MB to begin with
	if err != nil {
		log.Errorln("Cannot parse the multipart/form-data request: ", err)
		tools.ShowError(writer, "Invalid request, see log.")
		return
	}

	/*
		Storing a key without hashing is a bad practice, however this is not a password,
		but a random key and it needs to be changed anyway if an attacker was able to
		get access necessary to retreive it.
	*/
	if apik, ok := request.MultipartForm.Value["api-key"]; !ok {
		tools.ShowError(writer, "Invalid authentication!")
		return
	} else if len(apik) != 1 {
		tools.ShowError(writer, "Invalid authentication!")
		return
	} else if !tools.ConstantTimeCompare(apik[0], this.apiKey) {
		tools.ShowError(writer, "Invalid authentication!")
		return
	}

	var inFileHeader *multipart.FileHeader

	if inFileHeaders, ok := request.MultipartForm.File["upload"]; !ok {
		tools.ShowError(writer, "No file provided!")
		return
	} else if len(inFileHeaders) != 1 {
		tools.ShowError(writer, "No/too much file(s) provided!")
		return
	} else {
		inFileHeader = inFileHeaders[0]
	}

	// generate an ouput name
	outputName := filepath.Join(
		this.directory,
		fmt.Sprintf("backup_%s.tgz", time.Now().Format("2006-01-02")))

	if outFile, err := os.Stat(outputName); err == nil || (outFile != nil && outFile.Size() > 0) {
		tools.ShowError(writer, "Rate limit exceeded!")
		return
	}

	outFile, err := os.OpenFile(outputName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Errorln("Cannot open file for writing:", outFile, err)
		tools.ShowError(writer, "Internal error, see log.")
		return
	}
	defer outFile.Close()

	inFile, err := inFileHeader.Open()
	if err != nil {
		log.Errorln("Cannot open inFileHeader for reading:", outFile, err)
		tools.ShowError(writer, "Internal error, see log.")
		return
	}

	_, err = io.Copy(outFile, inFile)
	if err != nil {
		log.Errorln("Cannot COPY input file to destination:", outFile, err)
		tools.ShowError(writer, "Internal error, see log.")
		return
	}
	outFile.Sync()

	go this.executeCleanup()

	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.Write([]byte("OK"))
}
