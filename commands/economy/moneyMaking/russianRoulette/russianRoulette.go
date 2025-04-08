package russianroulette

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/economy/models"
	"github.com/RhykerWells/asbwig/common"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"slices"
)

var (
	activeGames = make(map[string]*RouletteGame)
)

type RouletteGame struct {
	Bet       int64
	PlayerIDs []string
	OwnerID   string
	IsActive  bool
}

var Command = &dcommand.AsbwigCommand{
	Command:     "russianroulette",
	Category: 	 dcommand.CategoryEconomy,
	Aliases:     []string{"rr"},
	Description: "Russian roulette with up to 6 people\nAll players must join with the same bet\nPayout for winners is `(<Bet>*Players)/winners`",
	Args: []*dcommand.Args{
		{Name: "Bet", Type: dcommand.Int},
	},
	Run: func(data *dcommand.Data) {
		embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: data.Author.Username, IconURL: data.Author.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
		guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", data.GuildID)).One(context.Background(), common.PQ)
		economyUser, err := models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = economyUser.Cash
		}
		if len(data.Args) <= 0 {
			embed.Description = "No `Bet` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		option := data.Args[0]
		if functions.ToInt64(option) <= 0 && option != "all" && option != "max" && option != "start" {
			embed.Description = "Invalid `Bet` argument provided"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if option == "start" {
			game, exists := activeGames[data.GuildID]
			if !exists {
				embed.Description = "There's no game to start."
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			if game.OwnerID != data.Author.ID {
				embed.Description = "Only the host can start the game."
				functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
				return
			}
			startGame(data.GuildID, data.ChannelID, data.Author.ID)
			return
		}
		bet := functions.ToInt64(option)
		if option == "all" {
			bet = cash
		} else if option == "max" {
			if guild.Maxbet > 0 {
				bet = guild.Maxbet
			} else {
				bet = cash
			}
		}
		if bet > cash {
			embed.Description = fmt.Sprintf("You can't bet more than you have in your hand. You currently have %s%s", guild.Symbol, humanize.Comma(cash))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if guild.Maxbet > 0 && bet > guild.Maxbet {
			embed.Description = fmt.Sprintf("You can't bet more than the servers limit. The limit is %s%s", guild.Symbol, humanize.Comma(guild.Maxbet))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		game, exists := activeGames[data.GuildID]
		if !exists {
			cooldown, err := models.EconomyCooldowns(qm.Where("guild_id=? AND user_id=? AND type='russianroulette'", data.GuildID, data.Author.ID)).One(context.Background(), common.PQ)
			if err == nil {
				if cooldown.ExpiresAt.Time.After(time.Now()) {
					embed.Description = fmt.Sprintf("You are on cooldown. You can start another game <t:%d:R>", (time.Now().Unix() + int64(time.Until(cooldown.ExpiresAt.Time).Seconds())))
					functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
					return
				}
			}
			activeGames[data.GuildID] = &RouletteGame{
				OwnerID:   data.Author.ID,
				Bet:       bet,
				PlayerIDs: []string{data.Author.ID},
				IsActive:  false,
			}
			embed.Description = fmt.Sprintf("A new game of russian roulette has started\n\nTo join, use the command `russianroulette %d` (1/6)\nTo start this game use the command `russianroulette start` or wait for more players\nThis game will automatically start in <t:%d:R> minutes if enough players join.", bet, (time.Now().Unix() + int64((2 * time.Minute).Seconds())))
			embed.Color = common.SuccessGreen
			cash = cash - bet
			userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
			userEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash"), boil.Infer())
			cooldowns := models.EconomyCooldown{GuildID: data.GuildID, UserID: data.Author.ID, Type: "russianroulette", ExpiresAt: null.Time{Time: time.Now().Add(300 * time.Second), Valid: true}}
			cooldowns.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id", "type"}, boil.Whitelist("expires_at"), boil.Infer())
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			time.AfterFunc(2*time.Minute, func() {startGame(data.GuildID, data.ChannelID, data.Author.ID)})
			return
		}
		if bet != game.Bet {
			embed.Description = fmt.Sprintf("All players must bet the same amount. This game's bet is %s%s", guild.Symbol, humanize.Comma(game.Bet))
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if game.IsActive {
			embed.Description = "You can't join or start a game while there is one running"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		if slices.Contains(game.PlayerIDs, data.Author.ID) {
			embed.Description = "You've already joined this game.\nIf you're the host, start the game with `russianroulette start`"
			functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
			return
		}
		game.PlayerIDs = append(game.PlayerIDs, data.Author.ID)
		if len(game.PlayerIDs) == 6 {
			game.IsActive = true
			startGame(data.GuildID, data.ChannelID, data.Author.ID)
			return
		}
		embed.Description = fmt.Sprintf("You've joined this game of russian roulette with a bet of %s%s (%d/6)", guild.Symbol, humanize.Comma(bet), len(game.PlayerIDs))
		embed.Color = common.SuccessGreen
		cash = cash - bet
		userEntry := models.EconomyUser{GuildID: data.GuildID, UserID: data.Author.ID, Cash: cash}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash"), boil.Infer())
		functions.SendMessage(data.ChannelID, &discordgo.MessageSend{Embed: embed})
	},
}

func startGame(guildID, channelID, ownerID string) {
	game, exists := activeGames[guildID]
	if !exists || game.OwnerID != ownerID {
		return
	}
	user, _ := functions.GetUser(game.OwnerID)
	embed := &discordgo.MessageEmbed{Author: &discordgo.MessageEmbedAuthor{Name: user.Username, IconURL: user.AvatarURL("256")}, Timestamp: time.Now().Format(time.RFC3339), Color: common.ErrorRed}
	guild, _ := models.EconomyConfigs(qm.Where("guild_id=?", guildID)).One(context.Background(), common.PQ)
	removeLeftPlayers(guildID, game)
	if len(game.PlayerIDs) <= 1 {
		functions.SendBasicMessage(channelID, "There weren't enough players to start russian roulette. Please start a new one.")
		delete(activeGames, guildID)
		return
	}
	functions.SendBasicMessage(channelID, "The russian roulette game has begun")
	loserCount := rand.Intn(len(game.PlayerIDs))
	var loser *discordgo.Member
	for playerNumber, player := range game.PlayerIDs {
		member, _ := functions.GetMember(guildID, player)
		if playerNumber != loserCount {
			functions.SendBasicMessage(channelID, fmt.Sprintf("> **%s** pulled the trigger and survived", member.Mention()))
			continue
		}
		loser = member
		functions.SendBasicMessage(channelID, fmt.Sprintf("**%s** pulled the trigger and died", member.Mention()))
		break
	}
	fields := []*discordgo.MessageEmbedField{}
	payout := (game.Bet * int64(len(game.PlayerIDs))) / int64(len(game.PlayerIDs)-1)
	for _, player := range game.PlayerIDs {
		if loser.User.ID == player {
			continue
		}
		economyUser, err := models.EconomyUsers(qm.Where("guild_id=? AND user_id=?", guildID, player)).One(context.Background(), common.PQ)
		var cash int64 = 0
		if err == nil {
			cash = economyUser.Cash
		}
		member, _ := functions.GetMember(guildID, player)
		field := &discordgo.MessageEmbedField{Name: member.User.Username, Value: member.Mention(), Inline: false}
		fields = append(fields, field)
		userEntry := models.EconomyUser{GuildID: guildID, UserID: player, Cash: cash + payout}
		userEntry.Upsert(context.Background(), common.PQ, true, []string{"guild_id", "user_id"}, boil.Whitelist("cash"), boil.Infer())
	}
	embed.Title = "Russian roulette winners"
	embed.Description = fmt.Sprintf("Payout is %s%s per winner", guild.Symbol, humanize.Comma(payout))
	embed.Fields = fields
	embed.Color = common.SuccessGreen
	delete(activeGames, guildID)
	functions.SendMessage(channelID, &discordgo.MessageSend{Embed: embed})
}

func removeLeftPlayers(guildID string, game *RouletteGame) {
	filtered := make([]string, 0, len(game.PlayerIDs))
	for _, playerID := range game.PlayerIDs {
		if _, err := functions.GetMember(guildID, playerID); err == nil {
			filtered = append(filtered, playerID)
		}
	}
	game.PlayerIDs = filtered
}