FROM nginx:alpine

# hapus default nginx
RUN rm -rf /usr/share/nginx/html/*

# copy semua file client ke nginx
COPY . /usr/share/nginx/html

# expose port
EXPOSE 80