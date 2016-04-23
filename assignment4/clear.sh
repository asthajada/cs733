#!/bin/sh
 
rm -rf logdirectory*
echo '{"Term":0,"VotedFor":-1}'> statestore1
cp statestore1 statestore2
cp statestore1 statestore3
cp statestore1 statestore4
cp statestore1 statestore5
