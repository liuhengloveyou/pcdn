#!/bin/bash

npm run build

scp -r dist/* root@101.37.182.58:/opt/pcdn/www/