# Mount-point change detection

## Introduction

This document presents the design of the mount-point change detection system in NDM.
The goal of mount-point change detection is to detect the changes in the mount-points
and the filesystem on the existing blockdevices discovered by NDM and trigger appropriate
action to update the blockdevice CRs.

## Design

The mount-points and filesystem info of a block device are found by the probe - *mountprobe*.
This probe uses the mounts file, which the Linux kernel provides in the *procfs* pseudo
filesystem. Reading the mounts file provides the status of all the mounted
filesystems (sample output below).

``` text
rootfs / rootfs rw 0 0
/dev/root / ext3 rw 0 0
/proc /proc proc rw 0 0 usbdevfs
/proc/bus/usb usbdevfs rw 0 0
/dev/sda1 /boot ext3 rw 0 0 none
/dev/pts devpts rw 0 0
/dev/sda4 /home ext3 rw 0 0 none
/dev/shm tmpfs rw 0 0 none
/proc/sys/fs/binfmt_misc binfmt_misc rw 0 0
```

Whenever a block device is (un)mounted or the fs type changes, the changes are reflected in the mounts file. This proposal introduces a change to the existing _**mount-probe**_ in NDM. Similar to how _udev-probe_ listens to udev events by starting a loop in its `Start()`, a loop is started by _mount-probe_ that watches for changes in the mounts file and triggers updation when a change is detected. The *epoll* API is used to watch the mounts file for changes. The *epoll* API is provided by the Linux kernel for the userspace programs to monitor file descriptors and get notifications about I/O events happening on them. Whenever the mounts file changes, the events `EPOLLPRI` and `EPOLLERR` are emitted. This behaviour has been documented [here](https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/commit/?id=31b07093c44a7a442394d44423e21d783f5523b8) (additional links - [\[1\]](https://lkml.org/lkml/2006/2/22/169), [\[2\]](http://lkml.iu.edu/hypermail/linux/kernel/1012.1/02246.html)).

A new package `libmount` is introduced for parsing the mounts file. The package `libmount` is a pure go implementation of the C library with same name (see [util-linux/libmount](https://github.com/karelzak/util-linux/tree/master/libmount)). This package also provides utility to compare two mount tables and get a diff data structure which can be used to tell the changes between the two tables. Initially on start-up *mount-probe* parses the mounts file and stores the mount table in memory. On receving an event from epoll, the mounts file is parsed again to get the new mount table. This new mount table is then compared with the older mount table stored to generate a diff, which is used to get the list of devices that changed mount-points or filesystems. An `EventMessage` is generated containing the list of changed devices and pushed to `udevevent.UdevEventMessageChannel`. The message contains information about what blockdevices to check (the list of changed device). The message also
has additional information regarding what probes are to be run and specifies that only _mount-probe_ needs to be run for the event. This is done since _mount-probe_ alone can fetch the new mounts and fs data for the blockdevices. Running the probes selectively helps us optimize the updation process.
The message is then received by the loop in `udevProbe.listen()` and sent further down to the `ProbeEvent` change handler.

For every blockdevice listed in the `EventMessage`, the change handler first fetches the latest copy of the blockdevice from the controller blockdevice cache (`controller.BDHierarchyCache`) and then runs the it though the requested probes which are also provided in the message. Once the blockdevice is run through all the probes, the cache is updated and an update request is send to the kuebrnetes api server to upate the corresponding blocdevice CR.

&nbsp;

``` text
+----------------------------------+
|                                  |
|                                  |
|                                  |
|                                  |
|             Epoll API            |
|                                  |
|                                  |
|                                  |
+-----------------+----------------+
                  |
                  |
                  |
                  |
                  |
                  |
                  |           Event
                  |   (EPOLLPRI & EPOLLERR)
                  |
                  |
                  |
                  |                                                                                                                 Updated       +------------------------+
                  |                                                                                                               Blockdevice     |                        |
                  |                                                                                                           +------------------->       Controller       |
                  |                                                                                                           |                   |    Blockdevice Cache   |
                  |                                                                                                           |                   |                        |
+-----------------v----------------+                    -----------------------------------------                 +-----------+-----------+       +------------------------+
|                                  |                                                                              |                       |
|                                  |                                                                              |                       |       +------------------------+
|                                  |                                                                              |                       +------->                        |
|            mount probe           |    EventMessage                                                EventMessage  |      Probe Event      |       |         Probes         |
|                                  +------------------->    udevevent.UdevEventMessageChannel   ------------------>                       |       |  - update Blockdevice  |
|            listen loop           |                                                                              |     Change Handler    <-------+                        |
|                                  |                                                                              |                       |       +------------------------+
|                                  |                                                                              |                       |
+-------------+------^-------------+                    -----------------------------------------                 +-----------+-----------+       +------------------------+
              |      |                                                                                                        |                   |                        |
              |      |                                                                                                        |                   |       Kubernetes       |
              |      |                                                                                                        +------------------->          etcd          |
              |      |                                                                                                               Updated      |                        |
              |      |                                                                                                            Blockdevice CR  +------------------------+
              |      |
              |      |
              |      |
              |      |
 +------------v------+--------------+
 |                                  |
 |                                  |
 |            libmount              |
 |                                  |
 |        - parse mounts file       |
 |        - generate diff           |
 |                                  |
 +----------------------------------+
```
