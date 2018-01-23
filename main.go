package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pfring"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func printStats(r *pfring.Ring) {
	for range time.Tick(1 * time.Second) {
		stats, err := r.Stats()
		if err != nil {
			log.Fatalln("pfring ring stats error:", err)
		}
		log.Printf(
			"recv/drop: %d/%d | drop%%: %f%%",
			stats.Received, stats.Dropped,
			float64(stats.Dropped)/float64(stats.Received))
	}
}

func handleKafkaEvents(producer *kafka.Producer) {
	for ev := range p.Events() {
		log.Println("kafka event:", ev)
	}
}

func main() {
	var err error

	iface := flag.String("i", "bond0", "Interface to read packets from")
	snaplen := flag.Int("s", 65536, "Sanp length (number of bytes max to read per packet")
	flag.Parse()

	ring, err := pfring.NewRing(*iface, uint32(*snaplen), pfring.FlagPromisc)
	if err != nil {
		log.Fatalln("pfring ring creation error:", err)
	}

	err = ring.SetSocketMode(pfring.ReadOnly)
	if err != nil {
		log.Fatalln("pfring SetSocketMode error:", err)
	}

	err = ring.Enable()
	if err != nil {
		log.Fatalln("pfring Enable error:", err)
	}

	go printStats(ring)

	producer, err := kafka.NewProducer(
		&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalln("producer creation error:", err)
	}

	source := gopacket.NewPacketSource(ring, layers.LayerTypeEthernet)
	for packet := range source.Packets() {
		producer.ProduceChannel() <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value: packet.Source(),
		}
	}
}
