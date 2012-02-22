#!/bin/sh

curl "http://127.0.0.1:8080/run?input=$1" 2>/dev/null  | awk 'BEGIN {flag=0}
	/^[ A-z]+/{flag=1}
	/^[^ A-z]+/{flag=0}
	{if (flag == 1) print}
	'



