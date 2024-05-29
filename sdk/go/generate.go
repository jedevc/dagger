package dagger

//go:generate sh -c "env; echo ---- /etc/hosts; cat /etc/hosts; echo ---- /etc/resolv.conf; cat /etc/resolv.conf; echo ---- pings; ping -c 4 dagger-engine; ping -c 4 $MYHOST; go -C ../../ run ./cmd/codegen --output ./sdk/go"
