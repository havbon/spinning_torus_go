package input

import (
	"github.com/eiannone/keyboard"
)

func GetKey(channel chan keyboard.KeyEvent) {
	for {
		r, key, err := keyboard.GetKey()
		if err != nil {
			continue
		}

		channel <- keyboard.KeyEvent{Key: key, Rune: r, Err: nil}
	}
}
