package toast

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"text/template"

	uuid "github.com/nu7hatch/gouuid"
)

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

type ToastManager struct {
	AppId       string
	DisplayName string
	Icon        string
}

func NewToastManager(appId string, displayName string, iconPath string) (*ToastManager, error) {
	fullIconPath, err := filepath.Abs(iconPath)
	if err != nil {
		return nil, err
	}

	toastManager := &ToastManager{
		AppId:       appId,
		DisplayName: displayName,
		Icon:        fullIconPath,
	}

	err = toastManager.registerToastManager()
	if err != nil {
		return nil, err
	}

	return toastManager, nil
}

func (n *ToastManager) registerToastManager() error {
	registerToastTemplate := template.New("register-toast")
	registerToastTemplate.Parse(`
		Add-Type -AssemblyName System.Windows.Forms

		$AppID = '{{if .AppId}}{{.AppId}}{{else}}com.windows.app{{end}}'
		$AppDisplayName = '{{if .DisplayName}}{{.DisplayName}}{{else}}Windows App{{end}}'
		$LogoImagePath = '{{.Icon}}'

		# Registry paths
		$regPathToastNotificationSettings = 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Notifications\Settings'
		$regPathToastApp = 'HKCU:\Software\Classes\AppUserModelId'

		# Create registry entries for notifications
		New-Item -Path "$regPathToastNotificationSettings\$AppID" -Force | Out-Null
		Set-ItemProperty -Path "$regPathToastNotificationSettings\$AppID" -Name 'ShowInActionCenter' -Value 1 -Force
		Set-ItemProperty -Path "$regPathToastNotificationSettings\$AppID" -Name 'Enabled' -Value 1 -Force

		# Create registry entries for the app
		New-Item -Path "$regPathToastApp\$AppID" -Force | Out-Null
		Set-ItemProperty -Path "$regPathToastApp\$AppID" -Name 'DisplayName' -Value $AppDisplayName -Force
		Set-ItemProperty -Path "$regPathToastApp\$AppID" -Name 'IconUri' -Value $LogoImagePath -Force
	`)

	var script bytes.Buffer
	err := registerToastTemplate.Execute(&script, n)
	if err != nil {
		return err
	}

	scriptPath, err := writeTempScript(script.String())
	if err != nil {
		return err
	}

	log.Println(script.String())
	return runScript(scriptPath)
}

type Action struct {
	Type      string
	Label     string
	Arguments string
}

func NewAction(Type string, Label string, Arguments string) Action { 
	return Action{
		Type:      Type,
		Label:     Label,
		Arguments: Arguments,
	}
}

type Toast struct {
	AppId               string
	Title               string
	Icon                string
	Message             string
	ActivationType      string
	ActivationArguments string
	Actions             []Action
	Audio               ToastAudio
	Loop                bool
	Duration            ToastDuration
}

func (toastManager *ToastManager) NewToast(Title string, Message string, ActivationType string, ActivationArguments string, Actions []Action, Audio ToastAudio, Loop bool, Duration ToastDuration) *Toast {
	return &Toast{
		AppId:               toastManager.AppId,
		Title:               Title,
		Icon:                toastManager.Icon,
		Message:             Message,
		ActivationType:      ActivationType,
		ActivationArguments: ActivationArguments,
		Actions:             Actions,
		Audio:               Audio,
		Loop:                Loop,
		Duration:            Duration,
	}
}

func (toastManager *ToastManager) NewSimpleToast(title string, message string) *Toast {
	return &Toast{
		AppId:   toastManager.AppId,
		Title:   title,
		Icon:    toastManager.Icon,
		Message: message,
		ActivationType:      "protocol",
		ActivationArguments: "",
		Audio:   Silent,
		Duration: Short,
		Loop:    false,
		Actions: []Action{},
	}
}

func (toast *Toast) Show() error {
	toastTemplate := template.New("toast")
	toastTemplate.Parse(`
		[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
		[Windows.UI.Notifications.ToastNotification, Windows.UI.Notifications, ContentType = WindowsRuntime] | Out-Null
		[Windows.Data.Xml.Dom.XmlDocument, Windows.Data.Xml.Dom.XmlDocument, ContentType = WindowsRuntime] | Out-Null

		$APP_ID = '{{if .AppId}}{{.AppId}}{{else}}Windows App{{end}}'
		$ACTIVATION_TYPE = '{{if .ActivationType}}{{.ActivationType}}{{else}}protocol{{end}}'
		$AUDIO = '{{if .Audio}}{{.Audio}}{{else}}Default{{end}}'
		$DURATION = '{{if .Duration}}{{.Duration}}{{else}}Short{{end}}'

		$template = @"
		<toast activationType="$ACTIVATION_TYPE" launch="{{.ActivationArguments}}" duration="$DURATION">
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
			<audio src="$AUDIO" loop="{{.Loop}}" />
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

	var script bytes.Buffer
	err := toastTemplate.Execute(&script, toast)
	if err != nil {
		return err
	}

	scriptPath, err := writeTempScript(script.String())
	if err != nil {
		return err
	}

	log.Println(script.String())
	return runScript(scriptPath)
}

func writeTempScript(script string) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	path := filepath.Join(os.TempDir(), id.String()+".ps1")
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	out := append(bomUtf8, []byte(script)...)
	return path, os.WriteFile(path, out, 0600)
}

// runScript runs the given script using PowerShell.
func runScript(scriptPath string) error {
	cmd := exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", "-File", scriptPath)

	// Capture the stdout and stderr
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	// Hide the window
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
