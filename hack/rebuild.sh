#!/usr/bin/env bash

set -e

if [[ -z "${TMPDIR}" ]]; then
    TMPDIR="/tmp"
fi

if [[ -z "${BUILD_PLATFORMS}" ]]; then
    BUILD_PLATFORMS="linux windows darwin"
fi

if [[ -z "${BUILD_ARCHS}" ]]; then
    BUILD_ARCHS="amd64 arm64"
fi

for os in $BUILD_PLATFORMS; do
    for arch in $BUILD_ARCHS; do
        # don't build for arm on windows
        if [[ "$os" == "windows" && "$arch" == "arm64" ]]; then
            continue
        fi
        echo "[INFO] Building for $os/$arch"
        if [[ $RACE == "yes" ]]; then
            echo "Building devspace with race detector"
            CGO_ENABLED=1 GOOS=$os GOARCH=$arch go build -race -ldflags "-s -w" -o test/devspace-cli-$os-$arch
        else
            CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "-s -w" -o test/devspace-cli-$os-$arch
        fi
    done
done
echo "[INFO] Built binaries for all platforms in test/ directory"

# if not true, install the binary
if [[ "${SKIP_INSTALL}" != "true" ]]; then
    if command -v sudo &> /dev/null; then
        go build -o test/devspace && sudo mv test/devspace /usr/local/bin/
    else 
        go install .
    fi
    echo "[INFO] Installed devspace binary to /usr/local/bin"
else 
    echo "[INFO] Skipping install of devspace binary"
fi

if [[ $BUILD_PLATFORMS == *"linux"* ]]; then
    cp test/devspace-cli-linux-amd64 test/devspace-linux-amd64 
    cp test/devspace-cli-linux-arm64 test/devspace-linux-arm64
fi

if [ -d "desktop/src-tauri/bin" ]; then
    if [[ $BUILD_PLATFORMS == *"linux"* ]]; then
        cp test/devspace-cli-linux-amd64 desktop/src-tauri/bin/devspace-cli-x86_64-unknown-linux-gnu
        cp test/devspace-cli-linux-arm64 desktop/src-tauri/bin/devspace-cli-aarch64-unknown-linux-gnu
    fi
    if [[ $BUILD_PLATFORMS == *"windows"* ]]; then
        cp test/devspace-cli-windows-amd64 desktop/src-tauri/bin/devspace-cli-x86_64-pc-windows-msvc.exe
    fi
    if [[ $BUILD_PLATFORMS == *"darwin"* ]]; then
        cp test/devspace-cli-darwin-amd64 desktop/src-tauri/bin/devspace-cli-x86_64-apple-darwin
        cp test/devspace-cli-darwin-arm64 desktop/src-tauri/bin/devspace-cli-aarch64-apple-darwin
    fi
echo "[INFO] Copied binaries to desktop/src-tauri/bin"
fi

if [[ $BUILD_PLATFORMS == *"linux"* ]]; then
    rm -R $TMPDIR/devspace-cache 2>/dev/null || true
    mkdir -p $TMPDIR/devspace-cache
    cp test/devspace-cli-linux-amd64 $TMPDIR/devspace-cache/devspace-linux-amd64
    cp test/devspace-cli-linux-arm64 $TMPDIR/devspace-cache/devspace-linux-arm64
    echo "[INFO] Copied binaries to $TMPDIR/devspace-cache"
fi
