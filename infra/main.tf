module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "6.6.0"

  name = var.cluster_name
  cidr = var.vpc_cidr

  azs = [
    "${var.region}a",
    "${var.region}b",
  ]

  private_subnets = [
    "10.0.1.0/24",
    "10.0.2.0/24",
  ]

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
  }

  public_subnets = [
    "10.0.101.0/24",
    "10.0.102.0/24",
  ]

  public_subnet_tags = {
    "kubernetes.io/role/elb" = 1
  }

  enable_nat_gateway = true
  single_nat_gateway = true

  tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
  }
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.0"

  name = var.cluster_name
  kubernetes_version = "1.34"
  enable_cluster_creator_admin_permissions = true
  
  // Disable EKS extended support (paid service)
  upgrade_policy = {support_type = "STANDARD"}

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets
  control_plane_subnet_ids = module.vpc.public_subnets
  eks_managed_node_groups = {
    default = {
      name = "ng-${var.cluster_name}"
      instance_types = ["t3.small"] # provision stuck with t3.nano

      min_size     = 1
      max_size     = 2
      desired_size = 1
    }
  }

  tags = {
    Environment = "dev"
  }

  addons = {
    coredns                = {}
    eks-pod-identity-agent = {
      before_compute = true
    }
    kube-proxy             = {}
    vpc-cni                = {
      before_compute = true
    }
  }
}
