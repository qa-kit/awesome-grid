package creator

import (
	"errors"
	"net/http"
	"time"

	"github.com/qa-kit/awesome-grid/cleaner"
	"github.com/qa-kit/awesome-grid/cluster"
	"github.com/qa-kit/awesome-grid/config"
	"github.com/qa-kit/awesome-grid/deploymentconfig"
	poolPkg "github.com/qa-kit/awesome-grid/pool"
	logger "github.com/sirupsen/logrus"
)

// Creator provides creating logic of pods
type Creator struct {
	config      *config.Config
	cluster     cluster.Cluster
	pool        *poolPkg.Pool
	cleaner     *cleaner.Cleaner
	healthcheck func(string) (*http.Response, error)
}

//New created new Creator
func New(config *config.Config,
	cluster cluster.Cluster,
	pool *poolPkg.Pool,
	cleaner *cleaner.Cleaner,
	healthcheck func(string) (*http.Response, error)) *Creator {

	return &Creator{config, cluster, pool, cleaner, healthcheck}
}

// Resolve creates new pod and resolve to it
func (c *Creator) Resolve(request *http.Request) (string, error) {
	// Creating config
	deploymentConfig, err := deploymentconfig.DeploymentFromTemplate(c.config.DeploymentTemplate)
	if err != nil {
		return "", errors.New("creating deployment, " + err.Error())
	}

	// Creating deployemnt
	logger.Info("creating new deployment in k8s")
	name, err := c.cluster.CreateDeployment(deploymentConfig)
	if err != nil {
		return "", errors.New("creating deployment, " + err.Error())
	}
	logger.Infof("deployment %s created", name)

	// Delayed removing
	go c.cleaner.DeleteDeployment(name, c.cluster, c.config.PodLifetime, c.pool)

	// Waiting for pod craeting
	var IP string
	IP, err = c.waitForPod(name)

	if err != nil {
		return "", errors.New("waiting pod ip, " + err.Error())
	}

	logger.Infof("deployment %s ip is %s", name, IP)

	// Addind to pool
	c.pool.AddPod(name, IP)

	return "http://" + IP + ":" + c.config.PodPort, nil
}

// waitForPod waits until timeout or pod is active
func (c *Creator) waitForPod(name string) (string, error) {
	type Pair struct {
		IP  string
		Err error
	}

	c1 := make(chan Pair, 1)

	go func() {
		for {
			select {
			case <-time.After(time.Duration(c.config.WaitForCreatingTimeout) * time.Second):
				c1 <- Pair{"", errors.New("timeout error")}
				return
			default:
				IP, _ := c.cluster.FindPodIP(name)
				if IP != "" {
					resp, _ := c.healthcheck("http://" + IP + ":" + c.config.PodPort)
					if resp != nil {
						c1 <- Pair{IP, nil}
						return
					}
				}

				time.Sleep(1 * time.Second)
			}
		}
	}()

	select {
	case p := <-c1:
		return p.IP, p.Err
	}
}
