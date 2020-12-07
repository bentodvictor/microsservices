FROM golang:1.15

WORKDIR /go/src/microsservices-payment

COPY . . 

RUN GOOS=linux go build -o microsservices-payment

EXPOSE 9091

ENTRYPOINT [ "./microsservices-payment" ]