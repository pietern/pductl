linux:
	GOOS=linux GOARCH=amd64 go build -o pductl .

copy: linux
	scp ./pductl 192.168.1.116:
