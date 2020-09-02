package main

func JoinCommand(ctx *Context) {
	ctx.Parties.Join(ctx.Guild.ID, ctx.Channel.ID, ctx.User)

	state := ctx.Guild.ID + "-" + ctx.Channel.ID + "-" + ctx.User.ID
	url := ctx.Auth.AuthURL(state)
	dm, _ := ctx.Session.UserChannelCreate(ctx.User.ID)
	ctx.Session.ChannelMessageSend(dm.ID, url)
}

func ListCommand(ctx *Context) {
	party := ctx.Parties.GetByGuild(ctx.Guild.ID)
	reply := "Party's list"
	for _, u := range party.users {
		if u.spotify == nil {
			reply += "\n - " + u.discord.Username + " - Not authenticated, please join the party."
			continue
		}
		s, err := u.spotify.CurrentUser()
		if err != nil {
			reply += "\n - " + u.discord.Username + " - Not authenticated, please join the party."
		} else {
			reply += "\n - " + u.discord.Username + " / " + s.DisplayName
		}
	}
	ctx.Reply(reply)
}

