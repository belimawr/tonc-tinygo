# TONC-TinyGo

This is my attempt to learn GameBoy Advance development in Go which
consists in going through the "[Tonc
tutorial](https://www.coranac.com/tonc/text/toc.htm)" and re-writing
them in Go using [TinyGo](https://tinygo.org/) to compile.

## Structure
The whole repository is a single Go project, each folder is one of
Tonc's chapters/sub-chapters. Chapters 1 and 2 are information about the hardware,
hence the first folder is `03.1-my-first-gba-demo`.

`sandbox` is me going rouge and having fun with the code of that
specific chapter, it has sub folders matching the chapter where I
decided to stop following the tutorial and went rogue.

## Dev env setup
### Download and install TinyGo
Get the latest [TinyGo](https://tinygo.org/) release from
[GitHub](https://github.com/tinygo-org/tinygo/releases/tag/v0.41.1)
directly. Then extract the archive and add its bin folder to your
path. You can test it by running 

```
$ tinygo version
tinygo version 0.41.1 linux/amd64 (using go version go1.26.5 and LLVM version 20.1.1)
```

### Configure Emacs (or your IDE)
That is as simple as setting the following environment variables
`GOROOT`, `GOOS`, `GOARCH`, `GOFLAGS`, `TINYGOROOT`.

`TINYGOROOT` is where you extracted the TinyGo archive, the folder
containing `bin`, `src`, `lib` and `targets`.

The other values you can get directly from the TinyGo binary:
```
% tinygo info gameboy-advance
LLVM triple:       armv4t-unknown-unknown-eabi
GOOS:              linux
GOARCH:            arm
build tags:        gameboyadvance arm7tdmi baremetal linux arm tinygo purego osusergo math_big_pure_go gc.conservative scheduler.none serial.none tinygo.unicore
garbage collector: conservative
scheduler:         none
cached GOROOT:     /home/belimawr/.cache/tinygo/goroot-82cf0cdee6e012bec20653743104fb8bf30730185f4a1879cf802dd09b205722
```

Set your `GOROOT` to the 'cached GOROOT'.

I'm using Emacs and lsp-mode for Go, so I automated that with a
`.dir-locals.el`:
```elisp
((go-mode
  . ((lsp-go-env . ((GOROOT . "/home/belimawr/.cache/tinygo/goroot-82cf0cdee6e012bec20653743104fb8bf30730185f4a1879cf802dd09b205722")
                    (TINYGOROOT . "/home/belimawr/bin/tinygo")
                    ;; (GOOS . "linux")
                    (GOARCH . "arm")
                    (GOFLAGS . "-tags=gameboyadvance,arm7tdmi,baremetal,linux,arm,tinygo,purego,osusergo,math_big_pure_go,gc.conservative,scheduler.none,serial.none,tinygo.unicore"))))))
```

`GOOS` is commented out because I'm already on Linux, but you might
need to set it.

## Emulator
Install [mGBA](https://mgba.io/), for Archlinux install `mgba-qt`, it
has a nice UI and many debug features.

## Compiling
```
tinygo build -o rom.gba -target gameboy-advance .
```

## Running
TinyGo might need `mgba` in your path, so if using `mgba-qt` create a symlink.
```
# Compile and run
tinygo run -target=gameboy-advance main.go

# Run in the emulator
mgba-qt rom.gba
```

Or open the emulator and load the rom

## Resources
 - [Tonc](https://www.coranac.com/tonc/text/toc.htm)
 - [GBATEK](https://problemkaputt.de/gbatek.htm)
 - Building GBA games in GO: [Video](https://www.youtube.com/watch?v=mrWJZSVSRVQ), [slides](https://engineering.getweave.com/talk/gba-games-in-go/go_users_group_slides.pdf)

