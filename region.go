package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

var ErrFailedGetRegion = errors.New("failed get region")

// Region is a region of the screen.
// 746,327 410x251
type Region struct {
	X, Y int
	W, H int
}

func (r *Region) String() string {
	return fmt.Sprintf("%d,%d %dx%d", r.X, r.Y, r.W, r.H)
}

// SelectRegion selects a region of the screen
func SelectRegion() (Region, error) {
	cmd := exec.Command("slurp")
	out, _ := cmd.CombinedOutput()

	if len(out) == 0 {
		return Region{}, ErrFailedGetRegion
	}

	// remove enters
	out = out[:len(out)-1]

	// is not a mistake
	if string(out) == "selection cancelled" {
		SendNotifyAndExit(NotifyNormal, "Selection cancelled", 0)
	}

	// explode by space
	slice := strings.Split(string(out), " ")
	// explode by comma
	xy := strings.Split(slice[0], ",")

	var (
		x, y, w, h int
		err        error
	)

	if x, err = strconv.Atoi(xy[0]); err != nil {
		return Region{}, ErrFailedGetRegion
	}

	if y, err = strconv.Atoi(xy[1]); err != nil {
		return Region{}, ErrFailedGetRegion
	}

	wh := strings.Split(slice[1], "x")

	if w, err = strconv.Atoi(wh[0]); err != nil {
		return Region{}, ErrFailedGetRegion
	}

	if h, err = strconv.Atoi(wh[1]); err != nil {
		return Region{}, ErrFailedGetRegion
	}

	return Region{
		X: x,
		Y: y,
		W: w,
		H: h,
	}, nil
}
