package main

import (
	"flag"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pfring"
	"log"
)

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

	source := gopacket.NewPacketSource(ring, layers.LayerTypeEthernet)
	for packet := range source.Packets() {
		log.Printf("Packet %d: %v", count, packet.Data())
	}
}
