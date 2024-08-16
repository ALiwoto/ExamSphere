#!/bin/bash

# Split the password into two parts
first_part="${1:0:${#1}/2}"
second_part="${1:${#1}/2}"

# Base64 encode the parts
encoded_first_part=$(echo -n "$first_part" | base64)
encoded_second_part=$(echo -n "$second_part" | base64)

# Concatenate the parts with the specific pattern and base64 encode the whole string
encoded_password=$(echo -n "passM66QFT_${encoded_first_part}_s5rS58O0O3ML_${encoded_second_part}_RendPassTS5S" | base64)

# Print the final encoded password
echo "$encoded_password"