package creator

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/qa-kit/awesome-grid/cleaner"
	"github.com/qa-kit/awesome-grid/cluster"
	"github.com/qa-kit/awesome-grid/config"
	"github.com/qa-kit/awesome-grid/deploymentconfig"
	poolPkg "github.com/qa-kit/awesome-grid/pool"
)

// Creator provides creating logic of pods
type Creator struct {
	config  *config.Config
	cluster cluster.Cluster
	pool    *poolPkg.Pool
	cleaner *cleaner.Cleaner
}

//New created new Creator
func New(config *config.Config,
	cluster cluster.Cluster,
	pool *poolPkg.Pool,
	cleaner *cleaner.Cleaner) *Creator {
	return &Creator{config, cluster, pool, cleaner}
}

// Resolve creates new pod and resolve to it
func (c *Creator) Resolve(request *http.Request) (string, error) {
	// Creating config
	deploymentConfig, err := deploymentconfig.DeploymentFromTemplate(c.config.DeploymentTemplate)
	if err != nil {
		return "", errors.New("creating deployment, " + err.Error())
	}

	// Creating deployemnt
	log.Println("creating new deployment in k8s")
	name, err := c.cluster.CreateDeployment(deploymentConfig)
	if err != nil {
		return "", errors.New("creating deployment, " + err.Error())
	}
	log.Println("deployment " + name + " created")

	// Delayed removing
	go c.cleaner.DeleteDeployment(name, c.cluster, c.config.PodLifetime, c.pool)
	time.Sleep(time.Duration(c.config.WaitForCreatingTimeout) * 1000 * time.Millisecond)

	// Getting ip of new pod
	IP, err := c.cluster.FindPodIP(name)
	if err != nil {
		return "", errors.New("finding pod ip, " + err.Error())
	}

	log.Println("deployment " + name + " ip is " + IP)

	// Addind to pool
	c.pool.AddPod(name, IP)

	return "http://" + IP + ":" + c.config.PodPort, nil
}
