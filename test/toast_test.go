package toast_test

import (
	"testing"

	"github.com/jimbertools/toast/pkg/toast"
)

func TestSimpleToast(t *testing.T) {
	toastManager, err := toast.NewToastManager("com.windows.app", "Windows App", "testdata/4201973.png")
	if err != nil {
		t.Error(err)
	}
	toast := toastManager.NewSimpleToast("Hello World", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.")

	err = toast.Show()
	if err != nil {
		t.Error(err)
	}
}

func TestToast(t *testing.T) {
	toastManager, err := toast.NewToastManager("Example App", "Hello World", "testdata/4201973.png")
	if err != nil {
		t.Error(err)
	}

	toast := toastManager.NewToast("Hello World", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "protocol", "https://google.com", []toast.Action{
		{Type: "protocol", Label: "I'm a button", Arguments: ""},
		{Type: "protocol", Label: "Me too!", Arguments: ""},
	}, toast.Silent, false, toast.Long)

	err = toast.Show()

	if err != nil {
		t.Error(err)
	}
}
