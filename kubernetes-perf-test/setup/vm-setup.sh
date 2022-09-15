set -e

sudo apt update && sudo apt -y upgrade
sudo apt -y install curl
sudo apt-get -y install zip


wget https://dlcdn.apache.org//jmeter/binaries/apache-jmeter-5.5.tgz
tar -xzf apache-jmeter-5.5.tgz
rm apache-jmeter-5.5.tgz


wget https://jmeter-plugins.org/files/packages/jpgc-casutg-2.10.zip
unzip jpgc-casutg-2.10.zip
rm jpgc-casutg-2.10.zip
cp -r lib/ apache-jmeter-5.5/
rm -r lib

# grpc plugin
wget https://jmeter-plugins.org/files/packages/jmeter-grpc-request-1.2.1.zip
unzip jmeter-grpc-request-1.2.1.zip
rm jmeter-grpc-request-1.2.1.zip
cp -r lib/ apache-jmeter-5.5/
rm -r lib

sudo apt-get install -y openjdk-11-jdk

# gRPC protos for ratelimit service
mkdir -p proto-lib/git-repos/
cd proto-lib/git-repos/

git clone https://github.com/envoyproxy/envoy.git
cp -r envoy/api/* ../

git clone https://github.com/cncf/udpa.git
cp -r udpa/* ../

git clone https://github.com/envoyproxy/protoc-gen-validate.git
cp -r protoc-gen-validate/* ../

cd ..
rm -rf git-repos/
