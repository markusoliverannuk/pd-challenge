output "certificate_arn" {
  value = aws_acm_certificate.api_challenge_cert.arn
}