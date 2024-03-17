cd academy_server
go build acade.go
cd ..
cd academy_client
sudo chmod 777 build.sh
./build.sh
1. start the server:
cd academy_server
./acade settings.yaml
2. start the client to be controlled
cd academy_client
./aca listenmode settings.yaml
3. start the client to control
cd academy_client
./aca [devicecode] settings.yaml [password]

1.向client端文件注入内容以使sed命令生效:
echo " ">>1.py
2.向client端1.py文件写内容:
sed -ri '1aprint("hello world")' 1.py
python3 1.py
3.向client端1.py文件删除第一行内容
sed -ri '1d' 1.py
4.查看client端1.py文件内容
cat 1.py

20240317-academy1.0beta
具备基本的远程与认证功能