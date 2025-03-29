package dcommand

type Args struct {
	Name	string
	Type	ArgumentType
}

type ArgumentType interface {
	Help() string
}

var (
	String =	&StringArg{}
	Int =		&IntArg{}
	User =		&UserArg{}
)

type StringArg struct{}
var _ ArgumentType = (*StringArg)(nil)
func (s *StringArg) Help() string {
	return "Text"
}

type IntArg struct{}
var _ ArgumentType = (*IntArg)(nil)
func (s *IntArg) Help() string {
	return "Integer"
}

type UserArg struct{}
var _ ArgumentType = (*UserArg)(nil)
func (s *UserArg) Help() string {
	return "ID/Mention"
}