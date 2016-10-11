package slack

type messageWrapper struct {
	message *Message
	status  Status
}

/*
	Compatibility with direct bot access (bot event reception)
	Actual event reception should be moved to subBots - the main "bot" should only handle delegation
	and subbot management.

	Easily addressed by passing through into the active bot.
	These are a sort of "default" scenario.
*/

// Internal direct handler reference.
func (bot *Bot) onEvent(event string, handler BotAction) {
	bot.activeBot.OnEvent(event, handler)
}

// OnEventWithSubtype registers handler to fire on the given type and subtype of event.
func (bot *Bot) OnEventWithSubtype(event, subtype string, handler BotAction) {
	bot.activeBot.OnEventWithSubtype(event, subtype, handler)
}

// Delegate to correct subBot
func (bot *Bot) handle(event map[string]interface{}) (wrappers []messageWrapper) {
	eventTarget, hasTarget := event["target"].(string)
	target := bot.activeBot // Obtain targeted bot
	if hasTarget {
		target = bot.GetTarget(eventTarget)
	}
	wrappers = target.handle(event)
	return
}
