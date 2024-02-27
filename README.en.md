cd Go_study
go build client.go
go build client2.go
go build server.go

1.start the server:
./server
2.start the client to be controlled
./client
3.start the client to control
./client2

input something in client2 cmd like:
echo " ">>1.py
#insert something in 1.py:
sed -ri '1aprint("hello world")' 1.py
python3 1.py
#delete something in 1.py
sed -ri '1d' 1.py
cat 1.py