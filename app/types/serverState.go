package types

type ServerState struct {
	Port             int
	Role             string // master | slave
	MasterHost       string // empty for master
	MasterPort       string // empty for master
	MasterReplId     string
	MasterReplOffset int // 0 for master
	Replicas         []Replica
	AckOffset        int // offset of the last acknowledged replication message (only for slaves)
	BytesSent        int // number of bytes sent to replicas (only for masters)

	DBDir      string
	DBFileName string
}
