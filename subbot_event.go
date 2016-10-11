package slack

// OnEvent registers handler to fire on the given type of event.
func (sBot *SubBot) OnEvent(event string, handler BotAction) {
	handlers, ok := sBot.Handlers[event]
	if !ok {
		handlers = make([]BotAction, 0)
	}
	handlers = append(handlers, handler)
	sBot.Handlers[event] = handlers
}

// OnEventWithSubtype registers handler to fire on the given type and subtype of event.
func (sBot *SubBot) OnEventWithSubtype(event, subtype string, handler BotAction) {
	subtypeMap, ok := sBot.Subhandlers[event]
	if !ok {
		subtypeMap = make(map[string]([]BotAction))
		sBot.Subhandlers[event] = subtypeMap
	}
	handlers, ok := sBot.Subhandlers[event][subtype]
	if !ok {
		handlers = make([]BotAction, 0)
	}
	handlers = append(handlers, handler)
	sBot.Subhandlers[event][subtype] = handlers
}

// Handle the event
func (sBot *SubBot) handle(event map[string]interface{}) (wrappers []messageWrapper) {
	eventType, hasType := event["type"].(string)
	eventSubtype, hasSubtype := event["subtype"].(string)

	if hasSubtype {
		subhandlerMap, ok := sBot.Subhandlers[eventType]
		if ok {
			subhandlers, ok := subhandlerMap[eventSubtype]
			if ok {
				for _, subhandler := range subhandlers {
					message, status := subhandler(sBot.Owner, event)
					wrappers = append(wrappers,
						messageWrapper{message, status})
				}
			}
		}
	}
	if hasType {
		handlers, ok := sBot.Handlers[eventType]
		if ok {
			for _, handler := range handlers {
				message, status := handler(sBot.Owner, event)
				wrappers = append(wrappers, messageWrapper{message, status})
			}
		}
	}
	return
}
