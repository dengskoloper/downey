package widevine

import "io/ioutil"

var DefaultPrivateKey string
var DefaultClientID []byte

func InitConstants() {
	DefaultPrivateKeyBuffer, err := ioutil.ReadFile("device_private_key")
	if err != nil {
		panic(err)
	}
	DefaultPrivateKey = string(DefaultPrivateKeyBuffer)

	DefaultClientID, err = ioutil.ReadFile("device_client_id_blob")
	if err != nil {
		panic(err)
	}
}