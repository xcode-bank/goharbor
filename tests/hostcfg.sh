#!/bin/bash
IP=`ip addr s eth0 |grep "inet "|awk '{print $2}' |awk -F "/" '{print $1}'`
#echo $IP
sudo sed "s/reg.mydomain.org/$IP/" -i Deploy/harbor.cfg
