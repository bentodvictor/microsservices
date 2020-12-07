FROM golang:1.15

WORKDIR /go/src/microsservices-checkout

COPY . . 

RUN GOOS=linux go build -o microsservices-checkout

EXPOSE 9090

ENTRYPOINT [ "./microsservices-checkout" ]