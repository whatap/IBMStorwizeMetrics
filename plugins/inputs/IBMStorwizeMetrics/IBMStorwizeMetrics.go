package IBMStorwizeMetrics

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

var sampleConfig string

type MetricConfig struct {
	Endpoint string   `toml:"endpoint"`
	Tags     []string `toml:"tags"`
	Fields   []string `toml:"fields"`
}

type IBMStorwizeMetrics struct {
	Endpoint           string         `toml:"endpoint"`
	AuthUsername       string         `toml:"auth_username"`
	AuthPassword       string         `toml:"auth_password"`
	InsecureSkipVerify bool           `toml:"insecure_skip_verify"`
	Metrics            []MetricConfig `toml:"metrics"`
	authCache          *AuthCache
}

func (sw *IBMStorwizeMetrics) Description() string {
	return "An input plugin based on IBM Spectrum Virtualize RESTful API."
}

func (sw *IBMStorwizeMetrics) SampleConfig() string {
	return sampleConfig
}

// Init is for setup, and validating config.
func (sw *IBMStorwizeMetrics) Init() error {
	sw.authCache = NewAuthCache(sw.Endpoint, sw.AuthUsername, sw.AuthPassword, sw.InsecureSkipVerify)
	return nil
}

func init() {
	inputs.Add("IBMStorwizeMetrics", func() telegraf.Input { return &IBMStorwizeMetrics{} })
}

func (sw *IBMStorwizeMetrics) Gather(acc telegraf.Accumulator) error {
	token, err := sw.authCache.GetToken()
	if err != nil {
		return err
	}

	for _, metric := range sw.Metrics {
		jsonData, err := sw.DoRequest(token, sw.Endpoint+metric.Endpoint)
		if err != nil {
			return err
		}

		for _, item := range jsonData {
			tags := make(map[string]string)
			fields := make(map[string]interface{})

			for _, tag := range metric.Tags {
				if value, ok := item[tag].(string); ok {
					tags[tag] = value
				}
			}

			for _, field := range metric.Fields {
				if value, ok := item[field]; ok {
					fields[field] = value
				}
			}

			// Add metric to the accumulator
			acc.AddFields("ibm_storwize_metrics", fields, tags)
		}
	}

	return nil
}

// Performs an HTTP request to the IBM Spectrum Virtualize API using an authenticated session
func (sw *IBMStorwizeMetrics) DoRequest(token, url string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %v", err)
	}

	req.Header.Set("X-Auth-Token", token)

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: sw.InsecureSkipVerify, // Note: Set to true for development purposes only
			},
		},
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Call failed: status code " + resp.Status)
	}

	var results []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("decoding response: %v", err)
	}

	return results, nil
}
