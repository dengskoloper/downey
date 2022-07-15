package main
import (
	"github.com/dengskoloper/downey/cdm"
	"github.com/dengskoloper/downey/util"
	"net/http"
	"bytes"
	"encoding/hex"
	"encoding/base64"
	"fmt"
	"io"
	"time"
	"os"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	LicenseServerURL string `long:"lic-server" description:"License Server URL"`
	AddHeaders bool `long:"add-headers" description:"Read HTTP headers from headers.json"`
	InitPSSH string `long:"pssh" description:"Override PSSH data"`
}

func init() {
	widevine.InitConstants()
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	var initData []byte
	if len(opts.InitPSSH) > 0 {
		initData, err = base64.StdEncoding.DecodeString(opts.InitPSSH)
		if err != nil {
			panic(err)
		}
	} else {
		file, err := os.Open("manifest.mpd")
		if err != nil {
			panic(err)
		}

		initData, err = widevine.InitDataFromMPD(file)
		if err != nil {
			panic(err)
		}
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

	if opts.AddHeaders {
		util.ReadHeadersFromJSON(&request.Header) 
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
	fmt.Println("\nDecryption keys: ", command)
}