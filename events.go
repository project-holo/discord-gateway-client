package main

import "github.com/bwmarrin/discordgo"

// Discord gateway event types (in order of list in README.md).
const (
	ready   = "READY"
	resumed = "RESUMED"

	channelCreate = "CHANNEL_CREATE"
	channelUpdate = "CHANNEL_UPDATE"
	channelDelete = "CHANNEL_DELETE"

	guildCreate = "GUILD_CREATE"
	guildUpdate = "GUILD_UPDATE"
	guildDelete = "GUILD_DELETE"

	guildBanAdd    = "GUILD_BAN_ADD"
	guildBanRemove = "GUILD_BAN_REMOVE"

	guildEmojisUpdate = "GUILD_EMOJIS_UPDATE"

	guildIntegrationsUpdate = "GUILD_INTEGRATIONS_UPDATE"

	guildMemberAdd    = "GUILD_MEMBER_ADD"
	guildMemberRemove = "GUILD_MEMBER_REMOVE"
	guildMemberUpdate = "GUILD_MEMBER_UPDATE"
	guildMembersChunk = "GUILD_MEMBERS_CHUNK"

	guildRoleCreate = "GUILD_ROLE_CREATE"
	guildRoleUpdate = "GUILD_ROLE_UPDATE"
	guildRoleDelete = "GUILD_ROLE_DELETE"

	messageCreate = "MESSAGE_CREATE"
	messageUpdate = "MESSAGE_UPDATE"
	messageDelete = "MESSAGE_DELETE"
	// messageDeleteBulk = "MESSAGE_DELETE_BULK"

	presenceUpdate = "PRESENCE_UPDATE"

	typingStart = "TYPING_START"

	userSettingsUpdate = "USER_SETTINGS_UPDATE"
	userUpdate         = "USER_UPDATE"

	voiceStateUpdate  = "VOICE_STATE_UPDATE"
	voiceServerUpdate = "VOICE_SERVER_UPDATE"
)

// Event handlers.
var eventHandlers = []interface{}{
	onReady,
	onResumed,

	onChannelCreate,
	onChannelUpdate,
	onChannelDelete,

	onGuildCreate,
	onGuildUpdate,
	onGuildDelete,

	onGuildBanAdd,
	onGuildBanRemove,

	onGuildEmojisUpdate,

	onGuildIntegrationsUpdate,

	onGuildMemberAdd,
	onGuildMemberRemove,
	onGuildMemberUpdate,
	onGuildMembersChunk,

	onGuildRoleCreate,
	onGuildRoleUpdate,
	onGuildRoleDelete,

	onMessageCreate,
	onMessageUpdate,
	onMessageDelete,

	onPresenceUpdate,

	onTypingStart,

	onUserSettingsUpdate,
	onUserUpdate,

	onVoiceStateUpdate,
	onVoiceServerUpdate,
}

func onReady(s *discordgo.Session, e *discordgo.Ready) {
	//s.UpdateStatus(0, "ProjectHOLO")
	serializeAndDispatchEvent(ready, *e)
}

func onResumed(s *discordgo.Session, e *discordgo.Resumed) {
	serializeAndDispatchEvent(resumed, *e)
}

func onChannelCreate(s *discordgo.Session, e *discordgo.ChannelCreate) {
	serializeAndDispatchEvent(channelCreate, *e.Channel)
}

func onChannelUpdate(s *discordgo.Session, e *discordgo.ChannelUpdate) {
	serializeAndDispatchEvent(channelUpdate, *e.Channel)
}

func onChannelDelete(s *discordgo.Session, e *discordgo.ChannelDelete) {
	serializeAndDispatchEvent(channelDelete, *e.Channel)
}

func onGuildCreate(s *discordgo.Session, e *discordgo.GuildCreate) {
	serializeAndDispatchEvent(channelCreate, *e.Guild)
}

func onGuildUpdate(s *discordgo.Session, e *discordgo.GuildUpdate) {
	serializeAndDispatchEvent(guildUpdate, *e.Guild)
}

func onGuildDelete(s *discordgo.Session, e *discordgo.GuildDelete) {
	serializeAndDispatchEvent(guildDelete, *e.Guild)
}

func onGuildBanAdd(s *discordgo.Session, e *discordgo.GuildBanAdd) {
	serializeAndDispatchEvent(guildBanAdd, *e)
}

func onGuildBanRemove(s *discordgo.Session, e *discordgo.GuildBanRemove) {
	serializeAndDispatchEvent(guildBanRemove, *e)
}

func onGuildEmojisUpdate(s *discordgo.Session, e *discordgo.GuildEmojisUpdate) {
	serializeAndDispatchEvent(guildEmojisUpdate, *e)
}

func onGuildIntegrationsUpdate(s *discordgo.Session, e *discordgo.GuildIntegrationsUpdate) {
	serializeAndDispatchEvent(guildIntegrationsUpdate, *e)
}

func onGuildMemberAdd(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	serializeAndDispatchEvent(guildMemberAdd, *e.Member)
}

func onGuildMemberRemove(s *discordgo.Session, e *discordgo.GuildMemberRemove) {
	serializeAndDispatchEvent(guildMemberRemove, *e.Member)
}

func onGuildMemberUpdate(s *discordgo.Session, e *discordgo.GuildMemberUpdate) {
	serializeAndDispatchEvent(guildMemberUpdate, *e.Member)
}

func onGuildMembersChunk(s *discordgo.Session, e *discordgo.GuildMembersChunk) {
	serializeAndDispatchEvent(guildMembersChunk, *e)
}

func onGuildRoleCreate(s *discordgo.Session, e *discordgo.GuildRoleCreate) {
	serializeAndDispatchEvent(guildRoleCreate, *e.GuildRole)
}

func onGuildRoleUpdate(s *discordgo.Session, e *discordgo.GuildRoleUpdate) {
	serializeAndDispatchEvent(guildRoleUpdate, *e.GuildRole)
}

func onGuildRoleDelete(s *discordgo.Session, e *discordgo.GuildRoleDelete) {
	serializeAndDispatchEvent(guildRoleDelete, *e)
}

func onMessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {
	serializeAndDispatchEvent(messageCreate, *e.Message)
}

func onMessageUpdate(s *discordgo.Session, e *discordgo.MessageUpdate) {
	serializeAndDispatchEvent(messageUpdate, *e.Message)
}

func onMessageDelete(s *discordgo.Session, e *discordgo.MessageDelete) {
	serializeAndDispatchEvent(messageDelete, *e.Message)
}

func onPresenceUpdate(s *discordgo.Session, e *discordgo.PresenceUpdate) {
	serializeAndDispatchEvent(presenceUpdate, *e)
}

func onTypingStart(s *discordgo.Session, e *discordgo.TypingStart) {
	serializeAndDispatchEvent(typingStart, *e)
}

func onUserSettingsUpdate(s *discordgo.Session, e *discordgo.UserSettingsUpdate) {
	serializeAndDispatchEvent(userSettingsUpdate, *e)
}

func onUserUpdate(s *discordgo.Session, e *discordgo.UserUpdate) {
	serializeAndDispatchEvent(userUpdate, *e.User)
}

func onVoiceStateUpdate(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	serializeAndDispatchEvent(presenceUpdate, *e.VoiceState)
}

func onVoiceServerUpdate(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
	serializeAndDispatchEvent(voiceServerUpdate, *e)
}
