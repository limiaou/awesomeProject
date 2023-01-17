package snow_flake

import (
	"errors"
	"sync"
	"time"
)

func GenerateServiceIdentities(service, length int) (identities []int64, err error) {
	worker := getSingleNode(int64(service))
	identities = make([]int64, 0)
	for i := 0; i < length; i++ {
		id := worker.Generate()
		identities = append(identities, id)
	}
	return
}
func getSingleNode(service int64) (singleAPINode *Node) {
	var once sync.Once
	once.Do(func() {
		node, err := NewNode(service)
		if err != nil {
			panic("snow flake error" + err.Error())
		}
		singleAPINode = node
	})
	return singleAPINode
}

func NewNode(node int64) (*Node, error) {
	if node < 0 || node > nodeMax {
		return nil, errors.New("Node number must be between 0 and 64")
	}
	return &Node{
		timestamp: 0,
		node:      node,
		step:      0,
	}, nil
}

func (n *Node) Generate() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.generate()
}

func (n *Node) generate() int64 {
	now := n.getMilliSeconds()
	if n.timestamp == now {
		n.step = (n.step + 1) & stepMax
		if n.step == 0 {
			for now <= n.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		n.step = 0
	}
	n.timestamp = now
	result := int64((now-Epoch)<<timeShift | (n.node << nodeShift) | (n.step))
	return result
}
func (n *Node) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}
