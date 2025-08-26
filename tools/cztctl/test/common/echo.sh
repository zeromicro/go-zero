#!/bin/bash

function console_red() {
	echo -e "\033[31m "$*" \033[0m"
}

function console_green() {
	echo -e "\033[32m "$*" \033[0m"
}

function console_yellow() {
	echo -e "\033[33m "$*" \033[0m"
}

function console_blue() {
	echo -e "\033[34m "$*" \033[0m"
}

function console_tip() {
   console_blue "========================== $* =============================="
}

function console_step() {
   console_blue "<<<<<<<<<<<<<<<< $* >>>>>>>>>>>>>>>>"
}
