package main

// This type provides a json encoding for the message body sent to messenger
type MessageAttach struct {
	MessagingType string `json:"messaging_type"`
	Recipient     Ident  `json:"recipient"`
	Mes           M      `json:"message"`
}

type M struct {
	Attach Attachment `json:"attachment"`
}

type Attachment struct {
	TypeOf string  `json:"type"`
	P      Payload `json:"payload"`
}

type Payload struct {
	TemplateType string    `json:"template_type"`
	Elem         []Element `json:"elements"`
}

type Element struct {
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Buttons  []Button `json:"buttons,omitempty"`
}

func NewElement(reminderText string, addButtons bool) (e Element) {
	e.Title = "Alert"
	e.Subtitle = "Reminder: " + reminderText
	if addButtons {
		e.Buttons = []Button{
			NewButton("Accept"),
			NewButton("Snooze"),
		}
	}
	return e
}

type Button struct {
	TypeOf string `json:"type"`
	Title  string `json:"title"`
	P      string `json:"payload"`
}

func NewButton(title string) (b Button) {
	b.TypeOf = "postback"
	b.P = title
	b.Title = title
	return b
}

func NewMessageAttach(to string, reminderText []string) *MessageAttach {
	m := new(MessageAttach)
	m.Recipient.ID = to
	m.MessagingType = "RESPONSE"
	m.Mes.Attach.TypeOf = "template"
	m.Mes.Attach.P.TemplateType = "generic"
	addButton := true
	for _, s := range reminderText {
		m.Mes.Attach.P.Elem = append(m.Mes.Attach.P.Elem, NewElement(s, addButton))
		addButton = false
	}
	return m
}

func NewWelcome(to string) *MessageAttach {
	m := new(MessageAttach)
	m.Recipient.ID = to
	m.MessagingType = "RESPONSE"
	m.Mes.Attach.TypeOf = "template"
	m.Mes.Attach.P.TemplateType = "generic"
	m.Mes.Attach.P.Elem = append(
		m.Mes.Attach.P.Elem,
		NewElement("Welcome to Reminder Bot", false),
	)
	m.Mes.Attach.P.Elem = append(
		m.Mes.Attach.P.Elem,
		NewElement("Just send a message and we'll take care of it", false),
	)
	return m
}
