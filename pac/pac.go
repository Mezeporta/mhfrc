package pac

import (
	"encoding/json"
	"mhfrc/util/byteframe"
	"mhfrc/util/sjis"
)

type Custom1 struct {
	Strings []string
	Data    [][]byte
}

type Pac struct {
	Ptr5Strings  []string
	Ptr6Strings  []string
	Ptr7Strings  []string
	Ptr8Strings  []string
	Ptr9Strings  []string
	Ptr10Strings []string
	Ptr11Strings []string
	Ptr12Strings []string
	Ptr13Strings []string
	Ptr14Data    [][]byte
	Ptr15Data    [][]byte
	Ptr16Data    [][]byte
	Ptr17Data    [][]byte
	Ptr18Data    []byte
	//Ptr19 unknown
	Ptr20Strings []string
	Ptr21Strings []string
	Ptr22Data    []Custom1
	Ptr23Data    []Custom1
	Ptr24Data    [][]byte
	Ptr25Data    [][]byte
	Ptr26Data    [][]byte
	Ptr27Data    [][]byte

	lenStructs []uint16
}

func DecompilePac(pacData []byte) ([]byte, error) {
	bf := byteframe.NewByteFrameFromBytes(pacData)
	bf.SetLE()
	_ = bf.ReadBytes(12) // 70 61 63 1A 0A 00 00 00 00 00 00 00
	_ = bf.ReadUint32()
	ptrStructLens := bf.ReadUint32()
	ptrUnk5 := bf.ReadUint32()  // strings type
	ptrUnk6 := bf.ReadUint32()  // strings type
	ptrUnk7 := bf.ReadUint32()  // strings type
	ptrUnk8 := bf.ReadUint32()  // strings type
	ptrUnk9 := bf.ReadUint32()  // strings type
	ptrUnk10 := bf.ReadUint32() // strings type
	ptrUnk11 := bf.ReadUint32() // strings type
	ptrUnk12 := bf.ReadUint32() // strings type
	ptrUnk13 := bf.ReadUint32() // strings type
	ptrUnk14 := bf.ReadUint32() // struct type
	ptrUnk15 := bf.ReadUint32() // struct type
	ptrUnk16 := bf.ReadUint32() // struct type
	ptrUnk17 := bf.ReadUint32() // struct type
	ptrUnk18 := bf.ReadUint32() // bytes type
	ptrUnk19 := bf.ReadUint32() // unknown
	ptrUnk20 := bf.ReadUint32() // strings type
	ptrUnk21 := bf.ReadUint32() // strings type
	ptrUnk22 := bf.ReadUint32() // custom struct type
	ptrUnk23 := bf.ReadUint32() // custom struct type
	ptrUnk24 := bf.ReadUint32() // struct type
	ptrUnk25 := bf.ReadUint32() // struct type
	ptrUnk26 := bf.ReadUint32() // struct type
	ptrUnk27 := bf.ReadUint32() // struct type

	pac := Pac{}

	_, _ = bf.Seek(int64(ptrStructLens), 0)
	for i := 0; i < 56; i++ {
		pac.lenStructs = append(pac.lenStructs, bf.ReadUint16())
	}

	pac.Ptr5Strings = readStringGroup(bf, ptrUnk5)
	pac.Ptr6Strings = readStringGroup(bf, ptrUnk6)
	pac.Ptr7Strings = readStringGroup(bf, ptrUnk7)
	pac.Ptr8Strings = readStringGroup(bf, ptrUnk8)
	pac.Ptr9Strings = readStringGroup(bf, ptrUnk9)
	pac.Ptr10Strings = readStringGroup(bf, ptrUnk10)
	pac.Ptr11Strings = readStringGroup(bf, ptrUnk11)
	pac.Ptr12Strings = readStringGroup(bf, ptrUnk12)
	pac.Ptr13Strings = readStringGroup(bf, ptrUnk13)
	pac.Ptr14Data = readBytesGroup(bf, ptrUnk14, pac.lenStructs[0], 36)
	pac.Ptr15Data = readBytesGroup(bf, ptrUnk15, pac.lenStructs[1], 116)
	pac.Ptr16Data = readBytesGroup(bf, ptrUnk16, pac.lenStructs[2], 8)
	pac.Ptr17Data = readBytesGroup(bf, ptrUnk17, pac.lenStructs[3], 8)

	_, _ = bf.Seek(int64(ptrUnk18), 0)
	pac.Ptr18Data = bf.ReadBytes(3240)

	_, _ = bf.Seek(int64(ptrUnk19), 0)
	// unknown format

	pac.Ptr20Strings = readStringGroup(bf, ptrUnk20)
	pac.Ptr21Strings = readStringGroup(bf, ptrUnk21)

	_, _ = bf.Seek(int64(ptrUnk22), 0)
	p22Ptrs := make([]uint32, 5)
	for i := 0; i < 5; i++ {
		p22Ptrs[i] = bf.ReadUint32()
	}
	for i := 0; i < 5; i++ {
		if p22Ptrs[i] == 0 {
			continue
		}
		_, _ = bf.Seek(int64(p22Ptrs[i]), 0)
		p22 := Custom1{}
		p22SubPtrs := make([]uint32, pac.lenStructs[16])
		for j := uint16(0); j < pac.lenStructs[16]; j++ {
			p22SubPtrs[j] = bf.ReadUint32()
		}
		for j := uint16(0); j < pac.lenStructs[16]; j++ {
			_, _ = bf.Seek(int64(p22SubPtrs[j]), 0)
			ptrStr := bf.ReadUint32()
			p22.Data = append(p22.Data, bf.ReadBytes(20))
			_, _ = bf.Seek(int64(ptrStr), 0)
			str := sjis.NewBytes(bf.ReadNullTerminatedBytes())
			p22.Strings = append(p22.Strings, str.String())
		}
		pac.Ptr22Data = append(pac.Ptr22Data, p22)
	}

	_, _ = bf.Seek(int64(ptrUnk23), 0)
	p23Ptrs := make([]uint32, 5)
	for i := 0; i < 5; i++ {
		p23Ptrs[i] = bf.ReadUint32()
	}
	for i := 0; i < 5; i++ {
		if p23Ptrs[i] == 0 {
			continue
		}
		_, _ = bf.Seek(int64(p23Ptrs[i]), 0)
		p23 := Custom1{}
		p23SubPtrs := make([]uint32, pac.lenStructs[16])
		for j := uint16(0); j < pac.lenStructs[16]; j++ {
			p23SubPtrs[j] = bf.ReadUint32()
		}
		for j := uint16(0); j < pac.lenStructs[16]; j++ {
			_, _ = bf.Seek(int64(p23SubPtrs[j]), 0)
			ptrStr := bf.ReadUint32()
			p23.Data = append(p23.Data, bf.ReadBytes(20))
			_, _ = bf.Seek(int64(ptrStr), 0)
			str := sjis.NewBytes(bf.ReadNullTerminatedBytes())
			p23.Strings = append(p23.Strings, str.String())
		}
		pac.Ptr23Data = append(pac.Ptr23Data, p23)
	}

	pac.Ptr24Data = readBytesGroup(bf, ptrUnk24, pac.lenStructs[17], 16)
	pac.Ptr25Data = readBytesGroup(bf, ptrUnk25, pac.lenStructs[18], 16)
	pac.Ptr26Data = readBytesGroup(bf, ptrUnk26, pac.lenStructs[19], 16)
	pac.Ptr27Data = readBytesGroup(bf, ptrUnk27, pac.lenStructs[20], 16)

	pacJson, err := json.MarshalIndent(pac, "", "\t")
	if err != nil {
		return nil, err
	}

	return pacJson, nil
}

func readStringGroup(bf *byteframe.ByteFrame, ptr uint32) []string {
	var i int64
	var strings []string
	for {
		_, _ = bf.Seek(int64(ptr)+i*4, 0)
		ptrString := bf.ReadUint32()
		if ptrString == 0 {
			break
		}
		_, _ = bf.Seek(int64(ptrString), 0)
		str := sjis.NewBytes(bf.ReadNullTerminatedBytes())
		strings = append(strings, str.String())
		i++
	}
	return strings
}

func readBytesGroup(bf *byteframe.ByteFrame, ptr uint32, lenStruct uint16, lenBytes uint) [][]byte {
	_, _ = bf.Seek(int64(ptr), 0)
	var data [][]byte
	for i := uint16(0); i < lenStruct; i++ {
		data = append(data, bf.ReadBytes(lenBytes))
	}
	return data
}
