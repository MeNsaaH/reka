# reka
> Never forget that instance running again

A Cloud Infrastructure Management Tool to stop, resume, clean and destroy resources based on tags

### Project Name
`REKA` is derived from a Native Nigerian Language, Igbo, meaning `Reap`.

### TODO
- [x] Bootstrap application architecture
- [x] Schedule resource refreshing
- [ ] Create Web Dashboard 
- [ ] generate Sample Yaml config and load config
- [ ] Allow users to specify tags/resources to track from reka UI with reaping Details
- [ ] Support Manual Trigger of resources reaping from Dashboard/CLI
- [ ] Allow authentication username and password set in config file
- [ ] Create Kubernetes Manifests and Helm Charts
- [ ] Create CLI

#### Supported Resources
- AWS: https://github.com/MeNsaaH/reka/issues/1 
- GCP: https://github.com/MeNsaaH/reka/issues/2 


## Development
```bash
# Web UI
cd web
go run main.go
```

using [air](https://github.com/cosmtrek/air) with autoreload UI features
```bash
    air -c .air.toml
```