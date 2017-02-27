package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/gmallard/stompngo"

	discordmanager "github.com/project-holo/discord-gateway-client/discord"
	stompmanager "github.com/project-holo/discord-gateway-client/stomp"
)

// Event message send parameters.
const (
	eventsDestination = "/events"
	consumeWindow     = time.Minute * 5 // 5 minute window
)

// Storage variables for Discord and STOMP connections.
var (
	discord *discordgo.Session
	stomp   *stompngo.Connection
)

// Default headers to send with all STOMP event messages.
var defaultEventMessageHeaders = stompngo.Headers{
	stompngo.HK_DESTINATION, eventsDestination,
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

// serializeAndDispatchEvent serializes and sends data to the events destination
// on the STOMP broker.
func serializeAndDispatchEvent(Type string, data interface{}) {
	// JSON encode data
	j, err := json.Marshal(event{
		Type:    Type,
		ShardID: discord.ShardID,
		Data:    data,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("Failed to serialize %v event", Type)
		return
	}

	// Construct message headers
	var body = []byte(j)
	h := stompngo.Headers{
		"expires", strconv.FormatInt(time.Now().Add(consumeWindow).UnixNano()/1000000, 10),
		stompngo.HK_CONTENT_LENGTH, strconv.Itoa(len(body)),
	}

	// Send the message
	err = stomp.SendBytes(defaultEventMessageHeaders.AddHeaders(h), body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("Failed to send a message to the STOMP broker")
	}
}

func main() {
	var (
		Debug      = flag.Bool("d", false, "Enable debug mode")
		ShardCount = flag.String("c", "0", "Shard count")
		ShardID    = flag.String("s", "0", "Shard ID")
		StompURI   = flag.String("b", "", "STOMP broker connection URI")
		Token      = flag.String("t", "", "Discord auth token")
		err        error
	)
	flag.Parse()

	// Debug mode
	if *Debug != false {
		log.SetLevel(log.DebugLevel)
	}

	// Print flags to debug
	log.WithFields(log.Fields{
		"Debug":      *Debug,
		"ShardCount": *ShardCount,
		"ShardID":    *ShardID,
		"StompURI":   *StompURI,
		"Token":      *Token,
	}).Debug("Flags")

	// Connect to STOMP broker
	s, err := stompmanager.CreateStompConnection(*StompURI)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create connection to STOMP broker")
	}
	stomp = s
	log.Debug("connected to STOMP broker")

	// Create Discord client
	shardID, _ := strconv.Atoi(*ShardID)
	shardCount, _ := strconv.Atoi(*ShardCount)
	d, me, err := discordmanager.CreateDiscordClient(*Token, shardID, shardCount)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create Discord client session")
		return
	}
	discord = d
	log.Debug("created Discord session")

	// Add Discord gateway event handlers
	for i := 0; i < len(eventHandlers); i++ {
		discord.AddHandler(eventHandlers[i])
	}
	log.Debugf("attached %v event handlers to Discord session", len(eventHandlers))

	// Connect to the Discord gateway
	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to connect to the Discord gateway")
		return
	}
	log.Debug("opened Discord gateway connection")

	// Ready
	log.Infof("%v is ready to rumble!", me.Username)

	// Wait for a signal to exit
	log.Info("Press CTRL+C to exit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
