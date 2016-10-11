package slack

import (
	"crypto/md5"
	"fmt"
)

/*
	Subordinate bots contained inside the main "bot"
*/

type SubBot struct {
	Owner       *Bot
	Name        string                                // Name for addressing SubBot
	BotID       string                                // Unique ID for access to SubBot and saving file
	Handlers    map[string]([]BotAction)              // SubBot-specific handlers
	Subhandlers map[string](map[string]([]BotAction)) // SubBot-specific subhandlers
}

// NewSubBot returns a SubBot under the specified Bot
func NewSubBot(owner *Bot, botID string, name string) *SubBot {
	return &SubBot{
		Owner:       owner,
		Name:        name,
		BotID:       botID,
		Handlers:    make(map[string]([]BotAction)),
		Subhandlers: make(map[string](map[string]([]BotAction))),
	}
}

func (bot *Bot) GetTarget(id string) *SubBot {
	subBot, hasSubBot := bot.SubBots[id]
	if !hasSubBot {
		subBot = bot.activeBot
	}
	return subBot
}

// Add a subbot.
func (bot *Bot) AddSubBot(name string) (index int, id string) {
	// generate unique ID
	id = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("BOTID_%v", bot.SubBotArr))))
	sBot := NewSubBot(bot, id, name)

	bot.SubBotArr[bot.SubBotNum] = sBot
	bot.SubBotName[sBot.Name] = sBot
	bot.SubBots[sBot.BotID] = sBot
	// If this is the first subBot bot (to be added), initialize it
	if bot.SubBotNum == 0 {
		bot.activeBot = sBot
	}
	bot.SubBotNum++
	index = bot.SubBotNum
	return
}
