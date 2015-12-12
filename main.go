package main

import (
	"os"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/jcooklin/snap-plugin-collector-cpuid/cpuid"
)

func main() {
	c := &cpuid.CPUID{}
	plugin.Start(
		cpuid.Meta(),
		c,
		os.Args[1],
	)
}
