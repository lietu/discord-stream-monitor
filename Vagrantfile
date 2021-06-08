Vagrant.configure("2") do |cfg|
    cfg.vm.box = "aspyatkin/ubuntu-20.04-server"
    cfg.vm.provision "shell", path: "setup.sh"

    # SSH agent forwarding makes life easier
    cfg.ssh.forward_agent = true

    vm_name = "dsm"
    cfg.vm.define :dsm do |s|
        s.vm.network "private_network", type: "dhcp"
        s.vm.hostname = vm_name

        s.vm.synced_folder ".", "/src/"

        s.vm.provider "virtualbox" do |vbox|
            vbox.name = vm_name
            vbox.customize ["modifyvm", :id, "--memory", "256"]
            vbox.customize ["modifyvm", :id, "--cpus", "1"]
            vbox.customize ["modifyvm", :id, "--ioapic", "on"]
        end
    end
end
