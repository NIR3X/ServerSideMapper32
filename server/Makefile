ifdef OS
	ifeq ($(OS),Windows_NT)
		RM = del /Q
		SLASH = \\
		EXE = .exe
	else
		RM = rm -f
		SLASH = /
		EXE = .out
	endif
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		RM = rm -f
		SLASH = /
		EXE = .out
	endif
endif

GO = go
DEST = ..$(SLASH)server$(EXE)
LDFLAGS = -s -w -extldflags -static
GOFLAGS = -trimpath -tags release -buildmode exe

build:
	$(GO) env -w CGO_ENABLED=1
	$(GO) build -o $(DEST) $(GOFLAGS) -ldflags "$(LDFLAGS)"

clean:
	$(RM) $(DEST)
