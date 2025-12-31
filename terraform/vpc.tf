module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "6.5.1"

  name="${var.cluster_name}-vpc"
  cidr=var.vpc_cidr

azs             = ["us-east-1a"]
private_subnets = ["10.0.1.0/24"]
public_subnets  = ["10.0.101.0/24"]

enable_nat_gateway = true
single_nat_gateway = true
enable_dns_hostnames = true
enable_dns_support = true

  public_subnet_tags = {
    "kubernetes.io/role/elb"                    = "1"
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb"           = "1"
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
  }

  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}
