FROM golang@1.22.5:alpine

WORKDIR /goserver
COPY . .

EXPOSE 8080

CMD ["go", "run", "/goapiserver/cmd/."]
