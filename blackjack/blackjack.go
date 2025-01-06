package blackjack

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Junior-Green/gophercises/deck"
)

type state uint8

const (
	dealing state = iota
	playerTurn
	dealerTurn
	end
)

type AI interface {
	DoubleDown(Hand) bool
	DecideHit(Hand) bool
	DecideBet() int
}

type BasicDealerStrategy struct{}

func (s BasicDealerStrategy) DecideHit(hand Hand) bool {
	return false
}

type StandardDealerStrategy struct{}

func (s StandardDealerStrategy) DecideHit(hand Hand) bool {
	handVal := hand.Value()

	if handVal < 17 {
		return true
	}

	for _, c := range hand.Hand {
		if c.Type == deck.ACE {
			return handVal == 17
		}
	}
	return false
}

type DealerStrategy interface {
	DecideHit(hand Hand) bool
}

type Hand struct {
	Hand []deck.Card
}

func (h *Hand) addCard(c deck.Card) {
	h.Hand = append(h.Hand, c)
}

func (h *Hand) Value() int {
	var aces, points int

	for _, card := range h.Hand {
		if card.Type == deck.ACE {
			points += 11
			aces++
		} else {
			points += cardValue(card)
		}
	}

	for i := 0; i < aces; i++ {
		if points <= 21 {
			return points
		}
		points -= 10
	}

	return points
}

type dealer struct {
	hand     Hand
	strategy DealerStrategy
}

func (d *dealer) play(g *game) {
	if !g.isSimulation {
		fmt.Println("\nDealer's turn")
		fmt.Println("-------------------------------")
		fmt.Println("\nDealer flips card.")
		d.printHand(g)
	}

	for d.strategy.DecideHit(d.hand) {
		d.draw(g)

		handVal := d.hand.Value()
		if handVal == 21 {
			fmt.Println("Dealer gets a Blackjack!")
			return
		} else if handVal > 21 {
			fmt.Println("Dealer busts!")
			return
		}
	}
	fmt.Println("\nDealer stands.")
}

func (d *dealer) draw(g *game) {
	card := g.drawCard()
	d.hand.addCard(card)

	if g.state == dealing && len(d.hand.Hand) == 2 {
		fmt.Println("Dealer draws a card face down")
		return
	}

	fmt.Printf("Dealer draws a %s\n", card)
}

func (d *dealer) printHand(g *game) {
	fmt.Println("\nDealer's Hand:")
	fmt.Printf("- %s\n", d.hand.Hand[0])

	if g.state != dealerTurn {
		fmt.Println("- <hidden>")
		return
	}

	for _, c := range d.hand.Hand[1:] {
		fmt.Printf("- %s\n", c)
	}
}

type player struct {
	name     string
	hand     Hand
	winnings int
	bet      int
	ai       AI
}

func (p *player) printHand() {
	fmt.Printf("\n%s's Hand:\n", p.name)
	for _, c := range p.hand.Hand {
		fmt.Printf("- %s\n", c)
	}
}

func (p *player) play(g *game) {
	if !g.isSimulation {
		fmt.Printf("\n%s's turn\n", p.name)
		fmt.Println("-------------------------------")
	}

	if p.hand.Value() == 21 {
		fmt.Printf("%s got a natural blackjack!\n", p.name)
		return
	}

	if g.isSimulation && p.ai.DoubleDown(p.hand) {
		fmt.Println(p.name, "double downs!")
		p.bet *= 2
	}

	for {
		var choice string

		if !g.isSimulation {
			g.dealer.printHand(g)
			p.printHand()
			choice = getUserInput("\n[1] HIT\n[2] STAND\nSelect an option: ", validateChooseOneOrTwo)
		}

		if choice == "2" || (g.isSimulation && !p.ai.DecideHit(p.hand)) {
			fmt.Printf("%s stands.\n", p.name)
			return
		}

		p.draw(g)

		handVal := p.hand.Value()
		if handVal == 21 {
			fmt.Printf("%s gets a Blackjack!\n", p.name)
			return
		} else if handVal > 21 {
			fmt.Printf("%s busts!\n", p.name)
			return
		}
	}
}

func (p *player) draw(g *game) {
	card := g.drawCard()
	fmt.Printf("%s draws a %s\n", p.name, card)
	p.hand.addCard(card)
}

type game struct {
	deck         []deck.Card
	players      []player
	dealer       dealer
	state        state
	rounds       int
	isSimulation bool
}

func (g *game) drawCard() deck.Card {
	c := g.deck[len(g.deck)-1]
	g.deck = g.deck[:len(g.deck)-1]
	return c
}

func (g *game) printPlayerWinnings() {
	for _, p := range g.players {
		fmt.Printf("%s winnings: %d\n", p.name, p.winnings)
	}
}

func (g *game) initPlayers(numPlayers int) {
	g.players = make([]player, 0, numPlayers)
	for i := 0; i < numPlayers; i++ {
		prompt := fmt.Sprintf("player %d enter your name: ", i+1)
		name := getUserInput(prompt, validateNonEmptyString)
		g.players = append(g.players, player{name: name})
	}
}

func (g *game) Start() {
	for i := 0; !g.isSimulation || i < g.rounds; i++ {
		g.reset()
		g.deal()
		g.play()
		g.finish()
	}

	if g.isSimulation {
		fmt.Printf("AI won/lost %d after %d rounds.\n", g.players[0].winnings, g.rounds)
	}
}

func SetupSimulation(strategy DealerStrategy, ai AI, rounds int) *game {
	p := player{ai: ai, name: "AI"}
	game := &game{
		players:      []player{p},
		dealer:       dealer{strategy: strategy},
		isSimulation: true,
		rounds:       rounds,
	}

	return game
}

func Setup(strategy DealerStrategy) *game {
	input := getUserInput("Enter number of players: ", validatePositiveInteger)

	numPlayers, err := strconv.Atoi(input)
	if err != nil {
		panic("Something unexpected occured")
	}

	game := &game{
		dealer:       dealer{strategy: strategy},
		isSimulation: false,
	}
	game.initPlayers(numPlayers)

	return game
}

func (g *game) deal() {
	g.state = dealing

	fmt.Println("Dealing cards...")
	for i, ppl := 0, 1+len(g.players); i < ppl*2; i++ {
		if i%ppl == 0 {
			g.dealer.draw(g)
		} else {
			g.players[i%ppl-1].draw(g)
		}
	}
}

func (g *game) reset() {
	g.deck = deck.NewDeck(deck.WithShuffle())

	//Empty everyone's hand
	g.dealer.hand = Hand{}
	for i := range g.players {
		bet := 0
		g.players[i].hand = Hand{}

		if !g.isSimulation {
			prompt := fmt.Sprintf("%s enter bet amount: ", g.players[i].name)
			input := getUserInput(prompt, validatePositiveInteger)
			num, _ := strconv.Atoi(input)
			bet = num

		} else {
			bet = g.players[i].ai.DecideBet()
		}
		g.players[i].bet = bet
	}
}

func (g *game) play() {
	g.state = playerTurn

	for i := range g.players {
		g.players[i].play(g)
	}
	g.state = dealerTurn
	g.dealer.play(g)
	g.state = end
}

func (g *game) finish() {

	for i := range g.players {
		result := getWinner(g.players[i].hand, g.dealer.hand)
		bet := g.players[i].bet
		if result > 0 {
			g.players[i].winnings += bet
			fmt.Printf("\n%s wins %d\n", g.players[i].name, bet)
		} else if result < 0 {
			g.players[i].winnings -= bet
			fmt.Printf("\n%s loses %d\n", g.players[i].name, bet)
		} else {
			fmt.Printf("\n%s ties with dealer\n", g.players[i].name)
		}
	}
	g.printPlayerWinnings()

	if !g.isSimulation && !continueGame() {
		os.Exit(0)
	}
}

func cardValue(c deck.Card) int {
	switch c.Type {
	case deck.JACK, deck.QUEEN, deck.KING:
		return 10
	default:
		return int(c.Type)
	}
}

// didPlayerWin returns a positive number if player wins, a negative
// number if the dealer wins, and 0 if it is a tie.
func getWinner(player, dealer Hand) int8 {
	playerVal, dealerVal := player.Value(), dealer.Value()

	if playerVal > 21 {
		return -1
	} else if dealerVal > 21 {
		return 1
	}

	if playerVal > dealerVal {
		return 1
	} else if dealerVal > playerVal {
		return -1
	}

	return 0
}

func continueGame() bool {
	input := getUserInput("Continue playing? (y/n): ", validateYesOrNo)

	return input[0] == 'y' || input[0] == 'Y'
}

func validatePositiveInteger(s string) bool {
	num, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	return num > 0
}

func validateYesOrNo(s string) bool {
	switch s {
	case "y", "Y", "n", "N", "yes", "no", "Yes", "No":
		return true
	}
	return false
}

func validateChooseOneOrTwo(s string) bool {
	return s == "1" || s == "2"
}

func validateNonEmptyString(s string) bool {
	return s != ""
}

func getUserInput(prompt string, validate func(string) bool) string {
	var input string
	for {
		fmt.Printf("\n%s", prompt)
		if _, err := fmt.Scanf("%s\n", &input); err != nil || !validate(input) {
			fmt.Print("Invalid input")
			continue
		}
		break
	}

	return input
}
