package controller

import (
	"github.com/tfriedel6/canvas"
)

type Suggestion struct {
	Message string
	Icon    string
}

func (s *Suggestion) Render(cv *canvas.Canvas) {

}

func GetHelperSuggestions(cp *ControlPanel) *Suggestion {
	if len(cp.C.Map.Countries[0].Towns[0].Farms) == 0 {

	} else if len(cp.C.Map.Countries[0].Towns[0].Workshops) == 0 {

	}
	return nil
}
