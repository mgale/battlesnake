package main

import (
	"log"

	"github.com/beefsack/go-astar"
	"github.com/google/go-cmp/cmp"
)

/*


 */

type Solo struct {
	path              []astar.Pather
	distance          float64
	found             bool
	firstDestination  Coord
	lastDestination   Coord
	rowLimiter        int
	smallBodySize     int
	ladderStage       int
	xMax              int
	yMax              int
	safeTrack         []Coord
	safeTrackPosition int
}

func generateSquare(xMax, yMax int) []Coord {
	return []Coord{
		{X: 0, Y: 0},
		{X: xMax, Y: 0},
		{X: xMax, Y: yMax},
		{X: 0, Y: yMax},
	}
}

func generateLadder(xMax, yMax int) []Coord {
	ladderCoords := []Coord{}

	xBuffer := xMax - 1
	xflip := 0
	for y := yMax; y > 0; y-- {
		ladderCoords = append(ladderCoords, Coord{
			X: xflip,
			Y: y,
		})
		xflip = xBuffer - xflip
	}

	ladderCoords = append(ladderCoords, []Coord{
		{
			X: xMax,
			Y: 1,
		},
		{
			X: xMax,
			Y: yMax / 2,
		},
		{
			X: xMax,
			Y: yMax,
		}}...,
	)

	return ladderCoords
}

func (s *Solo) Initialize(gameRequest GameRequest) {
	log.Println("Initialzing Solo strategy")
	log.Printf("Board Size: X:%d, Y:%d\n",
		gameRequest.Board.Width,
		gameRequest.Board.Height,
	)
	s.firstDestination.X = gameRequest.Board.Width / 2
	s.firstDestination.Y = gameRequest.Board.Height - 1
	s.lastDestination.X = gameRequest.You.Head.X
	s.lastDestination.Y = gameRequest.You.Head.Y

	s.xMax = gameRequest.Board.Width - 1
	s.yMax = gameRequest.Board.Height - 1

	s.safeTrackPosition = 0
	s.safeTrack = generateLadder(s.xMax, s.yMax)
}

func (s *Solo) GetMove(gameRequest GameRequest) string {
	w := createWorld(gameRequest)

	loopCounter := 0
	loopMax := 10
	for {
		loopCounter++
		if loopCounter > loopMax {
			break
		}
		var tmpDest Coord
		if s.safeTrackPosition >= len(s.safeTrack) {
			s.safeTrackPosition = 0
		}

		if cmp.Equal(gameRequest.You.Head, s.safeTrack[s.safeTrackPosition]) {
			s.safeTrackPosition++
			if s.safeTrackPosition >= len(s.safeTrack) {
				s.safeTrackPosition = 0
			}
		}

		tmpDest = s.safeTrack[s.safeTrackPosition]

		tmpTile := w.Tile(tmpDest.X, tmpDest.Y)
		if tmpTile.Kind == KindBlocker {
			foodTiles := []*Tile{}
			plainTiles := []*Tile{}
			log.Println("Error - Can't find path, dest failed:", tmpTile)
			directPaths := w.From().PathNeighbors()
			for _, path := range directPaths {
				checkTile := path.(*Tile)
				if checkTile.Kind == KindFood {
					foodTiles = append(foodTiles, checkTile)
				}
				if checkTile.Kind == KindPlain {
					plainTiles = append(plainTiles, checkTile)
				}
			}
			if len(foodTiles) > 0 {
				tmpDest.X = foodTiles[0].X
				tmpDest.Y = foodTiles[0].Y
			} else {
				tmpDest.X = plainTiles[0].X
				tmpDest.Y = plainTiles[0].Y
			}
		}

		w.SetTile(&Tile{Kind: KindTo}, tmpDest.X, tmpDest.Y)
		s.path, s.distance, s.found = astar.Path(w.From(), w.To())

		if s.found {
			break
		}

		//No found path, marking current dest tile as a blocker
		w.SetTile(&Tile{Kind: KindBlocker}, tmpDest.X, tmpDest.Y)
		s.safeTrackPosition++
	}

	/* Path is reverse order,
	postition 0 is the destination
	last position is our starting point or the snake head */

	pT := s.path[len(s.path)-2].(*Tile)
	destCoords := Coord{
		X: pT.X,
		Y: pT.Y,
	}

	move := determineDirection(gameRequest.You.Head, destCoords)

	log.Println("###########################")
	log.Printf("Current target: %v, Type: %v", s.lastDestination, string(KindRunes[pT.Kind]))
	log.Println("Estimated distance to dest:", s.distance)
	log.Printf("Head Coords: X:%d, Y:%d, New Coords: X:%d, Y:%d Move: %s MoveNum: %d\n",
		gameRequest.You.Head.X,
		gameRequest.You.Head.Y,
		pT.X,
		pT.Y,
		move,
		gameRequest.Turn,
	)
	return move

}
