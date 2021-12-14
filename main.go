package main

import (
	"fmt"
	"os"
	"packer-plugin-manifestory/builder/manifestory"
	manifestoryData "packer-plugin-manifestory/datasource/manifestory"
	manifestoryPP "packer-plugin-manifestory/post-processor/manifestory"
	manifestoryProv "packer-plugin-manifestory/provisioner/manifestory"
	manifestoryVersion "packer-plugin-manifestory/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("builder", new(manifestory.Builder))
	pps.RegisterProvisioner("provisioner", new(manifestoryProv.Provisioner))
	pps.RegisterPostProcessor("post-processor", new(manifestoryPP.PostProcessor))
	pps.RegisterDatasource("datasource", new(manifestoryData.Datasource))
	pps.SetVersion(manifestoryVersion.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
