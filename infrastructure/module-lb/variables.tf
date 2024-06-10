variable "vpc_id" {
  description = "the ID of the VPC"
  type        = string
}

variable "subnet1_id" {
  description = "id of first subnet"
  type        = string
}

variable "subnet2_id" {
  description = "id of second subnet"
  type        = string
}

variable "allowed_ip" {
  description = "id of second subnet"
  type        = string
}

variable "acm_certificate_arn" {
  description = "arn of tls cert from acm"
  type        = string
}
