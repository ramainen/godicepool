package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(D6())
	result := MakeAttacksRolls(10000, 3, 6)
	fmt.Println(result.CritsCount)
	fmt.Println(result.NonCritsCount)
	fmt.Println(result.Rolls)
}

func D6() int {
	return rand.Intn(6) + 1
}

type DicePoolResult struct {
	Rolls         [7]int
	CritsCount    int
	NonCritsCount int

	//FailedCount   int
	//SixesCount    int
	//OnesCount     int
}

type Weapon struct {
	Attacks int
	BS      int
	Lethal  int
	AP      int
	MW      int
	Rules   struct {
		Sniper bool
		Magic  bool
		Etc    bool
	}
}

func NewWeapon() Weapon {
	weapon := Weapon{}
	weapon.Lethal = 6
	weapon.AP = 0
	weapon.MW = 0
	return weapon
}

type Body struct {
	Wounds  int
	Defence int
	Save    int
	FNP     int
	Rules   struct {
		AllIsDust bool
		Xenos     bool
		Etc       bool
	}
}

func NewBody() Body {
	body := Body{}
	body.FNP = 0
	return body
}

//MakeAttacksRolls returns array of rolled dices, with marks about how many krits and not krits rolls
func MakeAttacksRolls(count int, bs int, critValue int) DicePoolResult {
	result := DicePoolResult{}
	for i := 1; i <= count; i++ {
		rollResult := D6()
		result.Rolls[rollResult]++
		if rollResult >= critValue {
			result.CritsCount++
		} else if rollResult >= bs {
			result.NonCritsCount++
		}
	}
	return result
}
