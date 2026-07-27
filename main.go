package main

import (
	v100 "bk-plugin-demo-go/bk_plugin/v100"
	v101 "bk-plugin-demo-go/bk_plugin/v101"
	v102 "bk-plugin-demo-go/bk_plugin/v102"
	v103 "bk-plugin-demo-go/bk_plugin/v103"
	_ "bk-plugin-demo-go/data_api"

	"github.com/TencentBlueKing/beego-runtime/runner"
	"github.com/TencentBlueKing/bk-plugin-framework-go/hub"
)

func registerPlugins() {
	hub.MustInstall(&v100.Plugin{}, v100.ContextInputs{}, v100.Outputs{}, v100.InputsFormJSON)
	hub.MustInstall(&v101.Plugin{}, v101.ContextInputs{}, v101.Outputs{}, v101.InputsFormJSON)
	hub.MustInstall(&v102.Plugin{}, v102.ContextInputs{}, v102.Outputs{}, v102.InputsFormJSON)
	hub.MustInstall(&v103.Plugin{}, v103.ContextInputs{}, v103.Outputs{}, v103.InputsFormJSON)
}

func main() {
	registerPlugins()
	runner.Run()
}
