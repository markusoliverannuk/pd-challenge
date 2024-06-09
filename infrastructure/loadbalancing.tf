resource "aws_lb" "app" {
  name               = "app-load-balancer"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.app_sg.id]
  subnets            = [aws_subnet.subnet1.id, aws_subnet.subnet2.id] # to modify l8r

  enable_deletion_protection = false
}



resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.app.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app-tg-http.arn
  }
}

resource "aws_lb_listener" "https_listener" {
  load_balancer_arn = aws_lb.app.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate.api_challenge_cert.arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app-tg-https.arn
  }
}


resource "aws_lb_target_group" "app-tg-http" {
  name        = "app-tg-http"
  port        = "80"
  protocol    = "HTTP"
  vpc_id      = aws_vpc.main.id

  # health_check {
  #   path               = "/"
  #   protocol           = "HTTP"
  #   port               = "80"
  #   interval           = 30
  #   timeout            = 10
  #   healthy_threshold  = 3
  #   unhealthy_threshold = 3
  # }
}

resource "aws_lb_target_group" "app-tg-https" {
  name        = "app-tg-https"
  port        = "443"
  protocol    = "HTTPS"
  vpc_id      = aws_vpc.main.id

  # health_check {
  #   path               = "/"
  #   protocol           = "HTTP"
  #   port               = "80"
  #   interval           = 30
  #   timeout            = 10
  #   healthy_threshold  = 3
  #   unhealthy_threshold = 3
  # }
}
