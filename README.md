# reka
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
- [ ] Add Resource Manager Dependencies
- [ ] Schedule tasks for destruction
- [ ] Expose API
- [ ] Create CLI to interact with API
- [ ] Create GUI to interact with API

#### AWS
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
