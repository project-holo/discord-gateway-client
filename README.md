# ProjectHOLO Discord Client
Lightweight Golang Discord client that receives incoming events from the Discord
gateway, serializes them into JSON, and passes them over a STOMP broker to the
worker nodes that process the incoming events.

### Usage
1. Build it: `go build -o discord-gateway-client`
2. `discord-gateway-client -t "Bot TOKEN_HERE" -b "stomp://login:passcode@host:port/host"`
3. ~~If you'd prefer to use environment variables: `DISCORD_TOKEN`, `STOMP_URI`, `DEBUG`~~ *soon*
4. If you want debug lines in console, `-d true`

### Serializes (in order of Discord documentation)
- [x] READY
- [x] RESUMED
- [x] CHANNEL_CREATE
- [x] CHANNEL_UPDATE
- [x] CHANNEL_DELETE
- [x] GUILD_CREATE
- [x] GUILD_UPDATE
- [x] GUILD_DELETE
- [x] GUILD\_BAN_ADD
- [x] GUILD\_BAN_REMOVE
- [x] GUILD\_EMOJIS_UPDATE
- [x] GUILD\_INTEGRATIONS_UPDATE
- [x] GUILD\_MEMBER_ADD
- [x] GUILD\_MEMBER_REMOVE
- [x] GUILD\_MEMBER_UPDATE
- [x] GUILD\_MEMBERS_CHUNK
- [x] GUILD\_ROLE_CREATE
- [x] GUILD\_ROLE_UPDATE
- [x] GUILD\_ROLE_DELETE
- [x] MESSAGE_CREATE
- [x] MESSAGE_UPDATE
- [x] MESSAGE_DELETE
- [ ] ~~MESSAGE\_DELETE_BULK~~ *DiscordGo doesn't allow handling this event type*
- [x] PRESENCE_UPDATE
- [x] TYPING_START
- [x] USER\_SETTINGS_UPDATE
- [x] USER_UPDATE
- [x] VOICE\_STATE_UPDATE
- [x] VOICE\_SERVER_UPDATE

Note: some of these events won't ever be serialized, as it's useless for this
application.

### License
A copy of the MIT license can be found in `LICENSE`.
