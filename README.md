# goweb


#### Setup Directories
```shell
mkdir ~/gowork
cd ~/gowork
mkdir bin
mkdir src
mkdir pkg
```

#### Setup Environment
```shell
cd ~/gowork
export GOPATH=$(pwd)
```

#### Install goweb
```shell
cd ~/gowork/src/
git clone https://github.com/satindergrewal/goweb.git
cd goweb

# install dependencies
go get -u github.com/kataras/iris
```

#### Run program
```shell
cd ~/gowork/src/goweb/
go run main.go
```
