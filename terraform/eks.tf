module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 21.10"

  # Cluster
  name               = var.cluster_name
  kubernetes_version = "1.29"

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  endpoint_public_access  = true
  endpoint_private_access = true

  # Managed Node Groups
  eks_managed_node_groups = {
    general = {
      desired_size = 1
      min_size     = 1
      max_size     = 2

      instance_types = ["t3.medium"]
      capacity_type  = "ON_DEMAND"

      labels = {
        role = "general"
      }

      tags = {
        Environment = var.environment
      }
    }

    application = {
      desired_size = 1
      min_size     = 1
      max_size     = 3

      instance_types = ["t3.medium"]
      capacity_type  = "ON_DEMAND"

      labels = {
        role = "application"
      }

      tags = {
        Environment = var.environment
      }
    }
  }

  # Enable the default node security group to ensure proper communication
  create_node_security_group = true

  # Enable admin permissions for the cluster creator
  enable_cluster_creator_admin_permissions = true

  # Addons
  addons = {
    coredns = {
      most_recent = true
    }

    kube-proxy = {
      most_recent = true
    }

    vpc-cni = {
      most_recent = true
    }

    aws-ebs-csi-driver = {
      most_recent = true
    }
  }

  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}
