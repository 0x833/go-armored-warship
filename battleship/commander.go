package battleship

type Commander interface {
	Fire(x, y int)
	Observe(shot Coordinate)
}

type Player struct {
	Name  string
	Ships []Ship 
}
