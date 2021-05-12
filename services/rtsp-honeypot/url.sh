#!/bin/sh

read url
echo $url >> /proc/1/fd/1
echo RTSP/1.0 404 Not Found
