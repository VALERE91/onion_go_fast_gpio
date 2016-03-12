package fast_gpio

import "testing"

func TestOpen(t *testing.T) {
	gpio := FastGPIO{}
	err := gpio.Open()

	if err != nil {
		t.Error("Error in open : ", err)
	}
}