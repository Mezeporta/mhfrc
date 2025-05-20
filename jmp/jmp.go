package jmp

import (
	"encoding/json"
	"errors"
	"mhfrc/util/byteframe"
	"mhfrc/util/sjis"
)

type Jmp struct {
	Jumps      []Jump
	Menus      []Menu
	Strings    []string
	ptrStrings []uint32

	ptrJumps      uint32
	ptrMenus      uint32
	ptrPtrStrings uint32
}

type Jump struct {
	Id             uint32
	Unk            uint32            // icon or colour
	StageIds       []uint16          // 4
	Destinations   []JumpDestination // 2
	Title          string
	Description    string
	ptrTitle       uint32
	ptrDescription uint32
}

type JumpDestination struct {
	Coordinates []float32
	Rotation    uint32
}

type Menu struct {
	Entries     []MenuEntry
	StageIds    []uint16
	ptrEntries  uint32
	ptrStageIds uint32
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
			for j := range jmp.Jumps[i].StageIds {
				jmp.Jumps[i].StageIds[j] = bf.ReadUint16()
			}
			jmp.Jumps[i].Destinations = make([]JumpDestination, 2)
			for j := range jmp.Jumps[i].Destinations {
				jmp.Jumps[i].Destinations[j].Coordinates = make([]float32, 3)
				for k := range jmp.Jumps[i].Destinations[j].Coordinates {
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
			_, _ = bf.Seek(int64(ptrJumps)+(int64(i+1)*56), 0)
		}
	}

	_, err = bf.Seek(int64(ptrMenus), 0)
	if err != nil {
		return []byte{}, err
	} else {
		jmp.Menus = make([]Menu, lenMenus)
		for i := range jmp.Menus {
			ptrMenuEntries := bf.ReadUint32()
			lenMenuEntries := bf.ReadUint32()
			ptrStageIds := bf.ReadUint32()
			jmp.Menus[i].Entries = make([]MenuEntry, lenMenuEntries)
			_, _ = bf.Seek(int64(ptrMenuEntries), 0)
			for j := range jmp.Menus[i].Entries {
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
			_, _ = bf.Seek(int64(ptrMenus)+(int64(i+1)*12), 0)
		}
	}

	_, err = bf.Seek(int64(ptrStrings), 0)
	if err != nil {
		return []byte{}, err
	} else {
		jmp.Strings = make([]string, lenStrings)
		for i := range jmp.Strings {
			ptrString := bf.ReadUint32()
			_, _ = bf.Seek(int64(ptrString), 0)
			strString := sjis.NewBytes(bf.ReadNullTerminatedBytes())
			jmp.Strings[i] = strString.String()
			_, _ = bf.Seek(int64(ptrStrings)+(int64(i+1)*4), 0)
		}
	}

	jmpJson, err := json.MarshalIndent(jmp, "", "\t")
	if err != nil {
		return []byte{}, err
	} else {
		return jmpJson, nil
	}
}

func CompileJmp(jmpJson []byte) ([]byte, error) {
	strings := make(map[string]uint32)

	jmp := Jmp{}
	err := json.Unmarshal(jmpJson, &jmp)
	if err != nil {
		return []byte{}, err
	}

	if len(jmp.Jumps) != 24 {
		return []byte{}, errors.New("must have 24 jumps")
	}

	bf := byteframe.NewByteFrame()
	bf.SetLE()
	bf.WriteUint32(0)
	bf.WriteUint32(0)
	bf.WriteUint32(0)
	bf.WriteUint32(0)
	bf.WriteUint32(0)

	for i := 0; i < len(jmp.Strings); i++ {
		if _, ok := strings[jmp.Strings[i]]; ok {
			jmp.ptrStrings = append(jmp.ptrStrings, strings[jmp.Strings[i]])
		} else {
			strings[jmp.Strings[i]] = uint32(bf.Index())
			jmp.ptrStrings = append(jmp.ptrStrings, uint32(bf.Index()))
			sjisString := sjis.NewString(jmp.Strings[i])
			bf.WriteBytes(sjisString.Bytes())
		}
	}

	for i := 0; i < len(jmp.Jumps); i++ {
		if _, ok := strings[jmp.Jumps[i].Title]; ok {
			jmp.Jumps[i].ptrTitle = strings[jmp.Jumps[i].Title]
		} else {
			strings[jmp.Jumps[i].Title] = uint32(bf.Index())
			jmp.Jumps[i].ptrTitle = uint32(bf.Index())
			sjisTitle := sjis.NewString(jmp.Jumps[i].Title)
			bf.WriteBytes(sjisTitle.Bytes())
		}
		if _, ok := strings[jmp.Jumps[i].Description]; ok {
			jmp.Jumps[i].ptrDescription = strings[jmp.Jumps[i].Description]
		} else {
			strings[jmp.Jumps[i].Description] = uint32(bf.Index())
			jmp.Jumps[i].ptrDescription = uint32(bf.Index())
			sjisDescription := sjis.NewString(jmp.Jumps[i].Description)
			bf.WriteBytes(sjisDescription.Bytes())
		}
	}

	for i := 0; i < len(jmp.Menus); i++ {
		jmp.Menus[i].ptrEntries = uint32(bf.Index())
		for j := 0; j < len(jmp.Menus[i].Entries); j++ {
			bf.WriteUint16(jmp.Menus[i].Entries[j].Index)
			bf.WriteUint16(jmp.Menus[i].Entries[j].Flags)
		}
		jmp.Menus[i].ptrStageIds = uint32(bf.Index())
		for j := 0; j < len(jmp.Menus[i].StageIds); j++ {
			bf.WriteUint16(jmp.Menus[i].StageIds[j])
		}
		bf.WriteUint16(0)
	}

	jmp.ptrPtrStrings = uint32(bf.Index())
	for i := 0; i < len(jmp.ptrStrings); i++ {
		bf.WriteUint32(jmp.ptrStrings[i])
	}

	jmp.ptrMenus = uint32(bf.Index())
	for i := 0; i < len(jmp.Menus); i++ {
		bf.WriteUint32(jmp.Menus[i].ptrEntries)
		bf.WriteUint32(uint32(len(jmp.Menus[i].Entries)))
		bf.WriteUint32(jmp.Menus[i].ptrStageIds)
	}

	jmp.ptrJumps = uint32(bf.Index())
	for i := 0; i < len(jmp.Jumps); i++ {
		bf.WriteUint32(jmp.Jumps[i].Id)
		bf.WriteUint32(jmp.Jumps[i].Unk)
		for j := 0; j < 4; j++ {
			bf.WriteUint16(jmp.Jumps[i].StageIds[j])
		}
		for j := 0; j < 2; j++ {
			for k := 0; k < 3; k++ {
				bf.WriteFloat32(jmp.Jumps[i].Destinations[j].Coordinates[k])
			}
			bf.WriteUint32(jmp.Jumps[i].Destinations[j].Rotation)
		}
		bf.WriteUint32(jmp.Jumps[i].ptrTitle)
		bf.WriteUint32(jmp.Jumps[i].ptrDescription)
	}

	_, _ = bf.Seek(0, 0)
	bf.WriteUint32(jmp.ptrJumps)
	bf.WriteUint32(jmp.ptrMenus)
	bf.WriteUint32(uint32(len(jmp.Menus)))
	bf.WriteUint32(jmp.ptrPtrStrings)
	bf.WriteUint32(uint32(len(jmp.Strings)))

	return bf.Data(), nil
}
