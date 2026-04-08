package v100

import (
	"log"

	"github.com/TencentBlueKing/bk-plugin-framework-go/kit"
)

var InputsFormJSON = []byte(`{
	"type": "object",
	"properties": {
		"hello": {
			"type": "string",
			"title": "Hello",
			"ui:component": {"name": "bk-input", "props": {"type": "textarea"}}
		}
	},
	"required": ["hello"]
}`)

type Inputs struct {
	Hello string `json:"hello" jsonschema:"title=Hello"`
}

type ContextInputs struct {
	Executor string `json:"executor" jsonschema:"title=任务执行人"`
}

type Outputs struct {
	World string `json:"world" jsonschema:"title=World"`
}

type Plugin struct{}

func (p *Plugin) Version() string {
	return "1.0.0"
}

func (p *Plugin) Desc() string {
	return "蓝鲸 Go 插件 Demo"
}

func (p *Plugin) Execute(c *kit.Context) error {
	var inputs Inputs
	if err := c.ReadInputs(&inputs); err != nil {
		return err
	}

	var contextInputs ContextInputs
	if err := c.ReadContextInputs(&contextInputs); err != nil {
		return err
	}

	log.Printf("executor: %s, hello: %s\n", contextInputs.Executor, inputs.Hello)

	outputs := Outputs{
		World: inputs.Hello,
	}
	return c.WriteOutputs(&outputs)
}
