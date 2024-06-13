#!/bin/bash

bracket="============================================="

# Fungsi untuk memutar file suara dengan volume yang disesuaikan
play_sound() {
    local volume="$1"
    local sound_file="$2"
    ffplay -nodisp -loglevel quiet -af "volume=$volume" -autoexit "$sound_file" >/dev/null 2>&1
}

# Cek apakah ada parameter masuk
if [ $# -eq 0 ]; then
    echo "Usage: $0 <commit_message>"
    play_sound 2 "./sound/error.mp3"
    exit 1
fi

# Memeriksa hasil build untuk kata "error"

build_output=$(yarn build 2>&1)
if echo "$build_output" | grep -qi "error"; then
    echo "$bracket"
    echo "# Kode ReactJS masih ada yang error !!"
    echo ""
    echo "$build_output"
    play_sound 2 "./sound/error.mp3"
    exit 1
fi

build_output=$(go build 2>&1)
if echo "$build_output" | grep -q "^# "; then
    echo "$bracket"
    echo "# Kode Golang masih ada yang error !!"
    echo ""
    echo "$build_output"
    play_sound 2 "./sound/error.mp3"
    exit 1
fi

echo ""
echo "$bracket"
echo "# Kode tidak ada masalah."
play_sound 1 "./sound/good-code.mp3"
echo ""

git add .
git commit -m "$*"
git push

echo ""
echo "$bracket"
echo "# Sukses push"
play_sound 1 "./sound/finish.mp3"

echo "$bracket" # Exit...
