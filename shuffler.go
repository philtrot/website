package main

import "fmt"
import "math/rand"
import "os"
import "sort"
import "strconv"
import "time"

const numCards = 40

// conclusion -> best on the long term is 9 piles!

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no")
		return
	}

	piles := make([]int, 0, len(os.Args)-1)
	for i := 1; i < len(os.Args); i++ {
        v, err := strconv.Atoi(os.Args[i])
		if err != nil {
			panic(err)
		}
		piles = append(piles, v)
	}
	fmt.Println("Piles:", piles)

	rand.Seed(time.Now().UnixNano())
	deck := make([]int, numCards)
	/* hardcoded deck1 -> flaw: perfect mix of second half penalizes pair number of piles
	// fill the deck, 1/4 terrain, 1/4 creature, 1/2 mixed
	for i := range deck {
		switch {
			case i < len(deck)/4:
				deck[i] = 1;
			case i < len(deck)/2:
				deck[i] = 2;
			default:
				deck[i] = i&1 + 1
		}
	}

	for _, v := range piles {
		fmt.println(deck)
		pshuffle(deck, v)
	}
	fmt.println(deck)

	fmt.println("score:", analyze(deck))
	*/

	/* hardcoded deck2: initial setup of all terrains at start! 
	// fill the deck, 20-60-20
	for i := range deck {
		switch {
			case i <= int(float32(len(deck))*.2):
				deck[i] = 1;
			case i >= int(float32(len(deck))*.8):
				deck[i] = 1;
			default:
				deck[i] = 2
		}
	}

	for _, v := range piles {
		fmt.Println(deck)
		pshuffle(deck, v)
	}
	fmt.Println(deck)

	fmt.Println("score:", analyze(deck))
	*/

	// random deck with 40% land
	for i := range deck {
		if i <= int(.4 * float32(len(deck))) {
			deck[i] = 1
			continue
		}
		deck[i] = 2
	}
	//fmt.Println(deck)

	resets := 500
	ngames := 31
	sum := 0
	tries := make([]int, 0, resets * ngames)
	for i := 0; i < resets; i++ {
		// reshuffle the deck perfectly
		rand.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})

		// play n games using only this shuffle
		for g := 0; g < ngames; g++ {
			// unbalance it semi-realistically: 50% of deck played, hand+bg(35); terrains; graveyard(65)
			numTerrains := 0
			for i := 0; i < len(deck)/2; i++ {
				if deck[i] == 1 {
					numTerrains++
				}
				deck[i] = 2
			}
			numOther := len(deck)/2 - numTerrains
			offset := int(.35 * float32(numOther))
			for i := offset; i < offset+numTerrains; i++ {
				deck[i] = 1
			}
			//fmt.Println(deck)

			// do the pile shuffle
			for _, v := range piles {
				pshuffle(deck, v)
			}
			//fmt.Println(deck)

			a := analyze(deck)
			sum += a
			tries = append(tries, a)
		}
	}

	sort.Ints(tries)
	//fmt.Println(tries)
	fmt.Println("Avg:", sum/len(tries))
	fmt.Println("80p:", tries[int(float32(len(tries))*.8)])
	fmt.Println("90p:", tries[int(float32(len(tries))*.9)])
}

func pshuffle(deck []int, n int) {
	pileSz := (len(deck) + n - 1) / n
	piles := make([][]int, n)
	for i := range piles {
		piles[i] = make([]int, 0, pileSz)
	}

	for i, v := range deck {
		idx := i % n
		piles[idx] = append(piles[idx], v)
	}

	deck = deck[:0]
	for _, p := range piles {
		for _, v := range p {
			deck = append(deck, v)
		}
	}
}

func analyze(deck []int) int {
	prev := -1
	cur := 0
	streaks := make([]int, 0, len(deck)/4)
	for _, v := range deck {
		if prev == v {
			cur++
			continue
		}
		prev = v
		if cur > 0 {
			streaks = append(streaks, cur)
			cur = 0
		}
	}
	if cur > 0 {
		streaks = append(streaks, cur)
	}

	score := 0
	for _, s := range streaks {
		score += s*s
	}
	return score
}
