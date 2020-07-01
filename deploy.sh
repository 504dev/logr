git pull origin master --recurse-submodules
make front
make build
sudo systemctl restart logr
journalctl -u logr -f