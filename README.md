# ServerSideMapper32 - Server-Side DLL Mapping Solution

This is a C++ and Go project that implements a server-side DLL manual mapping mechanism. The project includes a client-side component written in C++, a server-side component written in Go, and a "Hello World" DLL source.

## Installation

To use this project, you can clone the repository and compile it using C++ and Go compilers:

```bash
git clone https://github.com/NIR3X/ServerSideMapper32 --recurse-submodules
cd ServerSideMapper32
msys2_shell -mingw32 -defterm -here -no-start -lc 'mingw32-make -C client'
msys2_shell -mingw32 -defterm -here -no-start -lc 'mingw32-make -C helloworld32'
make -C server
```
