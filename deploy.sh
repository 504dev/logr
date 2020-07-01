git pull origin master
git submodule update --remote
make front
make build
sudo systemctl restart logr
journalctl -u logr -f