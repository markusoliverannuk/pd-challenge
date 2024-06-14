module "vpc" {
  source = "./module-vpc"
  vpc_cidr = var.vpc_cidr
}

module "dns" {
  source = "./module-dns"
  record_name = var.record_name
  lb_dns = module.lb.dns_name
  hosted_zone_id = var.hosted_zone_id
  lb_zone_id = module.lb.zone_id
}

module "iam" {
  source = "./module-iam"
}

module "compute" {
  source = "./module-compute"
  vpc_id = module.vpc.vpc_id
  allowed_ip = var.allowed_ip
  machine_ami = var.machine_ami
  instance_type = var.instance_type
  key_name = var.key_name
  instance_profile = module.iam.iam_instance_profile_ec2_ssm
}

module "lb" {
  source = "./module-lb"
  vpc_id = module.vpc.vpc_id
  subnet1_id = module.subnets.subnet1_id
  subnet2_id = module.subnets.subnet2_id
  allowed_ip = var.allowed_inbound_ips_lb
  acm_certificate_arn = module.acm.certificate_arn
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

    hosted_zone_id = var.hosted_zone_id
    record_name = var.record_name
}

module "asg" {
   source = "./module-asg"

    asg_desired_ec2 = var.asg_desired_ec2
    asg_min_ec2 = var.asg_min_ec2
    asg_max_ec2 = var.asg_max_ec2
    subnet1_id = module.subnets.subnet1_id
    subnet2_id = module.subnets.subnet2_id
    launch_template_id = module.compute.launch_template_id
    target_group_http_arn = module.lb.target_group_http_arn
    target_group_https_arn = module.lb.target_group_https_arn

}