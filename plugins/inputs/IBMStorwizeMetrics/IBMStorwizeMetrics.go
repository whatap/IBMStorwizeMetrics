package IBMStorwizeMetrics

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type IBMStorwize struct {
	Endpoint     string `toml:"endpoint"`
	AuthUsername string `toml:"auth_username"`
	AuthPassword string `toml:"auth_password"`
}

func (sw *IBMStorwize) Description() string {
	return "An input plugin based on IBM Spectrum Virtualize RESTful API."
}

func (sw *IBMStorwize) SampleConfig() string {
	return `
 ## Indicate if everything is fine
 ok = true
`
}

// Init is for setup, and validating config.
func (sw *IBMStorwize) Init() error {
	return nil
}

func (sw *IBMStorwize) Gather(acc telegraf.Accumulator) error {
	// TODO : WIP
	return nil
}

func init() {
	inputs.Add("IBMStorwizeMetrics", func() telegraf.Input { return &IBMStorwize{} })
}
