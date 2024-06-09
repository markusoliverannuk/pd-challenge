resource "aws_acm_certificate" "api_challenge_cert" {
  domain_name       = var.record_name
  validation_method = "DNS"

  tags = {
    Name = "api-challenge-cert"
  }
}

resource "aws_route53_record" "api_challenge_cert_validation" {
  for_each = {
    for dvo in aws_acm_certificate.api_challenge_cert.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  }

  zone_id = var.zone_id
  name    = each.value.name
  type    = each.value.type
  ttl     = 60
  records = [each.value.record]
}

resource "aws_acm_certificate_validation" "api_challenge_cert_validation" {
  certificate_arn         = aws_acm_certificate.api_challenge_cert.arn
  validation_record_fqdns = [for record in aws_route53_record.api_challenge_cert_validation : record.fqdn]
}

resource "aws_route53_record" "lb_record" {
  zone_id = var.zone_id  
  name    = var.record_name
  type    = "A"
  alias {
    name                   = aws_lb.app.dns_name
    zone_id                = aws_lb.app.zone_id
    evaluate_target_health = true
  }
}