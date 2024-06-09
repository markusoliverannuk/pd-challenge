#!/bin/bash

#installing necessities
sudo apt-get update
sudo apt-get install -y nginx certbot python3-certbot-nginx docker.io

#configuring nginx

sudo bash -c 'cat <<EOF > /etc/nginx/sites-available/default
server {
    listen 80;
    server_name api-challenge.techwithmarkus.com;

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

sudo certbot certonly --webroot -w /var/www/html -d api-challenge.techwithmarkus.com --non-interactive --agree-tos -m annukmarkusoliver@gmail.com

sudo bash -c 'cat <<EOF > /etc/nginx/sites-available/default
server {
    listen 80;
    server_name api-challenge.techwithmarkus.com;

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
    server_name api-challenge.techwithmarkus.com;

    ssl_certificate /etc/letsencrypt/live/api-challenge.techwithmarkus.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api-challenge.techwithmarkus.com/privkey.pem;

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

sudo docker run -d --name challenge -p 8050:8050 mannuk24/challenge:latest