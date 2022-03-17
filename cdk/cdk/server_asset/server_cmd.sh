cd home/ec2-user
touch test_command.txt

wget https://go.dev/dl/go1.18.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf go1.18.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a .bash_profile

mkdir app_folder
unzip app.zip -d app_folder