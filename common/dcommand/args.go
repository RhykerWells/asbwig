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