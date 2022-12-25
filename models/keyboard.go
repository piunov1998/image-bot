package models

type Key struct {
	Text string `json:"text"`
}

type Keyboard struct {
	Keyboard        [][]Key `json:"keyboard"`
	ResizeKeyboard  bool    `json:"resize_keyboard"`
	OneTimeKeyboard bool    `json:"one_time_keyboard"`
}
