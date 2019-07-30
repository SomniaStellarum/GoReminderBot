package main

// This type provides a json encoding for the message body sent to messenger
type MessageBody struct {
  MessagingType string `json:"messaging_type"`
  Recipient Ident `json:"recipient"`
  Mes Mess `json:"message"`
}

func NewResponseMessage(to, reply string) *MessageBody {
  m := new(MessageBody)
  m.Mes.Text = reply
  m.Recipient.ID = to
  m.MessagingType = "RESPONSE"
  return m
}
