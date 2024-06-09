resource "aws_lb" "app" {
  name               = "app-load-balancer"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.app_sg.id]
  subnets            = [aws_subnet.subnet1.id, aws_subnet.subnet2.id] # to modify l8r

  enable_deletion_protection = false
}

resource "aws_lb_target_group" "app" {
  name     = "app-tg-http"
  port     = 80
  protocol = "HTTP"
  vpc_id   = aws_vpc.main.id

    # health_check {
  #   healthy_threshold   = 2
  #   unhealthy_threshold = 2
  #   timeout             = 5
  #   interval            = 30
  #   path                = "/"
  #   matcher             = "200"
  # }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.app.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.app.arn
  }
}

