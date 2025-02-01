package keygen

import (
	"crypto/rand"
	"encoding/base32"
	log "github.com/sirupsen/logrus"
)

const TARGET_KEY_BINARY_LENGTH = 20

func GenerateApiKey() string {
	bkey := make([]byte, TARGET_KEY_BINARY_LENGTH)
	n, err := rand.Read(bkey)
	if err != nil {
		log.Fatalln("GenerateApiKey: error creating random key: ", err)
	}
	if n != TARGET_KEY_BINARY_LENGTH {
		log.Fatalln("GenerateApiKey: cannor read required number of bytes: ", n, TARGET_KEY_BINARY_LENGTH)
	}

	return base32.StdEncoding.EncodeToString(bkey)
}
