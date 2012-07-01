#!/bin/bash

rrdtool create goemon.rrd \
--step '6' \
'DS:Temperature:GAUGE:12:U:U' \
'DS:Phase1:GAUGE:12:U:U' \
'DS:Phase2:GAUGE:12:U:U' \
'DS:Phase3:GAUGE:12:U:U' \
'RRA:AVERAGE:0.5:1:600' \
'RRA:AVERAGE:0.5:300:1460' \
'RRA:AVERAGE:0.5:14400:365' \
'RRA:MIN:0.5:300:1460' \
'RRA:MAX:0.5:300:1460' \
'RRA:MIN:0.5:14400:365' \
'RRA:MAX:0.5:14400:365'
