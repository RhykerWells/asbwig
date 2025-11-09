package dcommand

// Args defines the structure to pass argument data with
type Args struct {
	Name     string
	Type     ArgumentType
	Optional bool
}

type ArgumentType interface {
	Help() string
}

var (
	String      = &StringArg{}
	Int         = &IntArg{}
	User        = &UserArg{}
	Member      = &MemberArg{}
	Bet         = &BetArg{}
	Duration    = &DurationArg{}
	CoinSide    = &CoinSideArg{}
	UserBalance = &BalanceArg{}
)

type StringArg struct{}

var _ ArgumentType = (*StringArg)(nil)

func (s *StringArg) Help() string {
	return "Text"
}

type IntArg struct{}

var _ ArgumentType = (*IntArg)(nil)

func (i *IntArg) Help() string {
	return "Whole number above 0"
}

type UserArg struct{}

var _ ArgumentType = (*UserArg)(nil)

func (u *UserArg) Help() string {
	return "Mention/ID"
}

type MemberArg struct{}

var _ ArgumentType = (*MemberArg)(nil)

func (m *MemberArg) Help() string {
	return "Mention/ID"
}

type BetArg struct{}

var _ ArgumentType = (*BetArg)(nil)

func (b *BetArg) Help() string {
	return "Whole integer|max|all"
}

type DurationArg struct{}

var _ ArgumentType = (*DurationArg)(nil)

func (d *DurationArg) Help() string {
	return "Duration"
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
