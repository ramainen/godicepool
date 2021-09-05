package main

import "fmt"

type AttackResult struct {
	Killed         bool
	AffectedWounds int
}

type AttackRoundResult struct {
	Killed         bool
	AffectedWounds int
	KilledInRound  int
	WoundsInRound  []int
}

type AttackSeriesResult struct {
	Killed        int   //Killed in all series
	KilledInRound []int //Killed in particular round
	MakedWounds   []int //How many wounds dealed (1 wound in 3000 simulations, 2 wounds in 500, etc)
}

func MakeAttack(weapon Weapon, body Body) (AttackResult, Body) {

	originalWounds := body.Wounds
	successFNP := 0
	result := AttackResult{}
	bs := weapon.BS
	attacks := weapon.Attacks
	critvalue := weapon.Lethal

	attackRolls := RollDicePool(attacks, bs, critvalue, weapon.Reroll)
	attackCritCount := attackRolls.CritsCount
	attackSuccessCount := attackRolls.NonCritsCount

	defence := body.Defence

	if weapon.P > 0 {
		if attackCritCount > 0 {
			weapon.AP = weapon.P
		}
	}

	if weapon.Rending > 0 {
		if attackCritCount > 0 {
			attackCritCount = attackCritCount + attackSuccessCount
			attackSuccessCount = 0
		}
	}

	//InvulnerableSave. Scenario: opponent has APs and invul better than usual saves. Behind this all depends to end user

	saveroll := body.Save
	if weapon.AP > 0 && body.InvulnerableSave > 0 {
		saveroll = body.InvulnerableSave
		weapon.AP = 0
	}

	//deal with AP
	if weapon.AP > 0 {
		defence = defence - weapon.AP
	}

	if body.InCover == 1 && defence > 0 {
		defence = defence - 1
	}

	saveRolls := RollDicePool(defence, saveroll, 6, body.Reroll)
	defenceCritCount := saveRolls.CritsCount
	defenceSuccessCount := saveRolls.NonCritsCount

	if body.InCover == 1 && defence > 0 {
		defenceSuccessCount++
	}
	//Mortal wounds works now, BEFORE chance to prevent them
	if weapon.MW > 0 && attackCritCount > 0 {
		if body.FNP > 0 {
			successFNP = XD6plus(weapon.MW*attackCritCount, body.FNP)
		}

		if successFNP < weapon.MW*attackCritCount {
			body.Wounds = body.Wounds - (weapon.MW*attackCritCount - successFNP)
			successFNP = 0
		}
		if body.Wounds <= 0 {
			body.Wounds = 0
			result.AffectedWounds = originalWounds
			result.Killed = true
			return result, body
		}

	}

	//Deal with crit and non crit rolls
	if attackCritCount > 0 {
		if defenceCritCount <= attackCritCount {
			attackCritCount = attackCritCount - defenceCritCount
			defenceCritCount = 0
		} else {
			attackCritCount = 0
			defenceCritCount = defenceCritCount - attackCritCount
		}
	}
	//Deal with attack crit rolls with pairs of defence rolls; Brutal do not allow this
	if attackCritCount > 0 && weapon.Brutal == 0 {
		if defenceSuccessCount <= attackCritCount*2 {
			//for example, defence success = 3, crits = 2, 3 <= 2*2, so...
			//..crits = 1, because only one is blocked

			resolvedCritAttacks := (int)(defenceSuccessCount / 2)
			attackCritCount = attackCritCount - resolvedCritAttacks

			//for another example, defence success = 3, crits = 99, 3 <= 99*2, so...
			defenceSuccessCount = defenceSuccessCount - resolvedCritAttacks*2
		} else {
			//attacks crits count is less than pairs of defence successes, for example 2 attack crits and 5 success defences, 1 remaining
			attackCritCount = 0
			defenceSuccessCount = defenceSuccessCount - attackCritCount*2
		}
	}
	//Deal with usual attacks
	//Deal with crit defence rolls
	if attackSuccessCount > 0 {
		if defenceCritCount <= attackSuccessCount {
			attackSuccessCount = attackSuccessCount - defenceCritCount
			defenceCritCount = 0
		} else {
			//defenceCrits is more than attacks

			//defenceCritCount = defenceCritCount - attackSuccessCount //this can be deleted
			attackSuccessCount = 0
		}
	}
	//Deal with usual defence rolls; Brutal down not allow this
	if attackSuccessCount > 0 && weapon.Brutal == 0 {
		if defenceSuccessCount <= attackSuccessCount {
			attackSuccessCount = attackSuccessCount - defenceSuccessCount
			defenceCritCount = 0
		} else {
			//defenceSuccessCount is more than attacks
			//defenceSuccessCount = defenceSuccessCount - attackSuccessCount //this can be deleted
			attackSuccessCount = 0
		}
	}

	//deal wounds!

	if body.FNP > 0 {
		successFNP = XD6plus(weapon.CriticalDamage*attackCritCount, body.FNP)
	}
	//FXIME!!!
	if successFNP < weapon.CriticalDamage*attackCritCount {
		body.Wounds = body.Wounds - (weapon.CriticalDamage*attackCritCount - successFNP)
	}
	if body.Wounds <= 0 {
		body.Wounds = 0
		result.AffectedWounds = originalWounds
		result.Killed = true
		return result, body
	}
	//FXIME!!!
	if body.FNP > 0 {
		successFNP = XD6plus(weapon.NormalDamage*attackSuccessCount, body.FNP)
	}
	if successFNP < weapon.NormalDamage*attackSuccessCount {
		body.Wounds = body.Wounds - (weapon.NormalDamage*attackSuccessCount - successFNP)
	}
	if body.Wounds <= 0 {
		body.Wounds = 0
		result.AffectedWounds = originalWounds
		result.Killed = true
		return result, body
	}
	result.AffectedWounds = originalWounds - body.Wounds
	if result.AffectedWounds < 0 {
		fmt.Println("originalWounds, body, ", originalWounds, body.Wounds)
		//result.AffectedWounds = 0
	}
	result.Killed = false
	return result, body
}

func MakeAttackRound(weapon Weapon, body Body, rounds int) AttackRoundResult {
	result := AttackRoundResult{}
	result.Killed = false
	if rounds < 1 {
		return result
	}
	result.WoundsInRound = make([]int, rounds+1)

	for i := 1; i <= rounds; i++ {
		roundResult := AttackResult{}
		roundResult, body = MakeAttack(weapon, body)
		result.AffectedWounds = result.AffectedWounds + roundResult.AffectedWounds
		result.WoundsInRound[i] = roundResult.AffectedWounds
		if roundResult.Killed {
			result.Killed = true
			result.KilledInRound = i
			return result
		}
	}
	return result
}
func MakeAttackSeries(count int, weapon Weapon, body Body, rounds int) AttackSeriesResult {
	result := AttackSeriesResult{}
	if count < 1 || rounds < 1 {
		return result
	}
	result.KilledInRound = make([]int, rounds+1)
	result.MakedWounds = make([]int, body.Wounds+1)

	for i := 1; i <= count; i++ {
		oneSimulationResult := MakeAttackRound(weapon, body, rounds)
		result.MakedWounds[oneSimulationResult.AffectedWounds]++
		if oneSimulationResult.Killed {
			result.Killed++
			result.KilledInRound[oneSimulationResult.KilledInRound]++
		}
	}
	return result
}
