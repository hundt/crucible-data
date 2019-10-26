package main

import (
	"encoding/binary"
	"log"
)

type Primitive int

const (
	PRIMITIVE_INT Primitive = iota
	PRIMITIVE_CHAR
	PRIMITIVE_STRING
	PRIMITIVE_BOOL
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
	case PRIMITIVE_BOOL:
		return binary.LittleEndian.Uint32(b) != 0
	case PRIMITIVE_STRING:
		addr := binary.LittleEndian.Uint32(b)
		if addr == 0 {
			return ""
		}
		addr -= exe.DataOffset
		if addr > uint32(len(exe.Bytes)) {
			return "INVALID"
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

func (s *Struct) Size() uint32 {
	offs := uint32(0)
	for _, f := range *s {
		for offs%f.Size() != 0 {
			offs++
		}
		offs += f.Size()
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
