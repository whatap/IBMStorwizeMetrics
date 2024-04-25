# IBMStorwizeMetrics

IBMStorwizeMetrics is a Telegraf input plugin designed to interact with IBM Spectrum Virtualize RESTful API. This plugin enables the gathering and monitoring of metrics from IBM storage solutions that use this API.

## Requirements

- **Go**: Ensure you have Go installed on your system to compile the project.
- **Telegraf**: This plugin is intended to be used as a part of Telegraf, which must be installed on your system.

## Project Initialization
To get started with the IBMStorwizeMetrics plugin, clone the repository to your local machine:
```bash
git clone https://github.com/alex-lata/IBMStorwizeMetrics.git
```

## Gather Dependencies
Navigate to the project directory and run the following command to resolve and tidy up the project's dependencies:
```bash
cd IBMStorwizeMetrics  
go mod tidy  # Ensures your project's dependencies are clean and up-to-date
```

## Build the Project
Compile the project using the Go compiler:
```bash
go build -o TelegrafIBMStoreWizeMetrics cmd/main.go
```
This command builds the project and outputs an executable named TelegrafIBMStoreWizeMetrics.

## Run Locally
To run the plugin locally, use the following command:
```bash
./TelegrafIBMStoreWizeMetrics --config plugins/inputs/IBMStorwizeMetrics/sample.conf
```
Replace sample.conf with your configuration file if you have a different setup.

## Debugging with VSCode
For those who prefer using Visual Studio Code for development, you can set up a launch.json configuration file in the .vscode folder to facilitate debugging:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Debug",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/main.go",
            "args": ["--config", "${workspaceFolder}/plugins/inputs/IBMStorwizeMetrics/sample.conf"],
        }
    ]
}
```

## Running as a Telegraf Plugin
To run IBMStorwizeMetrics as a Telegraf plugin, use the following command:
```bash
telegraf --config telegraf.conf
```
Make sure that telegraf.conf includes the configuration details for this plugin.

## Debugging as a Telegraf Plugin
For debugging, run Telegraf in debug mode:
```bash
telegraf --config telegraf.conf --debug
```
This command provides detailed debug output, useful for diagnosing issues or verifying that the plugin is functioning correctly.

## Sample Plugin Configuration
Below is a sample configuration for the IBMStorwizeMetrics plugin within Telegraf. This sample shows how to set up the plugin with basic authentication and endpoint details:
```json
[[inputs.IBMStorwizeMetrics]]
  endpoint = "https://IBM_URL:IBM_PORT/rest/v1"
  auth_username = "USERNAME"
  auth_password = "PASSWORD"
  insecure_skip_verify = false
  # Example of an endpoint with the mapptings from the response to tags and fields
  [[inputs.IBMStorwizeMetrics.metrics]]
    endpoint = "/lsnodestats"
    tags = ["node_id", "node_name"]
    fields = ["stat_current", "stat_name", "stat_peak", "stat_peak_time"]
  [[inputs.IBMStorwizeMetrics.metrics]]
    endpoint = "/lsvdisk"
    tags = ["name", "volume_name"]
    fields = ["capacity", "is_snapshot"]
```

## Sample Telegraf Configuration
Below is a sample Telegraf configuration for the IBMStorwizeMetrics plugin. 
```json
# Input Plugin: Execd
[[inputs.execd]]
  command = ["[PATH]/TelegrafIBMStoreWizeMetrics", "-config", "[PATH]/IBMStorwizeMetrics/plugins/inputs/IBMStorwizeMetrics/sample.conf"]
  signal = "none"
  interval = "30s"

# Output Plugin: Write metrics to a file
[[outputs.file]]
  files = ["[PATH]/IBMStorwizeMetrics/metrics.out"]
  data_format = "influx"
```

## Contributions
Contributions to the IBMStorwizeMetrics plugin are welcome. Please submit pull requests to the repository or report issues as needed.
