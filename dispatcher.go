package main

import (
	"encoding/json"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/gmallard/stompngo"
)

// consumeWindow is the amount of time the broker is permitted to queue the message for (without being consumed) before
// removing it from the queue (and possibly sending to dead letter queue).
const consumeWindow = time.Minute * 5 // 5 minute window

// Default headers to send with all STOMP event messages.
var defaultEventMessageHeaders = stompngo.Headers{
	"persistent", "true",
	"priority", "10",
	stompngo.HK_CONTENT_TYPE, "application/json; charset=utf8",
}

// Event data struct for all STOMP event messages.
type event struct {
	Type    string      `json:"type"`
	ShardID int         `json:"shard_id"`
	Data    interface{} `json:"data"`
}

// serializeAndDispatchEvent serializes and sends data to the events destination on the STOMP broker.
func serializeAndDispatchEvent(d *discordgo.Session, stomp *stompngo.Connection, eventType string, data interface{}) {
	// JSON encode data
	body, err := json.Marshal(event{
		Type:    eventType,
		ShardID: d.ShardID,
		Data:    data,
	})
	if err != nil {
		log.WithField("error", err).Errorf("failed to serialize %v event", eventType)
		return
	}

	// Construct message-specific headers
	h := stompngo.Headers{
		stompngo.HK_DESTINATION, eventsDestination,
		"expires", strconv.FormatInt(time.Now().Add(consumeWindow).UnixNano()/1000000, 10),
		stompngo.HK_CONTENT_LENGTH, strconv.Itoa(len(body)),
	}

	// Send the message
	err = stomp.SendBytes(defaultEventMessageHeaders.AddHeaders(h), body)
	if err != nil {
		log.WithField("error", err).Errorf("failed to send a message to the STOMP broker")
	}
}
