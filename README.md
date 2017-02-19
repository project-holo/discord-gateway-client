# ProjectHOLO Discord Client
Lightweight Golang Discord client that receives incoming events from the Discord
gateway, serializes them into JSON, and passes them over a STOMP broker to the
worker nodes that process the incoming events.

### Usage
1. Build it
2. `discord-gateway-client -t "Bot TOKEN_HERE" -b "stomp://login:passcode@host:port/host"`
3. ~~If you'd prefer to use environment variables: `DISCORD_TOKEN`, `STOMP_URI`, `DEBUG`~~ *soon*
4. If you want debug lines in console, `-d true`

### Serializes (in order of Discord documentation)
- [x] READY
- [ ] RESUMED
- [ ] CHANNEL_CREATE
- [ ] CHANNEL_UPDATE
- [ ] CHANNEL_DELETE
- [ ] GUILD_CREATE
- [ ] GUILD_UPDATE
- [ ] GUILD_DELETE
- [ ] GUILD\_BAN_ADD
- [ ] GUILD\_BAN_REMOVE
- [ ] GUILD\_EMOJIS_UPDATE
- [ ] GUILD\_INTEGRATIONS_UPDATE
- [ ] GUILD\_MEMBER_ADD
- [ ] GUILD\_MEMBER_REMOVE
- [ ] GUILD\_MEMBER_UPDATE
- [ ] GUILD\_MEMBERS_CHUNK
- [ ] GUILD\_ROLE_CREATE
- [ ] GUILD\_ROLE_UPDATE
- [ ] GUILD\_ROLE_DELETE
- [x] MESSAGE_CREATE
- [ ] MESSAGE_UPDATE
- [ ] MESSAGE_DELETE
- [ ] MESSAGE\_DELETE_BULK
- [ ] PRESENCE_UPDATE
- [ ] TYPING_START
- [ ] USER\_SETTINGS_UPDATE
- [ ] USER_UPDATE
- [ ] VOICE\_STATE_UPDATE
- [ ] VOICE\_SERVER_UPDATE

Note: some of these events won't ever be serialized, as it's useless for this
application.

### License
A copy of the MIT license can be found in `LICENSE`.
