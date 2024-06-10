variable "asg_desired_ec2" {
  description = "desired number of instances in the asg"
  type        = string
}

variable "asg_max_ec2" {
  description = "max number of isntances in the asg"
  type        = string
}

variable "asg_min_ec2" {
  description = "min number of isntances in the asg"
  type        = string
}

variable "subnet1_id" {
  description = "ID of the first subnet"
  type        = string
}

variable "subnet2_id" {
  description = "ID of the second subnet"
  type        = string
}

variable "launch_template_id" {
  description = "ID of the ec2 launch template"
  type        = string
}

variable "target_group_http_arn" {
  description = "arn of http target group"
  type        = string
}

variable "target_group_https_arn" {
  description = "arn of https target group"
  type        = string
}

