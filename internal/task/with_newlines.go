// Package task содержит реализации различных типов задач
package task

import (
	"github.com/qzeleza/ziva/internal/common"
)

// WithNewLinesInErrors реализация для SingleSelectTask
func (t *SingleSelectTask) WithNewLinesInErrors(preserve bool) common.Task {
	t.preserveErrorNewLines = preserve
	return t
}

// WithNewLinesInErrors реализация для MultiSelectTask
func (t *MultiSelectTask) WithNewLinesInErrors(preserve bool) common.Task {
	t.preserveErrorNewLines = preserve
	return t
}

// WithNewLinesInErrors реализация для YesNoTask
func (t *YesNoTask) WithNewLinesInErrors(preserve bool) common.Task {
	t.preserveErrorNewLines = preserve
	return t
}

// WithNewLinesInErrors реализация для InputTaskNew
func (t *InputTaskNew) WithNewLinesInErrors(preserve bool) common.Task {
	t.preserveErrorNewLines = preserve
	return t
}

// WithNewLinesInErrors реализация для FuncTask
func (t *FuncTask) WithNewLinesInErrors(preserve bool) common.Task {
	t.preserveErrorNewLines = preserve
	return t
}
