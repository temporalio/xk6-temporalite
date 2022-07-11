package temporalite

import (
	"context"

	"go.k6.io/k6/js/modules"

	"github.com/DataDog/temporalite"
	"github.com/temporalio/xk6-temporalite/metrics"
	"github.com/temporalio/xk6-temporalite/server"
)

func init() {
	modules.Register("k6/x/temporalite", new(RootModule))
}

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create `k6/x/temporal` module instances for each VU.
type RootModule struct{}

// ModuleInstance represents an instance of the module for every VU.
type ModuleInstance struct {
	vu            modules.VU
	customMetrics metrics.CustomMetrics
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Module   = &RootModule{}
	_ modules.Instance = &ModuleInstance{}
)

// NewModuleInstance implements the modules.Module interface to return
// a new instance for each VU.
func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu:            vu,
		customMetrics: metrics.RegisterServerMetrics(vu.InitEnv().Registry),
	}
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (temporal *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{Default: temporal}
}

// NewServer returns a new Temporalite Server.
func (m *ModuleInstance) NewTemporaliteServer() (*temporalite.Server, error) {
	metricsHandler := metrics.NewServerMetricsHandler(
		context.Background(),
		m.vu.State().Samples,
		m.vu.State().Tags.Clone(),
		m.customMetrics,
	)

	return server.NewTemporaliteServer(metricsHandler)
}
