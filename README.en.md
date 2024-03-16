cd academy_server
go build acade.go
cd ..
cd academy_client
sudo chmod 777 build.sh
./build.sh
1.start the server:
cd academy_server
./acade settings.yaml
2.start the client to be controlled
cd academy_client
./aca listenmode settings.yaml
3.start the client to control
cd academy_client
./aca [devicecode] settings.yaml [password]

input something in client2 cmd like:
echo " ">>1.py
#insert something in 1.py:
sed -ri '1aprint("hello world")' 1.py
python3 1.py
#delete something in 1.py
sed -ri '1d' 1.py
cat 1.py

20240317-academy1.0beta
Equipped with basic remote and authentication functions