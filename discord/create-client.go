package discord

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
)

// CreateDiscordClient creates a Discord client from a token and returns it.
func CreateDiscordClient(token string, shardID, shardCount int) (*discordgo.Session, *discordgo.User, error) {
	// Create session
	discord, err := discordgo.New(token)
	if err != nil {
		return nil, nil, err
	}

	// Fetch user information for the user, used to check for valid token and to
	// print user information to debug
	me, err := discord.User("@me")
	if err != nil {
		return nil, nil, err
	}
	log.Debug("fetched /users/@me from Discord API")
	log.Debugf("@me=%v", me.ID)
	log.Debugf("persona=%v#%v", me.Username, me.Discriminator)

	// Sharding parameters
	discord.ShardID = shardID
	discord.ShardCount = shardCount
	if discord.ShardCount <= 0 {
		discord.ShardCount = 1
	}
	log.Debugf("shard=%v/%v", discord.ShardID, discord.ShardCount)

	// Ready logger
	discord.AddHandlerOnce(func(s *discordgo.Session, e *discordgo.Ready) {
		log.Debug("READY")
	})

	return discord, me, nil
}
