package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Card struct {
	value int
	suit  int // 0 - spades, 1 - hearts, 2 - diamonds, 3 - clubs
}

func (card Card) getString() string {
	var suit string
	var value string

	switch card.suit {
	case 0:
		suit = "♠"
	case 1:
		suit = "♥"
	case 2:
		suit = "♦"
	case 3:
		suit = "♣"
	}

	switch card.value {
	case 11:
		value = "J"
	case 12:
		value = "Q"
	case 13:
		value = "K"
	case 1:
		value = "A"
	default:
		value = fmt.Sprintf("%d", card.value)
	}

	return value + suit
}

type Deck struct {
	cards []Card
}

// func (d *Deck) deal(num uint) []Card {

// }

func (game *Game) checkStatus(user string) string {
	if user == "Player" {
		var total int
		for i := 0; i < len(game.playerCards); i++ {
			total += int(game.playerCards[i].value)
		}
		if total > 21 {
			return "exit"
		} else if total == 21 {
			return "Winner"
		} else {
			return "S"
		}
	} else {
		var total int
		for i := 0; i < len(game.dealerCards); i++ {
			total += int(game.dealerCards[i].value)
		}
		if total > 21 {
			return "exit"
		} else if total < 17 {
			return "H"
		}
		return "S"
	}

}

func (d *Deck) create() {
	for i := 0; i < 4; i++ {
		for j := 1; j < 14; j++ {
			card := Card{}
			card.value = j
			card.suit = i
			d.cards = append(d.cards, card)
		}
	}
}

func (d *Deck) shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

type Game struct {
	deck        Deck
	playerCards []Card
	dealerCards []Card
}

func (game *Game) dealStartingCards() {
	card1 := game.deck.cards[0].getString()
	card2 := game.deck.cards[1].getString()
	game.playerCards = append(game.playerCards, game.deck.cards[0:2]...)

	total := game.deck.cards[0].value + game.deck.cards[1].value
	fmt.Printf("Your cards: %s + %s = %d \n", card1, card2, total)

	card3 := game.deck.cards[2].getString()
	game.dealerCards = append(game.dealerCards, game.deck.cards[2:4]...)
	fmt.Printf("Dealer card: %s\n", card3)

	game.deck.cards = game.deck.cards[4:]
}

func (game *Game) showCards(user string) {
	var total int
	fmt.Printf("Cards of %s :", user)
	if user == "Dealer" {
		for i := 0; i < len(game.dealerCards); i++ {
			card := game.dealerCards[i].getString()
			total += int(game.dealerCards[i].value)
			fmt.Printf(card + " + ")
		}
	} else if user == "Player" {
		for i := 0; i < len(game.playerCards); i++ {
			card := game.playerCards[i].getString()
			total += int(game.playerCards[i].value)
			fmt.Printf(card + " + ")
		}
	}
	fmt.Printf(": %d\n", total)
}

func (game *Game) play(bet float64) float64 {

	game.deck.create()
	game.deck.shuffle()
	game.dealStartingCards()

	var Poption, Doption string
	Poption = game.checkStatus("Player")
	Doption = game.checkStatus("Dealer")
	if Poption == "Winner" {
		fmt.Println("You won the bet")
		return bet
	} else if Poption == "exit" {
		fmt.Println("You are busted!!!.")
		return -bet
	} else if Doption == "Winner" {
		fmt.Println("Dealer Won")
		game.showCards("Dealer")
		return -bet
	} else if Doption == "exit" {
		fmt.Println("Dealer busted. You Won!!!")
		game.showCards("Dealer")
		return bet
	}
	Poption = "H"
	Doption = "H"
	option := "H"
	for option == "H" || Doption == "H" {
		fmt.Printf("Choose Hit (H) or Stay (S): ")
		option = enterString()
		if option == "H" {
			game.playerTurn()
			Poption = game.checkStatus("Player")
		}
		if Poption == "Winner" {
			fmt.Println("You won the bet")
			return bet
		} else if Poption == "exit" {
			fmt.Println("You are busted!!!.")
			return -bet
		}

		Doption = game.checkStatus("Dealer")
		if Doption == "H" {
			game.dealerTurn()
			Doption = game.checkStatus("Dealer")
		} else if Doption == "S" {
			fmt.Println("Dealer Stayed")
		}
		if Doption == "Winner" {
			fmt.Println("Dealer Won")
			game.showCards("Dealer")
			return -bet
		} else if Doption == "exit" {
			fmt.Println("Dealer busted. You Won!!!")
			game.showCards("Dealer")
			return bet
		}
	}
	var totalD int
	for i := 0; i < len(game.dealerCards); i++ {
		totalD += int(game.dealerCards[i].value)
	}
	var totalP int
	for i := 0; i < len(game.dealerCards); i++ {
		totalP += int(game.playerCards[i].value)
	}
	if (21 - totalD) < (21 - totalP) {
		fmt.Println("Winner is Dealer")
		game.showCards("Dealer")
		return -bet
	} else if (21 - totalD) > (21 - totalP) {
		fmt.Println("You are Winner")
		game.showCards("Player")
		return bet
	} else {
		fmt.Println("Draw")
		game.showCards("Player")
		game.showCards("Dealer")
		return 0
	}

}

func (game *Game) playerTurn() {

	game.playerCards = append(game.playerCards, game.deck.cards[0])
	game.deck.cards = game.deck.cards[1:]
	game.showCards("Player")
}

func (game *Game) dealerTurn() {
	game.dealerCards = append(game.dealerCards, game.deck.cards[0])
	card := game.deck.cards[0].getString()
	game.deck.cards = game.deck.cards[1:]
	fmt.Printf("Dealer hits and received: %s\n", card)
}

func enterString() string {
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occurred while reading input. Please try again", err)
		return ""
	}

	// remove the delimiter from the string
	input = strings.TrimSuffix(input, "\n")
	return input
}

func main() {
	balance := float64(100)

	for balance > 0 {
		fmt.Printf("Your balance is: $%.2f\n", balance)
		fmt.Printf("Enter your bet (q to quit): ")
		bet, err := strconv.ParseFloat(enterString(), 64)
		if err != nil {
			break
		}
		if bet > balance || bet <= 0 {
			fmt.Println("Invalid bet.")
			continue
		}

		game := Game{}
		balance += game.play(bet)
	}

	fmt.Printf("You left with: $%2.f\n", balance)
}
