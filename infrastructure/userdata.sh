#!/bin/bash

# Installing Docker
sudo apt-get update
sudo apt-get install -y nginx certbot python3-certbot-nginx docker.io

# Nginx configuration for ACME challenge
sudo bash -c 'cat <<EOF > /etc/nginx/sites-available/default
server {
    listen 80;
    server_name api-pd-challenge.techwithmarkus.com;

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
# Ensure the webroot directory exists
sudo mkdir -p /var/www/html

# Reload Nginx to apply the configuration
sudo systemctl reload nginx

# Obtain the SSL certificate using Certbot
sudo certbot certonly --webroot -w /var/www/html -d api-pd-challenge.techwithmarkus.com --non-interactive --agree-tos -m annukmarkusoliver@gmail.com

# Nginx configuration for HTTPS
sudo bash -c 'cat <<EOF > /etc/nginx/sites-available/default
server {
    listen 80;
    server_name api-pd-challenge.techwithmarkus.com;

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
    server_name api-pd-challenge.techwithmarkus.com;

    ssl_certificate /etc/letsencrypt/live/api-pd-challenge.techwithmarkus.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api-pd-challenge.techwithmarkus.com/privkey.pem;

    location / {
        proxy_pass http://localhost:8050;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF
'

# Reload Nginx to apply the HTTPS configuration
sudo systemctl reload nginx

# Stop and remove the existing Docker container if it exists
sudo docker stop challenge || true
sudo docker rm challenge || true

# Pull the latest Docker image
sudo docker pull mannuk24/challenge:latest

# Run the Docker container
sudo docker run -d --name challenge -p 8050:8050 mannuk24/challenge:latest
