# reka 
> Never forget that instance running again

<br>

[![Build Status](https://github.com/mensaah/reka/workflows/Test/badge.svg)](https://github.com/mensaah/reka/actions)
[![deploy](https://github.com/MeNsaaH/reka/workflows/deploy/badge.svg)](https://github.com/mensaah/reka/actions)


A Cloud Infrastructure Management Tool to stop, resume, clean and destroy resources based on tags. Reka uses a config to determine what actions should be taken on resources. It can prove to be a cost management tool where you can stop your tests environments during breaks, holidays and non-working hours. It can also be a nuke tool to nuke an account. It currently supports both AWS and GCP. A full list of supported resources can be found [here](./docs/supported-resources.md) 

#### What It can do
- Stop/Resume resources for example stopping an EC2 instance or resizing an EKS cluster to 0.
- Destroy/Terminate resources 
- Clean Up unused resources (such as EBS volumes, Elastic IPs)

#### Project Name
`REKA` is derived from a Native Nigerian Language, Igbo, meaning `Reap`|`Tear Down`.

### BEWARE!

This tool is **HIGHLY DESTRUCTIVE** and deletes cloud resources! This should be used in environments with **WITH CAUTION**.


## Table of Contents
- [Getting Started](#getting-started)
  - [Installation](#installation)
- [Usage](#usage)
  - [Authentication](#authentication)
  - [Rules](#rules)
  - [Excluding Resources](#excluding-resources)
- [RoadMap](#roadmap)
- [Contributing](#contributing)



## Getting Started
### Installation
#### Binary

The Reka binary can be downloaded from the [Releases](https://github.com/mensaah/reka/releases) and executed directly on its respective OS.

#### Docker

The reka image is also available on DockerHub. 
```bash
    docker pull mensaah/reka
```
If `config.yaml` is in the current directory, reka can be executed as:

```bash
    docker run -it -e AWS_ACCESS_KEY -e AWS_SECRET_ACCESS_KEY\
        -v `pwd`:/config mensaah/reka --config /config/config.yaml
```

#### Go
Reka can also be installed like a regular golang Binary if go is installed

```bash
go get -u github.com/mensaah/reka
```

## Usage

Reka uses a config file to know what resources to target and the actions to be taken for such resources. [Here](./docs/config.example.yaml) is an example configuration file that Reka uses.

```bash
    reka --config config.yaml

    # To run without destroying any instance, basically just stops, resumes resources
    reka --config config.yaml --disable-destroy

    # To see full range of commands that can executed with reka
    reka help

    # View supported resources
    reka resources
```
###  Authentication
- AWS
To use AWS Provider, you need to either have your aws credentials setup or export `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and optionally `AWS_SESSION_TOKEN` (optional).

- GCP
You can authenticate to GCP Provider by providing by setting the environment variable `GOOGLE_APPLICATION_CREDENTIALS` which is the path to the service account credentials. If `gcloud` is configured, the `gcloud` profile can also be used

### Rules
Reka supports different resource rules. These are usually in the form:

```yaml
rules:
  - name: <RULE_NAME>
    # Target specific resource types
    target:
      - resource-targets
    tags:
      tagKey1: tagValue1
      tagKey2: tagValue2
    condition:
        CONDITIONS
```
Conditions needs to be met before the action is taken on the resource.
- #### Stopping and Resuming instances within active Hours
This configuration sets that EC2 and EKS resources with tag `env = staging` needs to be active only within 7am - 7pm from Mondays to Fridays. When reka runs in any time outside that, the resource is stopped if stoppable and resumed once the condition is met.

```yaml
rules:
  - name: stop all staging instances after work and during weekends 
    tags:
      env: staging
    resources:
    - aws.ec2
    - aws.eks
    region: "us-east-2"
    condition:
      activeDuration: 
        startTime: "7:00"
        stopTime: "19:00"
        startDay: Monday
        stopDay: Friday
```
Specifying the `resources` list in any rule only applies the rules to those resources alone.

- #### Destroying instances after a particular time

```yaml
rules:
  - name: nuke all demo instances 2 weeks after demo october 9th (staging)
    tags:
      env: test
      project: proj-demo
    condition:
      terminationDate: "2021-10-23 01:00"
```

- #### Deleting unused Resources

```yaml
rules:
  - name: Delete all unused instances
    condition:
     terminationPolicy: unused
```


### Excluding Resources
You can additionally exclude resources which will make Reka not to act on those resources even when the satisfy a condition and a rule.
```yaml

exclude:
  - name: Exclude resources with prod tags in us-east-2
    region: "us-east-2"
    tags:
        env: prod

  - name: Exclude CI EC2 Instances on staging
    tags:
      env: staging
      ci-runner: true
    resources:
      - aws.ec2
```


## RoadMap
- [x] Schedule resource refreshing
- [x] generate Sample Yaml config and load config
- [x] [Persist state to remote sources](https://github.com/MeNsaaH/reka/issues/4)
- [ ] [Add More AWS Resources](https://github.com/MeNsaaH/reka/issues/1)
- [ ] [Add More GCP Resources](https://github.com/MeNsaaH/reka/issues/2)
- [ ] [Add MoreAzure Resources](https://github.com/MeNsaaH/reka/issues/6)
- [ ] [Create Web Dashboard to ease running reka](https://github.com/MeNsaaH/reka/issues/3)
- [ ] Tests 😆😆😆


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.