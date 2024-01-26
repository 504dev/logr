git pull origin master --tags --recurse-submodules
docker-compose up --build -d
docker system prune -a -f