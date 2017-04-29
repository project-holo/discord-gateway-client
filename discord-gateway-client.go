package main

import (
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/gmallard/stompngo"
)

// Configuration variables.
var (
	discordToken      string
	eventsDestination string
	shardCount        int
	shardID           int
	brokerURI         string
	debugMode         bool
)

func init() {
	// Load configuration from environment
	discordToken = os.Getenv("DISCORD_TOKEN")
	eventsDestination = os.Getenv("EVENTS_DESTINATION")
	shardCount, _ = strconv.Atoi(os.Getenv("SHARD_COUNT"))
	shardID, _ = strconv.Atoi(os.Getenv("SHARD_ID"))
	brokerURI = os.Getenv("BROKER_URI")
	debugMode, _ = strconv.ParseBool(os.Getenv("DEBUG"))

	// Parse configuration flags from command-line
	flag.StringVar(&discordToken, "token", "", "* Discord auth token")
	flag.StringVar(&eventsDestination, "events-dest", "/events", "Broker events destination")
	flag.IntVar(&shardCount, "shard-count", 0, "Shard count")
	flag.IntVar(&shardID, "shard", 0, "Shard ID")
	flag.StringVar(&brokerURI, "broker", "", "* Broker connection URI")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode")
	flag.Parse()

	// Debug mode
	if debugMode != false {
		log.SetLevel(log.DebugLevel)
	}

	// Print flags to debug
	log.WithFields(log.Fields{
		"discordToken": discordToken,
		"shardCount":   shardCount,
		"shardID":      shardID,
		"brokerURI":    brokerURI,
		"debugMode":    debugMode,
	}).Debug("Flags")
}

func main() {
	var err error

	// Connect to STOMP broker
	stomp, err := createStompConnection(brokerURI)
	if err != nil {
		log.WithField("error", err).Fatal("failed to create connection to STOMP broker")
	}
	log.Debug("connected to STOMP broker")

	// Create Discord client
	// Create session
	discord, err := discordgo.New(discordToken)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("failed to create Discord session")
	}

	// Fetch user information for the user, used to check for valid token
	// and to print user information to debug
	me, err := discord.User("@me")
	if err != nil {
		log.WithField("error", err).Fatal("failed to GET /users/@me from Discord API, token may be invalid")
	}
	log.Debug("fetched /users/@me from Discord API")
	log.Debugf("@me=%v", me.ID)
	log.Debugf("persona=%v#%v", me.Username, me.Discriminator)

	// Shard parameters
	discord.ShardID = shardID
	discord.ShardCount = shardCount
	if discord.ShardCount <= 0 {
		discord.ShardCount = 1
	}
	log.Debugf("shard=%v/%v", discord.ShardID, discord.ShardCount)

	// Ready logger
	discord.AddHandlerOnce(func(s *discordgo.Session, e *discordgo.Ready) {
		// s.UpdateStatus(0, "ProjectHOLO")
		log.Debug("READY")
	})

	// Add Discord gateway event handler
	discord.AddHandler(func(s *discordgo.Session, e *discordgo.Event) {
		if e.Operation != 0 || e.Type == "" {
			return
		}
		if e.Struct == nil {
			err := json.Unmarshal(e.RawData, &e.Struct)
			if err != nil {
				log.Warn("failed to unmarshal event without discordgo struct")
			}
		}

		serializeAndDispatchEvent(discord, stomp, e.Type, e.Struct)
	})

	// Connect to the Discord gateway
	err = discord.Open()
	if err != nil {
		log.WithField("error", err).Fatal("failed to connect to the Discord gateway")
		return
	}
	log.Debug("opened Discord gateway connection")

	// Ready
	log.Infof("%v is ready to rumble!", me.Username)

	// Wait for a SIGINT to exit
	log.Info("press CTRL+C to exit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c

	// Close connections to Discord and message broker, timeout after 5s
	log.Info("Exiting, press CTRL+C again to force exit.")
	go func() {
		err = discord.Close()
		if err != nil {
			log.WithField("error", err).Error("Failed to close Discord connection")
		}
		log.Debug("closed Discord connection")
		err = stomp.Disconnect(stompngo.Headers{})
		if err != nil {
			log.WithField("error", err).Error("Failed to disconnect from STOMP broker")
		}
		log.Debug("closed STOMP broker connection")
		c <- syscall.SIGINT
	}()
	time.AfterFunc(5*time.Second, func() {
		log.Debug("5s has passed, force exiting")
		c <- syscall.SIGINT
	})
	<-c
}
