#!/bin/bash

DOMAIN=$1
PATH=$2

/bin/cat /etc/letsencrypt/live/$DOMAIN/fullchain.pem /etc/letsencrypt/live/$DOMAIN/privkey.pem > $PATH/$DOMAIN
