package jmp

import (
	"encoding/json"
	"mhfrc/util/byteframe"
	"mhfrc/util/sjis"
)

type Jmp struct {
	Jumps   []Jump
	Menus   []Menu
	Strings []string
}

type Jump struct {
	Id           uint32
	Unk          uint32            // icon or colour
	StageIds     []uint16          // 4
	Destinations []JumpDestination // 2
	Title        string
	Description  string
}

type JumpDestination struct {
	Coordinates []float32
	Rotation    uint32
}

type Menu struct {
	Entries  []MenuEntry
	StageIds []uint16
}

type MenuEntry struct {
	Index uint16
	Flags uint16
}

func DecompileJmp(jmpData []byte) ([]byte, error) {
	bf := byteframe.NewByteFrameFromBytes(jmpData)
	bf.SetLE()
	ptrJumps := bf.ReadUint32()
	ptrMenus := bf.ReadUint32()
	lenMenus := bf.ReadUint32()
	ptrStrings := bf.ReadUint32()
	lenStrings := bf.ReadUint32()

	jmp := Jmp{}

	_, err := bf.Seek(int64(ptrJumps), 0)
	if err != nil {
		return []byte{}, err
	} else {
		jmp.Jumps = make([]Jump, 24)
		for i := 0; i < 24; i++ {
			jmp.Jumps[i].Id = bf.ReadUint32()
			jmp.Jumps[i].Unk = bf.ReadUint32()
			jmp.Jumps[i].StageIds = make([]uint16, 4)
			for j := 0; j < 4; j++ {
				jmp.Jumps[i].StageIds[j] = bf.ReadUint16()
			}
			jmp.Jumps[i].Destinations = make([]JumpDestination, 2)
			for j := 0; j < 2; j++ {
				jmp.Jumps[i].Destinations[j].Coordinates = make([]float32, 3)
				for k := 0; k < 3; k++ {
					jmp.Jumps[i].Destinations[j].Coordinates[k] = bf.ReadFloat32()
				}
				jmp.Jumps[i].Destinations[j].Rotation = bf.ReadUint32()
			}
			ptrTitle := bf.ReadUint32()
			ptrDescription := bf.ReadUint32()
			_, _ = bf.Seek(int64(ptrTitle), 0)
			strTitle := sjis.NewBytes(bf.ReadNullTerminatedBytes())
			jmp.Jumps[i].Title = strTitle.String()
			_, _ = bf.Seek(int64(ptrDescription), 0)
			strDescription := sjis.NewBytes(bf.ReadNullTerminatedBytes())
			jmp.Jumps[i].Description = strDescription.String()
			_, _ = bf.Seek(int64(ptrJumps)+int64(i*56), 0)
		}
	}

	_, err = bf.Seek(int64(ptrMenus), 0)
	if err != nil {
		return []byte{}, err
	} else {
		jmp.Menus = make([]Menu, lenMenus)
		for i := 0; i < int(lenMenus); i++ {
			ptrMenuEntries := bf.ReadUint32()
			lenMenuEntries := bf.ReadUint32()
			ptrStageIds := bf.ReadUint32()
			jmp.Menus[i].Entries = make([]MenuEntry, lenMenuEntries)
			_, _ = bf.Seek(int64(ptrMenuEntries), 0)
			for j := 0; j < int(lenMenuEntries); j++ {
				jmp.Menus[i].Entries[j].Index = bf.ReadUint16()
				jmp.Menus[i].Entries[j].Flags = bf.ReadUint16()
			}
			_, _ = bf.Seek(int64(ptrStageIds), 0)
			for {
				stageId := bf.ReadUint16()
				if stageId == 0 {
					break
				}
				jmp.Menus[i].StageIds = append(jmp.Menus[i].StageIds, stageId)
			}
			_, _ = bf.Seek(int64(ptrMenus)+int64(i*12), 0)
		}
	}

	_, err = bf.Seek(int64(ptrStrings), 0)
	if err != nil {
		return []byte{}, err
	} else {
		jmp.Strings = make([]string, lenStrings)
		for i := 0; i < int(lenStrings); i++ {
			ptrString := bf.ReadUint32()
			_, _ = bf.Seek(int64(ptrString), 0)
			strString := sjis.NewBytes(bf.ReadNullTerminatedBytes())
			jmp.Strings[i] = strString.String()
			_, _ = bf.Seek(int64(ptrStrings)+int64(i*4), 0)
		}
	}

	jmpJson, err := json.MarshalIndent(jmp, "", "\t")
	if err != nil {
		return []byte{}, err
	} else {
		return jmpJson, nil
	}
}
