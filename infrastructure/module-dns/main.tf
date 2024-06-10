resource "aws_route53_record" "lb_record" {
  zone_id = var.hosted_zone_id  
  name    = var.record_name
  type    = "A"
  alias {
    name                   = var.lb_dns
    zone_id                = var.lb_zone_id
    evaluate_target_health = true
  }
}