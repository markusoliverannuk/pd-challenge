output "load_balancer_dns" {
  description = "where we can access the LB :)"
  value       = aws_lb.app.dns_name
}