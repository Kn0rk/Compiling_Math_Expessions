package declareandassign

import (
	"bytes"
	"encoding/binary"
)

func translateTerm(num int) []byte {

	// push term on stack
	// we have to use the 64bit register for push/pop
	buf := new(bytes.Buffer)
	buf.WriteByte(0x68) // push
	binary.Write(buf, binary.LittleEndian, int32(num))
	byts := buf.Bytes()
	return byts
}

func translateOperation(token Operator) []byte {
	// args on stack
	// pop args RAX, RBX
	popped := []byte{0x58, 0x5B}

	if token == OpAdd {
		popped = append(popped,
			[]byte{
				0x48, 0x01, 0xd8}...)
	} else if token == OpMul {
		popped = append(popped,
			[]byte{
				0xf7, 0xeb}...)
	} else if token == OpSub {
		popped = append(popped,
			[]byte{
				0x48, 0x29, 0xd8}...)

	} else {
		// todo div is a diva and needs special attention
		panic(1)
	}
	// push   %rax for next operation
	popped = append(popped,
		[]byte{
			0x50}...)

	return popped
}

func addPrintResult() []byte {

	var ret = []byte{
		0xeb, 0x39,
		0xbb, 0x0a, 0x00, 0x00, 0x00,
		0x48, 0x83, 0xec, 0x20,
		0x48, 0x89, 0xe1,
		0x48, 0x31, 0xd2,
		0x48, 0xf7, 0xf3,
		0x80, 0xc2, 0x30,
		0x48, 0xff, 0xc9,
		0x88, 0x11,
		0x48, 0x85, 0xc0,
		0x75, 0xed,
		0x48, 0x89, 0xe2,
		0x48, 0x29, 0xca,
		0xb8, 0x01, 0x00, 0x00, 0x00,
		0xbf, 0x01, 0x00, 0x00, 0x00,
		0x48, 0x89, 0xce,
		0x0f, 0x05,
		0x48, 0x83, 0xc4, 0x20,
		0xc3}
	//   401000:       eb 39                   jmp    0x40103b
	//   401002:       bb 0a 00 00 00          mov    $0xa,%ebx
	//   401007:       48 83 ec 20             sub    $0x20,%rsp
	//   40100b:       48 89 e1                mov    %rsp,%rcx
	//   40100e:       48 31 d2                xor    %rdx,%rdx
	//   401011:       48 f7 f3                div    %rbx
	//   401014:       80 c2 30                add    $0x30,%dl
	//   401017:       48 ff c9                dec    %rcx
	//   40101a:       88 11                   mov    %dl,(%rcx)
	//   40101c:       48 85 c0                test   %rax,%rax
	//   40101f:       75 ed                   jne    0x40100e
	//   401021:       48 89 e2                mov    %rsp,%rdx
	//   401024:       48 29 ca                sub    %rcx,%rdx
	//   401027:       b8 01 00 00 00          mov    $0x1,%eax
	//   40102c:       bf 01 00 00 00          mov    $0x1,%edi
	//   401031:       48 89 ce                mov    %rcx,%rsi
	//   401034:       0f 05                   syscall
	//   401036:       48 83 c4 20             add    $0x20,%rsp
	//   40103a:       c3                      ret
	return ret
}
func call_print() []byte {
	return []byte{
		0xb9, 0x02, 0x10, 0x40, 0x00, 0xff, 0xd1,
	}
}

// 40104f:       b8 01 00 00 00          mov    $0x1,%eax
// 401054:       bf 01 00 00 00          mov    $0x1,%edi
// 6a 0a                   push   $0xa
//
//	40105b:       48 89 e6                mov    %rsp,%rsi
//	40105e:       ba 01 00 00 00          mov    $0x1,%edx
//	401063:       0f 05                   syscall
func newLine() []byte {
	return []byte{
		0xb8, 0x01, 0x00, 0x00, 0x00, 0xbf, 0x01, 0x00, 0x00, 0x00, 0x6a, 0x0a, 0x48, 0x89, 0xe6, 0xba, 0x01, 0x00, 0x00, 0x00, 0x0f, 0x05,
	}
}

func cleanExit() []byte {
	return []byte{
		// mov eax 1
		0xb8, 0x01, 0x0, 0x0, 0x0,
		// xor ebx
		0x31, 0xdb,
		// sys call
		0xcd, 0x80}
}

func fixedElfHeader() []byte {
	// ELF Header
	elfHeaderArr := make([]byte, 64)
	// 5B7x 2H5I 6H
	//4s
	binary.BigEndian.PutUint32(elfHeaderArr[0:], 0x7F454C46) // Magic number ascii == elf
	// 5B
	elfHeaderArr[4] = 2 // 64-bit
	elfHeaderArr[5] = 1 // Little-endian
	elfHeaderArr[6] = 1 // ELF version
	elfHeaderArr[7] = 0 // ABI
	elfHeaderArr[8] = 0 // ABI version
	// 16
	binary.LittleEndian.PutUint16(elfHeaderArr[16:], 2)    // Executable file
	binary.LittleEndian.PutUint16(elfHeaderArr[18:], 0x3e) // Intel 80386
	binary.LittleEndian.PutUint32(elfHeaderArr[20:], 1)
	binary.LittleEndian.PutUint64(elfHeaderArr[24:], 0x00401000) // Entry point
	binary.LittleEndian.PutUint64(elfHeaderArr[32:], 64)         // Start of program headers
	binary.LittleEndian.PutUint64(elfHeaderArr[40:], 0x2108)     // Start of section headers
	binary.LittleEndian.PutUint32(elfHeaderArr[48:], 0)          // Flags
	binary.LittleEndian.PutUint16(elfHeaderArr[52:], 64)         // Size of this header
	binary.LittleEndian.PutUint16(elfHeaderArr[54:], 56)         // Size of program headers
	binary.LittleEndian.PutUint16(elfHeaderArr[56:], 3)          // Number of program headers
	binary.LittleEndian.PutUint16(elfHeaderArr[58:], 64)         // Size of section headers
	binary.LittleEndian.PutUint16(elfHeaderArr[60:], 6)          // Number of section headers
	binary.LittleEndian.PutUint16(elfHeaderArr[62:], 5)          // Section header string table index
	return elfHeaderArr
}

type ProgramHeader struct {
	load                uint32
	flag                uint32
	file_offset         uint64
	virtual_address     uint64
	physical_address    uint64
	segment_byte_length uint64
	memory_size         uint64
	alignment           uint64
}

func (header ProgramHeader) toBytes() []byte {
	programHeader := make([]byte, 56)
	binary.LittleEndian.PutUint32(programHeader[0:], header.load)                 // Type (LOAD)
	binary.LittleEndian.PutUint32(programHeader[4:], header.flag)                 // Flags (R-X)
	binary.LittleEndian.PutUint64(programHeader[8:], header.file_offset)          // Offset
	binary.LittleEndian.PutUint64(programHeader[16:], header.virtual_address)     // Virtual address
	binary.LittleEndian.PutUint64(programHeader[24:], header.virtual_address)     // Physical address
	binary.LittleEndian.PutUint64(programHeader[32:], header.segment_byte_length) // File size
	// Memory size
	if header.memory_size == 0 {
		binary.LittleEndian.PutUint64(programHeader[40:], header.segment_byte_length)
	} else {
		binary.LittleEndian.PutUint64(programHeader[40:], header.memory_size)
	}
	// Alignment
	if header.alignment == 0 {
		binary.LittleEndian.PutUint64(programHeader[48:], 0x1000)
	} else {
		binary.LittleEndian.PutUint64(programHeader[48:], header.alignment)
	}

	return programHeader
}

type SectionHeader struct {
	name_offset     uint32
	section_type    uint32
	flags           uint64
	virtual_address uint64
	file_offset     uint64
	length          uint64
	link            uint32
	info            uint32
	align           uint64
	ent_size        uint64
}

func (sec SectionHeader) toBytes() []byte {
	programSection := make([]byte, 64)
	binary.LittleEndian.PutUint32(programSection[0:], sec.name_offset)      // Offset into string table
	binary.LittleEndian.PutUint32(programSection[4:], sec.section_type)     // Type e.g. symtab, strtab,...
	binary.LittleEndian.PutUint64(programSection[8:], sec.flags)            // Flags writable,...
	binary.LittleEndian.PutUint64(programSection[16:], sec.virtual_address) // Virtual address
	binary.LittleEndian.PutUint64(programSection[24:], sec.file_offset)     // file offset
	binary.LittleEndian.PutUint64(programSection[32:], sec.length)          // section size in file
	binary.LittleEndian.PutUint32(programSection[40:], sec.link)            // sh link
	binary.LittleEndian.PutUint32(programSection[44:], sec.info)            // sh info
	binary.LittleEndian.PutUint64(programSection[48:], sec.align)           // Alignment
	binary.LittleEndian.PutUint64(programSection[56:], sec.ent_size)        // ent size if s
	return programSection
}

func createBinary(
	program []byte,
	data []byte,
) []byte {

	program = append(addPrintResult(), program...)
	program = append(program, cleanExit()...)
	var programHeaders = make([]ProgramHeader, 0)
	programHeaders = append(programHeaders,
		ProgramHeader{
			load:                1,
			flag:                4,
			file_offset:         0,
			virtual_address:     0x00400000,
			segment_byte_length: 0xe8,
		})
	programHeaders = append(programHeaders, ProgramHeader{
		load:                1,
		flag:                5,
		file_offset:         0x1000,
		virtual_address:     0x00401000,
		segment_byte_length: 0x1f,
	})
	programHeaders = append(programHeaders, ProgramHeader{
		load:                1,
		flag:                6,
		file_offset:         0x2000,
		virtual_address:     0x00402000,
		segment_byte_length: 0xd,
	})

	var strtab = []byte{}
	strtab = append(strtab, append([]byte("min.asm"), 0)...)
	strtab = append(strtab, append([]byte("hello"), 0)...)
	strtab = append(strtab, append([]byte("__bss_start"), 0)...)
	strtab = append(strtab, append([]byte("_edata"), 0)...)
	strtab = append(strtab, append([]byte("_end"), []byte{0, 0}...)...)
	strtab = append(strtab, append([]byte(".symtab"), 0)...)
	strtab = append(strtab, append([]byte(".strtab"), 0)...)
	strtab = append(strtab, append([]byte(".shstrtab"), 0)...)
	strtab = append(strtab, append([]byte(".text"), 0)...)
	strtab = append(strtab, append([]byte(".data"), 0)...)

	var sections = make([]SectionHeader, 0)

	sections = append(sections,
		SectionHeader{
			// All 0
		})

	sections = append(sections, SectionHeader{
		name_offset:     0x1B,
		section_type:    0x1,
		flags:           6,
		virtual_address: 0x00401000,
		file_offset:     0x01000,
		align:           0x10,
		length:          uint64(len(program)),
	})

	sections = append(sections, SectionHeader{
		name_offset:     0x21,
		section_type:    0x1,
		flags:           3,
		virtual_address: 0x00402000,
		file_offset:     0x02000,
		align:           4,
		length:          uint64(len(data)) - 1,
	})

	sections = append(sections, SectionHeader{
		name_offset:     0x01, // symtab
		section_type:    0x2,  // symtab
		flags:           0,
		virtual_address: 0x00,
		file_offset:     0x02010,
		align:           8,
		length:          0xa8,
		ent_size:        0x18,
		info:            3,
		link:            4,
	})

	sections = append(sections, SectionHeader{
		name_offset:     0x9, // strtab
		section_type:    0x3, //strtab
		flags:           0,
		virtual_address: 0x0,
		file_offset:     0x020b8,
		align:           1,
		length:          0x27,
	})

	sections = append(sections, SectionHeader{
		name_offset:     0x11, // shstrtab
		section_type:    0x3,  //strtab
		flags:           0,
		virtual_address: 0x0,
		file_offset:     0x020df,
		align:           1,
		length:          0x27,
	})

	var bytes = make([]byte, 0)
	bytes = append(bytes, fixedElfHeader()...)
	for _, seg := range programHeaders {
		bytes = append(bytes, seg.toBytes()...)
	}
	bytes = append(bytes, make([]byte, 0x1000-len(bytes))...)

	bytes = append(bytes, program...)
	bytes = append(bytes, make([]byte, 0x2000-len(bytes))...)
	bytes = append(bytes, data...)
	bytes = append(bytes, make([]byte, 0x2010-len(bytes))...)

	bytes = append(bytes, make([]byte, 0x20b9-len(bytes))...)
	bytes = append(bytes, strtab...)
	// // skip symtab
	bytes = append(bytes, make([]byte, 0x2108-len(bytes))...)

	for _, sec := range sections {
		bytes = append(bytes, sec.toBytes()...)
	}
	return bytes
}
