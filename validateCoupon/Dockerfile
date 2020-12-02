FROM golang:1.15

WORKDIR /go/src/microsservices-validate-coupon

COPY . . 

RUN GOOS=linux go build -o microsservices-validate-coupon

EXPOSE 9092

ENTRYPOINT [ "./microsservices-validate-coupon" ]