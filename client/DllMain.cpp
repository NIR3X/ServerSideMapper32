#include "HttpRequest.cpp/HttpRequest.h"
#include "ServerSideMapper32.h"
#include <windows.h>

/**
 * @brief Entry point for the main thread.
 * 
 * This function is the entry point for the main thread of the application.
 * It retrieves the mmap data from a server, processes it, sends the processed data,
 * retrieves the mapped DLL from the server, and maps the DLL into the current process.
 * 
 * @param lpParam A pointer to the thread parameters (not used in this implementation).
 * @return DWORD The exit code of the thread.
 */
DWORD WINAPI mainThread(LPVOID lpParam){
	uint32_t statusCode = 0;
	auto mmapData = CHttpRequest::Request(L"http://127.0.0.1:8000/get_mmap_data", CHttpRequest::COptions(), statusCode);
	if (mmapData.size() == 0 || statusCode != 200) {
		MessageBoxA(NULL, "Failed to get mmap data", "Error", MB_OK);
		return 1;
	}

	std::vector<uint8_t> processedData = {};
	if (!CServerSideMapper32::ProcessData(mmapData, processedData)) {
		MessageBoxA(NULL, "Failed to process mmap data", "Error", MB_OK);
	}

	CHttpRequest::COptions options = {};
	options.data = processedData;
	auto mmapedDll = CHttpRequest::Request(L"http://127.0.0.1:8000/get_mapped_dll", options, statusCode);
	if (mmapedDll.size() == 0 || statusCode != 200) {
		MessageBoxA(NULL, "Failed to get mapped DLL", "Error", MB_OK);
		return 1;
	}

	if (!CServerSideMapper32::MMap32(mmapedDll)) {
		MessageBoxA(NULL, "Failed to mmap DLL", "Error", MB_OK);
		return 1;
	}

	return 0;
}

BOOL WINAPI DllMain(HMODULE hModule, DWORD dwReason, LPVOID lpReserved){
	if (dwReason == DLL_PROCESS_ATTACH) {
		CreateThread(NULL, 0, &mainThread, lpReserved, 0, NULL);
	}
	return TRUE;
}
