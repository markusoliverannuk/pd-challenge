
# Challenge PD

Hello! :)
I want to take a moment of your time to explain how this application works and how the whole infrastructure is built. 
I've divided this into 3 main sections:,<b>
 
1) Infrastructure<br>

2) EC2 Startup scripts<br>

3) Application code<br>

4) Running the application locally<br></b>

I'll try to move from the outermost layer all the way to the logic of our application.


## 1. Infrastructure

![App Screenshot](schemas/architecture.png)

The whole infrastructure for the API is written as code on <b>Terraform</b> for <b>AWS</b>. It consists of:<br>
• Hosted zones, records<br>
• IAM Roles, Policies, Attachments<br>
• Systems Manager Parameter Store as secure storage for the Pipedrive API key and Github Access Token<br>
• ACM Certs for AWS managed services (Application Load Balancer, Cloudfront CDN), Certbot on individual machines<br>
• VPC, Subnets, Security Groups, ACLs.<br>
• Internet Gateway<br>
• Application Load Balancer<br>
• Auto Scaling Group<br>
• EC2 instances running the latest version of the docker image from DockerHub, served through NGINX<br>
<br>
So by applying the code written for the infrastructure (I know what you might be asking "Where is TDD for Infrastructure?!". Yes, for this project I decided not to.), we are provisioning all the necessary resources on AWS and automatically fetching the necessary env variables from <b>Parameter Store</b>setting up the <b>Docker containers</b>, <b>requesting TLS certs</b> and <b>configuring NGINX</b> with newly requested certs on all active EC2 machines.

So let's get into it..

I apologize for all the incoming "Once that is done"s, "Then"s and "Next"s.

It starts by creating the VPC with a CIDR block of 10.0.0.0/16.

Once that is done we are requesting certificates for the load balancer (api-challenge-v2.techwithmarkus.com).

Once that is done it validates the ACM Certificates through DNS by creating records in our hosted zone.

Next, we get the internet gateway up and running.

Once that is done, it procceeds with the creation of two subnets in two separate availability zones, eu-north-1a & b (subnets.tf).

Next in line are two target groups, one intended for initial HTTP traffic for Certbot challenges and and the second one to serve all incoming application traffic. 

Next we are creating security groups for the target instances, which allow incoming connections from the load balancers security group. 

Now we have route tables and route tables associations with our subnets that we created earlier.

Once that is done, it moves on to creating our load balancer which forwards traffic to our 2 target groups depending on the rules. 

Next we create the launch template for our instances, currently I set them to run on t2.small Ubuntu 24.04 machines using the AMI ami-04b70fa74e45c3917.

Now it will be setting up the auto scaling group with a minimum amount of 1, desired amount of 2, and maximum amount of 3 instances spread across 2 availability zones (eu-north-1a & b).

Once that is done, it updates the Route53 records for api-challenge-v2.techwithmarkus.com to point to the IPv4 DNS of our Load Balancer.

Almost done! Now all that it has left to do is create the listeners for the load balancer (2).
<br>

## 2. Startup script on EC2 machines
![Instances](schemas/instances.png)
![Startup Script](schemas/startupscript.png.png)

Alright, so now that all resources are launched on AWS we have to have a way for the EC2 machine that are in the target groups, be able to serve our application. For that I've written a userdata script that each newly launched machine executes. I'll give you a quick explanation of what each part of it does :).

So it starts off by sleeping for 90 seconds, this pretty much gives enough time for the target groups and load balancer to finish configuring, otherwise we cannot respond to certbot challenges which sends a request which is eventually forwarded to the load balancer. So we need the load balancer to be in working order by that time. We can tweak it, currently set at 90 sec.

So next what we do is install the aws cli.

Once that is done we send requests to Systems Manager Parameter Store to get access to the Pipedrive API key and Github Access token, we query just the value and then echo them into envfile.env which we eventually will be passing into the docker container.


Then what it does, is it downloads package info from configured sources.

Next we install Certbot, Nginx, and Docker.

Now we start a new bash shell to insert a server configuration for NGINX. All it pretty much does is it tells the server to listen on port 80 (we will need this for the certbot challenge request).

After that is done we reload & start nginx. Now we are able to requests certs from certbot.

Moving on, we write a one-liner to request certificates for api-challenge-v2.techwithmarkus.com in non interactive mode, and we insert all the necessary details beforehand (like the agreement to tos and our email)

Now we start another shell to insert a new :443 listner into our NGINX configuration, provide it the certificates that certbot issued and tell it to forward traffic to the container port of the Docker container that we will be running shortly.

Now we will reload nginx for the changes to take effect.

Just in case we will make sure there is no matching container running (if there is we stop and remove it).

Next, we will pull the latest docker image of mannuk24/challenge from DockerHub with auto-restart enabled only in case of an error exit, pass it the environment file that we created earlier and run it, exposing port 8050.

![Inside machine](schemas/insidemachine.png)
![Pipedrive dashboard](schemas/pipedrive.png)

## 4. Running the application locally

If you wish to run the application locally, I'll now show you how to do so.

For local execution there are a few prerequisites:<br>

• You have a valid Github Access Token<br>
• You have a valid Pipedrive API key<br>
• You have golang installed on your machine<br>

Here are the steps you should take for local execution:
1) Clone this repository to a location on your machine using <b>git clone https://github.com/markusoliverannuk/challenge</b><br>
2) cd inside the directory that you just pulled. If for example you clowned the repository under 'Documents' and the repository name is <b>challenge</b>, you can execute <b>cd challenge</b> from your Documents folder. Pointed this out just in case :)<br>
3) Now make sure you have your Pipedrive API key and Github Access Token ready, we'll be using them quite soon.</br>
4) Now depending if you want to use the graphical user interface as well, or just go off of the JSON responses you can choose to either run only the golang api, or also the React client.<br>
5) Make sure you're currently inside the <b>challenge</b> directory and execute the following commands depending on your need.<br><br>
Option A) "I only want access to the response from the API without graphical user interface" - In your terminal enter <b>sh local-startup-api.sh</b>.<br>
- You should now be able the send POST requests to localhost:8050/user/{their Github username} and GET requests localhost:8050/trackedusers. You can do it either from a tool like Postman, curl, or for example from your web browser.<br><br>
Option B) "I want to view the results through the GUI" - Open 2 terminals. Make  sure both of them have the current location as the <b>challenge</b> directory. In one terminal window execute <b>local-startup-api.sh</b>. In the second terminal window execute <b>sh local-startup-ui.sh</b>.<br>
- You should be able to open the graphical user interface from your web browser by entering localhost:3000 which sends requests to our API and displays the response neatly.<br>
- If you wish, you may still view the API responses as pointed outt in Option A.<br>
<br>
When you now send a request, you will wait for a moment and then be presented with both "old gists" and "new gists" from the user. "new gists" are gists from that user that you haven't fetched before. "old gists" are gists that you have already seen through my application from the certain user.<br>
PS: You can view all the users you're already tracking by either clicking the button "View tracked users" in the GUI or by sending a GET request to localhost:8050/trackedusers.

If you have any questions regarding this project, then I'm always available via e-mail (<b>annukmarkusoliver@gmail.com</b>) or on LinkedIn.