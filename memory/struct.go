package memory

type NftablesRecord struct {
	ID        int
	EventUUID string
	SrcIP     string
	Nic       string
	DstIP     string
	DstPort   string
	Mac       string
	Protocol  string
	Timestamp int64
}
