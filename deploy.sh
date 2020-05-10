git pull origin master
make build
sudo systemctl restart logr
journalctl -u logr -f