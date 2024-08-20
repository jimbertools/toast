package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/jimbertools/toast/pkg/toast"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()

	app.Name = "toast"
	app.Usage = "Windows 10 toasts"
	app.Version = "v1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jacob Marshall",
			Email: "go-toast@jacobmarshall.co",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "app-id, id",
			Usage: "the app identifier (used for grouping multiple toasts)",
		},
		cli.StringFlag{
			Name:  "title, t",
			Usage: "the main toast title/heading",
		},
		cli.StringFlag{
			Name:  "message, m",
			Usage: "the toast's main message (new lines as separator)",
		},
		cli.StringFlag{
			Name:  "icon, i",
			Usage: "the app icon path (displays to the left of the toast)",
		},
		cli.StringFlag{
			Name:  "activation-type",
			Value: "protocol",
			Usage: "the type of action to invoke when the user clicks the toast",
		},
		cli.StringFlag{
			Name:  "activation-arg",
			Usage: "the activation argument",
		},
		cli.StringSliceFlag{
			Name:  "action",
			Usage: "optional action button",
		},
		cli.StringSliceFlag{
			Name:  "action-type",
			Usage: "the type of action button",
		},
		cli.StringSliceFlag{
			Name:  "action-arg",
			Usage: "the action button argument",
		},
		cli.StringFlag{
			Name:  "audio",
			Value: "silent",
			Usage: "which kind of audio should be played",
		},
		cli.BoolFlag{
			Name:  "loop",
			Usage: "whether to loop the audio",
		},
		cli.StringFlag{
			Name:  "duration",
			Value: "short",
			Usage: "how long the toast should display for",
		},
	}

	app.Action = func(c *cli.Context) error {
		appID := c.String("app-id")
		title := c.String("title")
		message := c.String("message")
		icon := c.String("icon")
		activationType := c.String("activation-type")
		activationArg := c.String("activation-arg")
		audio, _ := Audio(c.String("audio"))
		duration, _ := toast.Duration(c.String("duration"))
		loop := c.Bool("loop")

		var actions []toast.Action
		actionTexts := c.StringSlice("action")
		actionTypes := c.StringSlice("action-type")
		actionArgs := c.StringSlice("action-arg")

		for index, actionLabel := range actionTexts {
			var actionType string = "protocol"
			var actionArg string
			if len(actionTypes) > index {
				actionType = actionTypes[index]
			}
			if len(actionArgs) > index {
				actionArg = actionArgs[index]
			}
			actions = append(actions, toast.Action{
				Type:      actionType,
				Label:     actionLabel,
				Arguments: actionArg,
			})
		}

		notification := &toast.Notification{
			AppID:               appID,
			Title:               title,
			Message:             message,
			Icon:                icon,
			Actions:             actions,
			ActivationType:      activationType,
			ActivationArguments: activationArg,
			Audio:               audio,
			Loop:                loop,
			Duration:            duration,
		}

		if err := notification.Push(); err != nil {
			log.Fatalln(err)
		}

		return nil
	}

	app.Run(os.Args)
}



// Returns a toastAudio given a user-provided input (useful for cli apps).
//
// If the "name" doesn't match, then the default toastAudio is returned, along with ErrorInvalidAudio.
//
// The following names are valid;
//   - default
//   - im
//   - mail
//   - reminder
//   - sms
//   - loopingalarm
//   - loopimgalarm[2-10]
//   - loopingcall
//   - loopingcall[2-10]
//   - silent
//
// Handle the error appropriately according to how your app should work.
func Audio(name string) (toast.ToastAudio, error) {
	switch strings.ToLower(name) {
	case "default":
		return toast.Default, nil
	case "im":
		return toast.IM, nil
	case "mail":
		return toast.Mail, nil
	case "reminder":
		return toast.Reminder, nil
	case "sms":
		return toast.SMS, nil
	case "loopingalarm":
		return toast.LoopingAlarm, nil
	case "loopingalarm2":
		return toast.LoopingAlarm2, nil
	case "loopingalarm3":
		return toast.LoopingAlarm3, nil
	case "loopingalarm4":
		return toast.LoopingAlarm4, nil
	case "loopingalarm5":
		return toast.LoopingAlarm5, nil
	case "loopingalarm6":
		return toast.LoopingAlarm6, nil
	case "loopingalarm7":
		return toast.LoopingAlarm7, nil
	case "loopingalarm8":
		return toast.LoopingAlarm8, nil
	case "loopingalarm9":
		return toast.LoopingAlarm9, nil
	case "loopingalarm10":
		return toast.LoopingAlarm10, nil
	case "loopingcall":
		return toast.LoopingCall, nil
	case "loopingcall2":
		return toast.LoopingCall2, nil
	case "loopingcall3":
		return toast.LoopingCall3, nil
	case "loopingcall4":
		return toast.LoopingCall4, nil
	case "loopingcall5":
		return toast.LoopingCall5, nil
	case "loopingcall6":
		return toast.LoopingCall6, nil
	case "loopingcall7":
		return toast.LoopingCall7, nil
	case "loopingcall8":
		return toast.LoopingCall8, nil
	case "loopingcall9":
		return toast.LoopingCall9, nil
	case "loopingcall10":
		return toast.LoopingCall10, nil
	case "silent":
		return toast.Silent, nil
	default:
		return toast.Default, toast.ErrorInvalidAudio
	}
}