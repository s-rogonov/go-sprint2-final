package handlers

import (
	"time"

	"dbprovider"
	"dbprovider/models"
	renderPkg "github.com/unrolled/render"
)

var Render *renderPkg.Render

func init() {
	Render = renderPkg.New() // pass options if you want
}

func AccessDurations() (map[string]time.Duration, error) {
	tm := &models.Timings{}
	err := dbprovider.Manager.GetTimings(tm)
	if err != nil {
		return nil, err
	}

	return map[string]time.Duration{
		"+": tm.Addition,
		"-": tm.Subtraction,
		"*": tm.Multiplication,
		"/": tm.Division,
	}, nil
}
