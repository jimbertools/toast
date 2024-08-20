package toast

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"text/template"

	uuid "github.com/nu7hatch/gouuid"
)

var toastTemplate *template.Template

var (
	ErrorInvalidAudio    error = errors.New("toast: invalid audio")
	ErrorInvalidDuration       = errors.New("toast: invalid duration")
)

type ToastAudio string

const (
	Default        ToastAudio = "ms-winsoundevent:Notification.Default"
	IM             ToastAudio = "ms-winsoundevent:Notification.IM"
	Mail           ToastAudio = "ms-winsoundevent:Notification.Mail"
	Reminder       ToastAudio = "ms-winsoundevent:Notification.Reminder"
	SMS            ToastAudio = "ms-winsoundevent:Notification.SMS"
	LoopingAlarm   ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm"
	LoopingAlarm2  ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm2"
	LoopingAlarm3  ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm3"
	LoopingAlarm4  ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm4"
	LoopingAlarm5  ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm5"
	LoopingAlarm6  ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm6"
	LoopingAlarm7  ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm7"
	LoopingAlarm8  ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm8"
	LoopingAlarm9  ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm9"
	LoopingAlarm10 ToastAudio = "ms-winsoundevent:Notification.Looping.Alarm10"
	LoopingCall    ToastAudio = "ms-winsoundevent:Notification.Looping.Call"
	LoopingCall2   ToastAudio = "ms-winsoundevent:Notification.Looping.Call2"
	LoopingCall3   ToastAudio = "ms-winsoundevent:Notification.Looping.Call3"
	LoopingCall4   ToastAudio = "ms-winsoundevent:Notification.Looping.Call4"
	LoopingCall5   ToastAudio = "ms-winsoundevent:Notification.Looping.Call5"
	LoopingCall6   ToastAudio = "ms-winsoundevent:Notification.Looping.Call6"
	LoopingCall7   ToastAudio = "ms-winsoundevent:Notification.Looping.Call7"
	LoopingCall8   ToastAudio = "ms-winsoundevent:Notification.Looping.Call8"
	LoopingCall9   ToastAudio = "ms-winsoundevent:Notification.Looping.Call9"
	LoopingCall10  ToastAudio = "ms-winsoundevent:Notification.Looping.Call10"
	Silent         ToastAudio = "silent"
)

type ToastDuration string

const (
	Short ToastDuration = "short"
	Long  ToastDuration = "long"
)

func init() {
	toastTemplate = template.New("toast")
	toastTemplate.Parse(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

$APP_ID = '{{if .AppID}}{{.AppID}}{{else}}Windows App{{end}}'

$template = @"
<toast activationType="{{.ActivationType}}" launch="{{.ActivationArguments}}" duration="{{.Duration}}">
    <visual>
        <binding template="ToastGeneric">
            {{if .Icon}}
            <image placement="appLogoOverride" src="{{.Icon}}" />
            {{end}}
            {{if .Title}}
            <text><![CDATA[{{.Title}}]]></text>
            {{end}}
            {{if .Message}}
            <text><![CDATA[{{.Message}}]]></text>
            {{end}}
        </binding>
    </visual>
    {{if ne .Audio "silent"}}
	<audio src="{{.Audio}}" loop="{{.Loop}}" />
	{{else}}
	<audio silent="true" />
	{{end}}
    {{if .Actions}}
    <actions>
        {{range .Actions}}
        <action activationType="{{.Type}}" content="{{.Label}}" arguments="{{.Arguments}}" />
        {{end}}
    </actions>
    {{end}}
</toast>
"@

$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($template)
$toast = New-Object Windows.UI.Notifications.ToastNotification $xml
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($APP_ID).Show($toast)
    `)
}

// Notification
//
// The toast notification data. The following fields are strongly recommended;
//   - AppID
//   - Title
//
// If no toastAudio is provided, then the toast notification will be silent.
// You can set the toast to have a default audio by setting "Audio" to "toast.Default", or if your go app takes
// user-provided input for audio, call the "toast.Audio(name)" func.
//
// The AppID is shown beneath the toast message (in certain cases), and above the notification within the Action
// Center - and is used to group your notifications together. It is recommended that you provide a "pretty"
// name for your app, and not something like "com.example.MyApp".
//
// If no Title is provided, but a Message is, the message will display as the toast notification's title -
// which is a slightly different font style (heavier).
//
// The Icon should be an absolute path to the icon (as the toast is invoked from a temporary path on the user's
// system, not the working directory).
//
// If you would like the toast to call an external process/open a webpage, then you can set ActivationArguments
// to the uri you would like to trigger when the toast is clicked. For example: "https://google.com" would open
// the Google homepage when the user clicks the toast notification.
// By default, clicking the toast just hides/dismisses it.
//
// The following would show a notification to the user letting them know they received an email, and opens
// gmail.com when they click the notification. It also makes the Windows 10 "mail" sound effect.
//
//	toast := toast.Notification{
//	    AppID:               "Google Mail",
//	    Title:               email.Subject,
//	    Message:             email.Preview,
//	    Icon:                "C:/Program Files/Google Mail/icons/logo.png",
//	    ActivationArguments: "https://gmail.com",
//	    Audio:               toast.Mail,
//	}
//
//	err := toast.Push()
type Notification struct {
	// The name of your app. This value shows up in Windows 10's Action Centre, so make it
	// something readable for your users. It can contain spaces, however special characters
	// (eg. é) are not supported.
	AppID string

	// The main title/heading for the toast notification.
	Title string

	// The single/multi line message to display for the toast notification.
	Message string

	// An optional path to an image on the OS to display to the left of the title & message.
	Icon string

	// The type of notification level action (like toast.Action)
	ActivationType string

	// The activation/action arguments (invoked when the user clicks the notification)
	ActivationArguments string

	// Optional action buttons to display below the notification title & message.
	Actions []Action

	// The audio to play when displaying the toast
	Audio ToastAudio

	// Whether to loop the audio (default false)
	Loop bool

	// How long the toast should show up for (short/long)
	Duration ToastDuration
}

// Action
//
// Defines an actionable button.
// See https://msdn.microsoft.com/en-us/windows/uwp/controls-and-patterns/tiles-and-notifications-adaptive-interactive-toasts for more info.
//
// Only protocol type action buttons are actually useful, as there's no way of receiving feedback from the
// user's choice. Examples of protocol type action buttons include: "bingmaps:?q=sushi" to open up Windows 10's
// maps app with a pre-populated search field set to "sushi".
//
//	toast.Action{"protocol", "Open Maps", "bingmaps:?q=sushi"}
type Action struct {
	Type      string
	Label     string
	Arguments string
}

func (n *Notification) applyDefaults() {
	if n.ActivationType == "" {
		n.ActivationType = "protocol"
	}
	if n.Duration == "" {
		n.Duration = Short
	}
	if n.Audio == "" {
		n.Audio = Default
	}
}

func (n *Notification) buildXML() (string, error) {
	var out bytes.Buffer
	err := toastTemplate.Execute(&out, n)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

// Builds the Windows PowerShell script & invokes it, causing the toast to display.
//
// Note: Running the PowerShell script is by far the slowest process here, and can take a few
// seconds in some cases.
//
//	notification := toast.Notification{
//	    AppID: "Example App",
//	    Title: "My notification",
//	    Message: "Some message about how important something is...",
//	    Icon: "go.png",
//	    Actions: []toast.Action{
//	        {"protocol", "I'm a button", ""},
//	        {"protocol", "Me too!", ""},
//	    },
//	}
//	err := notification.Push()
//	if err != nil {
//	    log.Fatalln(err)
//	}
func (n *Notification) Push() error {
	xml, err := n.script()
	if err != nil {
		return err
	}

	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	filePath := filepath.Join(os.TempDir(), id.String()+".ps1")
	defer os.Remove(filePath)
	
	err = writeScript(filePath, xml)
	if err != nil {
		return err
	}

	return invokeScript(filePath)
}

func (notification *Notification) script() (xml string, err error) {
	notification.applyDefaults()

	return notification.buildXML()
}

func writeScript(filePath, xml string) error { 
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	out := append(bomUtf8, []byte(xml)...)
	err := os.WriteFile(filePath, out, 0600)
	return err
}

func invokeScript(filePath string) error {
	cmd := exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", "-File", filePath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
