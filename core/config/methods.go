package config

func (c *DaemonConfig) IsSudo(id int64) bool {
	if len(c.SudoUsers) == 0 {
		return false
	}

	for _, sudo := range c.SudoUsers {
		if sudo == id {
			return true
		}
	}

	return false
}
