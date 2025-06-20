
#!/usr/bin/env bash

killall -9 pcdn-server &>/dev/null

nohup ./pcdn-server &
