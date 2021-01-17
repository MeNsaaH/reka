# reka
> Never forget that instance running again

<br>

[![Build Status](https://github.com/mensaah/reka/workflows/Test/badge.svg)](https://github.com/mensaah/reka/actions)


A Cloud Infrastructure Management Tool to stop, resume, clean and destroy resources based on tags

### Project Name
`REKA` is derived from a Native Nigerian Language, Igbo, meaning `Reap`|`Tear Down`.

### BEWARE!

This tool is **HIGHLY DESTRUCTIVE** and can deletes all resources! This should be used in environments with **WITH CAUTION**.

## What Does this tool do?
- Stop/Resume resources based on configuration (activeDuration) for instance stopping an EKS cluster would mean resizing all nodegroups to 0 and resuming will be restoring back to original size
- Destroy resources that have specific tags/labels or after a certain Duration(terminationDate)
- Clean Up unused resources (such as EBS volumes, Elastic IPs)

### TODO
- [x] Schedule resource refreshing
- [x] generate Sample Yaml config and load config
- [x] [Persist state to remote sources](https://github.com/MeNsaaH/reka/issues/4)
- [ ] [Add More AWS Resources](https://github.com/MeNsaaH/reka/issues/1)
- [ ] [Add More GCP Resources](https://github.com/MeNsaaH/reka/issues/2)
- [ ] [Add More Azure Resources](https://github.com/MeNsaaH/reka/issues/6)
- [ ] [Create Web Dashboard](https://github.com/MeNsaaH/reka/issues/3)

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

###  Authentication
Reka uses the default authentication method for all providers