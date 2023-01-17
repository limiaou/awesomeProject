package snow_flake

import "sync"

type Node struct {
	mu        sync.Mutex
	timestamp int64
	node      int64
	step      int64
}

const (
	//nodeBits uint8 = 10
	nodeBits uint8 = 6
	//stepBits uint8 = 12
	stepBits  uint8 = 6
	nodeMax   int64 = -1 ^ (-1 << nodeBits)
	stepMax   int64 = -1 ^ (-1 << stepBits)
	timeShift uint8 = nodeBits + stepBits
	nodeShift uint8 = stepBits
)

var Epoch int64 = 1564588800000
