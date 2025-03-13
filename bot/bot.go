package bot

import "github.com/Ranger-4297/asbwig/internal"

func Run() {
	internal.Session.AddHandler(MessageCreate)
	internal.Session.AddHandler(GuildJoin)
	internal.Session.AddHandler(GuildLeave)
}