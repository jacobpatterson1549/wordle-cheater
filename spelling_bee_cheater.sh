#!/usr/bin/env sh
#
# A command to find "pangrams" for the New York Times game Spelling Bee.
#
# MIT License
# Copyright (c) 2025 Jacob Patterson

# example: May 17 2025: letters are zaenitl: finds 'initialize' and 'tantalize'
# $ echo zaenitl | sh spelling_bee_cheater.sh

read x \
  && echo $x \
  | sed -E "s/(.)/(?=.*\1)/g" \
  | sed -E "s/(.*)/'^\1.+$'/" \
  | xargs -I {} grep -P {} /usr/share/dict/words \
  | grep -E "^[$x]+$"

