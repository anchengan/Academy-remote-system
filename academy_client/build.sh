go build -o ./build/libsender.so -buildmode=c-shared sender.go
go build -o ./build/libreceiver.so -buildmode=c-shared receiver.go
gcc -o build/aca aca.c yamlreader.c macreader.c codereader.c sha256tool.c -lyaml  -lcrypto -L. -lsender -lreceiver
