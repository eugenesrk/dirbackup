package backup_client

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type BackupClient struct {
	directory string
	apiKey    string
	server    string
}

func NewBackupClient(directory, apiKey, server string) *BackupClient {
	bc := BackupClient{directory, apiKey, server}

	return &bc
}

func addToArchive(twr *tar.Writer, filename, nameInArchive string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	header.Name = nameInArchive

	err = twr.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(twr, file)
	if err != nil {
		return err
	}
	return nil
}

func (this *BackupClient) makeArchive() []byte {
	output := []byte{}
	outBuffer := bytes.NewBuffer(output)

	gzw := gzip.NewWriter(outBuffer)
	defer gzw.Close()
	twr := tar.NewWriter(gzw)
	defer twr.Close()

	log.Infoln("Walking: ", this.directory)
	err := filepath.Walk(this.directory, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			log.Warningln("INFO IS NILL!")
			return nil
		}
		if info.IsDir() {
			for _, exclusion := range exclusionDirs {
				if info.Name() == exclusion {
					log.Infoln("Skipping: ", path)
					return filepath.SkipDir
				}
			}
			// continue
			return nil
		}

		for _, exclusion := range exclusionExtensions {
			if strings.HasSuffix(info.Name(), exclusion) {
				return nil
			}
		}

		if info.Size() > maxFileSize {
			log.Infoln("Skipping HUGE file: ", path)
			return nil
		}

		log.Infoln("Adding: ", path)
		err = addToArchive(twr, path, strings.ReplaceAll(path, this.directory, ""))
		if err != nil {
			log.Fatalln("Cannot add a file to archive: ", err)
		}

		return nil
	})
	if err != nil {
		log.Fatalln("Unable to walk directory: ", err)
	}

	twr.Flush()
	gzw.Flush()

	if err = twr.Close(); err != nil {
		log.Fatalln("Unable to close archive1: ", err)
	}
	if err = gzw.Close(); err != nil {
		log.Fatalln("Unable to close archive2: ", err)
	}

	log.Warningln("Avail: ", outBuffer.Available())
	return outBuffer.Bytes()
}

func (this *BackupClient) CreateAndSend() {
	arch := this.makeArchive()

	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	err := mw.WriteField("api-key", this.apiKey)
	if err != nil {
		log.Fatalln("Cannot write a form field: ", err)
	}

	file, err := mw.CreateFormFile("upload", "upload.tgz")
	if err != nil {
		log.Fatalln("Cannot create a form field: ", err)
	}

	_, err = file.Write(arch)
	if err != nil {
		log.Fatalln("Cannot write a form field: ", err)
	}

	err = mw.Close()
	if err != nil {
		log.Fatalln("Cannot close a multipart writer: ", err)
	}

	req, err := http.NewRequest("PUT", this.server, &buf)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", mw.FormDataContentType())

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		bdat, _ := io.ReadAll(res.Body)
		log.Fatalln("Bad responce status code: ", res.Status, string(bdat))
	}
	log.Infoln("Backup uploaded!")
}
