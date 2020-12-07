FROM golang:1.15

WORKDIR /go/src/microsservices-validate-credcard

COPY . . 

RUN GOOS=linux go build -o microsservices-validate-credcard

EXPOSE 9093

ENTRYPOINT [ "./microsservices-validate-credcard" ]