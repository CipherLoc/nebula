## Cross Compile SafeLogic Libraries
There are many challenges with building nebula when using the SafeLogic libraries for crypto and ssl. For one (1), we need to use CGO to compile and link in the 'C' shared libs. Two (2), linking libraries built with differing/same names causes collisions with existing libraries. Three (3), each OS/ARCH is unique when building. Specific cross compilers must be installed for each OS/ARCH. All these differences causes many subtle errors.
This guide is only one solution to the problems of building for Linux amd64, Windows 386 and Windows amd64. It also, assumes we will be building all OS/ARCH from a Linux amd64 environment.

# Nebula Directory Structure
To build all the various outputs it will be necessary to structure the dependencies to limit the linking errors. For brevity I will only show the necessary affected directories. The general layout looks like:

```
.
└── nebula
    ├── build
    │   ├── linux_amd64
    │   │   └── nebula
    │   ├── windows_386
    │   │   └── nebula.exe
    │   └── windows_amd64
    │       └── nebula.exe
    ├── cmd
    │   └── nebula
    │       └── main.go
    └── libs
        ├── linux_amd64
        │   ├── engines
        │   │   ├── lib4758cca.so
        │   │   ├── libaep.so
        │   │   ├── libatalla.so
        │   │   ├── libcapi.so
        │   │   ├── libchil.so
        │   │   ├── libcswift.so
        │   │   ├── libgmp.so
        │   │   ├── libgost.so
        │   │   ├── libnuron.so
        │   │   ├── libpadlock.so
        │   │   ├── libsureware.so
        │   │   └── libubsec.so
        │   ├── libcrypto.so -> libcrypto.so.1.0.0
        │   ├── libcrypto.so.1.0.0
        │   ├── libssl.so -> libssl.so.1.0.0
        │   ├── libssl.so.1.0.0
        │   └── pkgconfig
        │       ├── libcrypto.pc
        │       ├── libssl.pc
        │       └── openssl.pc
        ├── windows_386
        │   ├── engines
        │   │   ├── 4758cca.dll
        │   │   ├── aep.dll
        │   │   ├── atalla.dll
        │   │   ├── capi.dll
        │   │   ├── chil.dll
        │   │   ├── cswift.dll
        │   │   ├── gmp.dll
        │   │   ├── gost.dll
        │   │   ├── nuron.dll
        │   │   ├── padlock.dll
        │   │   ├── sureware.dll
        │   │   └── ubsec.dll
        │   ├── libeay32.dll
        │   ├── libeay32.lib
        │   ├── ssleay32.dll
        │   └── ssleay32.lib
        └── windows_amd64
            ├── engines
            │   ├── 4758cca.dll
            │   ├── aep.dll
            │   ├── atalla.dll
            │   ├── capi.dll
            │   ├── chil.dll
            │   ├── cswift.dll
            │   ├── gmp.dll
            │   ├── gost.dll
            │   ├── nuron.dll
            │   ├── padlock.dll
            │   ├── sureware.dll
            │   └── ubsec.dll
            ├── libeay32.dll
            ├── libeay32.lib
            ├── ssleay32.dll
            └── ssleay32.lib
```

# Where To Find SafeLogic Libraries and DLLs
1. Checkout EngineV2 from the CipherLoc repo
2. Uncompress EngineV2/3rdParty/openssl/CompressedBinaries/openssl102w.zip
3. Uncompress the following files and copy the libraries/dlls to the build_dir/libs/[target]
    - openssl102w-linux-x86_64.tar.gz
    - openssl102w-windows-i386.zip
    - openssl102w-windows-x86_64.zip

# Installing Cross Platform Windows Compilers

`
sudo apt install gcc-mingw-w64 -y
`

This installs two (2) important directories:

1. /usr/i686-w64-mingw32
2. /usr/x86_64-w64-ming32

i686-w64-mingw32 contains the Windows 32 bit compiler and x86_64-w64-ming32 contains the Windows 64 bit compiler.

# Important Symlinks And Where To Create Them
Linking the appropriate SafeLogic libraries needed a symlink generated in the appropriate directory to successfully link the binaries.

1. For Linux_amd64
- Backup original libcrypto.so and libssl.so symlinks in /lib/x86_64-linux-gnu
            
```
    sudo mv /lib/x86_64-linux-gnu/libcrypto.so /lib/x86_64-linux-gnu/libcrypto.so.bak

    sudo mv /lib/x86_64-linux-gnu/libssl.co /lib/x86_64-linux-gnu/libssl.so.bak
```

- Create symlinks to Linux amd64 SafeLogic libraries

```
    sudo ln -s /path_to_nebula/libs/linux_amd64/libcrypto.so /lib/x86_64-linux-gnu/libcrypto.so

    sudo ln -s /path_to_nebula/libs/linux_amd64/libssl.so /lib/x86_64-linux-gnu/libssl.so
```

2. For Windows_386
#### NOTE: The library names are different. libeay32.lib == libcrypto.lib and ssleay32.lib == libssl.lib

```
    sudo ln -s /path_to_nebula/libs/windows_386/libeay32.lib /usr/i686-w64-mingw32/lib/libcrypto.lib

    sudo ln -s /path_to_nebula/libs/windows_386/ssleay32.lib /usr/i686-w64-mingw32/lib/libssl.lib
```

3. For Windows_amd64
#### NOTE: The library names are different. libeay32.lib == libcrypto.lib and ssleay32.lib == libssl.lib

```
    sudo ln -s /path_to_nebula/libs/windows_amd64/libeay32.lib /usr/x86_64-w64-mingw32/lib/libcrypto.lib

    sudo ln -s /path_to_nebula/libs/windows_amd64/ssleay32.lib /usr/x86_64-w64-mingw32/lib/libssl.lib
```

## Go Environment Settings
A set of go environment settings are needed for each build. It is important to understand to overrides for each build.

1. Linux

` GOOS=linux GOARCH=amd64 CGO_ENABLED=1 `

2. Windows_386

` GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-cpp`

3. Windows_amd64

` GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-cpp `

## Go Build Parameters

1. -trimpath
Cleans up the paths in the binary [NOT OPTIONAL]

2. -work
Prints the working path of the build files [OPTIONAL]

3. -v
Verbose - prints what is being compiles to the console [OPTIONAL]

4. -o
Output path and name of the binary

5. Entry point to build [cmd/nebula]
Go entry point to kick off the build [NOT OPTIONAL]

## Cleaning the Cache
I found it necessary to clean the cache between each build. Numerous errors occurred if I did not.

`
go clean -cache
`

## Example Makefile

```
linux_amd64: linux_clean
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -trimpath -work -v -o build/linux_amd64/nebula ./cmd/nebula

windows_386: windows_386_clean
	GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-cpp go build -trimpath -work -v -o build/windows-386/nebula.exe ./cmd/nebula

windows_amd64: windows_amd64_clean
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-cpp go build -trimpath -work -v -o build/windows-amd64/nebula.exe ./cmd/nebula

linux_clean:
	go clean -cache
	rm -f build/linux_amd64/nebula

windows_386_clean:
	go clean -cache
	rm -f build/windows-386/nebula.exe

windows_amd64_clean:
	go clean -cache
	rm -f buile/windows-amd64/nebula.exe

.PHONY: clean
```

