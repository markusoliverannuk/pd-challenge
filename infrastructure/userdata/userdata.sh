#!/bin/bash

#we will give a moment for the load balancer and target group to be configured, we can also use a while true loop to check, but this seems to work better
sleep 120

#installing the awscli
sudo snap install aws-cli --classic

GITHUB_AT=$(aws ssm get-parameter --name "/github/access_token" --region us-east-1 --with-decryption --query "Parameter.Value" --output text)
PIPEDRIVE_API_KEY=$(aws ssm get-parameter --name "/pipedrive/api_key" --region us-east-1 --with-decryption --query "Parameter.Value" --output text)

# creating an env file that I'll pass to docker at a later stage. injecting the keys into it
sudo echo "GITHUB_AT=${GITHUB_AT}" >> /home/ubuntu/envfile.env # using >> because we have multiple injections, otherwise it will overwrite
sudo echo "PIPEDRIVE_API_KEY=${PIPEDRIVE_API_KEY}" >> /home/ubuntu/envfile.env



#installing necessities
sudo apt-get update
sudo apt-get install -y nginx certbot python3-certbot-nginx docker.io

#configuring nginx

sudo bash -c 'cat <<EOF > /etc/nginx/sites-available/default
server {
    listen 80;
    server_name api-challenge-v2.techwithmarkus.com;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
        try_files \$uri =404;
    }

    location / {
        return 301 https://\$host\$request_uri;
    }
}
EOF
'


sudo systemctl start nginx

sudo mkdir -p /var/www/html

sudo systemctl reload nginx

#requesting certs

sudo certbot certonly --webroot -w /var/www/html -d api-challenge-v2.techwithmarkus.com --non-interactive --agree-tos -m annukmarkusoliver@gmail.com

sudo bash -c 'cat <<EOF > /etc/nginx/sites-available/default
server {
    listen 80;
    server_name api-challenge-v2.techwithmarkus.com;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
        try_files \$uri =404;
    }

    location / {
        return 301 https://\$host\$request_uri;
    }
}

server {
    listen 443 ssl;
    server_name api-challenge-v2.techwithmarkus.com;

    ssl_certificate /etc/letsencrypt/live/api-challenge-v2.techwithmarkus.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api-challenge-v2.techwithmarkus.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8050;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF
'

sudo systemctl reload nginx

sudo docker stop challenge || true
sudo docker rm challenge || true

sudo docker pull mannuk24/challenge:latest

sudo docker run -d --name challenge --env-file /home/ubuntu/envfile.env -p 8050:8050 --restart unless-stopped mannuk24/challenge:latest