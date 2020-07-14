package kafka

import (
	"log"

	"github.com/Shopify/sarama"
)

type msgConsumerGroupHandler struct {
	current chan *sarama.ConsumerMessage
}

func (msgConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (msgConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h msgConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
		h.current <- msg
	}
	return nil
}
