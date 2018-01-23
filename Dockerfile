FROM ubuntu:18.04

# Install dependencies
RUN apt-get update && apt-get install -y git golang librdkafka-dev pkg-config libpcap-dev

ENV GOPATH=/go
RUN go get github.com/google/gopacket
RUN go get github.com/confluentinc/confluent-kafka-go/kafka

RUN mkdir /opt/pcap-producer
WORKDIR /opt/pcap-producer
ADD main.go .
RUN go build main.go

ENTRYPOINT ["./main"]
