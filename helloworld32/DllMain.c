/**
 * @file DllMain.c
 * @brief This file contains the implementation of a hello world DLL.
 */

#include <windows.h>

/**
 * @brief Entry point for the main thread of the DLL.
 * @param lpParam A pointer to the thread parameter.
 * @return The thread exit code.
 */
DWORD WINAPI mainThread(LPVOID lpParam){
	MessageBoxA(NULL, "Hello World!", "helloworld32", MB_OK);
	return 0;
}

/**
 * @brief Entry point for the DLL.
 * @param hModule A handle to the DLL module.
 * @param dwReason The reason for calling the DLL entry point function.
 * @param lpReserved Reserved parameter.
 * @return TRUE if the DLL initialization is successful, otherwise FALSE.
 */
BOOL WINAPI DllMain(HMODULE hModule, DWORD dwReason, LPVOID lpReserved){
	if (dwReason == DLL_PROCESS_ATTACH) {
		CreateThread(NULL, 0, &mainThread, lpReserved, 0, NULL);
	}
	return TRUE;
}
