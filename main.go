package main

import (
	v100 "bk-plugin-demo-go/bk_plugin/v100"
	v101 "bk-plugin-demo-go/bk_plugin/v101"

	"github.com/TencentBlueKing/beego-runtime/runner"
	"github.com/TencentBlueKing/bk-plugin-framework-go/hub"
)

func registerPlugins() {
	hub.MustInstall(&v100.Plugin{}, v100.ContextInputs{}, v100.Outputs{}, v100.InputsFormJSON)
	hub.MustInstall(&v101.Plugin{}, v101.ContextInputs{}, v101.Outputs{}, v101.InputsFormJSON)
}

func main() {
	registerPlugins()
	runner.Run()
}
