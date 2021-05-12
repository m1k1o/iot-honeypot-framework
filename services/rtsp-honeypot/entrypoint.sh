#!/bin/sh

# https://ipvm.com/forums/video-surveillance/topics/how-to-write-an-rtsp-url-honeypot
echo RTSP honeypot is running
nc -ll -p 554 -e /usr/src/url.sh
