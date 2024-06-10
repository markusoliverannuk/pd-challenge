resource "aws_autoscaling_group" "app" {
# we can modify the desired, min and max values from variables.tf!
  desired_capacity     = var.asg_desired_ec2 
  max_size             = var.asg_max_ec2
  min_size             = var.asg_min_ec2
  vpc_zone_identifier  = [var.subnet1_id, var.subnet2_id]
  launch_template {
    id      = var.launch_template_id
    version = "$Latest"
  }
  target_group_arns = [var.target_group_http_arn, var.target_group_https_arn]

  tag {
    key                 = "Name"
    value               = "pd-challenge-machine"
    propagate_at_launch = true
  }
}