package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	templateFunctions = map[string]interface{}{
		"dict":  dict,
		"lower": lower,
		// Forms content
		"textInput":            textInput,
		"toggleSwitch":         toggleSwitch,
		"numberSelect":         numberSelection,
		"roleOptionsSingle":    roleOptionsSingle,
		"roleOptionsMulti":     roleOptionsMulti,
		"channelOptionsSingle": channelOptionsSingle,
	}
)

func dict(pairs ...interface{}) map[int]interface{} {
	result := make(map[int]interface{})
	for i := 0; i < len(pairs); i += 2 {
		key, _ := pairs[i].(int)
		result[key] = pairs[i+1]
	}
	return result
}

func lower(str string) string {
	return strings.ToLower(str)
}

// textInput generates a HTML object for the custom input.
// currentInput: the current input of the input
// uniqueID: string for the input ID (used to retrieve and store changed data)
func textInput(currentInput, uniqueID string) template.HTML {
	var menu strings.Builder
	menu.WriteString(`<div class="input-group mb-3">`)
	menu.WriteString(`<input type="text" class="textInput form-control text-light" name="` + uniqueID + `" id="` + uniqueID + `" autocomplete="off" value="` + currentInput + `">`)
	menu.WriteString("</div>")

	return template.HTML(menu.String())
}

// toggleSwitch generates a HTML object for the custom switch.
// currentState: the bool of the current state of the switch
// uniqueID: string for the input ID (used to retrieve and store changed data)
func toggleSwitch(currentState bool, uniqueID string) template.HTML {
	checked := ""
	if currentState {
		checked = " checked"
	}

	var menu strings.Builder
	menu.WriteString(`<label class="switch">`)
	menu.WriteString(`<input type="checkbox" name="` + uniqueID + `" id="` + uniqueID + `"` + checked + `/>`)
	menu.WriteString(`<span class="slider" style="left: 5px;"></span>`)
	menu.WriteString(`<span class="knob" style="left: 7px;"></span>`)
	menu.WriteString(`</label>`)

	return template.HTML(menu.String())
}

// numberSelection generates HTML options for a number picker
// min: lowest number possible
// max: highest number possible
// currentNumber: the current number
// uniqueID: string for the hidden input ID (used to retrieve and store changed data)
func numberSelection(min int64, max int64, currentNumber int64, uniqueID string) template.HTML {
	var menu strings.Builder
	menu.WriteString(`<div class="input-group mb-3">`)
	menu.WriteString(`<input type="number" name="` + uniqueID + `" id="` + uniqueID + `" min="` + strconv.FormatInt(min, 10) + `" step="1" max="` + strconv.FormatInt(max, 10) + `" class="form-control text-light" placeholder="0" style="background-color: var(--basePurple); border: 1px solid var(--accentGrey);" value="` + strconv.FormatInt(currentNumber, 10) + `">`)
	menu.WriteString(`<span class="input-group-text text-light" style="background-color: var(--primaryTetiaryPurple); border: 1px solid var(--accentGrey)">seconds</span>`)
	menu.WriteString(`</div>`)

	return template.HTML(menu.String())
}

// roleOptionsSingle generates HTML options for singular role selection
// roles: slice of Discord role objects
// selectedRoleID: string ID of currently selected role
// uniqueID: string for the hidden input ID (used to retrieve and store changed data)
// highestBotRolePosition: the position of the bots highest role
func roleOptionsSingle(roles []*discordgo.Role, selectedRoleID string, uniqueID string, highestBotRolePosition int) template.HTML {
	filteredRoles := make([]*discordgo.Role, 0, len(roles))
	for _, role := range roles {
		if role.Managed || role.Name == "@everyone" {
			continue
		}
		filteredRoles = append(filteredRoles, role)
	}
	sort.Slice(filteredRoles, func(i, j int) bool {
		return filteredRoles[i].Position > filteredRoles[j].Position
	})

	// Button label
	displayText := "Select role"
	if len(selectedRoleID) > 0 {
		label := ""
		for _, role := range filteredRoles {
			if selectedRoleID != role.ID {
				continue
			}
			label = role.Name
			break
		}
		if len(label) > 30 {
			displayText = "1 Selected"
		} else {
			displayText = label
		}
	}

	var menu strings.Builder
	menu.WriteString(`<div class="input-group mb-3">`)
	menu.WriteString(`
		<button class="btn dropdown-toggle text-start flex-grow-1 text-white" type="button" data-bs-toggle="dropdown" style="background-color: var(--basePurple); border: 1px solid var(--accentGrey); border-top-right-radius: var(--bs-btn-border-radius); border-bottom-right-radius: var(--bs-btn-border-radius);">
			<span id="` + uniqueID + `Label">` + template.HTMLEscapeString(displayText) + `</span>
		</button>
		<ul class="dropdown-menu w-100 overflow-auto" style="max-height: 250px;" aria-labelledby="` + uniqueID + `Dropdown">
		<li><a class="dropdown-item dropDownRoleSingleItem" data-value="">None</a></li>
	`)

	for _, role := range filteredRoles {
		disabled := ""
		disabledMsg := ""
		if highestBotRolePosition <= role.Position {
			disabled = " disabled"
			disabledMsg = " (bot higher than role)"
		}

		menu.WriteString(`<li>`)
		menu.WriteString(`<a class="dropdown-item dropDownRoleSingleItem` + disabled + `" data-value="` + role.ID + `">`)
		menu.WriteString(template.HTMLEscapeString(role.Name) + disabledMsg)
		menu.WriteString(`</a></li>`)
	}

	menu.WriteString(`</ul>`)
	menu.WriteString(`<input type="hidden" id="` + uniqueID + `" name="` + uniqueID + `" value="` + template.HTMLEscapeString(selectedRoleID) + `">`)
	menu.WriteString(`</div>`)
	return template.HTML(menu.String())
}

// roleOptionsMulti generates HTML options for multiple role selection
// roles: slice of Discord role objects
// selectedRoleIDs: slice of string IDs of currently selected roles
// uniqueID: string for the hidden input ID (used to retrieve and store changed data)
// highestBotRolePosition: the position of the bots highest role
func roleOptionsMulti(roles []*discordgo.Role, selectedRoleIDs interface{}, uniqueID string, highestBotRolePosition int) template.HTML {
	selectedMap := make(map[string]bool)
	if selectedRoleIDs != nil {
		if roleIDs, ok := selectedRoleIDs.([]string); ok {
			for _, id := range roleIDs {
				selectedMap[id] = true
			}
		}
	}

	filteredRoles := make([]*discordgo.Role, 0, len(roles))
	for _, role := range roles {
		if role.Managed || role.Name == "@everyone" {
			continue
		}
		filteredRoles = append(filteredRoles, role)
	}
	sort.Slice(filteredRoles, func(i, j int) bool {
		return filteredRoles[i].Position > filteredRoles[j].Position
	})

	var selectedNames []string
	for _, role := range filteredRoles {
		if selectedMap[role.ID] {
			selectedNames = append(selectedNames, role.Name)
		}
	}

	// Button label
	displayText := "Select roles"
	if len(selectedNames) > 0 {
		label := strings.Join(selectedNames, ", ")
		if len(selectedNames) > 3 || len(label) > 30 {
			displayText = fmt.Sprintf("%d Selected", len(selectedNames))
		} else {
			displayText = label
		}
	}

	var menu strings.Builder
	menu.WriteString(`<div class="input-group mb-3">`)
	menu.WriteString(`
		<button class="btn dropdown-toggle text-start flex-grow-1 text-white" type="button" data-bs-toggle="dropdown" data-bs-auto-close="outside" style="background-color: var(--basePurple); border: 1px solid var(--accentGrey); border-top-right-radius: var(--bs-btn-border-radius); border-bottom-right-radius: var(--bs-btn-border-radius);">
			<span id="` + uniqueID + `Label">` + template.HTMLEscapeString(displayText) + `</span>
		</button>
		<ul class="dropdown-menu w-100 overflow-auto" style="max-height: 250px;" aria-labelledby="` + uniqueID + `Dropdown">
	`)

	for _, role := range filteredRoles {
		checked := ""
		if selectedMap[role.ID] {
			checked = " checked"
		}
		disabled := ""
		disabledMsg := ""
		if highestBotRolePosition <= role.Position {
			disabled = " disabled"
			disabledMsg = " (bot higher than role)"
		}

		menu.WriteString(`<li>`)
		menu.WriteString(`<label class="dropdown-item` + disabled + `">`)
		menu.WriteString(`<input type="checkbox" class="dropDownRoleCheckbox me-2" value="` + role.ID + `"` + checked + disabled + `>`)
		menu.WriteString(template.HTMLEscapeString(role.Name) + disabledMsg)
		menu.WriteString(`</label></li>`)
	}
	menu.WriteString(`</ul>`)
	jsonVal, _ := json.Marshal(selectedRoleIDs)
	menu.WriteString(`<input type="hidden" id="` + uniqueID + `" name="` + uniqueID + `" value="` + template.HTMLEscapeString(string(jsonVal)) + `">`)
	menu.WriteString(`</div>`)
	return template.HTML(menu.String())
}

// channelOptionsSingle generates HTML options for singular channel selection
// channels: slice of Discord channel objects
// selectedChannelID: string ID of currently selected channel
// uniqueID: string for the hidden input ID (used to retrieve and store changed data)
func channelOptionsSingle(channels []*discordgo.Channel, selectedChannelID string, uniqueID string) template.HTML {
	filteredChannels := make([]*discordgo.Channel, 0, len(channels))
	for _, channel := range channels {
		if channel.Type != 0 {
			continue
		}
		filteredChannels = append(filteredChannels, channel)
	}
	sort.Slice(filteredChannels, func(i, j int) bool {
		return filteredChannels[i].Position > filteredChannels[j].Position
	})

	// Button label
	displayText := "Select channel"
	if len(selectedChannelID) > 0 {
		label := ""
		for _, channel := range filteredChannels {
			if selectedChannelID != channel.ID {
				continue
			}
			label = channel.Name
			break
		}
		if len(label) > 30 {
			displayText = "1 Selected"
		} else {
			displayText = label
		}
	}

	var menu strings.Builder
	menu.WriteString(`<div class="input-group mb-3">`)
	menu.WriteString(`
		<button class="btn dropdown-toggle text-start flex-grow-1 text-white" type="button" data-bs-toggle="dropdown" style="background-color: var(--basePurple); border: 1px solid var(--accentGrey); border-top-right-radius: var(--bs-btn-border-radius); border-bottom-right-radius: var(--bs-btn-border-radius);">
			<span id="` + uniqueID + `Label">` + template.HTMLEscapeString(displayText) + `</span>
		</button>
		<ul class="dropdown-menu w-100 overflow-auto" style="max-height: 250px;" aria-labelledby="` + uniqueID + `Dropdown">
		<li><a class="dropdown-item channelListItem" data-value="">None</a></li>
	`)

	for _, channel := range filteredChannels {
		menu.WriteString(`<li>`)
		menu.WriteString(`<a class="dropdown-item channelListItem" data-value="` + channel.ID + `">`)
		menu.WriteString(template.HTMLEscapeString(channel.Name))
		menu.WriteString(`</a></li>`)
	}

	menu.WriteString(`</ul>`)
	menu.WriteString(`<input type="hidden" id="` + uniqueID + `" name="` + uniqueID + `" value="` + template.HTMLEscapeString(selectedChannelID) + `">`)
	menu.WriteString(`</div>`)
	return template.HTML(menu.String())
}
