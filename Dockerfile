FROM golang:1.26

WORKDIR /app
ADD . /app
RUN go mod tidy

CMD ["sleep", "infinity"]
