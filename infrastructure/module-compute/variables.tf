variable "vpc_id" {
  description = "The ID of the VPC"
  type        = string
}

variable "allowed_ip" {
  description = "The IP address allowed to access the instances"
  type        = string
}

variable "machine_ami" {
  description = "The AMI ID for the instances"
  type        = string
}

variable "instance_type" {
  description = "The instance type for the instances"
  type        = string
}

variable "key_name" {
  description = "The key pair name for SSH access"
  type        = string
}