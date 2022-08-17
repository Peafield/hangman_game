package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"unicode"
)

func printHangedMan(i int) {
	hangedMan := []string{
		`-----------
|     |
|     
|    
|    
|
-----------`,
		`-----------
|     |
|     O
|    
|    
|
-----------`,
		`-----------
|     |
|     O
|     |
|    
|
-----------`,
		`-----------
|     |
|     O
|    /|
|    
|
-----------`,
		`-----------
|     |
|     O
|    /|\
|    
|
-----------`,
		`-----------
|     |
|     O
|    /|\
|    / 
|
-----------`,
		`-----------
|     |
|     O
|    /|\
|    / \
|
-----------`,
	}
	fmt.Println(hangedMan[i])
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getWord() []string {
	resp, err := http.Get("http://random-word-api.herokuapp.com/word")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	var wordString []string
	for _, ch := range body {
		if ch != 91 && ch != 34 && ch != 93 {
			wordString = append(wordString, string(ch))
		}
	}
	return wordString

}

func createUnderScores(l int) []string {
	var underScores []string
	i := 0
	for i < l {
		underScores = append(underScores, "_")
		i++
	}
	return underScores
}

func printBoardState(u []string) {
	for i, ch := range u {
		if i != len(u)-1 {
			fmt.Print(ch + " | ")
		} else {
			fmt.Print(ch)
		}
	}
	fmt.Println()
}

func guessALetter() string {
	fmt.Println("Guess a letter: ")
	var guessedLetter string
	fmt.Scanln(&guessedLetter)
	for _, r := range guessedLetter {
		if !unicode.IsLetter(r) {
			fmt.Println("That is not a valid guess!")
			guessedLetter = guessALetter()
		}
	}
	return strings.ToLower(guessedLetter)
}

func isUsedLetter(l string, u []string) bool {
	delim := true
	for _, ch := range u {
		if l != ch {
			delim = false
		}
	}
	return delim
}

func isCorrectGuess(g string, ws []string, u []string) ([]string, bool) {
	delim := false
	for i, ch := range ws {
		if g == ch {
			u[i] = g
			delim = true
		}
	}
	return u, delim
}

func gameOver(u []string) bool {
	game := true
	for _, ch := range u {
		if ch == "_" {
			game = false
		}
	}
	return game
}

func playAgain() {
	fmt.Println("Would you like to play again? Y or N:")
	var playAgain string
	fmt.Scanln(&playAgain)
	if strings.ToUpper(playAgain) == "Y" {
		clearScreen()
		play()
	} else {
		os.Exit(0)
	}
}

func play() {
	i := 0
	wordString := getWord()
	var valid bool
	var guessedLetter string
	var alreadyGuessed []string
	underScores := createUnderScores(len(wordString))
	for !gameOver(underScores) {
		clearScreen()
		if i == 6 {
			printHangedMan(i)
			fmt.Println("You Lost!")
			fmt.Println("The Word Was:")
			printBoardState(wordString)
			playAgain()
		}
		fmt.Println("Used Letters: ", alreadyGuessed[:])
		printHangedMan(i)
		printBoardState(underScores)
		guessedLetter = guessALetter()
		// TO DO: Fix this. Why only happens on first go? Could reuse isCorrectGuess here...?
		if len(alreadyGuessed) > 0 {
			for isUsedLetter(guessedLetter, alreadyGuessed) {
				fmt.Println("Already Used!")
				guessedLetter = guessALetter()
			}
		}
		alreadyGuessed = append(alreadyGuessed, guessedLetter)
		underScores, valid = isCorrectGuess(guessedLetter, wordString, underScores)
		if !valid {
			fmt.Println("Incorrect!")
			i++
		}
	}
	if gameOver(underScores) {
		clearScreen()
		printHangedMan(i)
		printBoardState(underScores)
		fmt.Println("You Won!")
		playAgain()
	}
}

func main() {
	play()
}

// TO DO: Add detection for already chosen letters
// TO DO: If lost, print out answer
// TO DO: Create game struct
