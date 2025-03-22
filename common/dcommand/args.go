package dcommand

type Args struct {
	Name	string
	Type	ArgumentType
}

type ArgumentType interface {
	Help() string
}

var (
	String = &StringArg{}
)

type StringArg struct{}
var _ ArgumentType = (*StringArg)(nil)
func (s *StringArg) Help() string {
	return "Text"
}