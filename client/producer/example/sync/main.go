package sync

import (
	"fmt"
	"log"

	"github.com/corrots/kafka/client/producer/client/sync"
)

func main() {
	p, err := sync.NewProducer()
	if err != nil {
		log.Fatal(err)
	}
	p.SendMessage()
	fmt.Println("done")
}
