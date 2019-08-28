package deploymentconfig

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// DeploymentConfig is alias for deployment data k8s client
type DeploymentConfig map[string]interface{}

// DeploymentFromTemplate creates deployment config with patterns
func DeploymentFromTemplate(template string) (d DeploymentConfig, err error) {
	rand.Seed(time.Now().UnixNano())
	// Craate patterns
	patterns := map[string]string{
		"%UNIQUE_ID%": strconv.Itoa(rand.Int()),
	}
	for key, value := range patterns {
		template = strings.ReplaceAll(template, key, value)
	}
	err = json.Unmarshal([]byte(template), &d)
	if err != nil {
		return nil, errors.New("processing deployment template, " + err.Error())
	}

	return d, nil
}
