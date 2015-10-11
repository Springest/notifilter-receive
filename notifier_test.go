package main

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
)

type LocalMessageNotifier struct {
	class     string
	message   []byte
	processed bool
}

func (mn *LocalMessageNotifier) sendMessage(class string, data []byte) NotifierResponse {
	mn.class = class
	mn.message = data
	mn.processed = true

	return NotifierResponse{}
}

func TestNewNotifier(t *testing.T) {
	n := Notifier{}
	assert.Equal(t, n.newNotifier(), &slackNotifier{})

	n.NotificationType = "email"
	assert.Equal(t, n.newNotifier(), &slackNotifier{})

	n.NotificationType = "slack"
	assert.Equal(t, n.newNotifier(), &slackNotifier{})
}

func TestNotifierCheckRulesSingle(t *testing.T) {
	var rules = types.JsonText(`[{"key": "number", "type": "number", "setting": "eq", "value": "12"}]`)
	n := Notifier{
		Id:               1,
		NotificationType: "email",
		Class:            "User",
		Template:         "name: {{.name}}",
		Rules:            rules,
	}

	var jt = types.JsonText(`{"active": true, "name": "Go", "number": "12"}`)
	s := Stat{"Mark", jt}

	assert.Equal(t, n.checkRules(&s), true)
}

func TestNotifierCheckRulesMultiple(t *testing.T) {
	var rules = types.JsonText(`[{"key": "number", "type": "number", "setting": "eq", "value": "12"},
	{"key": "name", "type": "string", "setting": null, "value": "Go"}]`)
	n := Notifier{
		Id:               1,
		NotificationType: "email",
		Class:            "User",
		Template:         "name: {{.name}}",
		Rules:            rules,
	}

	var jt = types.JsonText(`{"active": true, "name": "Go", "number": "12"}`)
	s := Stat{"Mark", jt}

	assert.Equal(t, n.checkRules(&s), true)
}

func TestNotifierCheckRulesSettingIsNull(t *testing.T) {
	var rules = types.JsonText(`[{"key": "name", "type": "string", "setting": null "value": "Go"}]`)
	n := Notifier{
		Id:               1,
		NotificationType: "email",
		Class:            "User",
		Template:         "name: {{.name}}",
		Rules:            rules,
	}

	var jt = types.JsonText(`{"active": true, "name": "Go", "number": 12}`)
	s := Stat{"Mark", jt}

	assert.Equal(t, n.checkRules(&s), true)
}

func TestNotifierCheckRulesSettingIsBlank(t *testing.T) {
	var rules = types.JsonText(`[{"key": "name", "type": "string", "setting": "", "value": "Go"}]`)
	n := Notifier{
		Id:               1,
		NotificationType: "email",
		Class:            "User",
		Template:         "name: {{.name}}",
		Rules:            rules,
	}

	var jt = types.JsonText(`{"active": true, "name": "Go", "number": "12"}`)
	s := Stat{"Mark", jt}

	assert.Equal(t, n.checkRules(&s), true)
}

func TestNotifierNotify(t *testing.T) {
	n := Notifier{
		Id:               1,
		NotificationType: "email",
		Class:            "User",
		Template:         "name: {{.name}}",
	}

	var jt = types.JsonText(`{"active": true, "name": "Go", "number": "12"}`)
	s := Stat{"Mark", jt}

	mn := &LocalMessageNotifier{}
	n.notify(&s, mn)

	assert.Equal(t, mn.class, "Mark")
	assert.Equal(t, mn.message, []byte("name: Go"))
	assert.Equal(t, mn.processed, true)
}

func TestNotifierNotifyReturnsEarlyIfRulesAreNotMet(t *testing.T) {
	var rules = types.JsonText(`[{"key": "number", "type": "number", "setting": "gt", "value": "1"}]`)
	n := Notifier{
		Id:               1,
		NotificationType: "email",
		Class:            "User",
		Template:         "name: {{.name}}",
		Rules:            rules,
	}

	var jt = types.JsonText(`{"active": true, "name": "Go", "number": "0"}`)
	s := Stat{"Mark", jt}

	mn := &LocalMessageNotifier{}
	n.notify(&s, mn)

	assert.Equal(t, mn.processed, false)
}

func TestNotifierRenderTemplate(t *testing.T) {
	n := Notifier{
		Id:               1,
		NotificationType: "email",
		Class:            "User",
		Template:         "name: {{.name}}",
	}

	var jt = types.JsonText(`{"active": true, "name": "Go", "number": "12"}`)
	s := Stat{"Mark", jt}

	result := n.renderTemplate(&s)
	expected := "name: Go"
	assert.Equal(t, result, expected)
}

func TestNotifierRenderTemplateWithLogic(t *testing.T) {
	template := `{{ if .active }}Active!{{ else }}inactive{{ end }}`
	n := Notifier{
		Id:               1,
		NotificationType: "email",
		Class:            "User",
		Template:         template,
	}

	var jt = types.JsonText(`{"active": true, "name": "Go", "number": "12"}`)
	s := Stat{"Mark", jt}

	result := n.renderTemplate(&s)
	expected := "Active!"
	assert.Equal(t, result, expected)

	jt = types.JsonText(`{"active": false, "name": "Go", "number": "12"}`)
	s = Stat{"Mark", jt}

	result = n.renderTemplate(&s)
	expected = "inactive"
	assert.Equal(t, result, expected)
}
