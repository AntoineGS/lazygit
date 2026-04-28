package oscommands

import (
	"os/exec"
	"strings"
	"sync"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/utils"
)

func getRunner() *cmdObjRunner {
	log := utils.NewDummyLog()
	return &cmdObjRunner{
		log:   log,
		guiIO: NewNullGuiIO(log),
	}
}

func toChanFn(f func(ct CredentialType) string) func(CredentialType) <-chan string {
	return func(ct CredentialType) <-chan string {
		ch := make(chan string)

		go func() {
			ch <- f(ct)
		}()

		return ch
	}
}

func TestProcessOutput(t *testing.T) {
	defaultPromptUserForCredential := func(ct CredentialType) string {
		switch ct {
		case Password:
			return "password"
		case Username:
			return "username"
		case Passphrase:
			return "passphrase"
		case PIN:
			return "pin"
		case Token:
			return "token"
		default:
			panic("unexpected credential type")
		}
	}

	scenarios := []struct {
		name                    string
		promptUserForCredential func(CredentialType) string
		output                  string
		expectedToWrite         string
	}{
		{
			name:                    "no output",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "",
			expectedToWrite:         "",
		},
		{
			name:                    "password prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Password:",
			expectedToWrite:         "password",
		},
		{
			name:                    "password prompt 2",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Bill's password:",
			expectedToWrite:         "password",
		},
		{
			name:                    "password prompt 3",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Password for 'Bill':",
			expectedToWrite:         "password",
		},
		{
			name:                    "username prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Username for 'Bill':",
			expectedToWrite:         "username",
		},
		{
			name:                    "passphrase prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Enter passphrase for key '123':",
			expectedToWrite:         "passphrase",
		},
		{
			name:                    "security key pin prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Enter PIN for key '123':",
			expectedToWrite:         "pin",
		},
		{
			name:                    "pkcs11 key pin prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Enter PIN for '123':",
			expectedToWrite:         "pin",
		},
		{
			name:                    "2FA token prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "testuser 2FA Token (citadel)",
			expectedToWrite:         "token",
		},
		{
			name:                    "username and password prompt",
			promptUserForCredential: defaultPromptUserForCredential,
			output:                  "Password:\nUsername for 'Alice':\n",
			expectedToWrite:         "passwordusername",
		},
		{
			name:                    "user submits empty credential",
			promptUserForCredential: func(ct CredentialType) string { return "" },
			output:                  "Password:\n",
			expectedToWrite:         "",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			runner := getRunner()
			reader := strings.NewReader(scenario.output)
			writer := &strings.Builder{}

			cmdObj := &CmdObj{task: gocui.NewFakeTask()}
			runner.processOutput(reader, writer, toChanFn(scenario.promptUserForCredential), func() error { return nil }, cmdObj)

			if writer.String() != scenario.expectedToWrite {
				t.Errorf("expected to write '%s' but got '%s'", scenario.expectedToWrite, writer.String())
			}
		})
	}
}

func TestRecordCommandOnTask_NilTaskNoop(t *testing.T) {
	cmdObj := &CmdObj{} // no task set
	cleanup := recordCommandOnTask(cmdObj)
	cleanup() // must not panic
}

type fakeRecorderTask struct {
	gocui.Task
	mu      sync.Mutex
	current string
}

func (r *fakeRecorderTask) SetCurrentCommand(cmd string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.current = cmd
}

func TestRecordCommandOnTask_StampsAndClears(t *testing.T) {
	rec := &fakeRecorderTask{}
	cmdObj := (&CmdObj{cmd: exec.Command("git", "pull")}).WithTask(rec)

	cleanup := recordCommandOnTask(cmdObj)
	rec.mu.Lock()
	got := rec.current
	rec.mu.Unlock()
	if got != "git pull" {
		t.Fatalf("expected current=%q after stamp, got %q", "git pull", got)
	}

	cleanup()
	rec.mu.Lock()
	got = rec.current
	rec.mu.Unlock()
	if got != "" {
		t.Fatalf("expected current=%q after cleanup, got %q", "", got)
	}
}

type fakePlainTask struct {
	gocui.Task
}

func TestRecordCommandOnTask_NonRecorderTaskNoop(t *testing.T) {
	cmdObj := (&CmdObj{cmd: exec.Command("git", "status")}).WithTask(&fakePlainTask{})
	cleanup := recordCommandOnTask(cmdObj)
	cleanup() // must not panic; nothing to assert beyond no-panic
}
