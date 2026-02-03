package mesh

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hashicorp/memberlist"
)

// Define o tipo de mensagem que circula no cluster
type BroadcastMessage struct {
	Type  string `json:"type"` // "kv_update"
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Node struct {
	List *memberlist.Memberlist
	Port int
	Name string
	// Removemos a Queue de Broadcast complexa para simplificar

	OnKVUpdate func(key, value string)
}

type NodeMetadata struct {
	Functions []string `json:"functions"`
	APIPort   int      `json:"api_port"`
}

type delegate struct {
	node *Node
	meta NodeMetadata
}

func (d *delegate) NodeMeta(limit int) []byte {
	b, _ := json.Marshal(d.meta)
	return b
}

func (d *delegate) NotifyMsg(b []byte) {
	// 📩 AQUI: Recebe a mensagem direta enviada por SendReliable
	var msg BroadcastMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		return
	}

	if msg.Type == "kv_update" && d.node.OnKVUpdate != nil {
		d.node.OnKVUpdate(msg.Key, msg.Value)
	}
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte { return nil }
func (d *delegate) LocalState(join bool) []byte                { return nil }
func (d *delegate) MergeRemoteState(buf []byte, join bool)     {}

func NewNode(bindPort int, advertisePort int, seeds []string, secretKey string, functions []string, apiPort int) (*Node, error) {
	config := memberlist.DefaultLANConfig()
	config.BindPort = bindPort
	config.AdvertisePort = advertisePort
	hostname, _ := config.Name, ""
	config.Name = fmt.Sprintf("%s-%d", hostname, bindPort)
	config.LogOutput = io.Discard

	if secretKey != "" {
		if len(secretKey) != 32 {
			return nil, fmt.Errorf("secret key must be 32 bytes")
		}
		config.SecretKey = []byte(secretKey)
	}

	node := &Node{
		Port: bindPort,
		Name: config.Name,
	}

	d := &delegate{
		node: node,
		meta: NodeMetadata{
			Functions: functions,
			APIPort:   apiPort,
		},
	}
	config.Delegate = d

	list, err := memberlist.Create(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create memberlist: %w", err)
	}
	node.List = list

	if len(seeds) > 0 {
		_, err := list.Join(seeds)
		if err != nil {
			log.Printf("⚠️ Mesh: Failed to join seeds %v: %v", seeds, err)
		} else {
			log.Printf("✅ Mesh: Joined cluster via %v. Members: %d", seeds, list.NumMembers())
		}
	}

	return node, nil
}

// BroadcastKV envia a atualização DIRETAMENTE para todos os pares conhecidos
func (n *Node) BroadcastKV(key, value string) {
	msg := BroadcastMessage{
		Type:  "kv_update",
		Key:   key,
		Value: value,
	}
	payload, _ := json.Marshal(msg)

	// Itera sobre todos os membros conhecidos
	for _, member := range n.List.Members() {
		// Não manda pra si mesmo
		if member.Name == n.Name {
			continue
		}

		// SendReliable (TCP) garante que a mensagem chegue e dispara NotifyMsg no destino
		// Para máxima performance poderíamos usar SendBestEffort (UDP), mas Reliable é melhor para KV.
		err := n.List.SendReliable(member, payload)
		if err != nil {
			log.Printf("⚠️ Failed to sync KV with %s: %v", member.Name, err)
		}
	}
}

func (n *Node) Shutdown() {
	if n.List != nil {
		n.List.Leave(time.Second)
		n.List.Shutdown()
	}
}

func (n *Node) GetPeers() []string {
	var peers []string
	for _, m := range n.List.Members() {
		peers = append(peers, fmt.Sprintf("%s (%s:%d)", m.Name, m.Addr, m.Port))
	}
	return peers
}
