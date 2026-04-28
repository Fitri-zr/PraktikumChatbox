# pakai Go official
FROM golang:1.22-alpine

# set working directory
WORKDIR /app

# copy semua file
COPY . .

# download dependency
RUN go mod tidy

# build aplikasi
RUN go build -o main .

# expose port
EXPOSE 8080

# jalankan
CMD ["./main"]