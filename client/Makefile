AR = ar
CXX = g++
CXXFLAGS = -std=c++2a -Wall -O2 -march=x86-64 -fno-rtti
LDFLAGS = -shared -s -static -lwinhttp -lntdll

SRC = DllMain.cpp ServerSideMapper32.cpp
OBJ = $(SRC:.cpp=.o)
DEST = ../client.dll

build: http_request $(OBJ)
	$(CXX) $(OBJ) HttpRequest.cpp/HttpRequest.a -o $(DEST) $(LDFLAGS)

%.o: %.cpp
	$(CXX) $(CXXFLAGS) -c -o $@ $<

http_request:
	$(MAKE) -C HttpRequest.cpp

clean:
	$(MAKE) -C HttpRequest.cpp clean
	rm -f $(OBJ) $(DEST)
