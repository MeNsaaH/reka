# reka
> Never forget that instance running again

A Cloud Infrastructure Management Tool to stop, resume, clean and destroy resources based on tags

### Project Name
`REKA` is derived from a Native Nigerian Language, Igbo, meaning `Reap`|`Tear Down`.

### BEWARE!

This tool is **HIGHLY DESTRUCTIVE** and can deletes all resources! This should be used in environments with **WITH CAUTION**.

### TODO
- [x] Bootstrap application architecture
- [x] Schedule resource refreshing
- [ ] Create Web Dashboard 
- [x] generate Sample Yaml config and load config
- [ ] Allow users to specify tags/resources to track from reka UI with reaping Details
- [ ] Support Manual Trigger of resources reaping from Dashboard/CLI
- [ ] Allow authentication username and password set in config file
- [ ] Create Kubernetes Manifests and Helm Charts
- [ ] Create CLI
- [ ] Save infra state and desired state in persistent storage

#### Supported Resources
- AWS: https://github.com/MeNsaaH/reka/issues/1 
- GCP: https://github.com/MeNsaaH/reka/issues/2 


## Development
Copy `config/config.example.yaml` to `config/config.yaml` and make all necessary changes
```bash
cp config/config.example.yaml config/config.yaml

# Start Web UI
cd web
go run main.go --config ../config/config.yaml
```

using [air](https://github.com/cosmtrek/air) with autoreload UI features
```bash
    air -c .air.toml
```