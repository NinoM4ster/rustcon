package main

// Data .
type Data struct {
	Message    string `json:"Message,omitempty"`
	Identifier int64  `json:"Identifier,omitempty"`
	Type       string `json:"Type,omitempty"`
	Stacktrace string `json:"Stacktrace,omitempty"`
	Name       string `json:"Name,omitempty"`
}

/* {
  "Message": "857815[857815] was killed by sentry.scientist.static (entity)",
  "Identifier": 0,
  "Type": "Generic",
  "Stacktrace": ""
} */

// Chat .
type Chat struct {
	Channel  int64  `json:"Channel,omitempty"`
	Message  string `json:"Message,omitempty"`
	UserID   string `json:"UserId,omitempty"`
	Username string `json:"Username,omitempty"`
	Color    string `json:"Color,omitempty"`
	Time     uint64 `json:"Time,omitempty"`
}

/* {
  "Channel": 2,
  "Message": "hi",
  "UserId": "0",
  "Username": "SERVER",
  "Color": "#eee",
  "Time": 1604453507
} */
