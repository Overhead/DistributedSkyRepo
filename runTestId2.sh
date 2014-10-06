xterm  -hold=true -e ./bin/node -localPort=2020 &
sleep 1s
xterm  -hold=true -e ./bin/node -localPort=2021 -remotePort=2020 &
sleep 1s
xterm  -hold=true -e ./bin/node -localPort=2022 -remotePort=2021 &
sleep 1s
xterm  -hold=true -e ./bin/node -localPort=2023 -remotePort=2022 &
sleep 1s
xterm  -hold=true -e ./bin/node -localPort=2024 -remotePort=2023 &
