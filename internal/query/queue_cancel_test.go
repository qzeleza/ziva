package query

import (
	te "github.com/charmbracelet/bubbletea"
	"testing"

	"github.com/qzeleza/ziva/internal/common"
)

func TestStripResultPrefixes_RemovesVerticalLineFromCancelMessage(t *testing.T) {
	input := "Задача\n  │────\n  │    отменено пользователем\n"
	got := stripResultPrefixes(input)
	want := "Задача\n   ────\n       отменено пользователем\n"
	if got != want {
		t.Fatalf("unexpected result: %q", got)
	}
}

type stubTask struct {
	final string
}

func (s stubTask) Title() string                           { return "" }
func (s stubTask) Run() te.Cmd                             { return nil }
func (s stubTask) Update(msg te.Msg) (common.Task, te.Cmd) { return s, nil }
func (s stubTask) View(width int) string                   { return "" }
func (s stubTask) IsDone() bool                            { return true }
func (s stubTask) FinalView(width int) string              { return s.final }
func (s stubTask) HasError() bool                          { return true }
func (s stubTask) Error() error                            { return nil }
func (s stubTask) StopOnError() bool                       { return true }
func (s stubTask) SetStopOnError(stop bool)                {}
func (s stubTask) WithNewLinesInErrors(bool) common.Task   { return s }

func TestFormatTaskResultCancelWithoutSummary(t *testing.T) {
	model := New("test")
	model.WithSummary(false)
	model.resultLineLength = 4
	task := stubTask{final: "Задача\n  │    отменено пользователем\n"}
	got := model.formatTaskResult(task, 50, true)
	want := "Задача\n     ────\n       отменено пользователем\n"
	if got != want {
		t.Fatalf("unexpected format: %q", got)
	}
}
