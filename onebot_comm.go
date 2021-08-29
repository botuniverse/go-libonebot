package libonebot

func (ob *OneBot) startCommMethods() {
	ob.commClosersLock.Lock()
	defer ob.commClosersLock.Unlock()
	ob.startCommMethodsHTTP()
	ob.startCommMethodsHTTPWebhook()
	ob.startCommMethodsWS()
	ob.startCommMethodsWSReverse()
}

func (ob *OneBot) startCommMethodsHTTP() {
	if ob.Config.CommMethods.HTTP == nil {
		return
	}
	for _, c := range ob.Config.CommMethods.HTTP {
		ob.commClosers = append(ob.commClosers, commStartHTTP(c, ob))
	}
}

func (ob *OneBot) startCommMethodsHTTPWebhook() {
	if ob.Config.CommMethods.HTTPWebhook == nil {
		return
	}
	for _, c := range ob.Config.CommMethods.HTTPWebhook {
		ob.commClosers = append(ob.commClosers, commStartHTTPWebhook(c, ob))
	}
}

func (ob *OneBot) startCommMethodsWS() {
	if ob.Config.CommMethods.WS == nil {
		return
	}
	for _, c := range ob.Config.CommMethods.WS {
		ob.commClosers = append(ob.commClosers, commStartWS(c, ob))
	}
}

func (ob *OneBot) startCommMethodsWSReverse() {
	if ob.Config.CommMethods.WSReverse == nil {
		return
	}
	for _, c := range ob.Config.CommMethods.WSReverse {
		ob.commClosers = append(ob.commClosers, commStartWSReverse(c, ob))
	}
}
