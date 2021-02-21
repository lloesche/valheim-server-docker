#!/bin/sh
# BepInEx running script
#
# This script is used to run a Unity game with BepInEx enabled.
#
# Usage: Configure the script below and simply run this script when you want to run your game modded.

# -------- SETTINGS --------
# ---- EDIT AS NEEDED ------

# EDIT THIS: The name of the Valheim Server executable
executable_name="valheim_server.x86_64"

# The rest is automatically handled by BepInEx

# Whether or not to enable Doorstop. Valid values: TRUE or FALSE
export DOORSTOP_ENABLE=TRUE

# What .NET assembly to execute. Valid value is a path to a .NET DLL that mono can execute.
export DOORSTOP_INVOKE_DLL_PATH="${PWD}/BepInEx/core/BepInEx.Preloader.dll"

# ----- DO NOT EDIT FROM THIS LINE FORWARD ------
# ----- (unless you know what you're doing) ------

doorstop_libs="${PWD}/doorstop_libs"
arch="x64"
executable_path="${PWD}/${executable_name}"
lib_postfix="so"

export LD_LIBRARY_PATH=./linux64:$LD_LIBRARY_PATH

doorstop_libname=libdoorstop_${arch}.${lib_postfix}
export LD_LIBRARY_PATH="${doorstop_libs}":${LD_LIBRARY_PATH}
export LD_PRELOAD=$doorstop_libname:$LD_PRELOAD
export DYLD_LIBRARY_PATH="${doorstop_libs}"
export DYLD_INSERT_LIBRARIES="${doorstop_libs}/$doorstop_libname"

export templdpath=$LD_LIBRARY_PATH
export SteamAppId=892970

"${PWD}/${executable_name}" "$@"

export LD_LIBRARY_PATH=$templdpath