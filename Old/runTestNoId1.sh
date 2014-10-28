xterm  -hold=true -e ./bin/node -ID=00 -localPort=2020 &
sleep 1s
xterm  -hold=true -e ./bin/node -ID=01 -localPort=2021 -remotePort=2020 &
sleep 1s
xterm  -hold=true -e ./bin/node -ID=02 -localPort=2022 -remotePort=2021 &
sleep 1s
xterm  -hold=true -e ./bin/node -ID=03 -localPort=2023 -remotePort=2022 &
sleep 1s
xterm  -hold=true -e ./bin/node -ID=04 -localPort=2024 -remotePort=2023 &
