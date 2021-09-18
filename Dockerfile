FROM golang:1.16.7-alpine3.14
RUN mkdir /anticap
ADD . /anticap
WORKDIR /anticap
#Build
RUN go build -o /anticap/anticap ./cmd/anticap/anticap.go
#Run
CMD ["/anticap/anticap"]