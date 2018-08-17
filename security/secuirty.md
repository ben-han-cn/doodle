# Security
## cryptography
The fundamental building block of most modern cryptography is a one-way function. 
A one-way function is a function that is easy to compute, but its inverse is hard 
to compute. (i.e. given f(x) = y, it is easy to calculate y given x, but to very hard
to calculate x given y)
There are two main ways that this is done:
1 Factorization
2 Elliptic Curve Logarithms


* Basic principle
  defense in depth
  need-to-know basis
  least privilege
  log everything 
  * mmap avoid the data copy from page cache to user application 
  ![""](mmap.png)

fadvise could be used to help the kernel to adjust and optimise the use of page cache
FADV_RANDOM: disable read ahead, greeing the memory 
FADV_SEQUENTIAL: encourage OS to read ahead
FADV_WILLNEED: notifies OS that the page will be needed in the near feature, the opposite
is FADV_DONTNEED

* DIO(direct IO)
  * ![""](direct_io.png)
* AIO/DIO(async Direct IO)
  ![""](io_tradeoffs.png)

 ```
  计世资讯（CCW Research）预计，按照销售额计算，2016年中国服务器虚拟化市场规模将达到21.7亿元。 
  到2020年，市场规模将达到44.1亿元。
 ```

## Virtualization Solution & CMP(Cloud Management Platform)

### VMWare

### OpenStack

## Why Virtualization & Private Cloud?
* Converged Infrastructure 
    * Virtual Machine/Compute
        * Cloud compution(On-demand self service, Resource pooling, Rapid elasticity, Measured service)
        * Multi-tenant isolation and resource limitation

    * Virtual Network
        * Central control over/knowledge of logical network topology
        * Decouple control and data plane
        * Network isolation
        * Virtualize network device(switch, router, load balance, firewall)
        * Programmatic integration with CMP

    * Virtual Storage
* Benefit
    * Wire once
    * Agility and flexibility 
    * Visibility 

* Virtual device
    * Distributed version(vDS), software is much easier to create abstraction

## Common Networking Challenges in Private Cloud Environment

* Manually Network configuration for VM is time-consuming and error-prone
* Solutions lack visibility and auditing capability
* Lack of centrialize IP address and DNS managelent

## Available Solution

### Infoblox Cloud Network Automation

* Support mainstream cloud management platform(CMP, 2014)
* Architecutre 
    * Adaptor
    * Cloud Platform Applicance
    * Grid Master

    Infoblox Cloud Platform Appliances are fully virtualized Infoblox Grid members
    that run on ESXi, Hyper-V, KVM or XenServer hypervisors. They deliver the full suite
    of Infoblox DNS, DHCP, and IPAM to cloud environments such as VMware,
    OpenStack, and Microsoft. These appliances, optimized for cloud deployments
    in the data center, also deliver a range of cloud-enabling functions including:
      * Automated IP address provisioning and reclamation when VMs are decommissioned
      * Automated DNS naming and reclamation when VMs are decommissioned
      * Automated DHCP lease assignment with fixed address support—especially
        important in OpenStack environments

### Our plan
  * Network view support/ virtual zdns slave
  * VMWare plugin
  * Openstack ipam agent
