package main

import (
	"flag"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func printStats(h *pcap.Handle) {
	for range time.Tick(1 * time.Second) {
		stats, err := h.Stats()
		if err != nil {
			log.Fatalln("pcap handle stats error:", err)
		}
		log.Printf(
			"recv/drop/ifdrop: %d/%d/%d",
			stats.PacketsReceived, stats.PacketsDropped,
			stats.PacketsIfDropped)
	}
}

func handleKafkaEvents(p *kafka.Producer) {
	for ev := range p.Events() {
		log.Println("kafka event:", ev)
	}
}

func main() {
	var err error

	iface := flag.String("i", "bond0", "Interface to read packets from")
	snaplen := flag.Int("s", 65536, "Sanp length (number of bytes max to read per packet")
	topic := flag.String("t", "pcap", "Kafka topic to write to")
	server := flag.String("k", "localhost:9092", "Kafka server to write to")
	flag.Parse()

	handle, err := pcap.OpenLive(*iface, int32(*snaplen), true, pcap.BlockForever)
	if err != nil {
		log.Fatalln("pcap handle opening error:", err)
	}

	go printStats(handle)

	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{"bootstrap.servers": *server})
	if err != nil {
		log.Fatalln("producer creation error:", err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range source.Packets() {
		producer.ProduceChannel() <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     topic,
				Partition: kafka.PartitionAny,
			},
			Value: packet.Data(),
		}
	}
}
