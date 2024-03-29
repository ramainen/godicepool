package main

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	js.Global().Set("MakeAttackSeries", WASMMakeAttackSeries())
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

		BrutalValue := args[0].Get("brutal")
		if !BrutalValue.IsUndefined() && !BrutalValue.IsNull() {
			weapon1.Brutal = BrutalValue.Int()
		}

		LethalValue := args[0].Get("lethal")
		if !LethalValue.IsUndefined() && !LethalValue.IsNull() {
			weapon1.Lethal = LethalValue.Int()
			if weapon1.Lethal == 0 {
				weapon1.Lethal = 6
			}
		}

		PValue := args[0].Get("p")
		if !PValue.IsUndefined() && !PValue.IsNull() {
			weapon1.P = PValue.Int()
		}

		RerollValue := args[0].Get("reroll")
		if !RerollValue.IsUndefined() && !RerollValue.IsNull() {
			weapon1.Reroll = RerollValue.Int()
		}

		RendingValue := args[0].Get("rending")
		if !RendingValue.IsUndefined() && !RendingValue.IsNull() {
			weapon1.Rending = RendingValue.Int()
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
		InvulnerableSaveValue := args[1].Get("invulnerable_save")
		if !InvulnerableSaveValue.IsUndefined() && !InvulnerableSaveValue.IsNull() {
			body1.InvulnerableSave = InvulnerableSaveValue.Int()
		}
		FNPValue := args[1].Get("fnp")
		if !FNPValue.IsUndefined() && !FNPValue.IsNull() {
			body1.FNP = FNPValue.Int()
		}

		InCoverValue := args[1].Get("in_cover")
		if !InCoverValue.IsUndefined() && !InCoverValue.IsNull() {
			body1.InCover = InCoverValue.Int()
		}

		BodyRerollValue := args[1].Get("reroll")
		if !BodyRerollValue.IsUndefined() && !BodyRerollValue.IsNull() {
			body1.Reroll = BodyRerollValue.Int()
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

		/*
			type AttackSeriesResult struct {
				Killed        int   //Killed in all series
				KilledInRound []int //Killed in particular round
				MakedWounds   []int //How many wounds dealed (1 wound in 3000 simulations, 2 wounds in 500, etc)
			}

		*/
		killedInRound := make([]interface{}, 0)
		for _, v := range attackSeriesResult.KilledInRound {
			killedInRound = append(killedInRound, v)
		}
		makedWounds := make([]interface{}, 0)
		for _, v := range attackSeriesResult.MakedWounds {
			makedWounds = append(makedWounds, v)
		}

		return map[string]interface{}{
			"killed":          attackSeriesResult.Killed,
			"killed_in_round": killedInRound,
			"maked_wounds":    makedWounds,
		}

	})
}
