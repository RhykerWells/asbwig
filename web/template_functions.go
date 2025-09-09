package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	templateFunctions = map[string]interface{}{
		"dict": dict,
		"seq": seq,
		"inSlice": inSlice,
		"toJson": toJson,
		"lower": lower,

	// Forms content
		"roleOptionsMulti": roleOptionsMulti,
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


func seq(start, end int) []int {
	var result []int
	for i := start; i <= end; i++ {
		result = append(result, i)
	}
	return result
}

func inSlice(val string, slice interface{}) bool {
	switch s := slice.(type) {
	case []string:
		for _, item := range s {
			if item == val {
				return true
			}
		}
	case []interface{}:
		for _, item := range s {
			if str, ok := item.(string); ok && str == val {
				return true
			}
		}
	}
	return false
}

func toJson(v interface{}) template.JS {
	b, _ := json.Marshal(v)
	return template.JS(b)
}

func lower(str string) string {
	return strings.ToLower(str)
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