#!/bin/sh
gotemplater -f json -d post.json -c post.txt post.html | \
gotemplater -f yaml -d site.yaml -c - site.html > final.html