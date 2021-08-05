package helpers

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
	"github.com/savannahghi/interserviceclient"
	"gopkg.in/yaml.v2"
)

// InitializeInterServiceClient initializes an external service in the correct environment given its name
func InitializeInterServiceClient(serviceName string) *interserviceclient.InterServiceClient {
	//os file and parse it to go type
	file, err := ioutil.ReadFile(filepath.Clean(interserviceclient.PathToDepsFile()))
	if err != nil {
		log.Errorf("error occurred while opening deps file %v", err)
		os.Exit(1)
	}
	var config interserviceclient.DepsConfig
	if err := yaml.Unmarshal(file, &config); err != nil {
		log.Errorf("failed to unmarshal yaml config file %v", err)
		os.Exit(1)
	}

	client, err := interserviceclient.SetupISCclient(config, serviceName)
	if err != nil {
		log.Panicf("unable to initialize inter service client for %v service: %s", err, serviceName)
	}
	return client
}
