# reka
> Never forget that instance running again

A Cloud Infrastructure Management Tool to stop, resume, clean and destroy resources based on tags

### Project Name
`REKA` is derived from a Native Nigerian Language, Igbo, meaning `Reap`|`Tear Down`.

### BEWARE!

This tool is **HIGHLY DESTRUCTIVE** and can deletes all resources! This should be used in environments with **WITH CAUTION**.

### TODO
- [x] Schedule resource refreshing
- [x] generate Sample Yaml config and load config
- [ ] Create Web Dashboard 
- [ ] Allow users to specify tags/resources to track from reka UI with reaping Details
- [ ] Persist state to remote sources (GCS, S3, if possible Databases)
- [ ] Allow web authentication username and password set in config file
- [ ] Create Kubernetes Manifests and Helm Charts

#### Supported Resources
- AWS: https://github.com/MeNsaaH/reka/issues/1 
- GCP: https://github.com/MeNsaaH/reka/issues/2 


## Development
Copy `config/config.example.yaml` to `config/config.yaml` and make all necessary changes
```bash
cp config/config.example.yaml config/config.yaml
# One time run
go run main.go --config ../config/config.yaml

# TODO Web Dashboard
go run main.go web --config ../config/config.yaml
```

## Installation
### Using go
```bash
go get -u github.com/mensaah/reka
```

## Usage

Reka loads default config from $HOME/.reka.yaml if --config param is not passed.
Copy `config/config.example.yaml` to `config/config.yaml` and make all necessary changes

```bash
cp config/config.example.yaml config/config.yaml
# Run reka using configuration. Stops stoppable resources, resume resumable resources and terminate
# resources dues for termination
reka --config ../config/config.yaml

# KABOOOOOM
# everything gone
reka nuke --config ../config/config.yaml
```