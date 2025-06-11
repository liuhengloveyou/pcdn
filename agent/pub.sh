#!/bin/bash


# ssh -p 58779 root@10.10.10.129
# te6ryiFCB9zc
make arm

# scp -P 58779 sn.sh pcdnagent.sh pcdnagent pcdnagent.service root@10.10.10.100:/opt/pcdnagent/
scp -P 55555 sn.sh pcdnagent.sh pcdnagent pcdnagent.service root@10.10.10.132:/opt/pcdnagent/
# scp -P 58779 sn.sh pcdnagent.sh pcdnagent pcdnagent.service root@10.10.10.110:/opt/pcdnagent/