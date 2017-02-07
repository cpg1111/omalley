package elect

import (
	"github.com/pullrequestrfb/omalley/action"
)

type Dispatcher interface {
	DispatchVote(vote *action.Action) error
}
