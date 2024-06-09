resource "aws_autoscaling_group" "app" {
# we can modify the desired, min and max values from variables.tf
  desired_capacity     = var.asg_desired_ec2 
  max_size             = var.asg_max_ec2
  min_size             = var.asg_min_ec2
  vpc_zone_identifier  = [aws_subnet.subnet1.id, aws_subnet.subnet2.id]
  launch_template {
    id      = aws_launch_template.app.id
    version = "$Latest"
  }
  target_group_arns = [aws_lb_target_group.app_https.arn]

  tag {
    key                 = "Name"
    value               = "pd-challenge-machine"
    propagate_at_launch = true
  }
}