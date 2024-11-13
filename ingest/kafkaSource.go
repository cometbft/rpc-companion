package ingest

import (
	"fmt"
	"os"

	"github.com/cometbft/cometbft/rpc/grpc/client"
	events_v1 "github.com/cometbft/rpc-companion/api/events/v1"
	proto "github.com/cosmos/gogoproto/proto"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	schemaregistry "github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/protobuf"
)

func ForwardData(data Job[client.BlockResults]) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": "0.0.0.0:9092",

		// Fixed properties
		"acks": "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	schemaRegistryClient, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://0.0.0.0:8017"))
	if err != nil {
		panic("Failed to create serializer" + err.Error())
	}

	sp, err := protobuf.NewSerializer(schemaRegistryClient, serde.ValueSerde, protobuf.NewSerializerConfig())
	if err != nil {
		panic("Failed to serialize " + err.Error())
	}
	//fmt.Println(sp)
	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	topic := "quickstart-events"

	txResult := data.cometType.TxResults
	height := data.cometType.Height

	for _, item := range txResult {
		for _, e := range item.Events {
			eventID := "x"
			for _, a := range e.Attributes {
				if !a.Index {
					continue
				}
				key := a.Key
				event := new(events_v1.CompanionEvent)
				event.Height = height
				event.Key = key

				event.Value = a.Value
				event.Type = e.Type
				event.EventId = eventID

				//value, err := event.ProtoMessage()
				// val, err := sp.Serialize("quickstart-events", event)

				value, err := proto.Marshal(event)
				if err != nil {
					panic("failed to marshal event attribute" + err.Error())
				}
				p.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
					Key:            []byte(key),
					Value:          value,
				}, nil)
			}
		}
	}

	blockResults := data.cometType.FinalizeBlockEvents
	for _, item := range blockResults {
		for idx, a := range item.Attributes {
			if !a.Index {
				continue
			}
			eventID := string(height) + ":" + string(idx)
			key := a.Key
			event := events_v1.CompanionEvent{
				Height:  height,
				Key:     key,
				Value:   a.Value,
				Type:    item.Type,
				EventId: eventID,
			}
			// 	 	  	 z z zz	value, err := event.Marshal()
			val, err := sp.Serialize("quickstart-events", event)
			//value, err := event.Marshal()
			if err != nil {
				panic("failed to marshal event attribute" + err.Error())
			}
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Key:            []byte(key),
				Value:          val,
			}, nil)
		}
	}

	// Wait for all messages to be delivered
	p.Flush(15 * 1000)
	p.Close()
}
