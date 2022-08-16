package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

func main() {
	i := 0
	wordString := getWord()
	fmt.Println(wordString)
	underScores := createUnderScores(len(wordString))
	printHangedMan(i)
	printBoardState(underScores)
	guessedLetter := guessALetter()
	for !gameOver(underScores) {
		underScores, valid := isCorrectGuess(guessedLetter, wordString, underScores)
		if !valid {
			if i == 7 {
				gameOver(underScores)
				fmt.Println("You Lost!")
				os.Exit(0)
			}
			fmt.Println("Incorrect!")
			i++
			printHangedMan(i)
			printBoardState(underScores)
			guessedLetter = guessALetter()
			underScores, valid = isCorrectGuess(guessedLetter, wordString, underScores)

		} else {
			fmt.Println(i)
			printHangedMan(i)
			printBoardState(underScores)
			guessedLetter = guessALetter()
			underScores, valid = isCorrectGuess(guessedLetter, wordString, underScores)
		}
	}
	if gameOver(underScores) {
		printHangedMan(i)
		printBoardState(underScores)
		fmt.Println("You Won!")
	}
}

// TO DO: Add detection for already chosen letters
