FROM golang:1.22.5

WORKDIR /goserver
COPY . .

EXPOSE 8080
EXPOSE 8081

RUN ["go", "build", "-o", "main", "/goserver/cmd/."]

CMD ["./main"]
