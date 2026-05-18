package main

import (
	"math/rand"
)

// Step represents a single instruction for the user.
type Step struct {
	Action string // "type" or "delete"
	Digit  int    // only meaningful when Action == "type"
	IsReal bool   // true if this is a real passcode digit
}

// generatePasscode returns 4 random digits (0–9).
func generatePasscode() [4]int {
	var code [4]int
	for i := range code {
		code[i] = rand.Intn(10)
	}
	return code
}

// generateSequence builds a step sequence using a tree-walk approach.
// A recursive random walk decides at each node whether to insert a
// distractor pair (type-fake + delete) or advance to the next real digit.
// Total steps: 8–12. The 4th real digit is always the very last step.
func generateSequence(passcode [4]int) []Step {
	budget := rand.Intn(5) + 8 // 8–12 total steps
	return walkTree(passcode, 0, budget)
}

// walkTree recursively walks the decision tree.
//   - realIdx: index of the next real digit to place (0–3)
//   - budget: remaining step slots
//
// At each node we either:
//   - Branch A (distract): emit type-fake + delete (2 steps), stay at same realIdx
//   - Branch B (advance): emit type-real (1 step), move to realIdx+1
//
// Constraint: realIdx==3 → must emit the 4th digit immediately (no distractors,
// because macOS auto-submits the instant the 4th character is entered).
func walkTree(passcode [4]int, realIdx int, budget int) []Step {
	if realIdx == 3 {
		// Terminal node: type the last digit, done.
		return []Step{{Action: "type", Digit: passcode[3], IsReal: true}}
	}

	remainingReals := 4 - realIdx
	slack := budget - remainingReals // how many extra steps we can spend

	// Can we afford a distractor pair (2 steps)?
	if slack >= 2 && rand.Float64() < distractProb(slack, budget) {
		fake := nearMissDigit(passcode[realIdx])
		var steps []Step
		steps = append(steps, Step{Action: "type", Digit: fake, IsReal: false})
		steps = append(steps, Step{Action: "delete", Digit: 0, IsReal: false})
		steps = append(steps, walkTree(passcode, realIdx, budget-2)...)
		return steps
	}

	// Advance: type the real digit, recurse to next node
	var steps []Step
	steps = append(steps, Step{Action: "type", Digit: passcode[realIdx], IsReal: true})
	steps = append(steps, walkTree(passcode, realIdx+1, budget-1)...)
	return steps
}

// distractProb returns the probability of taking a distractor branch.
// More slack → higher chance. Decays as budget shrinks.
func distractProb(slack, budget int) float64 {
	if slack <= 1 {
		return 0
	}
	// Base ~55%, scaled by how much room we have
	return 0.55 * float64(slack) / float64(budget)
}

// nearMissDigit returns a distractor digit that's often close to the real
// digit (±1 or ±2), making sequences harder to distinguish from real ones.
func nearMissDigit(realDigit int) int {
	if rand.Float64() < 0.6 {
		offsets := []int{-2, -1, 1, 2}
		off := offsets[rand.Intn(len(offsets))]
		near := (realDigit + off + 10) % 10
		return near
	}
	return randomDigitExcluding(realDigit)
}

// randomDigitExcluding returns a random digit 0–9 that is not equal to exclude.
func randomDigitExcluding(exclude int) int {
	for {
		d := rand.Intn(10)
		if d != exclude {
			return d
		}
	}
}
