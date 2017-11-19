```
sudo apt-get update
sudo apt-get upgrade

sudo snap install docker
sudo snap install go --classic

sudo curl -o /usr/local/bin/docker-compose -L "https://github.com/docker/compose/releases/download/1.15.0/docker-compose-$(uname -s)-$(uname -m)


mkdir ~/gopath
echo "export GOPATH=~/gopath" > ~/.bash_aliases
echo "export GOROOT=/snap/go/1016" >> ~/.bash_aliases
echo "export PATH=$PATH:$GOPATH/bin" >> ~/.bash_aliases

go get github.com/kardianos/govendor
go get github.com/cloudrimmers/imt2681-assignment3

cd $GOPATH/src/github.com/cloudrimmers/imt2681-assignment3

govendor sync
```
