
# Challenge PD

Hello! :)
I want to take a moment of your time to explain how this application works and how the whole infrastructure is built. 
I've divided this into 3 main sectors


## 1. Infrastructure

![App Screenshot](schemas/architecture.png)

The whole infrastructure is written as code on <b>Terraform</b> for AWS. It consists of (not in any particular order):<br>
• Hosted zones, records<br>
• ACM Certs for AWS managed services (Application Load Balancer, Cloudfront CDN), Certbot on individual machines<br>
• VPC, Subnets, Security Groups, ACLs.<br>
• Internet Gateway<br>
• Application Load Balancer<br>
• Auto Scaling Group<br>
• EC2 Images running the latest version of the docker image from DockerHub, served through NGINX<br>
• Cloudfront CDN<br>
• S3 Website Endpoint<br>
<br>
By applying the code written for the infrastructure, we are provisioning all the necessary resources on AWS and automatically setting up the docker containers, requesting TLS certs and configuring NGINX with newly requested certs on all active EC2 machines.

It starts by creating the VPC with a CIDR block of 10.0.0.0/16.

Once that is done we are requesting certificates for both the load balancer (api-challenge.techwithmarkus.com) and Cloudfront (challenge.techwithmarkus.com).

Once that is done it validates the ACM Certificates through DNS by creating records in our hosted zone. Next, we get the internet gateway up and running.

Once that is done, it procceeds with the creation of two subnets in two separate availability zones, us-east-1a & b (subnets.tf).

Next in line are two target groups, one intended for initial HTTP traffic for Certbot challenges and and the second one to serve all incoming application traffic. 

Next we are creating security groups for the target instances, which allow incoming connections from the load balancers security group. 

Now we have route tables and route tables associations with our subnets that we created earlier.

Once that is done, it moves on to creating our load balancer which forwards traffic to our 2 target groups depending on the rules. 

Next we create the launch template for our instances, currently I set them to run on t2.small Ubuntu 24.04 machines using the AMI ami-04b70fa74e45c3917.

Now it will be setting up the auto scaling group with a minimum amount of 1, desired amount of 2, and maximum amount of 3 instances spread across 2 availability zones (us-east-1a & b).

Once that is done, it updates the Route53 records for api-challenge.techwithmarkus.com to point to the IPv4 DNS of our Load Balancer.

Almost done! Now all that it has left to do is create the listeners for the load balancer (2).

