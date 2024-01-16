#pragma once

#include <windows.h>
#include <cstdint>
#include <vector>

class CServerSideMapper32 {
protected:
	struct EH4_SCOPETABLE_RECORD {
		INT EnclosingLevel;
		PVOID FilterFunc;
		PVOID HandlerFunc;
	};

	struct EH4_SCOPETABLE {
		INT GSCookieOffset;
		INT GSCookieXOROffset;
		INT EHCookieOffset;
		INT EHCookieXOROffset;
		struct EH4_SCOPETABLE_RECORD ScopeRecord[];
	};

	struct EH4_EXCEPTION_REGISTRATION_RECORD {
		PVOID SavedESP;
		EXCEPTION_POINTERS* ExceptionPointers;
		EXCEPTION_REGISTRATION_RECORD SubRecord;
		struct EH4_SCOPETABLE* EncodedScopeTable; //Xored with the __security_cookie
		UINT TryLevel;
	};

	static LONG NTAPI ExceptionHandlerBegin(PEXCEPTION_POINTERS exceptionInfo);
	static VOID ExceptionHandlerEnd();

public:
	static bool ProcessData(const std::vector<uint8_t>& mmap32Data, std::vector<uint8_t>& processedData);
	static bool MMap32(const std::vector<uint8_t>& mappedDll32);
};
