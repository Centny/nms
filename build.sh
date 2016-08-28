#!/bin/bash
##############################
#####Setting Environments#####
echo "Setting Environments"
set -e
export PATH=$PATH:$GOPATH/bin:$HOME/bin:$GOROOT/bin
##############################
######Install Dependence######
echo "Installing Dependence"
#go get github.com/go-sql-driver/mysql
#go get github.com/Centny/TDb
#go get code.google.com/p/go-uuid/uuid
##############################
#########Running Clear#########
if [ "$1" = "-u" ];then
 echo "Running Clear"
 rm -rf $GOPATH/src/github.com/Centny/fvm
 go get github.com/Centny/fvm/fvm
fi
#########Running Test#########
echo "Running Test"
pkgs="\
 github.com/Centny/nms/nmsrc\
"
# pkgs="\
# github.com/Centny/nms/task\
# github.com/Centny/nms/nmsrc\
# "
echo "mode: set" > a.out
for p in $pkgs;
do
 go test -v --coverprofile=c.out $p
 cat c.out | grep -v "mode" >>a.out
 go install $p
done
gocov convert a.out > coverage.json

##############################
#####Create Coverage Report###
echo "Create Coverage Report"
cat coverage.json | gocov-xml -b $GOPATH/src > coverage.xml
cat coverage.json | gocov-html coverage.json > coverage.html

######
go install github.com/Centny/fvm/fvm
