package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(D6())
	result := RollDicePool(10000, 3, 6)
	fmt.Println(result.CritsCount)
	fmt.Println(result.NonCritsCount)
	fmt.Println(result.Rolls)

	js.Global().Set("MakeAttackSeries", WASMMakeAttackSeries())

	body1 := NewBody()
	body1.Defence = 4
	body1.Wounds = 7
	body1.FNP = 5
	body1.Save = 5

	weapon1 := NewWeapon()
	weapon1.BS = 3
	weapon1.Attacks = 4
	weapon1.MW = 3
	weapon1.NormalDamage = 2
	weapon1.CriticalDamage = 3

	attackSeriesResult := MakeAttackSeries(1000, weapon1, body1, 2)
	fmt.Println(attackSeriesResult)

	//fmt.Println(MakeAttackRound(weapon1, body1, 2))
	select {}
}

func WASMMakeAttackSeries() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		body1 := NewBody()
		body1.Defence = 4
		body1.Wounds = 7
		body1.FNP = 5
		body1.Save = 5

		weapon1 := NewWeapon()
		weapon1.BS = 3
		weapon1.Attacks = 4
		weapon1.MW = 3
		weapon1.NormalDamage = 2
		weapon1.CriticalDamage = 3

		attackSeriesResult := MakeAttackSeries(1000, weapon1, body1, 2)
		return map[string]interface{}{
			"hello":  "world",
			"answer": 42,
			"killed": attackSeriesResult.Killed,
		}

	})
}
