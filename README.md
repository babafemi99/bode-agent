# Bọ̀dé Agent

Bọ̀dé Agent is the conversational layer for operating Bọ̀dé through messaging platforms.

The first integration is **Telegram**, with **Gemma 4** planned as the model powering natural-language interactions. The agent is a separate Go service and communicates with Bọ̀dé Core through its HTTP API.

## What We Have

The initial Telegram transport and bot lifecycle are being built around:

```text
Telegram
   ↓
Telegram Bot
   ↓
Receive
   ↓
Inbound Channel
   ↓
HandleMessage
   ↓
Command / Natural Language
```

Multiple Telegram bots can be connected to Bọ̀dé. A `Manager` keeps track of active bots and handles their lifecycle.

Each bot has:

* Bọ̀dé Bot ID
* Telegram Bot ID
* Telegram client
* inbound message channel
* outbound response channel

Bot identifiers are generated using the project's ULID-based `lid` package.

## Message Flow

Telegram updates are received by the bot and converted from Telegram-specific types into our own `Message` type.

```go
type Message struct {
    BotID     string
    UserID    int64
    ChatID    int64
    Text      string
    IsCommand bool
    Command   string
    Arguments string
}
```

This keeps the rest of the agent independent from the Telegram library.

Incoming messages flow through the inbound channel and are handled by:

```go
func (b *Bot) HandleMessage(msg Message)
```

Commands are handled directly by Go:

```text
/start
/help
/context
```

Normal text is passed to the conversational agent:

```text
"Show me the guests for the wedding"
        ↓
     Gemma 4
```

## Agent Direction

Gemma 4 will handle natural-language understanding and decide when a Bọ̀dé capability is required.

The model will **not** access Bọ̀dé's database or repositories directly.

Instead:

```text
User
 ↓
Telegram
 ↓
Bọ̀dé Agent
 ↓
Gemma 4
 ↓
Tool
 ↓
Bọ̀dé HTTP API
 ↓
Tool Result
 ↓
Gemma 4
 ↓
Telegram
```

Tools will provide capabilities such as:

* finding guests
* viewing events
* listing tables
* creating guests
* managing invitations
* retrieving event information

Bọ̀dé Core remains responsible for authorization, business rules and the final state of the system.

## Transport

The initial development path uses **Telegram long polling**.

The internal message boundary is designed so that Telegram webhooks can be introduced later without changing the agent logic.

```text
Telegram Transport
        ↓
     Message
        ↓
   Agent Logic
```

## Current Development Path

The project is being built incrementally:

1. Telegram bot connection
2. Multiple bot management
3. Receive / inbound message flow
4. Message handling and commands
5. Outbound message flow
6. Gemma 4 integration
7. Bọ̀dé API tools
8. Sessions and conversation context
9. Multi-step workflows and confirmations
10. Additional messaging platforms

The goal is simple: **make Bọ̀dé operable through conversation without coupling the AI layer to Bọ̀dé's internal implementation.**
