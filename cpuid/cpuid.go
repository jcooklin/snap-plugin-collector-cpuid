package cpuid

import (
	"os"
	"time"

	"github.com/intel-go/cpuid"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
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
func (c *CPUID) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	hostname, _ := os.Hostname()
	for idx, m := range mts {
		switch m.Namespace_[len(m.Namespace_)-2] {
		case "avx":
			mts[idx].Data_ = cpuid.EnabledAVX
		case "avx512":
			mts[idx].Data_ = cpuid.EnabledAVX512
		}
		mts[idx].Timestamp_ = time.Now()
		mts[idx].Source_ = hostname
	}
	return mts, nil
}

// GetMetricTypes returns the metrics this plugin provides through an array of plugin.PluginMetricType
func (c *CPUID) GetMetricTypes(_ plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	return []plugin.PluginMetricType{
		plugin.PluginMetricType{
			Namespace_: []string{"jcooklin", "cpuid", "avx", "enabled"},
			Version_:   1,
		},
		plugin.PluginMetricType{
			Namespace_: []string{"jcooklin", "cpuid", "avx512", "enabled"},
			Version_:   1,
		},
	}, nil
}

// GetConfigPolicy returns the policy for this plugin
func (c *CPUID) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}
