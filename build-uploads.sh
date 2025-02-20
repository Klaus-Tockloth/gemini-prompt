#!/bin/sh

# ------------------------------------
# Purpose:
# - Builds uploads (tar.gz or zip) for Github project repository (assets in release section).
#
# Releases:
# - v1.0.0 - 2025/02/20: initial release
# ------------------------------------

# set -o xtrace
set -o verbose

# recreate directory
rm -r ./uploads
mkdir ./uploads

# uploads 'darwin'
tar -cvzf ./uploads/macos-amd64_gemini-prompt.tar.gz ./binaries/darwin-amd64/gemini-prompt
tar -cvzf ./uploads/macos-arm64_gemini-prompt.tar.gz ./binaries/darwin-arm64/gemini-prompt

# uploads 'freebsd'
tar -cvzf ./uploads/freebsd-amd64_gemini-prompt.tar.gz ./binaries/freebsd-amd64/gemini-prompt
tar -cvzf ./uploads/freebsd-arm64_gemini-prompt.tar.gz ./binaries/freebsd-arm64/gemini-prompt

# uploads 'linux'
tar -cvzf ./uploads/linux-amd64_gemini-prompt.tar.gz ./binaries/linux-amd64/gemini-prompt
tar -cvzf ./uploads/linux-arm64_gemini-prompt.tar.gz ./binaries/linux-arm64/gemini-prompt

# uploads 'netbsd'
tar -cvzf ./uploads/netbsd-amd64_gemini-prompt.tar.gz ./binaries/netbsd-amd64/gemini-prompt
tar -cvzf ./uploads/netbsd-arm64_gemini-prompt.tar.gz ./binaries/netbsd-arm64/gemini-prompt

# uploads 'openbsd'
tar -cvzf ./uploads/openbsd-amd64_gemini-prompt.tar.gz ./binaries/openbsd-amd64/gemini-prompt
tar -cvzf ./uploads/openbsd-arm64_gemini-prompt.tar.gz ./binaries/openbsd-arm64/gemini-prompt

# uploads 'windows'
zip ./uploads/windows-amd64_gemini-prompt.zip ./binaries/windows-amd64/gemini-prompt.exe
zip ./uploads/windows-arm_gemini-prompt.zip ./binaries/windows-arm/gemini-prompt.exe
