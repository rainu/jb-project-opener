#!/bin/sh

PATH=$PATH:/opt/bin
SCRIPT_PATH=$(realpath "$0")
SCRIPT_HOME=$(dirname "$SCRIPT_PATH")

CACHE_FILE=$SCRIPT_HOME/.project-cache

if [[ -f "$CACHE_FILE" ]] && [[ "$(find $CACHE_FILE -mmin +1440 | wc -l)" == "0" ]]; then
  #use the cache
  CHOSEN=$(cat $CACHE_FILE | rofi -dmenu)
else
  #override the cache
  CHOSEN=$($SCRIPT_HOME/jb-project-opener | tee $CACHE_FILE | rofi -dmenu)
fi

if [ -z "$CHOSEN" ]; then
  exit 0
fi

PROJECT_PATH=$(echo $CHOSEN | sed 's/^[^:]*: //')
IDE=$(echo $CHOSEN | cut -d\: -f1)

case "$IDE" in
  GoLand)
    goland $PROJECT_PATH
    ;;
  WebStorm)
    webstorm $PROJECT_PATH
    ;;
  IntelliJIdea)
    idea $PROJECT_PATH
    ;;
  PyCharm)
    pycharm $PROJECT_PATH
    ;;
  DataGrip)
    datagrip $PROJECT_PATH
    ;;
  CLion)
    clion $PROJECT_PATH
    ;;
esac
