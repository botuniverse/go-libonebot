package libonebot

type commCloser func()

func (c commCloser) Close() {
	if c != nil {
		c()
	}
}
