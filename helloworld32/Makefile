CXX = g++
CXXFLAGS = -std=c++2a -Wall -O2 -march=x86-64 -fno-rtti
LDFLAGS = -shared -s -static

build:
	g++ -m32 $(CXXFLAGS) DllMain.c -o ../helloworld32.dll $(LDFLAGS)

clean:
	rm -f ../helloworld32.dll
