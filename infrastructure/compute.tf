resource "aws_security_group" "app_sg" {
  name        = "traffic_rules_for_challenge"
  description = "giving traffic access"
  vpc_id            = aws_vpc.main.id

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

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_launch_template" "app" {
  name_prefix = "app-template"

  image_id          = var.machine_ami # check description and additional info from variables.tf
  instance_type     = var.instance_type # we take the instance type from variables.tf where we declared it. currently set to t2.small
  key_name          = var.key_name
  security_groups = [aws_security_group.app_sg.id]
  network_interfaces {
    associate_public_ip_address = true
  }

  user_data = filebase64("userdata.sh")

  lifecycle {
    create_before_destroy = true
  }
}