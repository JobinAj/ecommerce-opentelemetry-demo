# Terraform Infrastructure

This directory contains Terraform configurations for deploying the infrastructure required for the e-commerce platform on AWS.

## Components

- **VPC**: Virtual Private Cloud with public and private subnets
- **EKS**: Elastic Kubernetes Service cluster for running containerized applications
- **RDS**: PostgreSQL database instance for data persistence

## Variables

The following variables need to be configured:

- `aws_region`: AWS region to deploy resources (default: us-east-1)
- `cluster_name`: Name of the EKS cluster (default: ecom-eks-cluster)
- `vpc_cidr`: CIDR block for the VPC (default: 10.0.0.0/16)
- `environment`: Environment name (default: development)
- `db_username`: Database username (required)
- `db_password`: Database password (required, sensitive)

## Usage

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. Review the execution plan:
   ```bash
   terraform plan
   ```

3. Apply the configuration:
   ```bash
   terraform apply
   ```

4. To destroy the infrastructure:
   ```bash
   terraform destroy
   ```

## Outputs

After applying the configuration, the following outputs will be available:

- `rds_endpoint`: The connection endpoint for the RDS instance
- `rds_port`: The port for the RDS instance
- `rds_db_name`: The name of the database
- `rds_username`: The username for the database
- `rds_security_group_id`: The ID of the RDS security group