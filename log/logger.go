package log

import (
	"os"
	"fmt"
	"path/filepath"
	"go.uber.org/zap"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func GetLogger() *zap.Logger {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binDir := filepath.Dir(path)
	configYaml, err := ioutil.ReadFile(fmt.Sprintf("%s/data/config.yaml", binDir))
	if err != nil {
		panic(err)
	}
	var config zap.Config
	if err := yaml.Unmarshal(configYaml, &config); err != nil {
		panic(err)
	}
	logger, _ := config.Build()

	return logger
}
