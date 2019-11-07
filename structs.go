package main

import (
	"encoding/binary"
	"fmt"
	"log"
)

type Primitive int

const (
	PRIMITIVE_INT Primitive = iota
	PRIMITIVE_CHAR
	PRIMITIVE_STRING
	PRIMITIVE_BOOL
	PRIMITIVE_USHORT
	PRIMITIVE_UCHAR = PRIMITIVE_CHAR
)

type Field struct {
	Type Primitive
	Name string
}

func CStr(b []byte) string {
	idx := 0
	for ; idx < len(b) && b[idx] != 0; idx++ {
	}
	return string(b[:idx])
}

func (f Field) Extract(b []byte, exe *Exe) interface{} {
	switch f.Type {
	case PRIMITIVE_INT:
		return binary.LittleEndian.Uint32(b)
	case PRIMITIVE_CHAR:
		return b[0]
	case PRIMITIVE_USHORT:
		return binary.LittleEndian.Uint16(b)
	case PRIMITIVE_BOOL:
		return binary.LittleEndian.Uint32(b) != 0
	case PRIMITIVE_STRING:
		addr := binary.LittleEndian.Uint32(b)
		if addr == 0 {
			return ""
		}
		addr -= exe.DataOffset
		if addr > uint32(len(exe.Bytes)) {
			return fmt.Sprintf("INVALID (%x)", binary.LittleEndian.Uint32(b))
		}
		return CStr(exe.Bytes[addr:])
	}
	log.Fatalf("Invalid type %#v", f.Type)
	return nil
}

func (f Field) Size() uint32 {
	switch f.Type {
	case PRIMITIVE_INT:
		return 4
	case PRIMITIVE_CHAR:
		return 1
	case PRIMITIVE_USHORT:
		return 2
	case PRIMITIVE_STRING:
		return 4
	case PRIMITIVE_BOOL:
		return 4
	}
	log.Fatalf("Invalid type %#v", f.Type)
	return 0
}

type Struct []Field

var ItemData = &Struct{
	Field{PRIMITIVE_INT, "iRnd"},
	Field{PRIMITIVE_CHAR, "iClass"},
	Field{PRIMITIVE_CHAR, "iLoc"},
	Field{PRIMITIVE_INT, "iCurs"},
	Field{PRIMITIVE_CHAR, "itype"},
	Field{PRIMITIVE_CHAR, "iItemId"},
	Field{PRIMITIVE_STRING, "iName"},
	Field{PRIMITIVE_STRING, "iSName"},
	Field{PRIMITIVE_CHAR, "iMinMLvl"},
	Field{PRIMITIVE_INT, "iDurability"},
	Field{PRIMITIVE_INT, "iMinDam"},
	Field{PRIMITIVE_INT, "iMaxDam"},
	Field{PRIMITIVE_INT, "iMinAC"},
	Field{PRIMITIVE_INT, "iMaxAC"},
	Field{PRIMITIVE_CHAR, "iMinStr"},
	Field{PRIMITIVE_CHAR, "iMinMag"},
	Field{PRIMITIVE_CHAR, "iMinDex"},
	// item_special_effect
	Field{PRIMITIVE_INT, "iFlags"},
	// item_misc_id
	Field{PRIMITIVE_INT, "iMiscId"},
	// spell_id
	Field{PRIMITIVE_INT, "iSpell"},
	Field{PRIMITIVE_BOOL, "iUsable"},
	Field{PRIMITIVE_INT, "iValue"},
	Field{PRIMITIVE_INT, "iMaxValue"},
}

var PowerData = &Struct{
	Field{PRIMITIVE_STRING, "PLName"},
	Field{PRIMITIVE_INT, "PLPower"},
	Field{PRIMITIVE_INT, "PLParam1"},
	Field{PRIMITIVE_INT, "PLParam2"},
	Field{PRIMITIVE_CHAR, "PLMinLvl"},
	Field{PRIMITIVE_INT, "PLIType"},
	Field{PRIMITIVE_INT, "PLGOE"},
	Field{PRIMITIVE_BOOL, "PLDouble"},
	Field{PRIMITIVE_BOOL, "PLOk"},
	Field{PRIMITIVE_INT, "PLMinVal"},
	Field{PRIMITIVE_INT, "PLMaxVal"},
	Field{PRIMITIVE_INT, "PLMultVal"},
}

var MonsterData = &Struct{
	Field{PRIMITIVE_INT, "width"},
	Field{PRIMITIVE_INT, "mImage"},
	Field{PRIMITIVE_STRING, "GraphicType"},
	Field{PRIMITIVE_BOOL, "has_special"},
	Field{PRIMITIVE_STRING, "sndfile"},
	Field{PRIMITIVE_BOOL, "snd_special"},
	Field{PRIMITIVE_BOOL, "has_trans"},
	Field{PRIMITIVE_STRING, "TransFile"},
	Field{PRIMITIVE_INT, "Frames[0]"},
	Field{PRIMITIVE_INT, "Frames[1]"},
	Field{PRIMITIVE_INT, "Frames[2]"},
	Field{PRIMITIVE_INT, "Frames[3]"},
	Field{PRIMITIVE_INT, "Frames[4]"},
	Field{PRIMITIVE_INT, "Frames[5]"},
	Field{PRIMITIVE_INT, "Rate[0]"},
	Field{PRIMITIVE_INT, "Rate[1]"},
	Field{PRIMITIVE_INT, "Rate[2]"},
	Field{PRIMITIVE_INT, "Rate[3]"},
	Field{PRIMITIVE_INT, "Rate[4]"},
	Field{PRIMITIVE_INT, "Rate[5]"},
	Field{PRIMITIVE_STRING, "mName"},
	Field{PRIMITIVE_CHAR, "mMinDLvl"},
	Field{PRIMITIVE_CHAR, "mMaxDLvl"},
	Field{PRIMITIVE_CHAR, "mLevel"},
	Field{PRIMITIVE_INT, "mMinHP"},
	Field{PRIMITIVE_INT, "mMaxHP"},
	Field{PRIMITIVE_CHAR, "mAi"},
	Field{PRIMITIVE_INT, "mFlags"},
	Field{PRIMITIVE_UCHAR, "mInt"},
	Field{PRIMITIVE_UCHAR, "mHit"}, // BUGFIX: Some monsters overflow this value on high difficulty
	Field{PRIMITIVE_UCHAR, "mAFNum"},
	Field{PRIMITIVE_UCHAR, "mMinDamage"},
	Field{PRIMITIVE_UCHAR, "mMaxDamage"},
	Field{PRIMITIVE_UCHAR, "mHit2"}, // BUGFIX: Some monsters overflow this value on high difficulty
	Field{PRIMITIVE_UCHAR, "mAFNum2"},
	Field{PRIMITIVE_UCHAR, "mMinDamage2"},
	Field{PRIMITIVE_UCHAR, "mMaxDamage2"},
	Field{PRIMITIVE_UCHAR, "mArmorClass"},
	Field{PRIMITIVE_CHAR, "mMonstClass"},
	Field{PRIMITIVE_USHORT, "mMagicRes"},
	Field{PRIMITIVE_USHORT, "mMagicRes2"},
	Field{PRIMITIVE_USHORT, "mTreasure"},
	Field{PRIMITIVE_CHAR, "mSelFlag"},
	Field{PRIMITIVE_USHORT, "mExp"},
}

var UniqueMonsterData = &Struct{
	Field{PRIMITIVE_CHAR, "mtype"},
	Field{PRIMITIVE_STRING, "mName"},
	Field{PRIMITIVE_STRING, "mTrnName"},
	Field{PRIMITIVE_UCHAR, "mlevel"},
	Field{PRIMITIVE_USHORT, "mmaxhp"},
	Field{PRIMITIVE_UCHAR, "mAi"},
	Field{PRIMITIVE_UCHAR, "mint"},
	Field{PRIMITIVE_UCHAR, "mMinDamage"},
	Field{PRIMITIVE_UCHAR, "mMaxDamage"},
	Field{PRIMITIVE_USHORT, "mMagicRes"},
	Field{PRIMITIVE_USHORT, "mUnqAttr"},
	Field{PRIMITIVE_UCHAR, "mUnqVar1"},
	Field{PRIMITIVE_UCHAR, "mUnqVar2"},
	Field{PRIMITIVE_INT, "mtalkmsg"},
}

func (s *Struct) Size() uint32 {
	offs := uint32(0)
	for _, f := range *s {
		for offs%f.Size() != 0 {
			offs++
		}
		offs += f.Size()
	}
	for offs%4 != 0 {
		offs++
	}
	return offs
}

type Exe struct {
	DataOffset uint32
	Bytes      []byte
}

type StructVal struct {
	Type   *Struct
	Exe    *Exe
	Offset uint32
}

func (v *StructVal) Get(field string) interface{} {
	offs := uint32(0)
	for _, f := range *v.Type {
		for offs%f.Size() != 0 {
			offs++
		}
		if f.Name == field {
			return f.Extract(v.Exe.Bytes[v.Offset+offs:], v.Exe)
		}
		offs += f.Size()
	}
	log.Fatalf("Invalid field %s", field)
	return nil
}
