package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"

	logger "github.com/sirupsen/logrus"
)

// Config struct stores config data for all app
type Config struct {
	DeploymentTemplatePath string `json:"deployment_template_path"`
	WaitForCreatingTimeout int    `json:"wait_pod_timeout"`
	PodLifetime            int    `json:"pod_lifetime"`
	Listen                 string `json:"listen"`
	Namespace              string `json:"namespace"`
	PodPort                string `json:"pod_port"`
	DeploymentTemplate     string
}

// Read reads config from yaml file
func (c *Config) Read(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("reading config " + path + ", " + err.Error())
	}

	err = json.Unmarshal(data, c)
	if err != nil {
		return errors.New("parsing config " + path + ", " + err.Error())
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return errors.New("error in get absolute path")
	}

	parentDir := filepath.Dir(absolutePath)
	c.DeploymentTemplatePath = parentDir + "/" + c.DeploymentTemplatePath

	data, err = ioutil.ReadFile(c.DeploymentTemplatePath)
	if err != nil {
		return errors.New("reading deployment template " + c.DeploymentTemplatePath + ", " + err.Error())
	}
	c.DeploymentTemplate = string(data)
	//TODO validate
	logger.Infof("config listing\n")
	logger.Infof("deployment template path: %s\n", c.DeploymentTemplatePath)
	logger.Infof("wait for creating timeout: %d\n", c.WaitForCreatingTimeout)
	logger.Infof("pod lifetime %d\n", c.PodLifetime)
	logger.Infof("listen: %s\n", c.Listen)
	logger.Infof("namespace: %s\n", c.Namespace)
	return nil
}
