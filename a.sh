#!/bin/bash

# These should be replaced with your actual values
# note that it is not good practice to store these here
# for purposes of this guide it's okay. 
CHANNEL_ID='<CHAT GROUP ID>'
BOT_TOKEN='<BOT TOKEN>'

# if the first argument is "-h" for help.
# display the usage information
if [ "$1" == "-h" ]; then
  echo "Usage: `basename $0` \"text message\""
  exit 0
fi

# Warn if the first argument is empty.
if [ -z "$1" ]; then
    echo "Add message text as the second argument"
    exit 0
fi

# Warn if more than one argument is passed to the script.
if [ "$#" -ne 1 ]; then
    echo "You can pass only one argument. For string with spaces, put it on quotes"
    exit 0
fi

# send a POST request to Telegram's API
# The '-s' = "silent" mode
# redirect outputs to /dev/null
curl -s --data "text=$1" --data "chat_id=$CHANNEL_ID" 'https://api.telegram.org/bot'$BOT_TOKEN'/sendMessage' 