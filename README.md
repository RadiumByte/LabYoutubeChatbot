# LabYoutubeChatbot

**LabYoutubeChatbot** is a broadcast interactivity component, which allows viewers to change video and audio content on broadcast (stream) using commands in chat. This project connects to the StreamServer (https://github.com/RadiumByte/StreamServer) as a HTTP client.

Chatbot will hook at the ewest broadcast on the user's channel.

## Features

- Processing YouTube Live chat.
- Support of command: list of cameras, select camera, active camera.
- Possible support of hardware management from chat.

## Requirements
- Any GNU/Linux distribution.
- OAuth2.0 token from Google Developer Console.

## How to install
1) Clone this repository to any directory.
2) Execute script /install/install_libs.sh
3) Execute script /install/build_project.sh
4) Download OAuth2.0 token file from Google Developer Console and place it in the same forlder with executable.
