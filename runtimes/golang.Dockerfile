FROM golang:1.14-alpine
 
RUN mkdir -p /usr/app
WORKDIR /usr/app
 
COPY . .
RUN go build -o main
 
EXPOSE 8000
CMD ["./main"]