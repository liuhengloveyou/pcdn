#!/bin/bash

pnpm run build

scp -r dist/* root@101.37.182.58:/opt/pcdn/www/

rm -fr dist