package main

import (
	"debug/pe"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/sanctuary/formats/image/cel"
)

type Item struct {
	_iPLToHit      int
	_iPLToHitRange Range
	_iPLDam        int
	_iPLAC         int
	_iAC           int
	_iPLFR         int
	_iPLLR         int
	_iPLMR         int
	_iSplLvlAdd    int
	_iCharges      int
	_iMaxCharges   int
	_iSpell        int
	_iFlags        int
	_iFMinDam      int
	_iFMaxDam      int
	_iLMinDam      int
	_iLMaxDam      int
	_iPLStr        int
	_iPLMag        int
	_iPLDex        int
	_iPLVit        int
	_iPLGetHit     int
	_iPLHP         int
	_iPLMana       int
	_iMaxDur       int
	_iDurability   int
	_iPLLight      int
	_iPLEnAc       int
	_iMinStr       int
	_iPLDamMod     int
	_iMinDam       int
	_iMaxDam       int
	_iLoc          int
	_iCurs         int
}

type Range struct {
	Min, Max int
}

func DescribePower(power *StructVal, r int, pLevel int) string {
	item := &Item{}
	if power.Type != PowerData {
		log.Fatal("Invalid power passed to UpdateItem")
	}
	param1 := int(power.Get("PLParam1").(uint32))
	param2 := int(power.Get("PLParam2").(uint32))
	pIdx := power.Get("PLPower").(uint32)
	switch pIdx {
	case IPL_TOHIT:
		item._iPLToHit += r
	case IPL_TOHIT_CURSE:
		item._iPLToHit -= r
	case IPL_DAMP:
		item._iPLDam += r
	case IPL_DAMP_CURSE:
		item._iPLDam -= r
	case IPL_TOHIT_DAMP:
		item._iPLDam += r
		var r2 Range
		if param1 == 20 {
			r2 = Range{1, 5}
		}
		if param1 == 36 {
			r2 = Range{6, 10}
		}
		if param1 == 51 {
			r2 = Range{11, 15}
		}
		if param1 == 66 {
			r2 = Range{16, 20}
		}
		if param1 == 81 {
			r2 = Range{21, 30}
		}
		if param1 == 96 {
			r2 = Range{31, 40}
		}
		if param1 == 111 {
			r2 = Range{41, 50}
		}
		if param1 == 126 {
			r2 = Range{51, 75}
		}
		if param1 == 151 {
			r2 = Range{76, 100}
		}
		item._iPLToHitRange = r2
	case IPL_TOHIT_DAMP_CURSE:
		item._iPLDam -= r
		var r2 Range
		if param1 == 25 {
			r2 = Range{-1, -5}
		}
		if param1 == 50 {
			r2 = Range{-6, -10}
		}
		item._iPLToHitRange = r2
	case IPL_ACP:
		item._iPLAC += r
	case IPL_ACP_CURSE:
		item._iPLAC -= r
	case IPL_SETAC:
		item._iAC = r
	case IPL_AC_CURSE:
		item._iAC -= r
	case IPL_FIRERES:
		item._iPLFR += r
	case IPL_LIGHTRES:
		item._iPLLR += r
	case IPL_MAGICRES:
		item._iPLMR += r
	case IPL_ALLRES:
		item._iPLFR += r
		item._iPLLR += r
		item._iPLMR += r
		if item._iPLFR < 0 {
			item._iPLFR = 0
		}
		if item._iPLLR < 0 {
			item._iPLLR = 0
		}
		if item._iPLMR < 0 {
			item._iPLMR = 0
		}
	case IPL_SPLLVLADD:
		item._iSplLvlAdd = r
	case IPL_CHARGES:
		item._iCharges *= param1
		item._iMaxCharges = item._iCharges
	case IPL_SPELL:
		item._iSpell = param1
		item._iCharges = param1
		item._iMaxCharges = param2
	case IPL_FIREDAM:
		item._iFlags |= ISPL_FIREDAM
		item._iFMinDam = param1
		item._iFMaxDam = param2
	case IPL_LIGHTDAM:
		item._iFlags |= ISPL_LIGHTDAM
		item._iLMinDam = param1
		item._iLMaxDam = param2
	case IPL_STR:
		item._iPLStr += r
	case IPL_STR_CURSE:
		item._iPLStr -= r
	case IPL_MAG:
		item._iPLMag += r
	case IPL_MAG_CURSE:
		item._iPLMag -= r
	case IPL_DEX:
		item._iPLDex += r
	case IPL_DEX_CURSE:
		item._iPLDex -= r
	case IPL_VIT:
		item._iPLVit += r
	case IPL_VIT_CURSE:
		item._iPLVit -= r
	case IPL_ATTRIBS:
		item._iPLStr += r
		item._iPLMag += r
		item._iPLDex += r
		item._iPLVit += r
	case IPL_ATTRIBS_CURSE:
		item._iPLStr -= r
		item._iPLMag -= r
		item._iPLDex -= r
		item._iPLVit -= r
	case IPL_GETHIT_CURSE:
		item._iPLGetHit += r
	case IPL_GETHIT:
		item._iPLGetHit -= r
	case IPL_LIFE:
		item._iPLHP += r << 6
	case IPL_LIFE_CURSE:
		item._iPLHP -= r << 6
	case IPL_MANA:
		item._iPLMana += r << 6
	case IPL_MANA_CURSE:
		item._iPLMana -= r << 6
	case IPL_DUR:
		// r2 := r * int(base.Get("iDurability").(uint32)) / 100
		// item._iMaxDur += r2
		// item._iDurability += r2
	case IPL_DUR_CURSE:
		item._iMaxDur -= r * item._iMaxDur / 100
		if item._iMaxDur < 1 {
			item._iMaxDur = 1
		}
		item._iDurability = item._iMaxDur
	case IPL_INDESTRUCTIBLE:
		// item._iDurability = DUR_INDESTRUCTIBLE
		// item._iMaxDur = DUR_INDESTRUCTIBLE
	case IPL_LIGHT:
		item._iPLLight += param1
	case IPL_LIGHT_CURSE:
		item._iPLLight -= param1
	case IPL_FIRE_ARROWS:
		item._iFlags |= ISPL_FIRE_ARROWS
		item._iFMinDam = param1
		item._iFMaxDam = param2
	case IPL_LIGHT_ARROWS:
		item._iFlags |= ISPL_LIGHT_ARROWS
		item._iLMinDam = param1
		item._iLMaxDam = param2
	case IPL_THORNS:
		item._iFlags |= ISPL_THORNS
	case IPL_NOMANA:
		item._iFlags |= ISPL_NOMANA
	case IPL_NOHEALPLR:
		item._iFlags |= ISPL_NOHEALPLR
	case IPL_ABSHALFTRAP:
		item._iFlags |= ISPL_ABSHALFTRAP
	case IPL_KNOCKBACK:
		item._iFlags |= ISPL_KNOCKBACK
	case IPL_3XDAMVDEM:
		item._iFlags |= ISPL_3XDAMVDEM
	case IPL_ALLRESZERO:
		item._iFlags |= ISPL_ALLRESZERO
	case IPL_NOHEALMON:
		item._iFlags |= ISPL_NOHEALMON
	case IPL_STEALMANA:
		if param1 == 3 {
			item._iFlags |= ISPL_STEALMANA_3
		}
		if param1 == 5 {
			item._iFlags |= ISPL_STEALMANA_5
		}
	case IPL_STEALLIFE:
		if param1 == 3 {
			item._iFlags |= ISPL_STEALLIFE_3
		}
		if param1 == 5 {
			item._iFlags |= ISPL_STEALLIFE_5
		}
	case IPL_TARGAC:
		item._iPLEnAc += r
	case IPL_FASTATTACK:
		if param1 == 1 {
			item._iFlags |= ISPL_QUICKATTACK
		}
		if param1 == 2 {
			item._iFlags |= ISPL_FASTATTACK
		}
		if param1 == 3 {
			item._iFlags |= ISPL_FASTERATTACK
		}
		if param1 == 4 {
			item._iFlags |= ISPL_FASTESTATTACK
		}

	case IPL_FASTRECOVER:
		if param1 == 1 {
			item._iFlags |= ISPL_FASTRECOVER
		}
		if param1 == 2 {
			item._iFlags |= ISPL_FASTERRECOVER
		}
		if param1 == 3 {
			item._iFlags |= ISPL_FASTESTRECOVER
		}
	case IPL_FASTBLOCK:
		item._iFlags |= ISPL_FASTBLOCK
	case IPL_DAMMOD:
		item._iPLDamMod += r
	case IPL_RNDARROWVEL:
		item._iFlags |= ISPL_RNDARROWVEL
	case IPL_SETDAM:
		item._iMinDam = param1
		item._iMaxDam = param2
	case IPL_SETDUR:
		item._iDurability = param1
		item._iMaxDur = param1
	case IPL_FASTSWING:
		item._iFlags |= ISPL_FASTERATTACK
	case IPL_ONEHAND:
		item._iLoc = ILOC_ONEHAND
	case IPL_DRAINLIFE:
		item._iFlags |= ISPL_DRAINLIFE
	case IPL_RNDSTEALLIFE:
		item._iFlags |= ISPL_RNDSTEALLIFE
	case IPL_INFRAVISION:
		item._iFlags |= ISPL_INFRAVISION
	case IPL_NOMINSTR:
		item._iMinStr = 0
	case IPL_INVCURS:
		item._iCurs = param1
	case IPL_ADDACLIFE:
		// item._iPLHP = (plr[myplr]._pIBonusAC + plr[myplr]._pIAC + plr[myplr]._pDexterity/5) << 6
	case IPL_ADDMANAAC:
		// item._iAC += (plr[myplr]._pMaxManaBase >> 6) / 10
	case IPL_FIRERESCLVL:
		item._iPLFR = 30 - pLevel
		if item._iPLFR < 0 {
			item._iPLFR = 0
		}
	}

	toHitMin := item._iPLToHit + item._iPLToHitRange.Min
	toHitMax := item._iPLToHit + item._iPLToHitRange.Max
	toHit := fmt.Sprintf("%+d%%", toHitMin)
	if toHitMin != toHitMax {
		toHit = fmt.Sprintf("%+d%% to %+d%% (random)", toHitMin, toHitMax)
	}

	switch pIdx {
	case IPL_TOHIT:
		fallthrough
	case IPL_TOHIT_CURSE:
		return fmt.Sprintf("chance to hit : %s", toHit)
	case IPL_DAMP:
		fallthrough
	case IPL_DAMP_CURSE:
		return fmt.Sprintf("%+d%% damage", item._iPLDam)
	case IPL_TOHIT_DAMP:
		fallthrough
	case IPL_TOHIT_DAMP_CURSE:
		return fmt.Sprintf("to hit: %s, %+d%% damage", toHit, item._iPLDam)
	case IPL_ACP:
		fallthrough
	case IPL_ACP_CURSE:
		return fmt.Sprintf("%+d%% armor", item._iPLAC)
	case IPL_SETAC:
		return fmt.Sprintf("armor class: %d", item._iAC)
	case IPL_AC_CURSE:
		return fmt.Sprintf("armor class: %d", item._iAC)
	case IPL_FIRERES:
		if item._iPLFR < 75 {
			return fmt.Sprintf("Resist Fire : %+d%%", item._iPLFR)
		}
		if item._iPLFR >= 75 {
			return fmt.Sprintf("Resist Fire : 75%% MAX")
		}
	case IPL_LIGHTRES:
		if item._iPLLR < 75 {
			return fmt.Sprintf("Resist Lightning : %+d%%", item._iPLLR)
		}
		if item._iPLLR >= 75 {
			return fmt.Sprintf("Resist Lightning : 75%% MAX")
		}
	case IPL_MAGICRES:
		if item._iPLMR < 75 {
			return fmt.Sprintf("Resist Magic : %+d%%", item._iPLMR)
		}
		if item._iPLMR >= 75 {
			return fmt.Sprintf("Resist Magic : 75%% MAX")
		}
	case IPL_ALLRES:
		if item._iPLFR < 75 {
			return fmt.Sprintf("Resist All : %+d%%", item._iPLFR)
		}
		if item._iPLFR >= 75 {
			return fmt.Sprintf("Resist All : 75%% MAX")
		}
	case IPL_SPLLVLADD:
		if item._iSplLvlAdd == 1 {
			return "spells are increased 1 level"
		}
		if item._iSplLvlAdd == 2 {
			return "spells are increased 2 levels"
		}
		if item._iSplLvlAdd < 1 {
			return "spells are decreased 1 level"
		}
	case IPL_CHARGES:
		return "Extra charges"
	case IPL_SPELL:
		return fmt.Sprintf("%d charges [TODO: what spell]", item._iMaxCharges)
		// return fmt.Sprintf("%d %s charges", item._iMaxCharges, spelldata[item._iSpell].sNameText)
	case IPL_FIREDAM:
		return fmt.Sprintf("Fire hit damage: %d-%d", item._iFMinDam, item._iFMaxDam)
	case IPL_LIGHTDAM:
		return fmt.Sprintf("Lightning hit damage: %d-%d", item._iLMinDam, item._iLMaxDam)
	case IPL_STR:
		fallthrough
	case IPL_STR_CURSE:
		return fmt.Sprintf("%+d to strength", item._iPLStr)
	case IPL_MAG:
		fallthrough
	case IPL_MAG_CURSE:
		return fmt.Sprintf("%+d to magic", item._iPLMag)
	case IPL_DEX:
		fallthrough
	case IPL_DEX_CURSE:
		return fmt.Sprintf("%+d to dexterity", item._iPLDex)
	case IPL_VIT:
		fallthrough
	case IPL_VIT_CURSE:
		return fmt.Sprintf("%+d to vitality", item._iPLVit)
	case IPL_ATTRIBS:
		fallthrough
	case IPL_ATTRIBS_CURSE:
		return fmt.Sprintf("%+d to all attributes", item._iPLStr)
	case IPL_GETHIT_CURSE:
		fallthrough
	case IPL_GETHIT:
		return fmt.Sprintf("%+d damage from enemies", item._iPLGetHit)
	case IPL_LIFE:
		fallthrough
	case IPL_LIFE_CURSE:
		return fmt.Sprintf("Hit Points : %+d", item._iPLHP>>6)
	case IPL_MANA:
		fallthrough
	case IPL_MANA_CURSE:
		return fmt.Sprintf("Mana : %+d", item._iPLMana>>6)
	case IPL_DUR:
		return "high durability"
	case IPL_DUR_CURSE:
		return "decreased durability"
	case IPL_INDESTRUCTIBLE:
		return "indestructible"
	case IPL_LIGHT:
		return fmt.Sprintf("+%d%% light radius", 10*item._iPLLight)
	case IPL_LIGHT_CURSE:
		return fmt.Sprintf("-%d%% light radius", -10*item._iPLLight)
	case IPL_FIRE_ARROWS:
		return fmt.Sprintf("fire arrows damage: %d-%d", item._iFMinDam, item._iFMaxDam)
	case IPL_LIGHT_ARROWS:
		return fmt.Sprintf("lightning arrows damage %d-%d", item._iLMinDam, item._iLMaxDam)
	case IPL_THORNS:
		return "attacker takes 1-3 damage"
	case IPL_NOMANA:
		return "user loses all mana"
	case IPL_NOHEALPLR:
		return "you can't heal"
	case IPL_ABSHALFTRAP:
		return "absorbs half of trap damage"
	case IPL_KNOCKBACK:
		return "knocks target back"
	case IPL_3XDAMVDEM:
		return "+200% damage vs. demons"
	case IPL_ALLRESZERO:
		return "All Resistance equals 0"
	case IPL_NOHEALMON:
		return "hit monster doesn't heal"
	case IPL_STEALMANA:
		if item._iFlags&ISPL_STEALMANA_3 != 0 {
			return "hit steals 3% mana"
		}
		if item._iFlags&ISPL_STEALMANA_5 != 0 {
			return "hit steals 5% mana"
		}
	case IPL_STEALLIFE:
		if item._iFlags&ISPL_STEALLIFE_3 != 0 {
			return "hit steals 3% life"
		}
		if item._iFlags&ISPL_STEALLIFE_5 != 0 {
			return "hit steals 5% life"
		}
	case IPL_TARGAC:
		return "damages target's armor"
	case IPL_FASTATTACK:
		if item._iFlags&ISPL_QUICKATTACK != 0 {
			return "quick attack"
		}
		if item._iFlags&ISPL_FASTATTACK != 0 {
			return "fast attack"
		}
		if item._iFlags&ISPL_FASTERATTACK != 0 {
			return "faster attack"
		}
		if item._iFlags&ISPL_FASTESTATTACK != 0 {
			return "fastest attack"
		}
	case IPL_FASTRECOVER:
		if item._iFlags&ISPL_FASTRECOVER != 0 {
			return "fast hit recovery"
		}
		if item._iFlags&ISPL_FASTERRECOVER != 0 {
			return "faster hit recovery"
		}
		if item._iFlags&ISPL_FASTESTRECOVER != 0 {
			return "fastest hit recovery"
		}
	case IPL_FASTBLOCK:
		return "fast block"
	case IPL_DAMMOD:
		return fmt.Sprintf("adds %d points to damage", item._iPLDamMod)
	case IPL_RNDARROWVEL:
		return "fires random speed arrows"
	case IPL_SETDAM:
		return fmt.Sprintf("unusual item damage")
	case IPL_SETDUR:
		return "altered durability"
	case IPL_FASTSWING:
		return "Faster attack swing"
	case IPL_ONEHAND:
		return "one handed sword"
	case IPL_DRAINLIFE:
		return "constantly lose hit points"
	case IPL_RNDSTEALLIFE:
		return "life stealing"
	case IPL_NOMINSTR:
		return "no strength requirement"
	case IPL_INFRAVISION:
		return "see with infravision"
	case IPL_INVCURS:
		return " "
	case IPL_ADDACLIFE:
		return "Armor class added to life"
	case IPL_ADDMANAAC:
		return "10% of mana added to armor"
	case IPL_FIRERESCLVL:
		if item._iPLFR <= 0 {
			return fmt.Sprintf(" ")
		} else if item._iPLFR >= 1 {
			return fmt.Sprintf("Resist Fire : %+d%%", item._iPLFR)
		}
	}
	return "Another ability (NW)"
}

// const ITEM_START uint32 = 0x91180

const (
	ITEM_START           uint32 = 0x8efa8
	ITEM_END             uint32 = 0x91df8
	PREFIX_START         uint32 = 0x7b688
	PREFIX_END           uint32 = 0x7c618
	SUFFIX_START         uint32 = 0x7c648
	SUFFIX_END           uint32 = 0x7d818
	MONSTER_START        uint32 = 0x977a8
	MONSTER_END          uint32 = 0x9afa8
	UNIQUE_MONSTER_START uint32 = 0x9b098
	UNIQUE_MONSTER_END   uint32 = 0x9bcb8
	// ITEM_START uint32 = 0x91134
	// ITEM_END          = 0x93f84
)

func ItemFlags(item *StructVal) uint32 {
	switch item.Get("itype").(byte) {
	case ITYPE_SWORD:
		fallthrough
	case ITYPE_AXE:
		fallthrough
	case ITYPE_MACE:
		return 0x1000
	case ITYPE_BOW:
		return 0x10
	case ITYPE_SHIELD:
		return 0x10000
	case ITYPE_LARMOR:
		fallthrough
	case ITYPE_HELM:
		fallthrough
	case ITYPE_MARMOR:
		fallthrough
	case ITYPE_HARMOR:
		return 0x100000
	case ITYPE_STAFF:
		return 0x100 // TODO: staff powers?
	case ITYPE_RING:
		fallthrough
	case ITYPE_AMULET:
		return 1
	}
	return 0
}

var LOCS = []string{
	"None",
	"One-hand",
	"Two-hand",
	"Armor",
	"Helm",
	"Ring",
	"Amulet",
	"Unequipable",
	"Belt",
}

// const ITEM_START = 0xe734

func main() {
	const exePath = "data/Crucible.exe"
	f, err := pe.Open(exePath)
	if err != nil {
		log.Fatalf("Error opening exe: %s", err)
	}
	baseAddr := f.OptionalHeader.(*pe.OptionalHeader32).ImageBase
	dataSec := f.Section(".data")
	if err != nil {
		log.Fatalf("Error reading data: %s", err)
	}
	f.Close()
	exeBytes, err := ioutil.ReadFile(exePath)
	if err != nil {
		log.Fatalf("Error reading executable: %s", err)
	}
	exe := &Exe{
		DataOffset: baseAddr + dataSec.VirtualAddress - dataSec.Offset,
		Bytes:      exeBytes,
	}

	// Searching code below here:
	// monsterSize := MonsterData.Size()
	// for offset := uint32(0); offset+monsterSize <= uint32(len(exeBytes)); offset++ {
	// 	monster := &StructVal{
	// 		Exe:    exe,
	// 		Type:   MonsterData,
	// 		Offset: offset,
	// 	}
	// 	if monster.Get("GraphicType").(string) == "Monsters\\Zombie\\Zombie%c.CL2" {
	// 		log.Printf("%#x %s", offset, monster.Get("mName").(string))
	// 	}
	// }
	// os.Exit(1)
	// uniqueMonsterSize := UniqueMonsterData.Size()
	// for offset := uint32(0); offset+uniqueMonsterSize <= uint32(len(exeBytes)); offset++ {
	// 	monster := &StructVal{
	// 		Exe:    exe,
	// 		Type:   UniqueMonsterData,
	// 		Offset: offset,
	// 	}
	// 	if monster.Get("mName").(string) == "Gharbad the Weak" {
	// 		log.Printf("%#x %s", offset, monster.Get("mName").(string))
	// 	}
	// }
	// uniqueMonsterSize := UniqueMonsterData.Size()
	// for offset := UNIQUE_MONSTER_START; offset < UNIQUE_MONSTER_END; offset += uniqueMonsterSize {
	// 	log.Printf("%#x", offset)
	// 	item := &StructVal{
	// 		Exe:    exe,
	// 		Type:   UniqueMonsterData,
	// 		Offset: offset,
	// 	}
	// 	for _, f := range *UniqueMonsterData {
	// 		log.Printf("  %s: %v", f.Name, item.Get(f.Name))
	// 	}
	// }
	// os.Exit(1)
	// os.Exit(1)
	// monsterSize := MonsterData.Size()
	// for offset := MONSTER_START; offset < MONSTER_END; offset += monsterSize {
	// 	log.Printf("%#x", offset)
	// 	item := &StructVal{
	// 		Exe:    exe,
	// 		Type:   MonsterData,
	// 		Offset: offset,
	// 	}
	// 	for _, f := range *MonsterData {
	// 		log.Printf("  %s: %v", f.Name, item.Get(f.Name))
	// 	}
	// }

	// itemSize := ItemData.Size()
	// for offset := ITEM_START; offset < ITEM_END; offset += itemSize {
	// 	log.Printf("%#x", offset)
	// 	item := &StructVal{
	// 		Exe:    exe,
	// 		Type:   ItemData,
	// 		Offset: offset,
	// 	}
	// 	for _, f := range *ItemData {
	// 		log.Printf("  %s: %v", f.Name, item.Get(f.Name))
	// 	}
	// }

	monsters := []*StructVal{}

	monsterSize := MonsterData.Size()
	mpqDir := "./data/diabdat"
	palPath := filepath.Join(mpqDir, "levels/towndata/town.pal")
	pal, err := cel.ParsePal(palPath)
	if err != nil {
		panic(err)
	}
	_ = pal
	fmt.Println("window.MONSTERS = [")
	id := 0
	monsterClasses := map[uint8]string{
		0: "Undead",
		1: "Demon",
		2: "Animal",
	}
	for offset := MONSTER_START; offset < MONSTER_END; offset += monsterSize {
		id++
		monster := &StructVal{
			Exe:    exe,
			Type:   MonsterData,
			Offset: offset,
		}
		monsters = append(monsters, monster)
		minDamage := monster.Get("mMinDamage").(uint8)
		maxDamage := monster.Get("mMaxDamage").(uint8)
		damage := fmt.Sprintf("%d", minDamage)
		if minDamage != maxDamage {
			damage = fmt.Sprintf("%d–%d", minDamage, maxDamage)
		}
		minDamage2 := monster.Get("mMinDamage2").(uint8)
		maxDamage2 := monster.Get("mMaxDamage2").(uint8)
		damage2 := fmt.Sprintf("%d", minDamage2)
		if maxDamage2 == 0 {
			damage2 = ""
		}
		if minDamage2 != maxDamage2 {
			damage2 = fmt.Sprintf("%d–%d", minDamage2, maxDamage2)
		}
		minHP := monster.Get("mMinHP").(uint32)
		maxHP := monster.Get("mMaxHP").(uint32)
		hp := fmt.Sprintf("%d", minHP)
		if minHP != maxHP {
			hp = fmt.Sprintf("%d–%d", minHP, maxHP)
		}
		// dstDir := fmt.Sprintf("/tmp/images/%d", id)
		// os.RemoveAll(dstDir)
		// if err := os.MkdirAll(dstDir, 0755); err != nil {
		// 	panic(err)
		// }
		// pngs := []string{}
		// anims := "nwahds"
		// seq := []int{1, 2, 1, 4}
		// if monster.Get("has_special").(bool) {
		// 	seq = []int{1, 2, 1, 5, 1, 4}
		// }
		// for _, animIdx := range seq {
		// 	animType := string(anims[animIdx])
		// 	if monster.Get(fmt.Sprintf("Frames[%d]", animIdx)).(uint32) <= 1 {
		// 		log.Printf("skipping %s for %s", animType, monster.Get("mName"))
		// 		continue
		// 	}
		// 	celPath := filepath.Join(mpqDir, strings.Replace(
		// 		strings.ToLower(strings.Replace(monster.Get("GraphicType").(string), "%c", animType, -1)),
		// 		"\\", "/", -1))
		// 	if celPath == "data/diabdat/monsters/golem/golemd.cl2" {
		// 		log.Printf("skipping %s", celPath)
		// 		continue
		// 	}
		// 	if celPath == "data/diabdat/monsters/golem/golems.cl2" {
		// 		log.Printf("skipping %s", celPath)
		// 		continue
		// 	}
		// 	log.Printf("decoding %s", celPath)
		// 	thisPal := pal
		// 	trnPath := monster.Get("TransFile").(string)
		// 	if trnPath != "" {
		// 		trnPath := filepath.Join(mpqDir, strings.ToLower(strings.Replace(trnPath, "\\", "/", -1)))
		// 		trn, err := cel.ParseTrn(trnPath)
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		thisPal = trn.Pal(pal)
		// 	}
		// 	imgs, err := cel.DecodeArchive(celPath, thisPal)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	for _, img := range imgs[7] {
		// 		pngName := fmt.Sprintf("%04d.png", len(pngs)+1)
		// 		pngPath := filepath.Join(dstDir, pngName)
		// 		if err := imgutil.WriteFile(pngPath, img); err != nil {
		// 			panic(err)
		// 		}
		// 		pngs = append(pngs, pngPath)
		// 	}
		// }
		// args := []string{
		// 	"-delay",
		// 	"5",
		// 	"-loop",
		// 	"0",
		// 	"-dispose",
		// 	"previous",
		// }
		// args = append(args, pngs...)
		// args = append(args, fmt.Sprintf("html/monsters/%d.gif", id))
		// cmd := exec.Command("/usr/bin/convert", args...)
		// out, err := cmd.CombinedOutput()
		// if err != nil {
		// 	panic(out)
		// }
		data := map[string]interface{}{
			"ID":         id,
			"Name":       monster.Get("mName"),
			"MinDLvl":    15*int32(int8(monster.Get("mMinDLvl").(uint8)))/30 + 1,
			"MaxDLvl":    15*int32(int8(monster.Get("mMaxDLvl").(uint8)))/30 + 1,
			"Damage":     damage,
			"Damage2":    damage2,
			"HP":         hp,
			"Experience": monster.Get("mExp"),
			"Resistance": monster.Get("mMagicRes"),
			"AC":         monster.Get("mArmorClass"),
			"Type":       monsterClasses[monster.Get("mMonstClass").(uint8)],
			"mLevel":     monster.Get("mLevel"),
		}
		buf, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Printf("  %s,\n", buf)
	}
	fmt.Println("];")

	fmt.Println("window.UNIQUE_MONSTERS = [")
	uniqueMonsterSize := UniqueMonsterData.Size()
	uid := 0
	for offset := UNIQUE_MONSTER_START; offset < UNIQUE_MONSTER_END; offset += uniqueMonsterSize {
		uid++
		unique := &StructVal{
			Exe:    exe,
			Type:   UniqueMonsterData,
			Offset: offset,
		}
		monsterType := unique.Get("mtype").(uint8)
		if monsterType == 255 {
			continue
		}
		baseMonster := monsters[unique.Get("mtype").(uint8)]

		minDamage := unique.Get("mMinDamage").(uint8)
		maxDamage := unique.Get("mMaxDamage").(uint8)
		damage := fmt.Sprintf("%d", minDamage)
		if minDamage != maxDamage {
			damage = fmt.Sprintf("%d–%d", minDamage, maxDamage)
		}

		hp := unique.Get("mmaxhp").(uint16)
		if hp < 64 {
			hp = 64
		}

		attr := unique.Get("mUnqAttr").(uint16)
		v1 := unique.Get("mUnqVar1").(uint8)

		ac := baseMonster.Get("mArmorClass").(uint8)
		if attr&8 > 0 {
			ac = v1
		}

		mLevel := unique.Get("mlevel").(byte)
		if mLevel == 0 {
			mLevel = baseMonster.Get("mLevel").(byte) + 5
		}

		data := map[string]interface{}{
			"ID":         uid,
			"Name":       unique.Get("mName"),
			"Damage":     damage,
			"Resistance": unique.Get("mMagicRes"),
			"HP":         hp,
			"Experience": baseMonster.Get("mExp").(uint16) * 2,
			"AC":         ac,
			"Type":       monsterClasses[baseMonster.Get("mMonstClass").(uint8)],
			"mLevel":     mLevel,
		}
		buf, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Printf("  %s,\n", buf)

		// dstDir := fmt.Sprintf("/tmp/images/%d", id)
		// os.RemoveAll(dstDir)
		// if err := os.MkdirAll(dstDir, 0755); err != nil {
		// 	panic(err)
		// }
		// pngs := []string{}
		// anims := "nwahds"
		// for _, animIdx := range []int{1, 2, 1, 4} {
		// 	animType := string(anims[animIdx])
		// 	if baseMonster.Get(fmt.Sprintf("Frames[%d]", animIdx)).(uint32) <= 1 {
		// 		log.Printf("skipping %s for %s", animType, baseMonster.Get("mName"))
		// 		continue
		// 	}
		// 	celPath := filepath.Join(mpqDir, strings.Replace(
		// 		strings.ToLower(strings.Replace(baseMonster.Get("GraphicType").(string), "%c", animType, -1)),
		// 		"\\", "/", -1))
		// 	if celPath == "data/diabdat/monsters/golem/golemd.cl2" {
		// 		log.Printf("skipping %s", celPath)
		// 		continue
		// 	}
		// 	log.Printf("decoding %s", celPath)
		// 	trnPath := fmt.Sprintf("Monsters\\Monsters\\%s.TRN", unique.Get("mTrnName").(string))
		// 	trnPath = filepath.Join(mpqDir, strings.ToLower(strings.Replace(trnPath, "\\", "/", -1)))
		// 	trn, err := cel.ParseTrn(trnPath)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	thisPal := trn.Pal(pal)
		// 	imgs, err := cel.DecodeArchive(celPath, thisPal)
		// 	if err != nil {
		// 		panic(err)
		// 	}

		// 	for _, img := range imgs[7] {
		// 		pngName := fmt.Sprintf("%04d.png", len(pngs)+1)
		// 		pngPath := filepath.Join(dstDir, pngName)
		// 		if err := imgutil.WriteFile(pngPath, img); err != nil {
		// 			panic(err)
		// 		}
		// 		pngs = append(pngs, pngPath)
		// 	}
		// }
		// args := []string{
		// 	"-delay",
		// 	"5",
		// 	"-loop",
		// 	"0",
		// 	"-dispose",
		// 	"previous",
		// }
		// args = append(args, pngs...)
		// args = append(args, fmt.Sprintf("html/unique-monsters/%d.gif", uid))
		// cmd := exec.Command("/usr/bin/convert", args...)
		// out, err := cmd.CombinedOutput()
		// if err != nil {
		// 	panic(out)
		// }
	}
	fmt.Println("];")

	prefixSize := PowerData.Size()
	fmt.Println("window.PREFIXES = [")
	for offset := PREFIX_START; offset < PREFIX_END; offset += prefixSize {
		power := &StructVal{
			Exe:    exe,
			Type:   PowerData,
			Offset: offset,
		}
		descriptions := []string{}
		for r := power.Get("PLParam1").(uint32); r <= power.Get("PLParam2").(uint32); r++ {
			descriptions = append(descriptions, DescribePower(power, int(r), 10))
		}
		data := map[string]interface{}{
			"Mod":             power.Get("PLPower"),
			"Name":            power.Get("PLName"),
			"ItemType":        power.Get("PLIType"),
			"MinParam":        power.Get("PLParam1"),
			"MaxParam":        power.Get("PLParam2"),
			"MinValue":        power.Get("PLMinVal"),
			"MaxValue":        power.Get("PLMaxVal"),
			"ValueMultiplier": power.Get("PLMultVal"),
			"Descriptions":    descriptions,
		}
		buf, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Printf("  %s,\n", buf)
	}
	fmt.Println("];")

	fmt.Println("window.SUFFIXES = [")
	for offset := SUFFIX_START; offset < SUFFIX_END; offset += prefixSize {
		power := &StructVal{
			Exe:    exe,
			Type:   PowerData,
			Offset: offset,
		}
		descriptions := []string{}
		for r := power.Get("PLParam1").(uint32); r <= power.Get("PLParam2").(uint32); r++ {
			descriptions = append(descriptions, DescribePower(power, int(r), 10))
		}
		data := map[string]interface{}{
			"Mod":             power.Get("PLPower"),
			"Name":            power.Get("PLName"),
			"ItemType":        power.Get("PLIType"),
			"MinParam":        power.Get("PLParam1"),
			"MaxParam":        power.Get("PLParam2"),
			"MinValue":        power.Get("PLMinVal"),
			"MaxValue":        power.Get("PLMaxVal"),
			"ValueMultiplier": power.Get("PLMultVal"),
			"Descriptions":    descriptions,
		}
		buf, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Printf("  %s,\n", buf)
	}
	fmt.Println("];")

	itemSize := ItemData.Size()
	fmt.Println("window.ITEMS = [")
	for offset := ITEM_START; offset < ITEM_END; offset += itemSize {
		item := &StructVal{
			Exe:    exe,
			Type:   ItemData,
			Offset: offset,
		}
		if item.Get("iLoc").(byte) == 0 || item.Get("iRnd").(uint32) == 0 {
			continue
		}
		minArmor := item.Get("iMinAC").(uint32)
		maxArmor := item.Get("iMaxAC").(uint32)
		armor := fmt.Sprintf("%d", minArmor)
		if minArmor != maxArmor {
			armor = fmt.Sprintf("%d–%d", minArmor, maxArmor)
		}
		minDamage := item.Get("iMinDam").(uint32)
		maxDamage := item.Get("iMaxDam").(uint32)
		damage := fmt.Sprintf("%d", minDamage)
		if minDamage != maxDamage {
			damage = fmt.Sprintf("%d–%d", minDamage, maxDamage)
		}
		data := map[string]interface{}{
			"Armor":      armor,
			"Damage":     damage,
			"Dexterity":  item.Get("iMinDex"),
			"Durability": item.Get("iDurability"),
			"Flags":      ItemFlags(item),
			"Loc":        LOCS[item.Get("iLoc").(byte)],
			"Magic":      item.Get("iMinMag"),
			"MaxDamage":  item.Get("iMaxDam"),
			"MinDamage":  item.Get("iMinDam"),
			"Name":       item.Get("iName"),
			"ShortName":  item.Get("iSName"),
			"Strength":   item.Get("iMinStr"),
			"Value":      item.Get("iValue"),
			"Icon":       item.Get("iCurs"),
		}
		buf, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		fmt.Printf("  %s,\n", buf)
	}
	fmt.Println("];")

	// Searching code below here:
	// itemSize := ItemData.Size()
	// for offset := uint32(0); offset+itemSize <= uint32(len(exeBytes)); offset++ {
	// 	item := &StructVal{
	// 		Exe:    exe,
	// 		Type:   ItemData,
	// 		Offset: offset,
	// 	}
	// 	if item.Get("iMinDam").(uint32) == 2 && item.Get("iMaxDam").(uint32) == 5 {
	// 		log.Printf("%#x %s", offset, item.Get("iName").(string))
	// 	}
	// }
	// prefixSize := PowerData.Size()
	// for offset := uint32(0); offset+prefixSize <= uint32(len(exeBytes)); offset++ {
	// 	item := &StructVal{
	// 		Exe:    exe,
	// 		Type:   PowerData,
	// 		Offset: offset,
	// 	}
	// 	if item.Get("PLParam1").(uint32) == 1 && item.Get("PLParam2").(uint32) == 2 {
	// 		log.Printf("%x %s", offset, item.Get("PLName").(string))
	// 	}
	// }
}
