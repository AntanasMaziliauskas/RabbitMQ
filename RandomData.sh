#!/bin/bash


echo 'How many documents do you want to insert?'
read amount
#Getting output from 'control listnodes' command into a output variable
#mapfile -t output <<< $(control listnodes)
#Going through the output variable and cutting the string into Node name only
#for i in "${output[@]}"
#do
    node="Node01"
    #node="$(cut -d ' ' -f 1 <<< $i)"
   # status="$(cut -d ' ' -f 4 <<< $i)"
#Checking to see if Node is ONLINE
    #if [ $status = 'ONLINE' ]
    #then
#This for loops through n times to generate random Name, Profession, Age and ID
#Runs 'control upsertonepersonnode' command to add random generated data into node
        for ((n=0;n<$amount;n++)) 
        do 
            NAME=$(cat /dev/urandom | tr -dc 'a-zA-Z' | fold -w 7 | head -n 1)
            PROFESSION=$(cat /dev/urandom | tr -dc 'a-zA-Z' | fold -w 9 | head -n 1)
            AGE=$(cat /dev/urandom | tr -dc '0-9' | fold -w 256 | head -n 1 | sed -e 's/^0*//' | head --bytes 2)
            if [ "$AGE" == "" ]; then
                AGE=0
            fi
#running 'randomObjectID' command to get random ObjectID
            ID=$(randomObjectID)
            go run cmd/server/main.go -node $node -person $ID.$NAME.$AGE.$PROFESSION upsertperson
        done
    #fi 
#done
