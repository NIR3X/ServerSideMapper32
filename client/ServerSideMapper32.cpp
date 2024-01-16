#include "ServerSideMapper32.h"
#include <ntdef.h>
#include <ntstatus.h>
#include <algorithm>
#include <sstream>
#include <string>
#include <thread>

extern "C" {
	NTSYSAPI NTSTATUS NTAPI LdrGetProcedureAddress(
		IN HMODULE              ModuleHandle,
		IN PANSI_STRING         FunctionName OPTIONAL,
		IN WORD                 Oridinal OPTIONAL,
		OUT LPVOID              *FunctionAddress
	);
}

#pragma GCC push_options
#pragma GCC optimize("O0")


static const UINT_PTR exceptionHandlerReplacement = 0x135772c0;

/**
 * @brief ExceptionHandlerBegin is an exception handler shell used to handle manual mapped DLL exceptions.
 * 
 * This function is designed to handle exceptions that occur within a manually mapped DLL. It is intended to be copied as byte code to the destination process, so compiler optimization is disabled to ensure the function can be accurately copied.
 * 
 * @param exceptionInfo A pointer to an EXCEPTION_POINTERS structure that contains information about the exception.
 * @return LONG The return value depends on the outcome of the exception handling. If the exception is successfully handled, the function returns EXCEPTION_CONTINUE_EXECUTION. If the exception is not handled, the function returns EXCEPTION_CONTINUE_SEARCH.
 */
LONG NTAPI CServerSideMapper32::ExceptionHandlerBegin(PEXCEPTION_POINTERS exceptionInfo) {
	/**
	 * @brief The starting address of the image to be patched.
	 * 
	 * The value stored in this variable, exceptionHandlerReplacement, will be replaced by the image begin address
	 * through bytecode patching.
	 */
	UINT_PTR imageBeginAddr = exceptionHandlerReplacement;

	/**
	 * @brief The ending address of the image to be patched.
	 * 
	 * The value stored in this variable, exceptionHandlerReplacement + 1, will be replaced by the image end address
	 * through bytecode patching.
	 */
	UINT_PTR imageEndAddr = exceptionHandlerReplacement + 1;

	if (!exceptionInfo) {
		return imageEndAddr - imageBeginAddr;
	}

	// // prevent the exception handler from being called for exceptions outside the image
	// PVOID exceptionAddress = exceptionInfo->ExceptionRecord->ExceptionAddress;
	// if ((UINT_PTR)exceptionAddress < imageBeginAddr || (UINT_PTR)exceptionAddress > imageEndAddr) {
	// 	return EXCEPTION_CONTINUE_SEARCH;
	// }

	EXCEPTION_REGISTRATION_RECORD* pFs =
		(EXCEPTION_REGISTRATION_RECORD*)__readfsdword(0); // get the current thread's exception handler list
	if((DWORD_PTR)pFs > 0x1000 && (DWORD_PTR)pFs < 0xFFFFFFF0) {
		struct EH4_EXCEPTION_REGISTRATION_RECORD* record = CONTAINING_RECORD(pFs, struct EH4_EXCEPTION_REGISTRATION_RECORD, SubRecord);
		EXCEPTION_ROUTINE* handler = record->SubRecord.Handler;

		if((UINT_PTR)handler > imageBeginAddr && (UINT_PTR)handler < imageEndAddr) {
			// call the original exception handler
			EXCEPTION_DISPOSITION exceptionDisposition = handler(
				exceptionInfo->ExceptionRecord, &record->SubRecord, exceptionInfo->ContextRecord, nullptr
			);
			if (exceptionDisposition == ExceptionContinueExecution) {
				return EXCEPTION_CONTINUE_EXECUTION;
			}
		}
	}
	return EXCEPTION_CONTINUE_SEARCH;
}

VOID CServerSideMapper32::ExceptionHandlerEnd() {
}

#pragma GCC pop_options

/**
 * @brief Loads DLL imports and functions by ordinals and names.
 * 
 * This function takes a vector of mmap32Data, which contains the data for loading DLL imports and functions.
 * It processes the data and populates the processedData vector with the loaded imports and functions.
 * 
 * @param mmap32Data The input vector containing the mmap32Data.
 * @param processedData The output vector to store the processed data.
 * @return Returns true if the data is successfully processed and loaded, false otherwise.
 */
bool CServerSideMapper32::ProcessData(const std::vector<uint8_t>& mmap32Data, std::vector<uint8_t>& processedData) {
	processedData = {};

	size_t mmap32DataSize = mmap32Data.size();
	if (mmap32DataSize < sizeof(ULONG) + sizeof(ULONG)) {
		return false;
	}

	uintptr_t offset = 0;
	ULONG imageBase = *(PULONG)&mmap32Data[offset]; offset += sizeof(ULONG);
	ULONG sizeOfImage = *(PULONG)&mmap32Data[offset]; offset += sizeof(ULONG);

	LPVOID vaImageBase = VirtualAlloc(nullptr, sizeOfImage, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
	if (!vaImageBase) {
		return false;
	}

	auto insertee = (uint8_t*)&vaImageBase;
	processedData.insert(processedData.end(), insertee, insertee + sizeof(ULONG));
	ULONG diff = (ULONG)vaImageBase - imageBase;
	insertee = (uint8_t*)&diff;
	processedData.insert(processedData.end(), insertee, insertee + sizeof(ULONG));

	HMODULE importHMod = nullptr;
	for (uintptr_t i = offset; i < mmap32DataSize; ++i) {
		switch (mmap32Data[i]) {
		case 255: {
				std::stringstream ss = {};
				bool isTruncated = true;
				for (uintptr_t j = i + 1; j < mmap32DataSize; ++j) {
					if (mmap32Data[j] == 0) {
						i = j;
						isTruncated = false;
						break;
					}
					ss << (char)mmap32Data[j];
				}
				if (isTruncated) {
					return false;
				}
				std::string importName = ss.str();
				importHMod = LoadLibraryA(importName.c_str());
				if (!importHMod) {
					MessageBoxA(NULL, "Failed to import library", "Error", MB_OK);
					return false;
				}
				auto insertee = (uint8_t*)&importHMod;
				processedData.insert(processedData.end(), insertee, insertee + sizeof(ULONG));
			}
			break;
		case 254: 
			if (i + sizeof(USHORT) < mmap32DataSize) {
				USHORT ordinal = *(PUSHORT)&mmap32Data[i + 1];
				LPVOID importProc = nullptr;
				if (LdrGetProcedureAddress(importHMod, nullptr, ordinal + 1 /* Adjusting for optional zero ordinal by incrementing */, &importProc) != STATUS_SUCCESS) {
					MessageBoxA(NULL, (std::string("Failed to import function by ordinal ") + std::to_string(ordinal)).c_str(), "Error", MB_OK);
					return false;
				}
				i += sizeof(USHORT);
				auto insertee = (uint8_t*)&importProc;
				processedData.insert(processedData.end(), insertee, insertee + sizeof(ULONG));
			} else {
				return false;
			}
			break;
		case 253: {
				std::stringstream ss = {};
				bool isTruncated = true;
				for (uintptr_t j = i + 1; j < mmap32DataSize; ++j) {
					if (mmap32Data[j] == 0) {
						i = j;
						isTruncated = false;
						break;
					}
					ss << (char)mmap32Data[j];
				}
				if (isTruncated) {
					return false;
				}
				std::string importProcName = ss.str();
				SIZE_T importProcNameLen = importProcName.length();
				if (importProcNameLen > 0xffff) {
					importProcNameLen = 0xffff;
				}
				ANSI_STRING importProcNameAnsi = {
					.Length = (USHORT)importProcNameLen,
					.MaximumLength = (USHORT)importProcNameLen,
					.Buffer = importProcName.data(),
				};
				LPVOID importProc = nullptr;
				if (LdrGetProcedureAddress(importHMod, &importProcNameAnsi, 0, &importProc) != STATUS_SUCCESS) {
					MessageBoxA(NULL, (std::string("Failed to import function ") + importProcName).c_str(), "Error", MB_OK);
					return false;
				}
				auto insertee = (uint8_t*)&importProc;
				processedData.insert(processedData.end(), insertee, insertee + sizeof(ULONG));
			}
			break;
		default:
			return false;
		}
	}
	return true;
}

using DllEntryPointType = BOOL(WINAPI*)(HINSTANCE, DWORD, LPVOID);

/**
 * @brief Injects an exception handler and a mapped DLL and calls its entry point.
 * 
 * This function takes a vector of uint8_t representing the mapped DLL and performs the following steps:
 * 1. Checks the size of the exception handler and ensures it is valid.
 * 2. Validates the size of the mapped DLL.
 * 3. Retrieves the image base address, size of the image, and address of the entry point from the mapped DLL.
 * 4. Copies the image data from the mapped DLL to the virtual address space.
 * 5. Replaces the exception handler placeholders in the exception handler shell code with the image base address and image end address.
 * 6. Allocates memory for the exception handler shell and copies the modified exception handler shell code to it.
 * 7. Adds the exception handler to the vectored exception handler chain.
 * 8. Creates a new thread and calls the DLL entry point with the image base address and the address of the entry point.
 * 
 * @param mappedDll32 A vector of uint8_t representing the mapped DLL.
 * @return True if the injection and execution were successful, false otherwise.
 */
bool CServerSideMapper32::MMap32(const std::vector<uint8_t>& mappedDll32) {
	SSIZE_T exceptionHandlerSize = (INT_PTR)&ExceptionHandlerEnd - (INT_PTR)&ExceptionHandlerBegin;
	if (exceptionHandlerSize <= 0) {
		MessageBoxA(NULL, "invalid exception handler size", "Error", MB_OK);
		return false;
	}

	size_t mappedDll32Size = mappedDll32.size();
	if (mappedDll32Size < sizeof(ULONG) + sizeof(ULONG) + sizeof(ULONG)) {
		return false;
	}
	uintptr_t offset = 0;
	ULONG vaImageBase = *(PULONG)&mappedDll32[offset]; offset += sizeof(ULONG);
	ULONG sizeOfImage = *(PULONG)&mappedDll32[offset]; offset += sizeof(ULONG);
	ULONG addressOfEntryPoint = *(PULONG)&mappedDll32[offset]; offset += sizeof(ULONG);
	if (mappedDll32Size - offset != sizeOfImage) {
		return false;
	}
	const uint8_t* image = &mappedDll32[offset];
	for (uintptr_t i = 0; i < sizeOfImage; ++i) {
		((PUINT8)vaImageBase)[i] = image[i];
	}

	UINT8 exceptionHandlerRaw[exceptionHandlerSize];
	std::copy_n((LPCBYTE)&ExceptionHandlerBegin, exceptionHandlerSize, exceptionHandlerRaw);
	INT replaced = 0;
	for (UINT_PTR i = 0; i <= exceptionHandlerSize - sizeof(UINT_PTR); ++i) {
		PUINT_PTR replacement = (PUINT_PTR)&exceptionHandlerRaw[i];
		switch (*replacement) {
		case exceptionHandlerReplacement:
			*replacement = vaImageBase;
			++replaced;
			break;
		case exceptionHandlerReplacement + 1:
			*replacement = vaImageBase + sizeOfImage;
			++replaced;
			break;
		}
	}
	if (replaced != 2) {
		MessageBoxA(NULL, "failed to replace exception handler addresses", "Error", MB_OK);
		return false;
	}
	LPVOID exceptionHandlerShell = VirtualAlloc(nullptr, exceptionHandlerSize, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
	if (!exceptionHandlerShell) {
		MessageBoxA(NULL, "failed to allocate memory for exception handler shell", "Error", MB_OK);
		return false;
	}
	std::copy_n(exceptionHandlerRaw, exceptionHandlerSize, (PUINT8)exceptionHandlerShell);
	AddVectoredExceptionHandler(1, (PVECTORED_EXCEPTION_HANDLER)exceptionHandlerShell);

	std::thread([](HINSTANCE hInst, DllEntryPointType dllEntryPoint) {
		dllEntryPoint(hInst, DLL_PROCESS_ATTACH, nullptr);
	}, (HINSTANCE)vaImageBase, (DllEntryPointType)(vaImageBase + addressOfEntryPoint)).detach();
	return true;
}
