# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"

  config.vm.provision "shell", inline: <<-SHELL
    curl -sf -o /tmp/go1.6.1.linux-amd64.tar.gz -L https://storage.googleapis.com/golang/go1.6.1.linux-amd64.tar.gz
    sudo mkdir -p /opt && cd /opt && sudo tar xfz /tmp/go1.6.1.linux-amd64.tar.gz && rm -f /tmp/go1.6.1.linux-amd64.tar.gz
    curl -s https://packagecloud.io/install/repositories/darron/consul/script.deb.sh | sudo bash
    curl -s https://packagecloud.io/install/repositories/darron/consul-cli/script.deb.sh | sudo bash
    sudo apt-get install -y consul git make graphviz consul-cli
    sudo mkdir -p /etc/consul.d /var/lib/consul /var/log/consul
    sudo ln -s /lib/init/upstart-job /etc/init.d/consul
    sudo tee /etc/consul.d/default.json << EOF
{
  "client_addr": "127.0.0.1",
  "data_dir": "/var/lib/consul",
  "server": true,
  "bootstrap": true,
  "recursor": "8.8.8.8",
  "bind_addr": "0.0.0.0",
  "log_level": "debug",
  "node_name": "sifter-consul"
}
EOF
    sudo cat > /etc/profile.d/go.sh << EOF
export GOROOT="/opt/go"
export GOPATH="/home/vagrant/gocode"
export PATH="/opt/go/bin://home/vagrant/gocode/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
export GOSHE_DEBUG=1
EOF
    cd /vagrant && source /etc/profile.d/go.sh
  SHELL
end
