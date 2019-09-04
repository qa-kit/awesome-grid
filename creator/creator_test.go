package creator

import (
	"net/http"
	"testing"

	"github.com/qa-kit/awesome-grid/cleaner"
	"github.com/qa-kit/awesome-grid/cluster"
	"github.com/qa-kit/awesome-grid/config"
	"github.com/qa-kit/awesome-grid/pool"
)

func TestCreatorResolve(t *testing.T) {
	expected := "http://127.0.0.1:4444"
	creator := Creator{
		config: &config.Config{
			WaitForCreatingTimeout: 0,
			PodLifetime:            0,
			DeploymentTemplate:     "{}",
			PodPort:                "4444",
		},
		cluster: cluster.FakeKubernetes{
			FindPodIPResult: "127.0.0.1",
		},
		pool:    &pool.Pool{},
		cleaner: &cleaner.Cleaner{},
		healthcheck: func(string) (*http.Response, error) {
			return &http.Response{}, nil
		},
	}
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("creating http request error %s", err.Error())
		return
	}
	IP, err := creator.Resolve(request)
	if err != nil {
		t.Errorf("expected empty error instead '%s'", err.Error())
	}
	if IP != expected {
		t.Errorf("expected ip '%s' instead '%s'", IP, expected)
	}
}
