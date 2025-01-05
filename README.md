# virest-storage-volume
## Go Package for Managing Virtualization Storage Volume
Storage pools are divided into storage volumes. Storage volumes are abstractions of physical partitions, LVM logical volumes, file-based disk images, and other storage types handled by libvirt. Storage volumes are presented to guest virtual machines as local storage devices regardless of the underlying hardware.

ViRest provides storage management on the physical host through storage pools and volumes. This Go package provides the REST API interface by utilizing Libvirt.

## Needed Package to Running Executable using Qemu/KVM Hypervisor
- qemu-kvm
- libvirt-daemon-system
- bridge-utils

## Needed Package for Development and Compiling
- libvirt-dev
- gcc

## Add User to Libvirt Group & KVM Group
```shell
sudo adduser '<username>' libvirt
```
```shell
sudo adduser '<username>' kvm
```

## Known error:
- Libvirt Go Binding methods undefined, please enable "cgo" with command:
    ```shell
    export CGO_ENABLED=1
    ```
- [Can't access storage, file permission denied](https://ostechnix.com/solved-cannot-access-storage-file-permission-denied-error-in-kvm-libvirt/)

#### References
- https://libvirt.org/storage.html