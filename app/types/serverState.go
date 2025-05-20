package types

type ServerState struct {
	Port             int
	Role             string // master | slave
	MasterHost       string // empty for master
	MasterPort       string // empty for master
	MasterReplId     string
	MasterReplOffset int // 0 for master
}
