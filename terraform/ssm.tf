data "aws_ami" "amazon_linux" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["al2023-ami-*-x86_64"]
  }
}

resource "aws_iam_role" "ssm_ec2_role" {
  name = "${var.cluster_name}-ssm-ec2-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect    = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })

  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

resource "aws_iam_role_policy_attachment" "ssm_core" {
  role       = aws_iam_role.ssm_ec2_role.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_instance_profile" "ssm_profile" {
  name = "${var.cluster_name}-ssm-profile"
  role = aws_iam_role.ssm_ec2_role.name
}

resource "aws_security_group" "ssm_ec2_sg" {
  name        = "${var.cluster_name}-ssm-sg"
  description = "Security group for SSM-only EC2"
  vpc_id      = module.vpc.vpc_id

egress {
  from_port   = 443
  to_port     = 443
  protocol    = "tcp"
  cidr_blocks = [module.vpc.vpc_cidr_block]
}

  tags = {
    Name        = "${var.cluster_name}-ssm-sg"
    Environment = var.environment
  }
}

resource "aws_instance" "ssm_admin" {
  ami                         = data.aws_ami.amazon_linux.id
  instance_type               = "t2.micro"
  subnet_id                   = module.vpc.private_subnets[0]
  iam_instance_profile        = aws_iam_instance_profile.ssm_profile.name
  vpc_security_group_ids      = [aws_security_group.ssm_ec2_sg.id]
  associate_public_ip_address = false

  tags = {
    Name        = "${var.cluster_name}-ssm-admin"
    Environment = var.environment
    Role        = "eks-admin"
  }
}
