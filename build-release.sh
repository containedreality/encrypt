#!/bin/sh
progs='digest encrypt decrypt'
rm -rf out && mkdir -p out

[ -z "${GOPARM}" ] && export GOARM=7
[ -z "${GOAMD64}" ] && export GOAMD64=v3
[ -z "${CGO_ENABLED}" ] && export CGO_ENABLED=0
[ -z "${ESOTERIC}" ] && export ESOTERIC=0

for prog in $progs; do
    for x in $(go tool dist list); do
        export GOOS=$(echo $x | cut -d'/' -f1)
        export GOARCH=$(echo $x | cut -d'/' -f2)

        if [ "${ESOTERIC}" = 0 ]; then
            case "$GOOS" in
            windows)
                ;;
            linux)
                ;;
            illumos)
                ;;
            *)
                continue
                ;;
            esac

            case "$GOARCH" in
            amd64)
                ;;
            arm64)
                ;;
            *)
                continue
                ;;
            esac
        fi

        [ "${GOARCH}" = 'wasm' ] && continue

        OUT="./out/${prog}-${GOOS}-${GOARCH}"

        if [ "${GOOS}" = 'windows' ]; then
            OUT="${OUT}.exe"
        fi

        go build -ldflags='-s -w' -trimpath -o "$OUT" "./$prog/"
    done
done