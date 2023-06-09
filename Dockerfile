FROM golang:1.19

WORKDIR /app
ADD . /app
RUN go mod tidy

CMD ["sleep", "infinity"]
