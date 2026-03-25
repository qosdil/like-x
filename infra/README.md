# like-x Infra (Terraform)

Terraform code for provisioning AWS VPC + EKS cluster for `like-x`.

## What this repo does

- Creates an AWS VPC with public/private subnets using `terraform-aws-modules/vpc/aws`.
- Creates an EKS cluster using `terraform-aws-modules/eks/aws`.
- EKS managed node group with `t3.small`, auto scaling 1-2 nodes.

## Prerequisites

- Terraform 1.5+ (or compatible with module versions)
- AWS CLI credentials configured (AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY or profile)
- IAM permissions to create VPC/EC2/EKS/Route53/ELB resources.

## Quick start

```bash
cd <this directory>
terraform init
terraform plan
terraform apply
```

## Outputs

- `cluster_name`
- `cluster_endpoint`

## Notes

- Existing `terraform.tfstate` and `terraform.tfstate.backup` indicate there may already be resources deployed.
- If `apply` hangs, check for stuck resources and avoid unnecessary name changes in-place (cluster name changes may trigger destroy/recreate).

## Cleanup

```bash
terraform destroy
```
