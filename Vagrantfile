# -*- mode: ruby -*-
# vi: set ft=ruby :
$script = <<SCRIPT
    apt-get update
    apt-get -y install git mc

    wget https://storage.googleapis.com/golang/go1.10.3.linux-amd64.tar.gz && tar -xvzf go1.10.3.linux-amd64.tar.gz; mv go /usr/local
    wget https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && mv dep-linux-amd64 /usr/local/bin/dep && chmod +x /usr/local/bin/dep

    export GOPATH=/home/vagrant; export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin:/home/vagrant/bin' >> .bash_profile
    echo 'export GOPATH=/home/vagrant' >> .bash_profile

    locale-gen ru_RU.UTF-8
SCRIPT

# Vagrantfile API/syntax version.
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

    config.vm.box = "ubuntu/xenial64"
    config.vm.network "private_network", ip: "192.168.50.249"

    config.vm.synced_folder "./", "/home/vagrant/src/github.com/Puppollo/di"

    config.vm.provider "virtualbox" do |v|
      v.name = "puppollo-di"
      v.memory = 1024
    end

    # provisioner config
    config.vm.provision "file", source: "~/.ssh/id_rsa.pub", destination: "~/.ssh/id_rsa.pub"
    config.vm.provision "file", source: "~/.ssh/id_rsa", destination: "~/.ssh/id_rsa"
    config.vm.provision "file", source: "~/.ssh/known_hosts", destination: "~/.ssh/known_hosts"
    config.vm.provision "shell", inline: $script
end
