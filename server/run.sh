
#!/usr/bin/env bash

killall -9 pcdn-server &>/dev/null
sleep 1

nohup ./pcdn-server &
