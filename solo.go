package main

import (
	"log"
	"math/rand"

	"github.com/beefsack/go-astar"
	"github.com/google/go-cmp/cmp"
)

type Solo struct {
	path             []astar.Pather
	distance         float64
	found            bool
	firstDestination Coord
	lastDestination  Coord
	rowLimiter       int
	smallBodySize    int
	ladderStage      int
	xMax             int
	yMax             int
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
	s.rowLimiter = gameRequest.Board.Width - 2
	s.smallBodySize = gameRequest.Board.Width - 2
	s.ladderStage = 0
	s.xMax = gameRequest.Board.Width - 1
	s.yMax = gameRequest.Board.Height - 1
}

// We need a new destination and food does not exist or we are too large
func (s *Solo) calculateDest(gameRequest GameRequest) Coord {

	//Snake placement completely unknown.
	if s.ladderStage == 0 {
		s.ladderStage = 1
		return Coord{
			X: s.xMax,
			Y: rand.Intn(s.yMax),
		}
	}

	if s.ladderStage == 1 {
		if gameRequest.You.Head.X == s.xMax {
			// We made it to the mountains
			s.ladderStage = 2
			return s.firstDestination
		} else {
			return Coord{
				X: s.xMax,
				Y: rand.Intn(s.yMax),
			}
		}
	}

	if s.ladderStage == 2 {
		//If we are near the bottom of the board go to our firstDestination
		if gameRequest.You.Head.Y <= 1 {
			s.ladderStage = 1
			return s.firstDestination
		}
		return Coord{
			X: s.xMax - gameRequest.You.Head.X,
			Y: gameRequest.You.Head.Y - 1,
		}
	}

	log.Println("Error - Unknown condition, generating random move")
	return Coord{
		X: rand.Intn(s.xMax),
		Y: rand.Intn(s.yMax),
	}
}

func (s *Solo) GetMove(gameRequest GameRequest) string {
	yMax := gameRequest.Board.Height
	xMax := gameRequest.Board.Width

	w := createWorld(gameRequest)
	for x := s.rowLimiter; x < xMax; x++ {
		for y := 0; y < yMax; y++ {
			w.SetTile(&Tile{
				Kind: KindMountain,
			}, x, y)
		}
	}

	for {
		var tmpDest Coord

		if cmp.Equal(gameRequest.You.Head, s.lastDestination) {
			//We have reached our target, goal, we need a new destination
			//Wiping current dest
			s.lastDestination = Coord{}
		}

		if (Coord{}) == s.lastDestination {
			// We need a new destination
			// Are we small and is there food
			tmpTile := w.FirstOfKind(KindRiver)
			if len(gameRequest.You.Body) < s.rowLimiter && tmpTile != nil {
				log.Println("Food Hunting")
				tmpDest.X = tmpTile.X
				tmpDest.Y = tmpTile.Y
			} else {
				log.Println("General Move")
				tmpDest = s.calculateDest(gameRequest)
			}
		}

		// Before code is path error handling
		// tmpDestTile := w.Tile(tmpDest.X, tmpDest.Y)
		// if tmpDestTile.Kind != KindPlain {
		// 	//If the destination is not safe loop again
		// 	continue
		// }

		w.SetTile(&Tile{Kind: KindTo}, tmpDest.X, tmpDest.Y)
		s.path, s.distance, s.found = astar.Path(w.From(), w.To())

		if s.found {
			s.lastDestination = tmpDest
			break
		}

		//No found path, marking current dest tile as a blocker
		w.SetTile(&Tile{Kind: KindBlocker}, tmpDest.X, tmpDest.Y)
		s.lastDestination = Coord{}
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
	log.Printf("Head Coords: X:%d, Y:%d, New Coords: X:%d, Y:%d Move: %s\n",
		gameRequest.You.Head.X,
		gameRequest.You.Head.Y,
		pT.X,
		pT.Y,
		move,
	)
	return move

}
