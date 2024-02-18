package ssm32

/*
#include <string.h>

#ifndef _WIN32
	#define UNALIGNED
#endif

typedef char CHAR;
typedef unsigned char UCHAR, * PUCHAR;
typedef unsigned short USHORT, * PUSHORT;
#define USHORT_SIZE sizeof(USHORT)
typedef int LONG;
typedef unsigned int ULONG, * PULONG;
#define ULONG_SIZE sizeof(ULONG)
typedef unsigned long long ULONGLONG, * PULONGLONG;
typedef void* PVOID;
typedef const char* PCSTR;

typedef struct _IMAGE_DOS_HEADER {
	USHORT e_magic;
	USHORT e_cblp;
	USHORT e_cp;
	USHORT e_crlc;
	USHORT e_cparhdr;
	USHORT e_minalloc;
	USHORT e_maxalloc;
	USHORT e_ss;
	USHORT e_sp;
	USHORT e_csum;
	USHORT e_ip;
	USHORT e_cs;
	USHORT e_lfarlc;
	USHORT e_ovno;
	USHORT e_res[4];
	USHORT e_oemid;
	USHORT e_oeminfo;
	USHORT e_res2[10];
	LONG e_lfanew;
} IMAGE_DOS_HEADER, * PIMAGE_DOS_HEADER;
#define IMAGE_DOS_HEADER_SIZE sizeof(IMAGE_DOS_HEADER)

typedef struct _IMAGE_FILE_HEADER {
	USHORT Machine;
	USHORT NumberOfSections;
	ULONG TimeDateStamp;
	ULONG PointerToSymbolTable;
	ULONG NumberOfSymbols;
	USHORT SizeOfOptionalHeader;
	USHORT Characteristics;
} IMAGE_FILE_HEADER, * PIMAGE_FILE_HEADER;

typedef struct _IMAGE_DATA_DIRECTORY {
	ULONG VirtualAddress;
	ULONG Size;
} IMAGE_DATA_DIRECTORY, * PIMAGE_DATA_DIRECTORY;

#define IMAGE_NUMBEROF_DIRECTORY_ENTRIES 16

typedef struct _IMAGE_OPTIONAL_HEADER32 {
	USHORT Magic;
	UCHAR MajorLinkerVersion;
	UCHAR MinorLinkerVersion;
	ULONG SizeOfCode;
	ULONG SizeOfInitializedData;
	ULONG SizeOfUninitializedData;
	ULONG AddressOfEntryPoint;
	ULONG BaseOfCode;
	ULONG BaseOfData;
	ULONG ImageBase;
	ULONG SectionAlignment;
	ULONG FileAlignment;
	USHORT MajorOperatingSystemVersion;
	USHORT MinorOperatingSystemVersion;
	USHORT MajorImageVersion;
	USHORT MinorImageVersion;
	USHORT MajorSubsystemVersion;
	USHORT MinorSubsystemVersion;
	ULONG Win32VersionValue;
	ULONG SizeOfImage;
	ULONG SizeOfHeaders;
	ULONG CheckSum;
	USHORT Subsystem;
	USHORT DllCharacteristics;
	ULONG SizeOfStackReserve;
	ULONG SizeOfStackCommit;
	ULONG SizeOfHeapReserve;
	ULONG SizeOfHeapCommit;
	ULONG LoaderFlags;
	ULONG NumberOfRvaAndSizes;
	IMAGE_DATA_DIRECTORY DataDirectory[IMAGE_NUMBEROF_DIRECTORY_ENTRIES];
} IMAGE_OPTIONAL_HEADER32, * PIMAGE_OPTIONAL_HEADER32;

typedef struct _IMAGE_NT_HEADERS32 {
	ULONG Signature;
	IMAGE_FILE_HEADER FileHeader;
	IMAGE_OPTIONAL_HEADER32 OptionalHeader;
} IMAGE_NT_HEADERS32, * PIMAGE_NT_HEADERS32;
#define IMAGE_NT_HEADERS32_SIZE sizeof(IMAGE_NT_HEADERS32)

#define IMAGE_NT_OPTIONAL_HDR32_MAGIC 0x10b

#define IMAGE_DIRECTORY_ENTRY_IMPORT 1
#define IMAGE_DIRECTORY_ENTRY_BASERELOC 5

typedef struct _IMAGE_SECTION_HEADER {
	UCHAR Name[8];
	ULONG Misc;
	ULONG VirtualAddress;
	ULONG SizeOfRawData;
	ULONG PointerToRawData;
	ULONG PointerToRelocations;
	ULONG PointerToLinenumbers;
	USHORT NumberOfRelocations;
	USHORT NumberOfLinenumbers;
	ULONG Characteristics;
} IMAGE_SECTION_HEADER, * PIMAGE_SECTION_HEADER;
#define IMAGE_SECTION_HEADER_SIZE sizeof(IMAGE_SECTION_HEADER)

typedef struct _IMAGE_BASE_RELOCATION {
	ULONG VirtualAddress;
	ULONG SizeOfBlock;
} IMAGE_BASE_RELOCATION;
typedef IMAGE_BASE_RELOCATION UNALIGNED * PIMAGE_BASE_RELOCATION;
#define IMAGE_BASE_RELOCATION_SIZE sizeof(IMAGE_BASE_RELOCATION)

#define IMAGE_REL_BASED_ABSOLUTE 0
#define IMAGE_REL_BASED_HIGHLOW 3

typedef struct _IMAGE_IMPORT_DESCRIPTOR {
	ULONG u1;
	// union {
	//     ULONG Characteristics;
	//     ULONG OriginalFirstThunk;
	// };
	ULONG TimeDateStamp;
	ULONG ForwarderChain;
	ULONG Name;
	ULONG FirstThunk;
} IMAGE_IMPORT_DESCRIPTOR;
typedef IMAGE_IMPORT_DESCRIPTOR UNALIGNED * PIMAGE_IMPORT_DESCRIPTOR;
#define IMAGE_IMPORT_DESCRIPTOR_SIZE sizeof(IMAGE_IMPORT_DESCRIPTOR)

typedef struct _IMAGE_THUNK_DATA32 {
	ULONG u1;
	// union {
	//     ULONG ForwarderString;      // PBYTE
	//     ULONG Function;             // PULONG
	//     ULONG Ordinal;
	//     ULONG AddressOfData;        // PIMAGE_IMPORT_BY_NAME
	// } u1;
} IMAGE_THUNK_DATA32, * PIMAGE_THUNK_DATA32;
#define IMAGE_THUNK_DATA32_SIZE sizeof(IMAGE_THUNK_DATA32)

#define IMAGE_ORDINAL_FLAG32 0x80000000

typedef struct _IMAGE_IMPORT_BY_NAME {
	USHORT Hint;
	CHAR   Name[1];
} IMAGE_IMPORT_BY_NAME, * PIMAGE_IMPORT_BY_NAME;

void* memset_ull(void* dest, int c, ULONGLONG n) {
	return memset(dest, c, n);
}
*/
import "C"
import (
	"encoding/binary"
	"serversidemapper32/rawwrapper"
	"unsafe"
)

type MMap32Context struct {
	dataWrapper                                                                                   *rawwrapper.RawWrapper
	elfaNew                                                                                       uintptr
	ntHeaders                                                                                     C.PIMAGE_NT_HEADERS32
	offsetToSectionHeader, sizeOfImage, sizeOfHeaders, baseRelocVA, importDirVA, numberOfSections uintptr
	vaImitationDataWrapper                                                                        *rawwrapper.RawWrapper
	vaImitationImportDir                                                                          C.PIMAGE_IMPORT_DESCRIPTOR
}

// GetMMap32Data takes the raw data of a DLL and returns a uint8 array containing
// the image base, size of the image, names of imports (DLLs, functions by ordinals,
// and names), along with a context object that stores extracted data required for
// calling the MMap32 function.
//
// Parameters:
//   - peRawData: The raw data of the DLL as a []uint8 array.
//
// Returns:
//   - A []uint8 array containing the image base, size of the image, and names of imports.
//   - A pointer to the MMap32Context object that stores extracted data.
func GetMMap32Data(peRawData []uint8) ([]uint8, *MMap32Context) {
	var ctx = new(MMap32Context)
	ctx.dataWrapper = rawwrapper.NewRawWrapper(peRawData)
	ctx.elfaNew = uintptr(C.PIMAGE_DOS_HEADER(ctx.dataWrapper.At(0)).e_lfanew)
	ctx.ntHeaders = C.PIMAGE_NT_HEADERS32(ctx.dataWrapper.At(ctx.elfaNew))
	if ctx.ntHeaders.OptionalHeader.Magic != C.IMAGE_NT_OPTIONAL_HDR32_MAGIC {
		return nil, nil
	}

	ctx.offsetToSectionHeader = ctx.elfaNew + C.IMAGE_NT_HEADERS32_SIZE
	ctx.sizeOfImage = uintptr(ctx.ntHeaders.OptionalHeader.SizeOfImage)
	ctx.sizeOfHeaders = uintptr(ctx.ntHeaders.OptionalHeader.SizeOfHeaders)
	ctx.baseRelocVA = uintptr(ctx.ntHeaders.OptionalHeader.DataDirectory[C.IMAGE_DIRECTORY_ENTRY_BASERELOC].VirtualAddress)
	ctx.importDirVA = uintptr(ctx.ntHeaders.OptionalHeader.DataDirectory[C.IMAGE_DIRECTORY_ENTRY_IMPORT].VirtualAddress)
	ctx.numberOfSections = uintptr(ctx.ntHeaders.FileHeader.NumberOfSections)
	ctx.vaImitationDataWrapper = rawwrapper.NewRawWrapper(make([]uint8, ctx.sizeOfImage))
	ctx.vaImitationImportDir = C.PIMAGE_IMPORT_DESCRIPTOR(ctx.vaImitationDataWrapper.At(ctx.importDirVA))

	copy(ctx.vaImitationDataWrapper.Data, ctx.dataWrapper.Data[:ctx.sizeOfHeaders])
	for i := uintptr(0); i < ctx.numberOfSections; i++ {
		var sectionHeader = C.PIMAGE_SECTION_HEADER(ctx.dataWrapper.At(ctx.offsetToSectionHeader + C.IMAGE_SECTION_HEADER_SIZE*i))
		copy(ctx.vaImitationDataWrapper.Data[sectionHeader.VirtualAddress:], ctx.dataWrapper.Data[sectionHeader.PointerToRawData:][:sectionHeader.SizeOfRawData])
	}

	var (
		mmap32Data             = make([]uint8, 0)
		vaImitationNtHeaders   = C.PIMAGE_NT_HEADERS32(ctx.vaImitationDataWrapper.At(ctx.elfaNew))
		vaImitationImportDirIt = ctx.vaImitationImportDir
	)
	mmap32Data = binary.LittleEndian.AppendUint32(mmap32Data, uint32(vaImitationNtHeaders.OptionalHeader.ImageBase))
	mmap32Data = binary.LittleEndian.AppendUint32(mmap32Data, uint32(ctx.sizeOfImage))
	for vaImitationImportDirIt.u1 /* Characteristics */ != 0 {
		var (
			origFirstThunk = C.PIMAGE_THUNK_DATA32(ctx.vaImitationDataWrapper.At(uintptr(vaImitationImportDirIt.u1 /* OriginalFirstThunk */)))
			imgCSTRName    = C.PCSTR(ctx.vaImitationDataWrapper.At(uintptr(vaImitationImportDirIt.Name)))
		)

		mmap32Data = append(mmap32Data, 255)
		mmap32Data = append(mmap32Data, []uint8(C.GoString(imgCSTRName))...)
		mmap32Data = append(mmap32Data, 0)

		for origFirstThunk.u1 /* AddressOfData */ != 0 {
			if (origFirstThunk.u1 /* Ordinal */ & C.IMAGE_ORDINAL_FLAG32) == C.IMAGE_ORDINAL_FLAG32 {
				mmap32Data = append(mmap32Data, 254)
				mmap32Data = binary.LittleEndian.AppendUint16(mmap32Data, uint16(C.USHORT(origFirstThunk.u1 /* Ordinal */ &0xffff)))
			} else {
				mmap32Data = append(mmap32Data, 253)
				mmap32Data = append(mmap32Data, []uint8(C.GoString(C.PCSTR(unsafe.Pointer(&C.PIMAGE_IMPORT_BY_NAME(ctx.vaImitationDataWrapper.At(uintptr(origFirstThunk.u1 /* AddressOfData */))).Name[0]))))...)
				mmap32Data = append(mmap32Data, 0)
			}

			origFirstThunk = C.PIMAGE_THUNK_DATA32(unsafe.Pointer(uintptr(unsafe.Pointer(origFirstThunk)) + C.IMAGE_THUNK_DATA32_SIZE))
		}

		vaImitationImportDirIt = C.PIMAGE_IMPORT_DESCRIPTOR(unsafe.Pointer(uintptr(unsafe.Pointer(vaImitationImportDirIt)) + C.IMAGE_IMPORT_DESCRIPTOR_SIZE))
	}

	return mmap32Data, ctx
}

// MMap32 gets the DLL into the VA state, fixes its relocations and imports using processedData.
// It takes the processedData, which is a byte slice containing the necessary information for fixing the DLL,
// and the ctx, which is a pointer to the MMap32Context struct containing additional context information.
// It returns the modified byte slice representing the DLL in the VA state.
func MMap32(processedData []uint8, ctx *MMap32Context) []uint8 {
	if len(processedData) < C.ULONG_SIZE*2 {
		return nil
	}

	vaImageBase := binary.LittleEndian.Uint32(processedData)
	diff := binary.LittleEndian.Uint32(processedData[C.ULONG_SIZE:])
	relocs := processedData[C.ULONG_SIZE*2:]

	var vaImitationDataWrapperData = make([]uint8, 0)
	vaImitationDataWrapperData = binary.LittleEndian.AppendUint32(vaImitationDataWrapperData, uint32(vaImageBase))
	vaImitationDataWrapperData = binary.LittleEndian.AppendUint32(vaImitationDataWrapperData, uint32(ctx.sizeOfImage))
	vaImitationDataWrapperData = binary.LittleEndian.AppendUint32(vaImitationDataWrapperData, uint32(ctx.ntHeaders.OptionalHeader.AddressOfEntryPoint))
	vaImitationDataWrapperData = append(vaImitationDataWrapperData, make([]uint8, ctx.sizeOfImage)...)

	var (
		baseRelocVA            = uintptr(ctx.ntHeaders.OptionalHeader.DataDirectory[C.IMAGE_DIRECTORY_ENTRY_BASERELOC].VirtualAddress)
		vaImitationDataWrapper = rawwrapper.NewRawWrapper(vaImitationDataWrapperData[int(unsafe.Sizeof(uint32(0))*3):])
		ntHeaders              = C.PIMAGE_NT_HEADERS32(vaImitationDataWrapper.At(ctx.elfaNew))
		vaImitationImportDir   = C.PIMAGE_IMPORT_DESCRIPTOR(ctx.vaImitationDataWrapper.At(ctx.importDirVA))
		relocsPos              = 0
	)

	copy(vaImitationDataWrapper.Data, ctx.dataWrapper.Data[:ctx.sizeOfHeaders])
	for i := uintptr(0); i < ctx.numberOfSections; i++ {
		var sectionHeader = C.PIMAGE_SECTION_HEADER(vaImitationDataWrapper.At(ctx.offsetToSectionHeader + C.IMAGE_SECTION_HEADER_SIZE*i))
		copy(vaImitationDataWrapper.Data[sectionHeader.VirtualAddress:], ctx.dataWrapper.Data[sectionHeader.PointerToRawData:][:sectionHeader.SizeOfRawData])
		sectionHeader.SizeOfRawData = 0
		sectionHeader.PointerToRawData = 0
		C.memset(unsafe.Pointer(&sectionHeader.Name[0]), 0, 8)
	}

	C.memset_ull(vaImitationDataWrapper.At(C.IMAGE_DOS_HEADER_SIZE), 0, C.ULONGLONG(ctx.elfaNew)-C.IMAGE_DOS_HEADER_SIZE)
	ntHeaders.FileHeader.Machine = 0x01C0
	ntHeaders.OptionalHeader.DataDirectory[C.IMAGE_DIRECTORY_ENTRY_BASERELOC].VirtualAddress = 0
	ntHeaders.OptionalHeader.DataDirectory[C.IMAGE_DIRECTORY_ENTRY_IMPORT].VirtualAddress = 0
	ntHeaders.OptionalHeader.AddressOfEntryPoint = 0
	ntHeaders.OptionalHeader.ImageBase = 0
	ntHeaders.OptionalHeader.SizeOfImage = 0
	ntHeaders.OptionalHeader.SizeOfHeaders = 0

	var vaImitationBaseRelocIt = C.PIMAGE_BASE_RELOCATION(vaImitationDataWrapper.At(baseRelocVA))
	for vaImitationBaseRelocIt.VirtualAddress > 0 {
		var relocItem = C.PUSHORT(unsafe.Pointer(uintptr(unsafe.Pointer(vaImitationBaseRelocIt)) + C.IMAGE_BASE_RELOCATION_SIZE))
		var numItems = uintptr((vaImitationBaseRelocIt.SizeOfBlock - C.IMAGE_BASE_RELOCATION_SIZE) / C.USHORT_SIZE)

		for i := uintptr(0); i < numItems; i++ {
			switch *relocItem >> 12 {
			case C.IMAGE_REL_BASED_ABSOLUTE:
				break
			case C.IMAGE_REL_BASED_HIGHLOW:
				*C.PULONG(vaImitationDataWrapper.At(uintptr(vaImitationBaseRelocIt.VirtualAddress) + uintptr(*relocItem&0xfff))) += C.ULONG(diff)
				break
			default:
				return nil
			}
			relocItem = C.PUSHORT(unsafe.Pointer(uintptr(unsafe.Pointer(relocItem)) + C.USHORT_SIZE))
		}

		vaImitationBaseRelocIt = C.PIMAGE_BASE_RELOCATION(unsafe.Pointer(uintptr(unsafe.Pointer(vaImitationBaseRelocIt)) + uintptr(vaImitationBaseRelocIt.SizeOfBlock)))
	}

	var vaImitationImportDirIt = vaImitationImportDir
	for vaImitationImportDirIt.u1 /* Characteristics */ != 0 {
		var (
			origFirstThunk = C.PIMAGE_THUNK_DATA32(vaImitationDataWrapper.At(uintptr(vaImitationImportDirIt.u1 /* OriginalFirstThunk */)))
			firstThunk     = C.PIMAGE_THUNK_DATA32(vaImitationDataWrapper.At(uintptr(vaImitationImportDirIt.FirstThunk)))
		)

		relocsPos += 4

		for origFirstThunk.u1 /* AddressOfData */ != 0 {
			if relocsPos >= len(relocs) {
				return nil
			}
			origFirstThunk.u1 = 0
			firstThunk.u1 = C.ULONG(binary.LittleEndian.Uint32(relocs[relocsPos:]))
			relocsPos += 4
			origFirstThunk = C.PIMAGE_THUNK_DATA32(unsafe.Pointer(uintptr(unsafe.Pointer(origFirstThunk)) + C.IMAGE_THUNK_DATA32_SIZE))
			firstThunk = C.PIMAGE_THUNK_DATA32(unsafe.Pointer(uintptr(unsafe.Pointer(firstThunk)) + C.IMAGE_THUNK_DATA32_SIZE))
		}

		vaImitationImportDirIt = C.PIMAGE_IMPORT_DESCRIPTOR(unsafe.Pointer(uintptr(unsafe.Pointer(vaImitationImportDirIt)) + C.IMAGE_IMPORT_DESCRIPTOR_SIZE))
	}

	return vaImitationDataWrapperData
}
