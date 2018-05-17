package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"math"
	"time"
	"bufio"
	"os"
	"strings"
)

/*
he players use red and black checkers. Red goes first. Players alternate turns until one player has four checkers in a row horizontally, vertically, or diagonally.
It is possible for the game to end in a draw if no player can achieve four in a row.

The board has 6 rows and 7 columns per row. To make the game easier for human users number the columns 1 to 7.
The actual columns in your 2d array are numbered 0 to 6 so you will have to adjust the user input.
Error check the input to ensure it is an int (that method is already provided), that the selected column is a valid column,
and that there is an open spot in that column. A checker cannot be placed in a column that is full.
*/
const (
	// height is actually the Y coordinate but it goes first in array access - since we look at rows as vertical, we put that first
	width  = 8 // we will ignore 0th column
	height = 6
)

const (
	HUMAN    = "h"
	COMPUTER = "c"
)

func main() {
	fmt.Println("Welcome to connect4")
	runGame()
}

var (
	reader = bufio.NewReader(os.Stdin)

)

func runGame() {
	// board starts out empty
	var (
		board     [height][width]string
		finished  bool
		curPlayer = HUMAN
		round = 1
	)
	rand.Seed(int64(time.Now().Second()))
	numSpotsFilled :=0
	for ; !finished; {
		printBoard(&board, round)
		curMove := getMove(curPlayer)
		// checking that move is within bounds - does not apply to computer player
		fmt.Printf("next move for player [%v] is column [%s] in round %d\n", curPlayer, curMove, round)

		lostAMove, height, width := fillBoardWithMove(&board, curPlayer, curMove, round)
		if lostAMove {
			fmt.Printf("Player [%s] chose a full column [%v] and coin rolled off\n", curPlayer, curMove)
			// flip game order
			curPlayer = flipGameOrder(curPlayer)
			continue
		}
		numSpotsFilled++
		finished, isDraw, desc := checkForGameEnd(board, curPlayer, height, width, numSpotsFilled)
		if finished {
			if isDraw {
				fmt.Printf("game is a draw\n")
				printBoard(&board, round)
				return
			}
			fmt.Printf("**** player [%s] WON via [%s]!\n", curPlayer, desc)
			fmt.Println("-------------------------------")
			printBoard(&board, round)
			return;
		}

		// flip game order
		curPlayer = flipGameOrder(curPlayer)
		round++
	}
}

func flipGameOrder(curPlayer string) string {
	// flip game order
	if curPlayer == HUMAN {
		return COMPUTER
	} else {
		return HUMAN
	}
}

// passing in x.y (ie height/width) of newly placed coin
// returns isFinished, isDraw
func checkForGameEnd(board [height][width]string, curPlayer string, inHeight int, inWidth int, numSpotsFilled int) (bool, bool, string) {
	isFinished := checkVertical(board, curPlayer, inHeight, inWidth)
	if isFinished {
		return true, false, "vertical"
	}

	isFinished = checkHorizontal(&board, curPlayer, inHeight, inWidth)
	if isFinished {
		return isFinished, false, "Horizontal"
	}

	// need to check diagonals
	isFinished = checkLeftDiagonal(&board, curPlayer, inHeight, inWidth)
	if isFinished {
		return isFinished, false, "left diagonal"
	}

	isFinished = checkRightDiagonal(&board, curPlayer, inHeight, inWidth)
	if isFinished {
		return isFinished, false, "rigth diagonal"
	}


	// if board is full and nothing is a conncect4 we have a draw
	if numSpotsFilled == (width-1)*height {
		fmt.Printf("have %d spots filled, detecting a draw\n", numSpotsFilled)
		return true, true, "DRAW"

	}
	return false, false, "kep playing"
}

// go down and to the left 3 and see if you have same player
func checkLeftDiagonal(board *[height][width]string, curPlayer string, inHeight int, inCol int) bool {
	fmt.Printf("Checking left diagonal around {%v, %v}\n", inHeight, inCol)
	numTimesSamePlayer := 0
	newHeight := inHeight
	newCol := inCol
	for i:=1; i<=3; i++ {
		newHeight = newHeight+1
		if newHeight >= height {
			return false
		}

		newCol = newCol -1
		if newCol <1 { // want to cap at 1 since we are skipping 0th column
			return false
		}

		fmt.Printf("in loop: checking left diagonal around {%v, %v}\n", newHeight, newCol)
		if board[newHeight][newCol] == curPlayer {
			numTimesSamePlayer++
		} else {
			return false
		}
	}
	return numTimesSamePlayer==3
}

// go down and to the right 3 and see if you have same player
func checkRightDiagonal(board *[height][width]string, curPlayer string, inHeight int, inCol int) bool {
	fmt.Printf("Checking right diagonal around {%v, %v}\n", inHeight, inCol)
	numTimesSamePlayer := 0
	newHeight := inHeight
	newCol := inCol
	for i:=0; i<=3; i++ {
		newHeight = newHeight+1
		if newHeight >= height {
			return false
		}

		newCol = newCol+1
		if newCol >=width {
			return false
		}

		if board[newHeight][newCol] == curPlayer {
			numTimesSamePlayer++
		} else {
			return false
		}
	}
	return numTimesSamePlayer==3
}


// need to check horizontal, all combos where incoming width column may be anywhere in Max(0, width) and Min(7, width)
// so we go from -3...width...+3
// ie board[height][Max(0, width-3)]...board[height, width]...board[height][Math.Min(7, width+3)]
// height stays constant, and we vary the column
func checkHorizontal(board *[height][width]string, curPlayer string, height int, width int) bool {
	fmt.Printf("Checking horizontal from position {%d,%d} for %v\n", height, width, curPlayer)
	sameInARow := 0
	row := board[height]
	fmt.Printf("horizontal check analyzes row %v in row [%v] from {%v to %v}\n",
		row, height, int(math.Max(0, float64(width-3))), int(math.Min(7, float64(width+3))))
	for i:=int(math.Max(0, float64(width-3))); i<= int(math.Min(7, float64(width+3))); i++ {
		//fmt.Printf("checking for i %v\n", int(i))
		if row[i] == curPlayer {
			sameInARow++
		} else {
			sameInARow = 0
		}
		if sameInARow >=4 {
			return true;
		}
	}
	return sameInARow >= 4
}

// checks to see if this new coin has 3 of same coin underneath
// so need to go down 3 from height
func checkVertical(board [height][width]string, curPlayer string, inHeightOfCoin int, width int) bool {
	fmt.Printf("checking vertical for column %d for 3 down from height %d\n", width, inHeightOfCoin)
	// check vertical 3 down
	for i:=1; i<=3; i++ {
		if inHeightOfCoin+ i >= height || board[inHeightOfCoin+i][width] != curPlayer {
			return false;
		}
	}
	return true
}

// keep falling until we hit something. if we had an empty spot before, go to i-1 and set it to curPlayer
func fillBoardWithMove(board *[height][width]string, curPlayer string, column string, round int) (bool, int, int) {
	theCol, _ := strconv.Atoi(column)
	fmt.Printf("Player [%v] is dropping into column [%d]\n", curPlayer, theCol)
	emptyAtHeight := -1
	for i:=0; i<height; i++ {
		//fmt.Printf("Current value in {%d, %d} is [%v]\n", i, theCol, board[i][theCol])
		if len(board[i][theCol]) == 0 { // if spot is empty, keep falling to next one
			emptyAtHeight = i
		} else {
			break
		}
	}

	if emptyAtHeight != -1 {
		board[emptyAtHeight][theCol] = curPlayer
		fmt.Printf("board[%d][%d] becomes %s\n", emptyAtHeight, theCol, curPlayer)
		return false, emptyAtHeight, theCol
	}

	fmt.Printf("column [%d] is full and nothing got placed, player [%s] lost a move\n", theCol, curPlayer)
	printBoard(board, round)
	return true, -1, theCol // height/width doesn't matter for false result
}

// returns the column being chosen
func getMove(player string) string {
	var column string
	if player == HUMAN {
		for
		{
			fmt.Print("Enter column: ")
			column, _ = reader.ReadString('\n')
			column = strings.Trim(column, "\n")
			//fmt.Printf("read [%v]\n", column)
			colInt, err := strconv.Atoi(column)
			if err != nil {
				fmt.Printf("You entered an invalid colInt, please try again: %v\n", err)
				continue
			}
			if colInt < 1 || colInt >= width {
				fmt.Printf("%d is not a valid colInt, please pick one [1,7]\n", colInt)
				continue
			}
			// valid colInt, break
			break
		}
	} else {
		column = strconv.Itoa(rand.Intn(width-1)+1) // it's exclusive, [0,width)
	}
	fmt.Printf("Player [%v] chooses column [%s]\n", player, column)
	return column
}

/*
arrays in Go are wide x down, and are considered down-dimensional
ie a 3x2 array is 3 down, 2 wide and is [3][2]int
 */

func printBoard(board *[height][width]string, round int) {
	fmt.Printf("Current board for round %d\n", round)
	fmt.Println("    1 2 3 4 5 6 7")
	fmt.Println("    -------------")
	for i, _ := range board {
		fmt.Printf("%d | ", i)
		for j:=1; j<len(board[i]); j++ {
			if len(board[i][j]) > 0 {
				fmt.Printf("%s ", board[i][j])
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}
