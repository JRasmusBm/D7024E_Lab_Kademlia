#/bin/sh
trap_ctrlc() {
  (rm .temp.o) &> /dev/null
  (cd ../..) &> /dev/null
  exit 2
}

trap "trap_ctrlc" 2

cd cli/src 
clear 
go build -o ../../.temp.o . 
cd ../../
./.temp.o
rm .temp.o &> /dev/null
