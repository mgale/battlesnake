package main

import (
	"fmt"
	"math/rand"

	"github.com/beefsack/go-astar"
	"github.com/google/go-cmp/cmp"
)

// Board is 11, max row is 10, row counting starts at 0
// X axis is horizontal
// Y axis is vertical
// Coords start in bottom left, x increasing to the right
// Y increasing going up
// Direction is based on top view of board not based on
// snake direction.
type GameState struct {
	currentTarget Coord
	gameRequest   GameRequest
	strategy      Strategy
}

var gameStates map[string]GameState

func determineDirection(from Coord, to Coord) string {
	if from.X == to.X {
		// Y axis is changing, left or right
		if from.Y > to.Y {
			return "down"
		} else {
			return "up"
		}
	} else {
		// X axis is changing, up or down
		if from.X > to.X {
			return "left"
		} else {
			return "right"
		}
	}
}

func createWorld(gameRequest GameRequest) World {
	yMax := gameRequest.Board.Height
	xMax := gameRequest.Board.Width

	// Create World
	w := World{}
	for x := 0; x < xMax; x++ {
		for y := 0; y < yMax; y++ {
			w.SetTile(&Tile{
				Kind: KindPlain,
			}, x, y)
		}
	}

	for _, coord := range gameRequest.Board.Food {
		w.SetTile(&Tile{
			Kind: KindFood,
		}, coord.X, coord.Y)
	}

	for _, coord := range gameRequest.Board.Hazards {
		w.SetTile(&Tile{
			Kind: KindBlocker,
		}, coord.X, coord.Y)
	}

	for _, coord := range gameRequest.You.Body {
		w.SetTile(&Tile{
			Kind: KindBlocker,
		}, coord.X, coord.Y)
	}

	for _, snake := range gameRequest.Board.Snakes {
		for _, coord := range snake.Body {
			w.SetTile(&Tile{
				Kind: KindBlocker,
			}, coord.X, coord.Y)
		}
	}

	w.SetTile(&Tile{
		Kind: KindFrom,
	}, gameRequest.You.Head.X, gameRequest.You.Head.Y)

	return w
}

func determineMove(gameState GameState) (string, GameState) {
	// Choose a random direction to move in
	//possibleMoves := []string{"up", "down", "left", "right"}
	//move := possibleMoves[rand.Intn(len(possibleMoves))]

	fmt.Println("##################################")
	fmt.Printf("GS: Turn: %v, Head: %v, Health: %v, Body: %v\n",
		gameState.gameRequest.Turn,
		gameState.gameRequest.You.Head,
		gameState.gameRequest.You.Health,
		gameState.gameRequest.You.Body,
	)

	w := createWorld(gameState.gameRequest)

	var path []astar.Pather
	var distance float64
	var found bool
	goRandom := false
	for {
		var destinationCords Coord

		if goRandom {
			destinationCords = Coord{
				X: rand.Intn(gameState.gameRequest.Board.Width),
				Y: rand.Intn(gameState.gameRequest.Board.Height),
			}

		} else {
			goRandom = true
			if cmp.Equal(gameState.gameRequest.You.Head, gameState.currentTarget) {
				// Create new target
				destinationCords = Coord{
					X: gameState.gameRequest.Board.Width - gameState.gameRequest.You.Head.X - 1,
					Y: gameState.gameRequest.Board.Height - gameState.gameRequest.You.Head.Y - 1,
				}
			} else {
				destinationCords = gameState.currentTarget
			}
		}
		destinationTile := w.Tile(destinationCords.X, destinationCords.Y)
		if destinationTile.Kind != KindPlain {
			continue
		}

		w.SetTile(&Tile{
			Kind: KindTo,
		}, destinationCords.X, destinationCords.Y)

		// t1 and t2 are *Tile objects from inside the world.
		path, distance, found = astar.Path(w.From(), w.To())

		if found {
			gameState.currentTarget = destinationCords
			break
		}

		fmt.Println("Could not find path")
	}

	// for _, p := range path {
	// 	pT := p.(*Tile)

	// 	r, _ := KindRunes[pT.Kind]
	// 	fmt.Println("Tile Type:", string(r))
	// 	fmt.Printf("Tile: X:%d, Y:%d\n", pT.X, pT.Y)
	// }
	pT := path[len(path)-2].(*Tile)

	fmt.Println("Current target:", gameState.currentTarget)
	fmt.Println("Estimated distance to dest:", distance)
	fmt.Println("Tile Type:", string(KindRunes[pT.Kind]))
	fmt.Printf("New Coords: X:%d, Y:%d\n", pT.X, pT.Y)

	destCoords := Coord{
		X: pT.X,
		Y: pT.Y,
	}

	move := determineDirection(gameState.gameRequest.You.Head, destCoords)

	return move, gameState
}
