#!/usr/bin/env sh

# Usage: prompt_var <var_name> <prompt_msg> [default_value]
# Prompts for a variable with a default value if provided.

prompt_var() {
  local var_name=$1
  local prompt_msg=$2
  local default=$3
  if [ -n "$default" ]; then
    read -p "$prompt_msg [$default]: " var_value
    if [ -z "$var_value" ]; then
      var_value=$default
    fi
  else
    read -p "$prompt_msg: " var_value
  fi
  eval "$var_name=\"\$var_value\""
}