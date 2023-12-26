package console

import (
	"Dur4nC2/misc/protobuf/clientpb"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"log"
	"time"
)

const (
	// ANSI Colors
	Normal    = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Blue      = "\033[34m"
	Cyan      = "\033[36m"
	Bold      = "\033[1m"
	Clearln   = "\r\x1b[2K"
	Underline = "\033[4m"

	// Info - Display colorful information
	Info = Bold + Cyan + "[*] " + Normal
	// Warn - Warn a user
	Warn = Bold + Red + "[!] " + Normal
	// Success - Diplay success
	Success = Bold + Green + "[+] " + Normal
)

// Observer - A function to call when the sessions change
type Observer func(*clientpb.Beacon)
type ActiveTarget struct {
	beacon     *clientpb.Beacon
	observers  map[int]Observer
	observerID int
}
type ServerConsoleClient struct {
	App          *grumble.App
	ActiveTarget *ActiveTarget
	IsServer     bool
}
type BindCmds func(console *ServerConsoleClient)

func Start(cmds BindCmds) error {
	con := &ServerConsoleClient{
		App: grumble.New(&grumble.Config{
			Name:                  "Server",
			Description:           "Server Client",
			PromptColor:           color.New(),
			HelpHeadlineColor:     color.New(),
			HelpHeadlineUnderline: true,
			HelpSubCommands:       true,
		}),
		ActiveTarget: &ActiveTarget{
			observers:  map[int]Observer{},
			observerID: 0,
		},
	}
	con.App.SetPrintASCIILogo(func(_ *grumble.App) {
		con.PrintLogo()
	})
	con.App.SetPrompt(con.GetPrompt())
	cmds(con)
	con.ActiveTarget.AddObserver(func(_ *clientpb.Beacon) {
		con.App.SetPrompt(con.GetPrompt())
	})
	err := con.App.Run()
	if err != nil {
		log.Printf("Run loop returned error: %v", err)
	}
	return err
}

// *** Console Print utilities ***

func (con *ServerConsoleClient) GetPrompt() string {
	prompt := Underline + "dur4nc2" + Normal

	if con.ActiveTarget.GetBeacon() != nil {
		prompt += fmt.Sprintf(Bold+Blue+" (%s)%s", con.ActiveTarget.GetBeacon().Name, Normal)
	}

	prompt += " > "
	return Clearln + prompt
}

func (con *ServerConsoleClient) PrintLogo() {
	con.Println("[+] Dur4nC2 Server console")
}

func (con *ServerConsoleClient) Printf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(con.App.Stdout(), format, args...)
}

func (con *ServerConsoleClient) Println(args ...interface{}) (n int, err error) {
	return fmt.Fprintln(con.App.Stdout(), args...)
}

func (con *ServerConsoleClient) PrintInfof(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(con.App.Stdout(), Clearln+Info+format, args...)
}

func (con *ServerConsoleClient) PrintSuccessf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(con.App.Stdout(), Clearln+Success+format, args...)
}

func (con *ServerConsoleClient) PrintWarnf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(con.App.Stdout(), Clearln+"⚠️  "+Normal+format, args...)
}

func (con *ServerConsoleClient) PrintErrorf(format string, args ...interface{}) (n int, err error) {
	return fmt.Fprintf(con.App.Stderr(), Clearln+Warn+format, args...)
}
func (con *ServerConsoleClient) FormatDateDelta(t time.Time, includeDate bool, color bool) string {
	nextTime := t.Format(time.UnixDate)

	var interval string

	if t.Before(time.Now()) {
		if includeDate {
			interval = fmt.Sprintf("%s (%s ago)", nextTime, time.Since(t).Round(time.Second))
		} else {
			interval = time.Since(t).Round(time.Second).String()
		}
		if color {
			interval = fmt.Sprintf("%s%s%s", Bold+Red, interval, Normal)
		}
	} else {
		if includeDate {
			interval = fmt.Sprintf("%s (in %s)", nextTime, time.Until(t).Round(time.Second))
		} else {
			interval = time.Until(t).Round(time.Second).String()
		}
		if color {
			interval = fmt.Sprintf("%s%s%s", Bold+Green, interval, Normal)
		}
	}
	return interval
}

// *** Active Target ***

func (s *ActiveTarget) GetBeacon() *clientpb.Beacon {
	return s.beacon
}

func (s *ActiveTarget) Set(beacon *clientpb.Beacon) {
	if beacon == nil {
		//panic("cannot set nil beacon")
		return
	}
	s.beacon = beacon
	for _, observer := range s.observers {
		observer(s.beacon)
	}
	return
}

// AddObserver - Observers to notify when the active session changes
func (s *ActiveTarget) AddObserver(observer Observer) int {
	s.observerID++
	s.observers[s.observerID] = observer
	return s.observerID
}

func (s *ActiveTarget) RemoveObserver(observerID int) {
	delete(s.observers, observerID)
}

// Background - Background the active session
func (s *ActiveTarget) Background() {
	s.beacon = nil
	for _, observer := range s.observers {
		observer(nil)
	}
}
