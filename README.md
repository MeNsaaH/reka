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
- [ ] [Create Web Dashboard](https://github.com/MeNsaaH/reka/issues/3)
- [ ] [Persist state to remote sources](https://github.com/MeNsaaH/reka/issues/4)
- [ ] [Dockerize application](https://github.com/MeNsaaH/reka/issues/5)
- [ ] [Add AWS Resources](https://github.com/MeNsaaH/reka/issues/1)
- [ ] [Add GCP Resources](https://github.com/MeNsaaH/reka/issues/2)
- [ ] [Add Azure Resources](https://github.com/MeNsaaH/reka/issues/6)

#### Supported Resources
Here is a list of all [supported resources](./supported-resources.md) 

## Development
Copy `config/config.example.yaml` to `config/config.yaml` and make all necessary changes
```bash
cp config/config.example.yaml config/config.yaml
# One time run
go run main.go --config ../config/config.yaml
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
reka --config config/config.yaml
```
