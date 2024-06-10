output "dns_name" {
  value = aws_lb.app.dns_name
}

output "zone_id" {
    value = aws_lb.app.zone_id
}