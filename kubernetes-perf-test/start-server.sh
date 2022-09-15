HEAP="-Xms12g -Xmx12g"
cd ~/apache-jmeter-5.5/bin/
nohup ./jmeter-server &
tail -f nohup.out
