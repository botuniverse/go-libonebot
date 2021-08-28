package onebot

func (ob *OneBot) startCommMethods() {
	ob.commClosers = append(ob.commClosers, commStartHTTP("127.0.0.1", 5700, ob))
	ob.commClosers = append(ob.commClosers, commStartWS("127.0.0.1", 6700, ob))
	ob.commClosers = append(ob.commClosers, commStartHTTPWebhook("http://127.0.0.1:8080", ob))
}
