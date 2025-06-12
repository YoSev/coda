package coda

func (c *Coda) debug(msg string) {
	c.Logs = append(c.Logs, msg)
}
