# Toast

A go package for Windows 10 toast notifications.

As seen in [jacobmarshall/pokevision-cli](https://github.com/jacobmarshall/pokevision-cli).

## CLI

As well as using go-toast within your Go projects, you can also utilise the CLI - for any of your projects.

```ps1
./toast.exe
--app-id "Example App"
--title "Hello World"
--message "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
--icon "C:\Users\Wannes\Desktop\Work\Jimber\Go\toast\test\testdata\4201973.png"
--audio "default"
--duration "long"
--activation-arg "https://google.com"
--action "Open maps" --action-arg "bingmaps:?q=sushi"
--action "Open browser" -action-arg ""
```

## Example

```go
package main

import (
    "log"

    "github.com/jimbertools/toast/pkg/toast"
)

func main() {
  toastManager, err := toast.NewToastManager("com.windows.app", "Windows App", "C:Path/to/your/image.png")
	if err != nil {
		log.Println(err)
	}

	toast := toastManager.NewSimpleToast("Hello World", "Lorem ipsum dolor sit amet, consectetur adipiscing elit.")

	err = toast.Show()
	if err != nil {
		log.Println(err)
	}
}
```

