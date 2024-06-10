module "vpc" {
  source = "./module-vpc"
  # Pass necessary variables
  vpc_cidr = var.vpc_cidr
}

module "dns" {
  source = "./module-dns"
  # Pass necessary variables
  record_name = var.record_name
  lb_dns = module.lb.dns_name
  lb_zone_id = module.lb.zone_id
}

module "compute" {
  source = "./module-compute"
  vpc_id = module.vpc.vpc_id
  allowed_ip = var.allowed_ip
  machine_ami = var.machine_ami
  instance_type = var.instance_type
  key_name = var.key_name
}

module "lb" {
  source = "./module-lb"
  # Pass necessary variables
  vpc_id = module.vpc.vpc_id
  subnet1_id = module.subnets.subnet1_id
  subnet2_id = module.subnets.subnet2_id
  allowed_ip = var.allowed_inbound_ips_lb
}

module "subnets" {
   source = "./module-subnets"

    vpc_id = module.vpc.vpc_id
    igw_id = module.igw.igw_id
}

module "igw" {
   source = "./module-igw"

    vpc_id = module.vpc.vpc_id
}

module "acm" {
   source = "./module-acm"

    vpc_id = module.vpc.vpc_id
    igw_id = module.igw.igw_id
}