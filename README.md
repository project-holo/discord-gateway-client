# ProjectHOLO Discord Gateway Client

Lightweight Golang Discord client that receives incoming events from the
Discord gateway and sends them to a STOMP broker. Useful for creating
distributed and highly-available bots for Discord.

Rolling restarts can be performed on worker code without affecting shard
availability, and events can be balanced between all workers.

Events are serialized into the JSON format outlined below and passed
over a STOMP connection to a message broker (such as ActiveMQ or
RabbitMQ) to be processed by your own workers.

If you are looking for a compatible worker framework, this organization
maintains
[https://github.com/project-holo/discord-worker-framework](project-holo/discord-worker-framework),
which consumes the format that this project produces.

### Usage

1. Build it: `go build -o discord-gateway-client`
2. `discord-gateway-client --token "Bot TOKEN_HERE" --broker
   "stomp://login:passcode@host:port/host"`
3. If you'd prefer to use environment variables or would like shard
   configuration, check `init()` in
   [discord-gateway-client.go](discord-gateway-client.go)
4. If you want debug lines in console, set `-d true` or `DEBUG=true`
5. Consume events from `/events` on the virtual host of the broker
   (unless you changed the destination)

### Serialization Format (JSON)

```js
{
  "type": "string" // event type (e.g. MESSAGE_CREATE)
  "shard_id": int, // shard ID
  "data": { ... } // event data, whatever is serialized from DiscordGo
}
```

### Serializes the Following Events

- [x] `READY`
- [x] `RESUMED`
- [x] `CHANNEL_CREATE`
- [x] `CHANNEL_UPDATE`
- [x] `CHANNEL_DELETE`
- [x] `GUILD_CREATE`
- [x] `GUILD_UPDATE`
- [x] `GUILD_DELETE`
- [x] `GUILD_BAN_ADD`
- [x] `GUILD_BAN_REMOVE`
- [x] `GUILD_EMOJIS_UPDATE`
- [x] `GUILD_INTEGRATIONS_UPDATE`
- [x] `GUILD_MEMBER_ADD`
- [x] `GUILD_MEMBER_REMOVE`
- [x] `GUILD_MEMBER_UPDATE`
- [x] `GUILD_MEMBERS_CHUNK`
- [x] `GUILD_ROLE_CREATE`
- [x] `GUILD_ROLE_UPDATE`
- [x] `GUILD_ROLE_DELETE`
- [x] `MESSAGE_CREATE`
- [x] `MESSAGE_UPDATE`
- [x] `MESSAGE_DELETE`
- [ ] ~~`MESSAGE_DELETE_BULK`~~ (Discord doesn't send this anymore)
- [x] `PRESENCE_UPDATE`
- [x] `TYPING_START`
- [x] `USER_SETTINGS_UPDATE`
- [x] `USER_UPDATE`
- [x] `VOICE_STATE_UPDATE`
- [x] `VOICE_SERVER_UPDATE`

### License

A copy of the MIT license can be found in [LICENSE](LICENSE).
