#!/bin/bash

function run() {
  go run main.go seeds -f && go test http-api/tests -v -count=1
}

function confirm() {
cat <<EOF
🔥🔥🔥这是集成测试，期间会强制重置数据库的全部数据，请注意!!!
🚫 请注意!!!
🚫 请注意!!!
🚫 请注意!!!
EOF
read -r -p "你清楚自己在干什么并确定还是要这么做? [y/N] " response
case "$response" in
    [yY][eE][sS]|[yY])
    run
         ;;
    *)
        echo "拜拜！！！"
        ;;
esac
}

if [[ $1 == -y ]]; then
    run
else
    confirm
fi

