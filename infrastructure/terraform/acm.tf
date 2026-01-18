# ACM Certificate for hexops.online
# This provides SSL/TLS for HTTPS access

resource "aws_acm_certificate" "hexops" {
  domain_name               = var.domain_name
  subject_alternative_names = ["*.${var.domain_name}"]
  validation_method         = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name        = "${var.cluster_name}-acm-cert"
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# Output the DNS validation records
# These need to be added to GoDaddy manually
output "acm_certificate_validation_records" {
  description = "DNS records to add to GoDaddy for ACM certificate validation"
  value = {
    for dvo in aws_acm_certificate.hexops.domain_validation_options : dvo.domain_name => {
      name  = dvo.resource_record_name
      type  = dvo.resource_record_type
      value = dvo.resource_record_value
    }
  }
}

output "acm_certificate_arn" {
  description = "ARN of the ACM certificate for use in ALB Ingress"
  value       = aws_acm_certificate.hexops.arn
}

output "acm_certificate_status" {
  description = "Status of the ACM certificate (PENDING_VALIDATION or ISSUED)"
  value       = aws_acm_certificate.hexops.status
}
