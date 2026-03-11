package review

import (
	"bytes"
	"strings"
)

// DimensionID identifies one of the four evaluation dimensions (O005/R007, T5.6).
type DimensionID string

const (
	DimSignalState    DimensionID = "signal_state"
	DimIteration      DimensionID = "iteration"
	DimScopeConverge  DimensionID = "scope_convergence"
	DimSubjectiveDone DimensionID = "subjective_completion"
)

// DimensionResult holds feedback for one dimension.
type DimensionResult struct {
	ID       DimensionID
	Label    string
	Feedback string
	OK       bool
}

// Embedded reviewer instructions (Ralph-owned); not user-editable.
// Evaluates prompt text against the four dimensions and returns structured feedback.
func evaluateDimensions(promptContent []byte) []DimensionResult {
	text := string(promptContent)
	lower := strings.ToLower(text)
	lines := strings.Split(text, "\n")

	var out []DimensionResult

	// 1. Signal and state: success/failure signals, statefulness compatible with loop
	signalFeedback, signalOK := evalSignalState(lower, text)
	out = append(out, DimensionResult{
		ID:       DimSignalState,
		Label:    "Signal and state",
		Feedback: signalFeedback,
		OK:       signalOK,
	})

	// 2. Iteration awareness: multi-iteration, fresh process each time
	iterFeedback, iterOK := evalIterationAwareness(lower, text)
	out = append(out, DimensionResult{
		ID:       DimIteration,
		Label:    "Iteration awareness",
		Feedback: iterFeedback,
		OK:       iterOK,
	})

	// 3. Scope and convergence: defined scope, checkable completion criteria
	scopeFeedback, scopeOK := evalScopeConvergence(lower, text, lines)
	out = append(out, DimensionResult{
		ID:       DimScopeConverge,
		Label:    "Scope and convergence",
		Feedback: scopeFeedback,
		OK:       scopeOK,
	})

	// 4. Subjective completion: escape techniques when "done" is subjective
	subjFeedback, subjOK := evalSubjectiveCompletion(lower, text)
	out = append(out, DimensionResult{
		ID:       DimSubjectiveDone,
		Label:    "Subjective completion criteria",
		Feedback: subjFeedback,
		OK:       subjOK,
	})

	return out
}

func evalSignalState(lower, text string) (feedback string, ok bool) {
	hasSuccess := hasAny(lower, "success", "exit 0", "exit code 0", "done", "complete", "succeed", "promise", "signal")
	hasFailure := hasAny(lower, "failure", "fail", "exit 1", "exit code 1", "error", "signal")
	if hasSuccess && hasFailure {
		return "Prompt defines or implies both success and failure signals. Compatible with Ralph's loop model.", true
	}
	if hasSuccess {
		return "Success signal or completion notion present; consider adding explicit failure signal so the loop can detect when to retry or exit.", false
	}
	if hasFailure {
		return "Failure signal present; consider defining success/done so the loop knows when to exit successfully.", false
	}
	return "Success and failure signals are missing or unclear. Add explicit markers (e.g. exit codes, SUCCESS/FAILURE lines) so Ralph can detect completion and failure.", false
}

func evalIterationAwareness(lower, text string) (feedback string, ok bool) {
	hasLoop := hasAny(lower, "loop", "iteration", "each run", "every run", "re-run", "rerun", "again")
	hasFresh := hasAny(lower, "fresh", "new process", "each time", "re-read", "reread", "state", "context")
	if hasLoop && hasFresh {
		return "Prompt acknowledges multi-iteration execution and re-reading state; suitable for Ralph's execution model.", true
	}
	if hasLoop {
		return "Loop or iteration mentioned; consider noting that each run is a fresh process so the AI re-reads state and emits signals each time.", false
	}
	if hasFresh {
		return "State or re-read mentioned; consider explicitly stating that execution is multi-iteration with a fresh process per run.", false
	}
	return "Iteration awareness missing or unclear. The prompt should acknowledge that execution is multi-iteration with a fresh process each time, so the AI can re-read state and emit signals.", false
}

func evalScopeConvergence(lower, text string, lines []string) (feedback string, ok bool) {
	hasScope := hasAny(lower, "scope", "task", "goal", "objective", "deliverable", "criteria", "check", "verify")
	hasConverge := hasAny(lower, "converge", "complete", "done", "finish", "max iteration", "max iterations", "limit")
	if hasScope && (hasConverge || len(lines) >= 3) {
		return "Scope or completion criteria are present; supports convergence of the loop.", true
	}
	if hasScope {
		return "Scope or task defined; consider adding checkable completion criteria or limits so the loop can converge.", false
	}
	if hasConverge {
		return "Completion or limits mentioned; consider defining a clear scope so the loop has a bounded task.", false
	}
	if len(text) < 100 {
		return "Very short prompt; scope and convergence criteria are unclear. Define task scope and how to verify completion.", false
	}
	return "Scope and convergence are missing or unclear. Define a bounded task and checkable completion criteria so the loop can converge.", false
}

func evalSubjectiveCompletion(lower, text string) (feedback string, ok bool) {
	hasSubjective := hasAny(lower, "good enough", "reads well", "sounds good", "looks good", "subjective", "judgment", "quality")
	hasEscape := hasAny(lower, "alternative", "variation", "step back", "challenge", "assumption", "explore", "try different", "another approach")
	if hasSubjective && hasEscape {
		return "Subjective completion is acknowledged with variation or stepping-back techniques; helps avoid getting stuck.", true
	}
	if hasSubjective {
		return "Subjective 'done' present; consider adding techniques to escape local optima (e.g. variation, step back, challenge assumptions) so the AI does not get stuck in small tweaks.", false
	}
	if hasEscape {
		return "Variation or stepping-back techniques present; supports convergence when completion is subjective.", true
	}
	return "No explicit subjective completion or escape techniques. If 'done' is subjective, add guidance for variation or stepping back to avoid repetitive tweaks.", true
}

func hasAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// narrativeFromDimensions formats dimension results as report narrative (R002, R007).
func narrativeFromDimensions(results []DimensionResult) string {
	var b bytes.Buffer
	for _, r := range results {
		b.WriteString("**")
		b.WriteString(r.Label)
		b.WriteString("**\n")
		b.WriteString(r.Feedback)
		b.WriteString("\n\n")
	}
	return strings.TrimSpace(b.String())
}

// suggestedRevisionFromDimensions returns a revised prompt that addresses dimension gaps.
// For T5.6 we add short inline suggestions as comments or append a "Review suggestions" block.
func suggestedRevisionFromDimensions(promptContent []byte, results []DimensionResult) string {
	base := string(promptContent)
	var additions []string
	for _, r := range results {
		if r.OK {
			continue
		}
		switch r.ID {
		case DimSignalState:
			additions = append(additions, "\n# Ralph: Define success (e.g. emit SUCCESS or exit 0) and failure (e.g. FAILURE or exit 1) so the loop can detect outcome.")
		case DimIteration:
			additions = append(additions, "\n# Ralph: This runs in a loop; each run is a fresh process. Re-read state and emit signals each time.")
		case DimScopeConverge:
			additions = append(additions, "\n# Ralph: Define scope and checkable completion criteria so the loop can converge.")
		case DimSubjectiveDone:
			additions = append(additions, "\n# Ralph: If 'done' is subjective, add variation or step-back techniques to avoid getting stuck.")
		}
	}
	if len(additions) == 0 {
		return base
	}
	// Append a single suggestions block to avoid duplicating per dimension
	return base + "\n\n---\n# Review suggestions (consider incorporating):\n" + strings.Join(additions, "\n")
}
