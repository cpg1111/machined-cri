# -*- mode: ruby -*-
# vi: set ft=ruby :

CLOUD_CONFIG_PATH = File.join(File.dirname(__FILE__), "cloud-config.yml")

$update_channel='stable'
$image_version='current'

def vm_gui
  $vb_gui.nil? ? $vm_gui : $vb_gui
end

def vm_memory
  $vb_memory.nil? ? $vm_memory : $vb_memory
end

def vm_cpus
  $vb_cpus.nil? ? $vm_cpus : $vb_cpus
end

Vagrant.configure("2") do |config|
  config.vm.box = "coreos-#{$update_channel}"
  config.vm.network "private_network", ip: "172.22.22.22"
  config.vm.box_url = "https://storage.googleapis.com/#{$update_channel}.release.core-os.net/amd64-usr/#{$image_version}/coreos_production_vagrant.json"
  ["vmware_fusion", "vmware_workstation"].each do |vmware|
    config.vm.provider vmware do |v, override|
      v.check_guest_additions = false
      override.vm.box_url = "http:///storage.googleapis.com/#{$update_channel}.release.core-os.net/amd64-usr/#{$image_version}/coreos_production_vagrant_vmware_fusion.json"
    end
  end
  config.vm.provider :virtualbox do |v|
    v.check_guest_additions = false
    v.functional_vboxsf = false
  end
  config.vm.provision :file, :source => "#{CLOUD_CONFIG_PATH}", :destination => "/tmp/vagrantfile-user-data"
  config.vm.provision :shell, :inline => "mv /tmp/vagrantfile-user-data /var/lib/coreos-vagrant/", :privileged => true
  config.vm.synced_folder "./", "/opt/src/", :mount_options => ['nolock,vers=3,udp']
end
