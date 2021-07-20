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

		/*
			window.MakeAttackSeries({'attacks': 4, 'bs': 3, 'normal_damage': 4,'critical_damage':5},{'wounds':7, 'save':5, 'defence': 3},{'count': 1000, 'rounds':5})
		*/
		weapon1 := NewWeapon()
		weapon1.BS = 9
		weapon1.Attacks = 1
		weapon1.NormalDamage = 0
		weapon1.CriticalDamage = 0

		attacksValue := args[0].Get("attacks")
		if !attacksValue.IsUndefined() && !attacksValue.IsNull() {
			weapon1.Attacks = attacksValue.Int()
		}

		NormalDamageValue := args[0].Get("normal_damage")
		if !NormalDamageValue.IsUndefined() && !NormalDamageValue.IsNull() {
			weapon1.NormalDamage = NormalDamageValue.Int()
		}

		CriticalDamageValue := args[0].Get("critical_damage")
		if !CriticalDamageValue.IsUndefined() && !CriticalDamageValue.IsNull() {
			weapon1.CriticalDamage = CriticalDamageValue.Int()
		}
		BSValue := args[0].Get("bs")
		if !BSValue.IsUndefined() && !BSValue.IsNull() {
			weapon1.BS = BSValue.Int()
		}
		APValue := args[0].Get("ap")
		if !APValue.IsUndefined() && !APValue.IsNull() {
			weapon1.AP = APValue.Int()
		}

		MWValue := args[0].Get("mw")
		if !MWValue.IsUndefined() && !MWValue.IsNull() {
			weapon1.MW = MWValue.Int()
		}

		fmt.Println(weapon1)
		body1 := NewBody()
		body1.Defence = 9
		body1.Wounds = 1
		body1.Save = 9

		DefenceValue := args[1].Get("defence")
		if !DefenceValue.IsUndefined() && !DefenceValue.IsNull() {
			body1.Defence = DefenceValue.Int()
		}
		WoundsValue := args[1].Get("wounds")
		if !WoundsValue.IsUndefined() && !WoundsValue.IsNull() {
			body1.Wounds = WoundsValue.Int()
		}

		SaveValue := args[1].Get("save")
		if !SaveValue.IsUndefined() && !SaveValue.IsNull() {
			body1.Save = SaveValue.Int()
		}
		FNPValue := args[1].Get("fnp")
		if !FNPValue.IsUndefined() && !FNPValue.IsNull() {
			body1.FNP = FNPValue.Int()
		}

		simulationsCount := 1000
		if len(args) > 2 {
			CountValue := args[2].Get("count")
			if !CountValue.IsUndefined() && !CountValue.IsNull() {
				simulationsCount = CountValue.Int()
			}
		}

		roundsCount := 1
		if len(args) > 2 {
			roundsCountValue := args[2].Get("rounds")
			if !roundsCountValue.IsUndefined() && !roundsCountValue.IsNull() {
				roundsCount = roundsCountValue.Int()
			}
		}

		attackSeriesResult := MakeAttackSeries(simulationsCount, weapon1, body1, roundsCount)
		return map[string]interface{}{
			"killed": attackSeriesResult.Killed,
		}

	})
}
