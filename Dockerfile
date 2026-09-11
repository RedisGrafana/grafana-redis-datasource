FROM golang:1.25.6

WORKDIR /app
ADD . /app
RUN go mod tidy

CMD ["sleep", "infinity"]
