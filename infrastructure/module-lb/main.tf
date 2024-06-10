module "acm" {
  source      = "../module-acm"
}

resource "aws_lb" "app" {
  name               = "app-load-balancer"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.lb_sg.id]
  subnets            = [var.subnet1_id, var.subnet2_id] # to modify l8r

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
  certificate_arn   = module.acm.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app-tg-https.arn
  }
}


resource "aws_lb_target_group" "app-tg-http" {
  name        = "app-tg-http"
  port        = "80"
  protocol    = "HTTP"
  vpc_id      = var.vpc_id

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
  vpc_id      = var.vpc_id

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

resource "aws_security_group" "lb_sg" {
  name        = "traffic_rules_for_lb"
  description = "giving traffic access"
  vpc_id            = var.vpc_id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.allowed_ip]## check @ variables.tf
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = [var.allowed_ip] ## check @ variables.tf
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = [var.allowed_ip] ## check @ variables.tf
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}