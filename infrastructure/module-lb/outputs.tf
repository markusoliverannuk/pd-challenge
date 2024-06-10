output "dns_name" {
  value = aws_lb.app.dns_name
}

output "zone_id" {
    value = aws_lb.app.zone_id
}

output "target_group_http_arn" {
    value = aws_lb_target_group.app-tg-http.arn
}

output "target_group_https_arn" {
    value = aws_lb_target_group.app-tg-https.arn
}

