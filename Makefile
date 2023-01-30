BINDIR=bin

#.PHONY: pbs

all: a i m w64
#
#pbs:
#	cd pbs/ && $(MAKE)
#

m:
	# build the Go library for ARM

	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 \
	CGO_CFLAGS=-mmacosx-version-min=11.0 \
	CGO_LDFLAGS=-mmacosx-version-min=11.0 \
	MIN_SUPPORTED_MACOSX_DEPLOYMENT_TARGET=11.0 \
	go build --buildmode=c-archive -o $(BINDIR)/stamp-arm64.a pc/*.go

	# build the Go library for AMD
	CGO_CFLAGS=-mmacosx-version-min=11.0 \
	CGO_LDFLAGS=-mmacosx-version-min=11.0 \
	MIN_SUPPORTED_MACOSX_DEPLOYMENT_TARGET=11.0 \
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build --buildmode=c-archive -o $(BINDIR)/stamp-amd64.a pc/*.go

	lipo -create $(BINDIR)/stamp-arm64.a $(BINDIR)/stamp-amd64.a -o $(BINDIR)/stamp.a

	rm $(BINDIR)/stamp-amd64.h  $(BINDIR)/stamp-amd64.a $(BINDIR)/stamp-arm64.a

	mv $(BINDIR)/stamp-arm64.h $(BINDIR)/stamp.h
	cp pc/callback.h $(BINDIR)/

w64:#compile this on windows device and change gcc bin path to mingw32
	cp -f pc/callback.h  $(BINDIR)/

	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	go build -ldflags "-w -s -H windowsgui" -v -o $(BINDIR)/export.a -buildmode=c-archive  pc/*.go
	gcc -m64 -shared -pthread -o $(BINDIR)/libstamp.dll pc/export.c $(BINDIR)/export.a -lWinMM -lntdll -lWS2_32
	gzip -f $(BINDIR)/libstamp.dll

w32:#compile this on windows device and change gcc bin path to mingw32
	cp -f pc/callback.h  $(BINDIR)/

	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=386 \
	go build -ldflags "-w -s -H windowsgui" -v -o $(BINDIR)/export.a -buildmode=c-archive  pc/*.go
	gcc -m32 -shared -pthread -o $(BINDIR)/libstamp.dll pc/export.c $(BINDIR)/export.a -lWinMM -lntdll -lWS2_32
	gzip -f $(BINDIR)/libstamp.dll

a:
	 gomobile bind -v -o $(BINDIR)/stamp.aar -target=android -ldflags=-s github.com/blockchainstamp/go-stamp-uilib/mobile
i:
	gomobile bind -v -o $(BINDIR)/stamp.xcframework -target=ios  -ldflags="-w" -ldflags=-s github.com/blockchainstamp/go-stamp-uilib/mobile

clean:
	gomobile clean
	rm $(BINDIR)/*
