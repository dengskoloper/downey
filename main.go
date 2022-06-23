package main
import (
	"github.com/dengskoloper/downey/cdm"
	"net/http"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"time"
	"os"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	LicenseServerURL string `long:"lic-server" description:"License Server URL"`
	InitPSSH string `long:"pssh-data" description:"Override PSSH Init data from MPD"`
}

func init() {
	widevine.InitConstants()
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	file, err := os.Open("manifest.mpd")
	if err != nil {
		panic(err)
	}

	initData, err := widevine.InitDataFromMPD(file)
	if err != nil {
		panic(err)
	}

	cdm, err := widevine.NewDefaultCDM(initData)
	if err != nil {
		panic(err)
	}
	licenseRequest, err := cdm.GetLicenseRequest()
	if err != nil {
		panic(err)
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	request, err := http.NewRequest(http.MethodPost, opts.LicenseServerURL, bytes.NewReader(licenseRequest))

	if err != nil {
		panic(err)
	}
	request.Close = true
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	licenseResponse, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	keys, err := cdm.GetLicenseKeys(licenseRequest, licenseResponse)
	if err != nil {
		panic(err)
	}

	command := ""
	for _, key := range keys {
		if key.Type == widevine.License_KeyContainer_CONTENT {
			command += "\n" + hex.EncodeToString(key.ID) + ":" + hex.EncodeToString(key.Value)
		}
	}
	fmt.Println("Decryption keys: ", command)
}