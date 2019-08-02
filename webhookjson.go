package main

type Object struct {
	Entries []Entry `json:"entry"`
}

type Entry struct {
	Messages []Message `json:"messaging"`
}

type Message struct {
	Sender    Ident    `json:"sender"`
	Recipient Ident    `json:"recipient"`
	Mes       Mess     `json:"message,omitempty"`
	Pb        Postback `json:"postback,omitempty"`
}

type Ident struct {
	ID string `json:"id"`
}

type Mess struct {
	Text string `json:"text"`
}

type Postback struct {
	Title string `json:"title"`
	P     string `json:"payload"`
}
