server {

    listen 80;
    server_name fileserver.s3apis.com;

    return 301 https://fileserver.s3apis.com$request_uri;
}


upstream fileserverblue  {
   server localhost:5000  weight=10 max_fails=3  fail_timeout=15s;
}

upstream fileservergreen  {
   server localhost:5001  weight=10 max_fails=3  fail_timeout=15s;
}

server { 
    #listen 80;
    listen 443 ssl;
    server_name fileserver.s3apis.com;

    include fileserver-bluegreen.conf;
    access_log /var/log/nginx/fileserver.log;


    location / {

      #proxy_set_header        Host $host;
      proxy_set_header        Host $http_host;
      proxy_set_header        X-Real-IP $remote_addr;
      proxy_set_header        X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header        X-Forwarded-Proto $scheme;
      #proxy_redirect off;

      # Fix the “It appears that your reverse proxy set up is broken" error.
      proxy_pass              http://$activeBackend;

      proxy_read_timeout  90;

    }
 
    ssl_certificate /etc/letsencrypt/live/fileserver.s3apis.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/fileserver.s3apis.com/privkey.pem; # managed by Certbot
}

