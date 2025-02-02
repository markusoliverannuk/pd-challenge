variable "aws_region" {
    default = "eu-north-1"
    description = "where we create our resources"
}

variable "instance_type" {
    default ="t2.small"
    description = "the size of the instances that will be inside the ASG"
}

variable "asg_desired_ec2" {
  description = "desired number of instances"
  default     = 1
}

variable "hosted_zone_id" {
  description = "hosted zone id"
  default     = "Z09754511Y8SW428ZP8R3"

}

variable "record_name" {
  description = "hosted zone record name / subdomain value for certs"
  default     = "api-challenge-v2.techwithmarkus.com"

}

variable "asg_min_ec2" {
  description = "minimum number of instances"
  default     = 1
}

variable "asg_max_ec2" {
  description = "maximum number of instances"
  default     = 3
}

variable "vpc_cidr" {
  description = "maximum number of instances"
  type = string
  default = "10.0.0.0/16"
}

variable "allowed_ip" {
  description = "address allowed to connect to the instance"
  default     = "0.0.0.0/0"
}

variable "key_name" {
  description = "The name of the key pair"
  default = "pd-challenge-kp"
}

variable "machine_ami" {
  description = "the id of the machine image which we'll be using for our instnaces (ubuntu 24:04 64bit (x86))"
  default     = "ami-04b70fa74e45c3917"
}

variable "allowed_inbound_ips_lb" {
  default = "0.0.0.0/0"
}