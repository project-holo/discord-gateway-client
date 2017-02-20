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

	// Add Discord gateway event handlers (in order of list in events.go)
	discord.AddHandler(onReady)   // READY
	discord.AddHandler(onResumed) // RESUMED

	discord.AddHandler(onChannelCreate) // CHANNEL_CREATE
	discord.AddHandler(onChannelUpdate) // CHANNEL_UPDATE
	discord.AddHandler(onChannelDelete) // CHANNEL_DELETE

	discord.AddHandler(onGuildCreate) // GUILD_CREATE
	discord.AddHandler(onGuildUpdate) // GUILD_UPDATE
	discord.AddHandler(onGuildDelete) // GUILD_DELETE

	discord.AddHandler(onGuildBanAdd)    // GUILD_BAN_ADD
	discord.AddHandler(onGuildBanRemove) // GUILD_BAN_REMOVE

	discord.AddHandler(onGuildEmojisUpdate) // GUILD_EMOJIS_UPDATE

	discord.AddHandler(onGuildIntegrationsUpdate) // GUILD_INTEGRATIONS_UPDATE

	discord.AddHandler(onGuildMemberAdd)    // GUILD_MEMBER_ADD
	discord.AddHandler(onGuildMemberRemove) // GUILD_MEMBER_REMOVE
	discord.AddHandler(onGuildMemberUpdate) // GUILD_MEMBER_UPDATE
	discord.AddHandler(onGuildMembersChunk) // GUILD_MEMBERS_CHUNK

	discord.AddHandler(onGuildRoleCreate) // GUILD_ROLE_CREATE
	discord.AddHandler(onGuildRoleUpdate) // GUILD_ROLE_UPDATE
	discord.AddHandler(onGuildRoleDelete) // GUILD_ROLE_DELETE

	discord.AddHandler(onMessageCreate) // MESSAGE_CREATE
	discord.AddHandler(onMessageUpdate) // MESSAGE_UPDATE
	discord.AddHandler(onMessageDelete) // MESSAGE_DELETE

	discord.AddHandler(onPresenceUpdate) // PRESENCE_UPDATE

	discord.AddHandler(onTypingStart) // TYPING_START

	discord.AddHandler(onUserSettingsUpdate) // USER_SETTINGS_UPDATE
	discord.AddHandler(onUserUpdate)         // USER_UPDATE

	discord.AddHandler(onVoiceStateUpdate)  // VOICE_STATE_UPDATE
	discord.AddHandler(onVoiceServerUpdate) // VOICE_SERVER_UPDATE

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
