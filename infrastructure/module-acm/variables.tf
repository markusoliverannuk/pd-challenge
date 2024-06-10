

variable "hosted_zone_id" {
  description = "hosted zone id"
}

variable "record_name" {
  description = "hosted zone record name / subdomain value for certs"

}

variable "asg_min_ec2" {
  description = "minimum number of instances"
  default     = 1
}

variable "asg_max_ec2" {
  description = "maximum number of instances"
  default     = 3
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