#！/bin/bash

for((i=0;i<10;i++)); do
    sleep 3
    echo "${i}"
    if (( i%2 == 0  )); then
       ln  -sf  case2 example
    else
       ln  -sf  case1 example
    fi
    
done
