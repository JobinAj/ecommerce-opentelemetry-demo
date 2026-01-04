output "rds_endpoint" {
  description = "The connection endpoint for the RDS instance"
  value       = aws_db_instance.postgres_db.endpoint
}

output "rds_port" {
  description = "The port for the RDS instance"
  value       = aws_db_instance.postgres_db.port
}

output "rds_db_name" {
  description = "The name of the database"
  value       = aws_db_instance.postgres_db.db_name
}

output "rds_username" {
  description = "The username for the database"
  value       = aws_db_instance.postgres_db.username
}

output "rds_security_group_id" {
  description = "The ID of the RDS security group"
  value       = aws_security_group.rds_security_group.id
}