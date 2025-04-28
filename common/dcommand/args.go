package dcommand

type Args struct {
	Name     string
	Type     ArgumentType
	Optional bool
}

type ArgumentType interface {
	Help() string
}

var (
	Any    = &AnyArg{}
	String = &StringArg{}
	Int    = &IntArg{}
	User   = &UserArg{}
	Channel = &ChannelArg{}
	Bet = &BetArg{}
	Money = &MoneyArg{}
	CoinSide = &CoinSideArg{}
	UserBalance = &BalanceArg{}
)

type AnyArg struct{}
var _ ArgumentType = (*AnyArg)(nil)
func (a *AnyArg) Help() string {
	return "Any"
}

type StringArg struct{}
var _ ArgumentType = (*StringArg)(nil)
func (s *StringArg) Help() string {
	return "Text"
}

type IntArg struct{}
var _ ArgumentType = (*IntArg)(nil)
func (i *IntArg) Help() string {
	return "Whole number"
}

type UserArg struct{}
var _ ArgumentType = (*UserArg)(nil)
func (u *UserArg) Help() string {
	return "Mention/ID"
}

type ChannelArg struct{}
var _ ArgumentType = (*ChannelArg)(nil)
func (c *ChannelArg) Help() string {
	return "Mention/ID"
}

type BetArg struct{}
var _ ArgumentType = (*BetArg)(nil)
func (b *BetArg) Help() string {
	return "Whole integer|max|all"
}

type MoneyArg struct{}
var _ ArgumentType = (*MoneyArg)(nil)
func (m *MoneyArg) Help() string {
	return "Whole integer above 0"
}

type CoinSideArg struct{}
var _ ArgumentType = (*CoinSideArg)(nil)
func (c *CoinSideArg) Help() string {
	return "Heads/Tails"
}

type BalanceArg struct{}
var _ ArgumentType = (*BalanceArg)(nil)
func (b *BalanceArg) Help() string {
	return "Bank/Cash"
}