# redis-stack-server with systemd
sudo systemctl daemon-reload
sudo systemctl start redis-stack-server
sudo systemctl stop redis-stack-server
sudo systemctl status redis-stack-server
sudo vim /etc/systemd/system/redis-stack-server.service
sudo vim /etc/redis-stack.conf

# View logs
journalctl -u redis-stack-server

# kill process binding to port
sudo kill -9 $(sudo lsof -t -i:8080)

# run ansible with local inventory
ansible-playbook ec2_redis_stack.ansible.yml -i local_inventory.yml
