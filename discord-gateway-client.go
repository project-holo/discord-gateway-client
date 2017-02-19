package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	stompgo "github.com/go-stomp/stomp"

	discordmanager "github.com/project-holo/discord-gateway-client/discord"
	stompmanager "github.com/project-holo/discord-gateway-client/stomp"
)

const (
	messageCreate = "MESSAGE_CREATE"
	ready         = "READY"
)

var (
	discord           *discordgo.Session
	stomp             *stompgo.Conn
	eventsDestination = "/events"
)

type event struct {
	Type    string      `json:"type"`
	ShardID int         `json:"shard_id"`
	Data    interface{} `json:"data"`
}

// serializeAndDispatchEvent serializes and sends data to the events destination
// on the STOMP broker.
func serializeAndDispatchEvent(Type string, data interface{}) {
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
	t := stomp.Begin()
	err = t.Send(eventsDestination, "application/json", []byte(j))
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("Failed to send a message to the STOMP broker")
		return
	}
	err = t.Commit()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Errorf("Failed to commit transaction to STOMP broker")
		return
	}
}

// onMessageCreate serializes and sends incoming messages to the STOMP events
// destination.
func onMessageCreate(s *discordgo.Session, m *discordgo.Message) {
	serializeAndDispatchEvent(messageCreate, *m)
}

// onReady serializes and sends ready packets to the STOMP events destination.
func onReady(s *discordgo.Session, e *discordgo.Ready) {
	//s.UpdateStatus(0, "ProjectHOLO")
	serializeAndDispatchEvent(ready, *e)
}

func main() {
	var (
		Debug      = flag.Bool("d", false, "Debug mode (bool)")
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

	// Add event handlers
	discord.AddHandler(onMessageCreate) // MESSAGE_CREATE
	discord.AddHandler(onReady)         // READY

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
