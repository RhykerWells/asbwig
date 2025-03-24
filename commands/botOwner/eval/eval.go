package eval

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RhykerWells/asbwig/bot/functions"
	"github.com/RhykerWells/asbwig/commands/util"
	"github.com/RhykerWells/asbwig/common/dcommand"
	"github.com/bwmarrin/discordgo"
	piston "github.com/milindmadhukar/go-piston"
)

var Command = &dcommand.AsbwigCommand{
	Command:      []string{"eval"},
	Description:  "Evaluates Go code",
	ArgsRequired: 1,
	Args: []*dcommand.Args{
		{Name: "Code", Type: dcommand.String},
	},
	Run: util.OwnerCommand(eval),
}

func eval(data *dcommand.Data) {
	codeBlock := codeBlockExtractor(data)
	if codeBlock == "" {
		return
	}
	output, err := exec(codeBlock, data.Message.Reference())
	if err != nil {
		functions.SendBasicMessage(data.Message.ChannelID, "Something went wrong.")
		return
	}
	functions.SendBasicMessage(data.Message.ChannelID, "```"+output.Run.Output+"```")
}

func exec(code string, messageReference *discordgo.MessageReference) (*piston.PistonExecution, error) {
	// execute code using piston library
	pclient := piston.CreateDefaultClient()
	output, err := pclient.Execute("go", "",
		[]piston.Code{
			{
				Name:    fmt.Sprintf("%s-code", messageReference.MessageID),
				Content: code,
			},
		},
	)
	return output, err
}

func codeBlockExtractor(data *dcommand.Data) string {
	messageContent := strings.Join(data.Args, " ")
	codeblockRegex, _ := regexp.Compile("```go.*")
	c := strings.Split(messageContent, "\n")
	for bi, bb := range c {
		if codeblockRegex.MatchString(bb) {
			var codeBlock string
			endBlockRegx, _ := regexp.Compile("```")
			sa := c[bi+1:]
			for ei, eb := range sa {
				if endBlockRegx.Match([]byte(eb)) {
					codeBlock = strings.Join(sa[:ei], "\n")
					return codeBlock
				}
			}
		}
	}
	return ""
}
