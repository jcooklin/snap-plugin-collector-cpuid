package cpuid

import (
	"time"

	"github.com/intel-go/cpuid"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
)

const (
	// Name of plugin
	name = "cpuid"
	// Version of plugin
	version = 1
	// Type of plugin
	pluginType = plugin.CollectorPluginType
)

type CPUID struct{}

// Meta returns plugin meta which is sent to snapd
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		name,
		version,
		pluginType,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType})

}

// CollectMetrics returns []plugin.PluginMetrictType
func (c *CPUID) CollectMetrics(mts []plugin.MetricType) ([]plugin.MetricType, error) {
	for idx, m := range mts {
		switch m.Namespace().Strings()[len(m.Namespace())-2] {
		case "avx":
			mts[idx].Data_ = cpuid.EnabledAVX
		case "avx512":
			mts[idx].Data_ = cpuid.EnabledAVX512
		}
		mts[idx].Timestamp_ = time.Now()
	}
	return mts, nil
}

// GetMetricTypes returns the metrics this plugin provides through an array of plugin.PluginMetricType
func (c *CPUID) GetMetricTypes(_ plugin.ConfigType) ([]plugin.MetricType, error) {
	return []plugin.MetricType{
		plugin.MetricType{
			Namespace_: core.NewNamespace("jcooklin", "cpuid", "avx", "enabled"),
			Version_:   1,
		},
		plugin.MetricType{
			Namespace_: core.NewNamespace("jcooklin", "cpuid", "avx512", "enabled"),
			Version_:   1,
		},
	}, nil
}

// GetConfigPolicy returns the policy for this plugin
func (c *CPUID) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}
