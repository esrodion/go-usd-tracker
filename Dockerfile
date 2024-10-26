FROM golang:1.22.5-alpine

WORKDIR /goserver
COPY . .

EXPOSE 8080

RUN ["go", "mod", "tidy"]
RUN ["go", "build", "-o", "main", "/goserver/cmd/."]

CMD ["./main"]
