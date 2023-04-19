[ ! -d "./filenode" ] && exit 0

cd ./filenode
go build
cd ../nodeorc
go build
cd ../cfggen
go build

cd ..
[ ! -d "./cluster" ] && mkdir "./cluster"
[ ! -d "./cluster/orc" ] && mkdir "./cluster/orc"
[ ! -d "./cluster/node1" ] && mkdir "./cluster/node1"
[ ! -d "./cluster/node2" ] && mkdir "./cluster/node2"
[ ! -d "./cluster/node3" ] && mkdir "./cluster/node3"
[ ! -d "./cluster/node4" ] && mkdir "./cluster/node4"
[ ! -d "./cluster/node5" ] && mkdir "./cluster/node5"

pwd=$(pwd)

cp ./cfggen/orccfgtemplate/app.cfg ./cluster/orc
./cfggen/cfggen 1 8081 $pwd"/cluster/node1" $pwd"/cluster/node1" "./cfggen"
./cfggen/cfggen 2 8082 $pwd"/cluster/node2" $pwd"/cluster/node2" "./cfggen"
./cfggen/cfggen 3 8083 $pwd"/cluster/node3" $pwd"/cluster/node3" "./cfggen"
./cfggen/cfggen 4 8084 $pwd"/cluster/node4" $pwd"/cluster/node4" "./cfggen"
./cfggen/cfggen 5 8085 $pwd"/cluster/node5" $pwd"/cluster/node5" "./cfggen"

cp ./filenode/filenode $pwd"/cluster/node1"
cp ./filenode/filenode $pwd"/cluster/node2"
cp ./filenode/filenode $pwd"/cluster/node3"
cp ./filenode/filenode $pwd"/cluster/node4"
cp ./filenode/filenode $pwd"/cluster/node5"
cp ./nodeorc/nodeorc $pwd"/cluster/orc"

# TODO доработать kill процессов для повторных запусков, пока можно запускать из терминала
#cd $pwd"/cluster/node1" && "./filenode" &
#cd $pwd"/cluster/node2" && "./filenode" &
#cd $pwd"/cluster/node3" && "./filenode" &
#cd $pwd"/cluster/node4" && "./filenode" &
#cd $pwd"/cluster/node5" && "./filenode" &
#cd $pwd"/cluster/orc" && "./nodeorc" &