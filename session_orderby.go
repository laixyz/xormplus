package xormplus

// GetOrderBy get order by string.
func (session *Session) GetOrderBy() string {
	return session.statement.OrderStr
}
