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

## License
[![GNU AGPLv3 Image](https://www.gnu.org/graphics/agplv3-155x51.png)](https://www.gnu.org/licenses/agpl-3.0.html)  

This program is Free Software: You can use, study share and improve it at your
will. Specifically you can redistribute and/or modify it under the terms of the
[GNU Affero General Public License](https://www.gnu.org/licenses/agpl-3.0.html) as
published by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
