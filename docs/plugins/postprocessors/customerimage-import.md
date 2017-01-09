# Customer image import post-processor

The customer image import post-processor converts a local VMWare (`.vmx`) virtual machine into OVF (`.ovf`) format, uploads it to CloudControl, and then imports it as a customer image.

## Settings

* `mcp_region` (Required) is the CloudControl region code (e.g. AU, NA, EU, etc).
* `mcp_user` (Required) is the CloudControl user name.  
Can also be specified via the `MCP_USER` environment variable.
* `mcp_password` (Required) is the CloudControl password.  
Can also be specified via the `MCP_PASSWORD` environment variable.
* `datacenter` (Required) is the Id of the datacenter where the image will be imported (must be MCP 2.0).
* `target_image` (Required) is the name of the customer image to create.
* `ovf_package_prefix` (Optional) is the prefix used to name the OVF package files.  
If not specified, `target_image` is used.

## Sample configurations

### Create a new VMWare virtual machine locally, and import it into Cloud Control as a customer image

`build.json`:

```json
{
	"builders": [
		{
			"type": "vmware-iso",
			"iso_url": "./iso/ubuntu-16.04.1-server-amd64.iso",
			"iso_checksum": "d2d939ca0e65816790375f6826e4032f",
			"iso_checksum_type": "md5",
			"output_directory": "./vmx",
			"ssh_port": 22,
			"ssh_username": "packer",
			"ssh_password": "packer",
			"ssh_wait_timeout": "10000s",
			"boot_command": [
				"<enter><wait><f6><esc><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
				"<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
				"<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
				"<bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs><bs>",
				"/install/vmlinuz<wait>",
				" auto<wait>",
				" console-setup/ask_detect=false<wait>",
				" console-setup/layoutcode=us<wait>",
				" console-setup/modelcode=pc105<wait>",
				" debconf/frontend=noninteractive<wait>",
				" debian-installer=en_US<wait>",
				" fb=false<wait>",
				" initrd=/install/initrd.gz<wait>",
				" kbd-chooser/method=us<wait>",
				" keyboard-configuration/layout=USA<wait>",
				" keyboard-configuration/variant=USA<wait>",
				" locale=en_US<wait>",
				" netcfg/get_domain=vm<wait>",
				" netcfg/get_hostname=packer<wait>",
				" grub-installer/bootdev=/dev/sda<wait>",
				" noapic<wait>",
				" preseed/url=http://{{ .HTTPIP }}:{{ .HTTPPort }}/preseed.cfg",
				" -- <wait>",
				"<enter><wait>"
			],
			"boot_wait": "10s",
			"disk_size": 81920,
			"vm_name": "ubuntu-from-packer",
			"vmx_data": {
				"cpuid.coresPerSocket": "1",
				"memsize": "1024",
				"numvcpus": "1"
			},
			"guest_os_type": "ubuntu-64",
			"http_directory": ".",
			"shutdown_command": "echo 'packer'|sudo -S shutdown -P now",
			"headless":"true"
		}
	],
	"provisioners": [
		{
			"type": "shell",
			"inline": [
				"echo 'packer' | sudo -S true",
				"sudo apt-get update",
				"sudo apt-get install -y open-vm-tools"
			]
		}
	],
	"post-processors": [
		[
			{
				"type": "ddcloud-customerimage-import",
				"mcp_region": "AU",
				"mcp_user": "my_mcp_user",
				"mcp_password": "my_mcp_password",
				"datacenter": "AU9",
				"target_image": "ubuntu-from-packer",
				"ovf_package_prefix": "ubuntu-from-packer"
			}
		]
	]
}
```

`preseed.cfg`:

```
choose-mirror-bin mirror/http/proxy string
d-i base-installer/kernel/override-image string linux-server
d-i clock-setup/utc boolean true
d-i clock-setup/utc-auto boolean true
d-i finish-install/reboot_in_progress note
d-i grub-installer/only_debian boolean true
d-i grub-installer/with_other_os boolean true
d-i partman-auto-lvm/guided_size string max
d-i partman-auto/choose_recipe select atomic
d-i partman-auto/method string lvm
d-i partman-lvm/confirm boolean true
d-i partman-lvm/confirm boolean true
d-i partman-lvm/confirm_nooverwrite boolean true
d-i partman-lvm/device_remove_lvm boolean true
d-i partman/choose_partition select finish
d-i partman/confirm boolean true
d-i partman/confirm_nooverwrite boolean true
d-i partman/confirm_write_new_label boolean true
d-i pkgsel/include string openssh-server cryptsetup build-essential libssl-dev libreadline-dev zlib1g-dev linux-source dkms nfs-common
d-i pkgsel/install-language-support boolean false
d-i pkgsel/update-policy select none
d-i pkgsel/upgrade select full-upgrade
d-i time/zone string UTC
tasksel tasksel/first multiselect standard, ubuntu-server

d-i console-setup/ask_detect boolean false
d-i keyboard-configuration/layoutcode string us
d-i keyboard-configuration/modelcode string pc105
d-i debian-installer/locale string en_US

# Create packer user account.
d-i passwd/user-fullname string packer
d-i passwd/username string packer
d-i passwd/user-password password packer
d-i passwd/user-password-again password packer
d-i user-setup/allow-password-weak boolean true
d-i user-setup/encrypt-home boolean false
d-i passwd/user-default-groups packer sudo
d-i passwd/user-uid string 900
```
