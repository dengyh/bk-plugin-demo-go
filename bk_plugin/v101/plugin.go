package v101

import (
	"fmt"

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
	TenantID string `json:"tenant_id" jsonschema:"title=租户 ID"`
}

type Outputs struct {
	World string `json:"world" jsonschema:"title=World"`
}

type Plugin struct{}

func (p *Plugin) Version() string {
	return "1.0.1"
}

func (p *Plugin) Desc() string {
	return "蓝鲸 Go 插件 Demo - 租户 ID 验证"
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
	if contextInputs.TenantID == "" {
		c.WithField("tenant_id", contextInputs.TenantID).Error("tenant context missing")
		return fmt.Errorf("tenant_id is required")
	}

	c.WithField("tenant_id", contextInputs.TenantID).Info("tenant context received")

	outputs := Outputs{
		World: inputs.Hello,
	}
	return c.WriteOutputs(&outputs)
}
