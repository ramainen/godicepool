package main

import "math/rand"

const REROLL_NONE int = 0
const REROLL_ONES int = 1
const REROLL_ONE_ROLL int = 2
const REROLL_ALL int = 3

func D6() int {
	return rand.Intn(6) + 1
}
func XD6plus(numberOfDices int, plus int) int {
	result := 0
	if numberOfDices < 1 {
		return 0
	}
	for i := 1; i <= numberOfDices; i++ {
		if D6() >= plus {
			result++
		}
	}
	return result
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
	Attacks        int
	NormalDamage   int
	CriticalDamage int
	BS             int
	Lethal         int
	AP             int
	MW             int
	Brutal         int
	P              int
	Reroll         int
	Rending        int
}

func NewWeapon() Weapon {
	weapon := Weapon{}
	weapon.Lethal = 6  //DONE
	weapon.AP = 0      //DONE
	weapon.MW = 0      //DONE
	weapon.Brutal = 0  //DONE
	weapon.Reroll = 0  //DONE
	weapon.P = 0       //DONE
	weapon.Rending = 0 //DONE
	return weapon
}

type Body struct {
	Wounds           int
	Defence          int
	Save             int
	InvulnerableSave int
	FNP              int
	InCover          int
	Reroll           int
}

func NewBody() Body {
	body := Body{}
	body.FNP = 0
	body.InvulnerableSave = 0
	body.InCover = 0
	body.Reroll = 0
	return body
}

//RollDicePool returns array of rolled dices, with marks about how many krits and not krits rolls
func RollDicePool(count int, bs int, critValue int, rerolls int) DicePoolResult {
	/*
		rerolls=
			0 - none
			1 - only ones
			2 - only one failed
			3 - all failed

	*/
	result := DicePoolResult{}
	atMostOneFailed := false
	switch rerolls {
	case 0:
		for i := 1; i <= count; i++ {
			rollResult := D6()
			result.Rolls[rollResult]++
			if rollResult == 1 {
				//Nothing happens
			} else if rollResult >= critValue || rollResult == 6 {
				result.CritsCount++
			} else if rollResult >= bs {
				result.NonCritsCount++
			}
		}
	case 1:
		//reroll ones
		for i := 1; i <= count; i++ {
			rollResult := D6()
			if rollResult == 1 {
				rollResult = D6()
			}
			result.Rolls[rollResult]++
			if rollResult == 1 {
				//Nothing happens
			} else if rollResult >= critValue || rollResult == 6 {
				result.CritsCount++
			} else if rollResult >= bs {
				result.NonCritsCount++
			}
		}
	case 2:
		//only one failed
		for i := 1; i <= count; i++ {
			rollResult := D6()
			if !atMostOneFailed && rollResult < bs {
				rollResult = D6()
				atMostOneFailed = true
			}
			result.Rolls[rollResult]++
			if rollResult == 1 {
				//Nothing happens
			} else if rollResult >= critValue || rollResult == 6 {
				result.CritsCount++
			} else if rollResult >= bs {
				result.NonCritsCount++
			}
		}
	case 3:
		//all failed
		for i := 1; i <= count; i++ {
			rollResult := D6()
			if rollResult < bs {
				rollResult = D6()
			}
			result.Rolls[rollResult]++
			if rollResult == 1 {
				//Nothing happens
			} else if rollResult >= critValue || rollResult == 6 {
				result.CritsCount++
			} else if rollResult >= bs {
				result.NonCritsCount++
			}
		}
	default:
		//This is an error... I do no not know;
		result.CritsCount = 0
		result.NonCritsCount = 0
		result.Rolls = [7]int{1, 1, 1, 1, 1, 1, 1}
	}

	return result
}
