#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export cpwd=`pwd`
export LD_LIBRARY_PATH=/usr/local/lib:/usr/lib
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
o_dir=build
if [ "$2" != "" ];then
	o_dir=$2/nms
fi
rm -rf $o_dir
mkdir -p $o_dir

#### Package ####
n_srv=nms
v_srv=0.0.1
##
d_srv="$n_srv"d
o_srv=$o_dir/$n_srv
mkdir $o_srv
mkdir $o_srv/conf
mkdir $o_srv/conf/clients
mkdir $o_srv/www
if [ "$IG_TEST" == "1" ];then
	echo "Build order ig test executor..."
	go test -c -i -cover -o $o_srv/$n_srv".test" github.com/Centny/nms/nms
	cp -rf *.sh $o_srv
else
	echo "Build order normal executor..."
	go build -o $o_srv/$n_srv github.com/Centny/nms/nms
fi
cp *.properties $o_srv/conf
cp nmstask/*.properties $o_srv/conf
cp -rf www/* www

###
if [ "$1" != "" ];then
	curl -o $o_srv/srvd $1/srvd
	curl -o $o_srv/srvd_i $1/srvd_i
	chmod +x $o_srv/srvd
	chmod +x $o_srv/srvd_i
	echo "./srvd_i \$1 $n_srv \$2 \$3" >$o_srv/install.sh
	chmod +x $o_srv/install.sh
fi
cd $o_dir
zip -r $n_srv.zip $n_srv
cd ../
echo "Package $n_srv..."