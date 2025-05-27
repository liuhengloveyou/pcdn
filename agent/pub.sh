#!/bin/bash


# ssh -p 58779 root@10.10.10.102
make arm
scp -P 58779 sn.sh pcdnagent.sh pcdnagent pcdnagent.service root@10.10.10.102:/opt/pcdnagent/