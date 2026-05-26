# dev.frontend.Dockerfile  //frontend arayuzu render ıcın 
FROM nginx:alpine
COPY . /usr/share/nginx/html
