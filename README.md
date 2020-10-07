# reka
> Never forget that instance running again

A Cloud Infrastructure Management Tool to stop, resume, clean and destroy resources based on tags

### Project Name
`REKA` is derived from a Native Nigerian Language, Igbo, meaning `Reap`.

Currently Supports:
#### AWS
- EC2
- S3

### TODO
- [x] Bootstrap application architecture
- [x] get resources with specified tags for destruction
- [ ] Schedule tasks for destruction
- [ ] Create Web Dashboard 
- [ ] Allow users to specify tags/resources to track from reka UI with reaping Details
- [ ] Support Manual Trigger of resources reaping from Dashboard/CLI
- [ ] Allow authentication username and password set in config file
- [ ] Create Kubernetes Manifests and Helm Charts
- [ ] Create CLI

#### AWS
- [x] EC2: Stop|Resume| Destroy
- [x] S3: Destroy
- [ ] EKS : Stop|Resume: by Resizing | Destroy
- [ ] RDS : Stop | resume | Destroy
- [ ] EBS : Destroy | Unused EBS Volumes
- [ ] Elastic IPs : Destroy | Unused IPs
- [ ] VPCs : Destroy | Unused VPCs

#### GCP
- [ ] Compute : Stop|Resume: by Resizing | Destroy
- [ ] Cloud SQL: Stop | resume | Destroy
- [ ] EBS : Destroy | Unused EBS Volumes
- [ ] Cloud Storage : Destroy
- [ ] IPs : Destroy | Unused IPs
- [ ] VPCs : Destroy | Unused VPCs
