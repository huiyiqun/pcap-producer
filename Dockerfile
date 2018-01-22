FROM ubuntu:16.04

# Enable repo of ntop
RUN apt-get update && apt-get upgrade -y && apt-get install -y wget lsb-release
RUN wget http://apt-stable.ntop.org/16.04/all/apt-ntop-stable.deb && dpkg -i apt-ntop-stable.deb && rm apt-ntop-stable.deb

# Install latest go
RUN wget https://dl.google.com/go/go1.9.2.linux-amd64.tar.gz -O- | tar xvz -C /opt

# Install dependencies
RUN apt-get update && apt-get install -y pfring git
ENV GOPATH=/go
RUN /opt/go/bin/go get github.com/google/gopacket/examples/pfdump
RUN /opt/go/bin/go get github.com/segmentio/kafka-go

RUN mkdir /opt/pcap-producer
WORKDIR /opt/pcap-producer
ADD main.go .
RUN /opt/go/bin/go build main.go

ENTRYPOINT ["./main"]
