package bbsteps

import (
	"braid-game/bot/bbprefab"
	"braid-game/bot/bbprefab/bbcards"

	"github.com/pojol/gobot/prefab"
)

// NewGuestLoginStep guest
func NewGuestLoginStep(md *bbprefab.BotData) *prefab.Step {

	step := prefab.NewStep()

	step.AddCard(bbcards.NewGuestLoginCard(md))

	return step
}

// NewLoginOutStep login out
func NewLoginOutStep(md *bbprefab.BotData) *prefab.Step {
	step := prefab.NewStep()

	step.AddCard(bbcards.NewLoginOutCard(md))

	return step
}
