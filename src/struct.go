package main

// Data .
type Data struct {
	Message    string `json:"Message,omitempty"`
	Identifier int64  `json:"Identifier,omitempty"`
	Type       string `json:"Type,omitempty"`
	Stacktrace string `json:"Stacktrace,omitempty"`
	Name       string `json:"Name,omitempty"`
}

// Chat .
type Chat struct {
	Channel  int64
	Message  string
	UserID   string
	Username string
	Color    string
	Time     uint64
}

/* [-1] {
  "Channel": 2,
  "Message": "hi",
  "UserId": "0",
  "Username": "SERVER",
  "Color": "#eee",
  "Time": 1604453507
} */
