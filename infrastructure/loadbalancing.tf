resource "aws_lb" "app" {
  name               = "app-load-balancer"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.app_sg.id]
  subnets            = [aws_subnet.subnet1.id, aws_subnet.subnet2.id] # to modify l8r

  enable_deletion_protection = false
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.app.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "<your-acm-certificate-arn>"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app_https.arn
  }
}

resource "aws_lb_target_group" "app_https" {
  name     = "app-tg-https"
  port     = 443
  protocol = "HTTPS"
  vpc_id   = aws_vpc.main.id

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 5
    interval            = 30
    path                = "/"
    matcher             = "200"
  }
}