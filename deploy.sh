git pull origin master --tags --recurse-submodules
make front
make build
sudo systemctl restart logr
journalctl -u logr -f