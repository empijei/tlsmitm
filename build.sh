mkdir -p bin
for os in darwin freebsd linux
do 
	for arch in amd64 386
	do 
		GOOS=$os GOARCH=$arch go build -o "bin/${os}_${arch}_tlsmitm"
	done
done

for arch in amd64 386
do 
	GOOS=windows GOARCH=$arch go build -o "bin/windows_${arch}_tlsmitm.exe"
done
