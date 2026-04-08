package main

import (
	v100 "bk-plugin-demo-go/bk_plugin/v100"

	"github.com/TencentBlueKing/beego-runtime/runner"
	"github.com/TencentBlueKing/bk-plugin-framework-go/hub"
)

func main() {
	hub.MustInstall(&v100.Plugin{}, v100.ContextInputs{}, v100.Outputs{}, v100.InputsFormJSON)
	runner.Run()
}
