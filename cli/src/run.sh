#/bin/sh
trap_ctrlc() {
  (rm .temp.o) &> /dev/null
  exit 2
}

trap "trap_ctrlc" 2

(clear && go build -o .temp.o . && ./.temp.o) ; (rm .temp.o &> /dev/null)
