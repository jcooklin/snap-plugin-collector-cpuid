
## Creating a Go based plugin for [snap](https://github.com/intelsdi-x/snap) through a stripped down example  

1.  Choose a repository and plugin name
  *  It is suggested that your name include the following convention snap-plugin-[*plugin-type*]-[*plugin-name*]
  *  For example: snap-plugin-collector-foo
  *  Plugin types: collector, processor, publisher
 
2. Implement the [CollectorPlugin](https://github.com/intelsdi-x/snap/blob/master/control/plugin/collector.go#L26-L30) interface (cpuid.go)
 ```go
// Collector plugin
type CollectorPlugin interface {
	Plugin
	CollectMetrics([]PluginMetricType) ([]PluginMetricType, error)
	GetMetricTypes(PluginConfigType) ([]PluginMetricType, error)
}```
```go
type Plugin interface {
	GetConfigPolicy() (*cpolicy.ConfigPolicy, error)
}
```
  * ```GetMetricTypes([]PluginMetricType) ([]PluginMetricType, error)``` communicates what metrics
 the plugin collects through returning an array of [plugin.PluginMetricType](https://github.com/intelsdi-x/snap/blob/master/control/plugin/metric.go#L90-L118)
    * Through the `plugin.PluginMetricType` we are providing the `namespace` and `version` of the metric(s) our plugin will provide    
 ```go
 // GetMetricTypes returns the metrics this plugin provides through an array of plugin.PluginMetricType
func (c *CPUID) GetMetricTypes(_ plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	return []plugin.PluginMetricType{
		plugin.PluginMetricType{
			Namespace_: []string{"jcooklin", "cpuid", "VendorIdentificationString"},
			Version_: 1,
		},
	}, nil
}
```
  
  * ```CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error)``` returns the requested metrics
    * The argument to `CollectMetrics` is an array of `plugin.PluginMetricType` where the `Namespace` field communicates 
	what metric is being requested
	* The result of `CollectMetrics` is also an array of `plugin.PluginMetricType` where the fields `Data`, `TimeStamp` and `Source`
	communicate the value, time of collection and source of the metric
	
	```go
	// CollectMetrics returns []plugin.PluginMetrictType
func (c *CPUID) CollectMetrics(mts []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	hostname, _ := os.Hostname()
	for idx, m := range mts {
		switch m.Namespace_[len(m.Namespace_)-1] {
		case "VendorIdentificatorString":
			mts[idx].Data_ = cpuid.VendorIdentificatorString
			mts[idx].Timestamp_ = time.Now()
			mts[idx].Source_ = hostname
		}
	}
	return mts, nil
}
	```
	
  * `GetConfigPolicy() (*cpolicy.ConfigPolicy, error)` returns the required config policy, if any, for the metrics provided by the plugin
    * This plugin does not require a config policy so an empty cpolicy is returned
	```go
	// GetConfigPolicy returns the policy for this plugin
func (c *CPUID) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}
	```
	* *Note: a config policy provides the plugin author a mechanism which forces the users of the plugin to provide configuration detail
	such as username, password, address, port, etc when their plugin is used*

* Start the plugin (main.go)
  * `plugin.Start(m *PluginMeta, c Plugin, requestString string)` starts the plugin
  ```go
    func main() {
	c := &cpuid.CPUID{}
	plugin.Start(
		cpuid.Meta(),
		c,
		os.Args[1],
	)
}
``` 
    * `cpuid.Meta()` returns `*PluginMeta`    	   
  	```go
	func Meta() *plugin.PluginMeta {
		return plugin.NewPluginMeta(
			name,
			version,
			pluginType,
			[]string{plugin.SnapGOBContentType},
			[]string{plugin.SnapGOBContentType})
	}
```  

### Download, build and use the plugin

* Prerequisite
  * Go is properly [installed](https://golang.org/doc/install) and configured
  * Snap binaries `snapd`,  `snapctl` and the plugin file publisher plugin are available
    * Download the [latest snap release](https://github.com/intelsdi-x/snap/releases) or build [snap]((https://github.com/intelsdi-x/snap) yourself	

  
1. Use 'go get' to get the plugin
```
go get github.com/jcooklin/snap-plugin-collector-cpuid
```
*Note* The rest of the commands below are run from the root of the snap-plugin-collector-cpuid which should now be located at $GOPATH/src/github.com/jcooklin/

* Build the plugin
```
make
```

* Start snapd
```
snapd -t 0 -l 1
```
  * `-t 0` disables plugin trust
  * `-l 1` starts snapd with debug logging

* Load plugins
```
snapctl plugin load build/rootfs/snap-plugin-collector-cpuid
snapctl plugin load ../../intelsdi-x/snap/build/plugin/snap-publisher-file
``` 
*Note*: You may need to reference a different path for the snap-publisher-file.  For instance, if you downloaded a [release](https://github.com/intelsdi-x/snap/releases). 
* Start task 
```
snapctl task create -t example/task.json
```

* Watch task

```
snapctl task watch <TASK_ID>
```

![img](http://i.giphy.com/26tP7WQsLepN4Rooo.gif) 