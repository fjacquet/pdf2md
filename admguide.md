 

POWERFLEX 4.X ADMINISTRATION 

 

# PARTICIPANT** **COURSE GUIDE

 

# COURSE GUIDE

![Image](images/admguide_Image12.png)

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 2  |

 

# Table of Contents 

| **PowerFlex 4.x Administration** |  **9**  |
| PowerFlex 4.x Administration |  9  |
| Course Content Key Identifiers |  9  |
| **PowerFlex Introduction** |  **10**  |
| What Is PowerFlex? |  10  |
| How PowerFlex Scales |  10  |
| **Model Breakdown and Positioning** |  **12**  |
| Model Breakdown and Positioning |  12  |
| PowerFlex Offerings |  12  |

Knowledge Check - PowerFlex Offerings Case Study 1 15 Knowledge Check - PowerFlex Offerings Case Study 2 15 Knowledge Check - PowerFlex Model Case Study 3 16 

Knowledge Check - PowerFlex Model Case Study 4 16 

| PowerFlex Hardware |  17  |
| Knowledge Check - Hardware Case Study 1 |  19  |
| Knowledge Check - Hardware Case Study 2 |  20  |
| Knowledge Check - Hardware Case Study 3 |  20  |
| Knowledge Check - Hardware Case Study 4 |  21  |
| **Software Architecture** |  **21**  |
| Software Architecture |  21  |
| Component Communication |  27  |
| Deployment Options |  32  |
| Knowledge Check - Deployment Case Study |  34  |
| **Management Interfaces** |  **35**  |
| PowerFlex Manager |  35  |
| Explore PowerFlex Manager |  38  |
| SCLI Management Interface |  38  |
| Query Command Examples |  39  |

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 3  |

 

| PowerAPI |  40  |
| PowerAPI Example |  42  |
| **Deployment Process** |  **43**  |
| Outline of the Implementation Process |  43  |
| Implementation Process |  43  |
| Implementation Process |  44  |
| Implementation Process |  45  |
| Implementation Process |  46  |
| Implementation Process |  47  |
| **Post-Installation Tasks** |  **48**  |

Getting Started Wizard: Download a Compliance File 48 Getting Started Wizard: Define the Networks 49 Getting Started Wizard: Discover Resources 50 

Getting Started Wizard: Manage Deployed Resources 51 Getting Started Wizard: Deploy Resources 51 Deploy the Production Cluster - CSV Method 52 

| CSV File Formatting |  53  |
| Apply a PowerFlex License |  54  |
| **Upgrade a PowerFlex Cluster** |  **56**  |
| Upgrade a PowerFlex 4.x System |  56  |

# PowerFlex** **Templates and Resource Groups

**58** 

| Templates |  58  |
| Resource Groups |  62  |

Working With Templates and Resource Groups: Demonstration 65 

| **Configure Protection Domains** |  **66**  |
| About Protection Domains |  66  |

Configure a Protection Domain with PowerFlex Manager 67 

| Configure a Protection Domain |  70  |
| **Configure Fault Sets** |  **71**  |

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 4  |

 

| Fault Sets |  71  |
| Create and Delete Fault Sets |  75  |
| **Configure Storage Pools** |  **76**  |
| Acceleration Pools |  77  |
| Configure Storage Pool Settings |  78  |
| Enable the Background Device Scanner |  80  |

Configure I/O Priorities (Rebuild/Rebalance) 81 

| Configure Storage Pools |  85  |

# Create, Modify,** **and Delete Volumes    **86** 

| Adding Volumes Using PowerFlex Manager |  86  |

Increase Volume Size Using PowerFlex Manager 86 Map Volumes to a Host Using PowerFlex Manager 86 Set Volume Bandwidth and IOPS Limits Using PowerFlex Manager 87 

| vTree Migration Using PowerFlex Manager |  87  |

Unmap Volumes from a Host Using PowerFlex Manager 88 Removing Volumes Using PowerFlex Manager 88 

| **Create a NAS Server (NFS and SMB)** |  **89**  |
| Overview of a File System Storage |  89  |

PowerFlex 4.5.x File Service System Overview 93 NAS Capabilities of a PowerFlex Cluster (System) 95 Create a NAS Server in PowerFlex Manager 101 

| **Export/Share Filesystems** |  **104**  |

Configure the Settings of an Existing NAS Server 104 

| Create a File System for NFS Export |  107  |
| Create a File System for SMB Share |  107  |
| PowerFlex File System Global Name Space |  108  |

Create a Global Namespace in PowerFlex 4.5.x 109 

| **Protect NAS Servers** |  **115**  |

File System and NAS Server Protection Options 115 

| Protect NAS Servers |  117  |

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 5  |

 

| **Configure NAS Quotas** |  **118**  |

File System Quotas in a PowerFlex Cluster 118 

| Quota Limits for a File System |  119  |

Configure NAS Quotas for File Systems in PowerFlex Manager 123 

| **Enter and Exit Maintenance Modes** |  **124**  |
| Types of Maintenance Modes |  125  |
| Abort Protected Maintenance Mode |  126  |
| Put an SDS into Maintenance Mode |  127  |
| Put a Node into Service Mode |  127  |
| **Perform a Node Expansion** |  **128**  |
| Expansion Process Overview |  128  |
| PowerFlex 4.x rack Expansion Types |  128  |
| Validate System Health |  132  |
| Verify Cluster Configuration Details |  135  |
| Cabling Pre-Requisites |  137  |
| PowerFlex rack Node Port Locations |  138  |
| Switch Cabling |  142  |
| Verify Cabling |  144  |

Perform Expansion Using PowerFlex Manager 145 

| Performing Expansion Using CSV File |  148  |

# Manage Users and Passwords in a PowerFlex Cluster

**149** Change PowerFlex Cluster System Passwords 149 

| Manage User Roles |  150  |

Change Passwords for Resources Discovered in PowerFlex Manager 151 Manage the Secadmin User in the CloudLink Center and PowerFlex 151 

| **Protect Volumes using SnapShots** |  **153**  |
| PowerFlex Snapshot |  153  |

Difference Between Regular Snapshot and Secure Snapshot 154 

| Create, Modify, and Delete Snapshots |  155  |
| Create a Volume Snapshot |  158  |

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 6  |

 

| Managing Volume Snapshots |  159  |
| Snapshot Policies |  159  |
| Managing Snapshot Policies |  160  |
| Working With Snapshot Policies |  160  |
| **Replicate Volumes Between Clusters** |  **161**  |
| Native Asynchronous Replication |  161  |
| Configure Replication |  162  |
| Swap Root Certificates - Exercise |  163  |
| Configure Replication |  164  |
| Journal Configuration - Knowledge Check |  165  |
| Configure Replication |  165  |
| Configure Replication |  167  |
| Replication I/O Flow |  169  |
| **Configure Storage Data Servers** |  **170**  |

Prerequisites for Adding a Storage Data Server 170 

| Add a Storage Data Server (SDS) |  171  |
| Add an SDS - Exercise |  172  |
| Add an SDS - Advanced Options |  172  |
| Add a Device - Advanced Options |  173  |
| **Reconfiguring MDMs** |  **175**  |
| Reconfiguring MDM Roles |  175  |

Reconfigure MDM Roles in a PowerFlex Cluster 175 

| A PowerFlex Cluster Refresh |  175  |

Enable and Disable SDC Automatic Update with MDM IP Addresses 176 

| How the Automated Solution Works |  177  |

Automated SDC Update with MDM IP Addresses Scenarios 178 

| **NVMe/TCP** |  **181**  |
| Configure Storage Target (SDT) Service |  181  |

Add an NVMe Target Using PowerFlex Manager 182 

| Register NVMe Host Initiator |  182  |

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 7  |

 Add an NVMe Host Using PowerFlex Manager 183 

# Add CloudLink to a PowerFlex Cluster    **184** 

Role of Cloudlink in a PowerFlex Cluster 184 Add a CloudLink License to the PowerFlex Cluster 185 

| Discover and Deploy CloudLink |  186  |

Deploy a CloudLink Cluster Using PowerFlex Manager 187 

| **Configure PowerFlex Alerting** |  **188**  |
| Monitoring Menu |  188  |
| Enabling SupportAssist |  190  |
| Configure SupportAssist |  191  |
| Configure SupportAssist |  194  |
| Connecting PowerFlex to CloudIQ |  195  |

Configure Events and Alert Notifications to External Systems Overview 196 Configure External Source and Destination 197 

| **PowerFlex in the Cloud (APEX)** |  **198**  |
| What Is Dell APEX? |  198  |

The Dell APEX as-a-Service Portfolio Offerings 199 

| APEX Subscription Model |  200  |
| PowerFlex Cloud |  203  |
| PowerFlex Public Cloud Environments |  205  |
| **Appendix** |  **207**  |
| **Glossary** |  **233**  |
| **Knowledge Check** **Answer** **Key** |  **242**  |

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 8  |

PowerFlex 4.x Administration  

# PowerFlex 4.x Administration 

## PowerFlex 4.x Administration 

Upon completion of this course, you should be able to: 

- $•$  Identify the PowerFlex hardware and software components and 
  - how they interact with each other. 

$•$  Manage storage and network resources. $•$  Configure and manage data protection and security. $•$  Create NAS servers, configure NAS quotas for filesystems and configure exports / shares for NFS and SMB filesystems on a 

PowerFlex system. $•$  Backup and restore various components of a PowerFlex system. $•$  Perform basic log collection and troubleshooting. 

## Course Content Key Identifiers 

This course uses the following guidelines:   Text appearing in the UI element like wizard titles, window titles, and labels appear in **bold** font. 

  Command-line commands appear in a **command name** style format. For example: review the output of the **ipconfig** command. 

  Definitions appear as a Glossary Term.  

The following terms are used interchangeably between this course and other documentation: 

$•$  PowerFlex cluster and PowerFlex system. $•$  File controller and File Server Node. 

$•$  PFMC, Controller Nodes, and Management Nodes. $•$  CloudIQ and APEX AIOps Observability (AAIOO) 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 9  |

PowerFlex Introduction 

# PowerFlex Introduction 

## What Is PowerFlex? 

 To learn more, review the[ Dell PowerFlex 4.5.x Technical Overview ](https://www.dell.com/support/manuals/en-us/scaleio/flex-software-to-45x/introduction?) documentation located on dell.com/support. 

 PowerFlex is a storage virtualization solution that uses a node-based architecture for scaling resources including storage and compute. 

It can have separate, dedicated, compute and storage nodes; can combine them into hyperconverged nodes, or have a combination of both. 

PowerFlex is considered Software Defined Storage. Software Defined Storage means that physical storage from aggregate nodes is pooled together to create volumes. Volumes are mapped to targets, allowing an application to have access to the storage. 

PowerFlex uses multiple IP-based protocols for communication between hosts and storage. These include file-level TCP/IP protocols, NVMe/TCP, and a proprietary TCP/IP connection unique to PowerFlex. 

*Video:* *The web version* *of this content contains a* *video.* 

## How PowerFlex Scales 

PowerFlex can start scaling with as few as four nodes, and the platform can scale to thousands of nodes linearly without any disruption. When more nodes are added to the cluster, the resource pool grows as shown. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 10  |

![Image](images/admguide_Image452.jpeg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path8.svg)

PowerFlex Introduction PowerFlex can consolidate heterogeneous workloads while maintaining performance and throughput levels by aggregating resources from 

underlying nodes in the cluster. 

 PowerFlex scaling is achieved by building a large and distributed pool that eliminates I/O bottlenecks present in any single node by itself. Each node 

contributes bandwidth, which is combined to deliver high and scalable throughput performance at low latency. 

# What happens when** **you** **remove** **and** **add** **nodes?

PowerFlex administrators can add and remove nodes to the system with no impact on application availability and performance with PowerFlex. If a node must be upgraded to newer-generation hardware, or more compute resources added, nodes can be replaced with a new node. 

Take a scenario where a node has reached its end of service life. The administrator can remove the node, and the system rapidly rebalances. 

The administrator can then deploy a new node, and the system balances again, ensuring all applications remain online. When a planned or unplanned outage occurs, or a node requires maintenance, it can be removed from the cluster. The administrator can add it back when the 

maintenance is complete, and the node is online. 

Through this node-based architecture, applications and workloads are insulated from infrastructure issues. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 11  |

![Image](images/admguide_Image455.jpegimages/admguide_Path1.svg)

Model Breakdown and Positioning 

# Model Breakdown and Positioning 

## Model Breakdown and Positioning 

## PowerFlex Offerings 

PowerFlex offerings consist of PowerFlex rack, PowerFlex appliance, PowerFlex custom node, and a PowerFlex software only. Each offering comes with PowerFlex software defined storage, and is designed to meet differing customer requirements and expectations. 

To learn more about each offering, scroll through each tab below. 

PowerFlex rack 

 PowerFlex rack is a fully engineered rack-scale system with integrated networking. It has a larger starting point, and scales to hundreds of nodes 

and racks. It comes with licensing for a PowerFlex system and a unified management platform. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 12  |

![Image](images/admguide_Image458.png)

Model Breakdown and Positioning PowerFlex appliance 

 PowerFlex appliance has a smaller starting point than rack, but scales to hundreds of nodes. It allows customers to use a broad set of supported 

networking options and can be added to existing networking infrastructures. However, it is limited to supported networking configurations. It comes with licensing for a PowerFlex system and a unified management platform. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 13  |

![Image](images/admguide_Image461.png)

Model Breakdown and Positioning PowerFlex custom node 

 PowerFlex Custom Nodes are servers that are configured to support a PowerFlex system. The servers are tested and certified to deliver 

predictable performance for a wide variety of workloads. This offering is ideal for customers who prefer to build their own environments and have their own management services. It also allows broader networking options than the appliance. It comes with licensing for a PowerFlex system and a 

unified management platform. 

PowerFlex software only 

 PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 14  |

![Image](images/admguide_Image464.png)

![Image](images/admguide_Image465.jpeg)

Model Breakdown and Positioning PowerFlex software only is intended for large enterprise implementations where robust onsite hardware support is readily available. Customers who 

want to pursue a software-only implementation must commit to licensing in the hundreds to thousands of nodes over a predefined time period. 

# Knowledge Check - PowerFlex Offerings Case Study 1 

A customer is expanding their data center to include large scale virtual storage and compute capabilities. They have recently upgraded their networking infrastructure with the latest Dell switches. The customer wants to increase their compute, storage, and networking capacity with a 

simple to deploy, scalable multi-node solution that includes a unified management platform. They also want to be able to maintain their existing network infrastructure. 

# 1.  Which PowerFlex offering best suits these requirements? 

a.  PowerFlex rack b.  PowerFlex appliance 

c.  PowerFlex custom node d.  PowerFlex software only 

Answer Key 

# Knowledge Check - PowerFlex Offerings Case Study 2 

A customer has a large data center that has been recently updated with hundreds of new servers. They have learned about PowerFlex and want to leverage the scalability of software-defined storage on their existing servers without investing in new hardware. 

# 2.  Which PowerFlex offering best suits these requirements? 

a.  PowerFlex rack b.  PowerFlex appliance 

c.  PowerFlex custom node d.  PowerFlex software only 

Answer Key PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 15  |

Model Breakdown and Positioning 

# Knowledge Check - PowerFlex Model Case Study 3 

A customer wants to build their own virtualized storage clusters. They are not interested in upgrading their network infrastructure or managing it through a shared interface. They do want hardware that Dell supports and is configured for PowerFlex. 

# 3.  Which PowerFlex offering best suits these requirements? 

a.  PowerFlex rack b.  PowerFlex appliance 

c.  PowerFlex custom node d.  PowerFlex software only 

Answer Key 

# Knowledge Check - PowerFlex Model Case Study 4 

A customer wants to switch their current data center over to virtual storage and compute capabilities. The current bare-metal data center infrastructure consists of stand-alone servers with internal storage. The data center has led to underperformance, limited scalability, and 

management complexity. The customer is looking for a software-defined storage and compute solution that has Dell provided fully integrated network devices and a unified management platform. The customer also wants the solution to give them room to scale as the company grows. 

# 4.  Which PowerFlex offering best suits these requirements? 

a.  PowerFlex rack b.  PowerFlex appliance 

c.  PowerFlex custom node d.  PowerFlex software only 

Answer Key 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 16  |

Model Breakdown and Positioning 

# PowerFlex Hardware 

PowerFlex custom node, appliance, and rack configurations are all built around the PowerEdge Server platform. The PowerEdge platform provides the latest processing, memory, and storage resources that form the foundation of the PowerFlex system. Several models are offered, each 

providing distinct capabilities. 

To learn about the PowerEdge models offered, scroll through each page below. 

PowerEdge 650/R660 

 The PowerEdge R650 platform is a single rack-unit. As such, it can provide great CPU core density for usage in a compute-only, or 

hyperconverged node use-case. It provides ten drive slots for an overall maximum capacity of 76.8 TB. 

The 16G PowerEdge version is the R660 server. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 17  |

![Image](images/admguide_Image476.jpeg)

Model Breakdown and Positioning PowerEdge R6525/R6625/7525/7625 

 

The 15G PowerEdge R6525 (single rack-unit) and R7525 (two rack-unit) are similar in specification to other node types but feature AMD processors. This difference allows customers choice in computing platform, while having the same functionality other node types can 

provide. 

The 16G version PowerEdge servers are the R6625 (single rack-unit) and R7625 (two rack-unit). 

Of note, PowerEdge servers with AMD processors are only supported in a PowerFlex Cluster as compute-only nodes. 

PowerEdge 750/760 

 PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 18  |

![Image](images/admguide_Image478.jpeg)

![Image](images/admguide_Image479.png)

![Image](images/admguide_Image480.jpeg)

Model Breakdown and Positioning The 15G PowerEdge R750 platform is a two rack-unit server. This form factor provides greater storage density, making it an ideal choice for 

storage-only nodes. Its maximum capacity is 128 TB with all 24 drive slots populated. The 2RU configuration also provides greater expansion in the form of full-height PCI cards such as double-slot GPUs. These features make the R750 a good candidate for hyperconverged nodes as well. 

The 16G PowerEdge version is the two-rack-unit R760 server. 

PowerEdge R860 

 The 16G PowerEdge R860 rounds out the hardware platform offerings. The R860 with a 2U air-cooled form factor and four Intel Xeon® Scalable 

Processors enables businesses to accelerate a wide variety of core and business-critical applications including databases, analytics, and other data-driven applications. 

**Go to:**  For additional information about the PowerEdge rack servers, refer to the[ Dell PowerEdge servers by ](https://www.dell.com/support/kbdoc/en-us/000137343/how-to-identify-which-generation-your-dell-poweredge-server-belongs-to) [generation ](https://www.dell.com/support/kbdoc/en-us/000137343/how-to-identify-which-generation-your-dell-poweredge-server-belongs-to)documentation page located in dell.com/support.  

 

# Knowledge Check - Hardware Case Study 1 

Platform Needs Analysis A customer is planning a two-layer PowerFlex implementation, with dedicated storage and compute nodes. Their applications run exclusively 

on AMD processors. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 19  |

![Image](images/admguide_Image484.jpeg)

![Image](images/admguide_Image485.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Model Breakdown and Positioning 

# 5.  Which PowerEdge platform best suits these requirements? 

a.  R650 b.  R6525 

c.  R750 d.  R860 

Answer Key 

# Knowledge Check - Hardware Case Study 2 

Platform Needs Analysis A customer is planning a hyperconverged PowerFlex solution, where all nodes will perform both compute and storage roles. They want a single 

hardware platform for all nodes that can provide the greatest amount of computing power. The customer also wants the greatest storage-density and expansion capabilities. Cost is not a factor. 

# 6.  Which PowerEdge platform best suits these requirements? 

a.  R650 b.  R6525 

c.  R750 d.  R860 

Answer Key 

# Knowledge Check - Hardware Case Study 3 

Platform Needs Analysis A customer requires a compute-dense solution that minimizes the amount of rack space for their PowerFlex cluster. The nodes should be capable of 

providing hyperconverged functionality. Cost is not a factor. 

# 7.  Which PowerEdge platform best suits these requirements? 

a.  R650 PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 20  |

Software Architecture b.  R6525 c.  R750 

d.  R860 Answer Key 

## Knowledge Check - Hardware Case Study 4 

Platform Needs Analysis A customer is planning a two-layer PowerFlex implementation, with dedicated storage and compute nodes. They are looking for the most cost-

effective solution for their storage nodes, with room for expansion. Cost to storage density is their greatest driving factor in choosing a hardware platform for the storage nodes. 

# 8.  Which PowerEdge platform best suits these requirements? 

a.  R650 b.  R6525 

c.  R750 d.  R860 

Answer Key  

# Software Architecture 

## Software Architecture 

Other PowerFlex Cluster components: 

- $•$  Lightweight Installer Agent (LIA) 
- $•$  Nodes 
- $•$  Storage Media 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 21  |

![Image](images/admguide_Path3.svg)

Software Architecture The PowerFlex cluster is a combination of components that enable the consumption, virtualization, and presentation of local storage devices as a 

shared, software-defined storage solution. 

To learn about these components, scroll through each page below. 

Meta Data Manager (MDM) 

 The MDM is the authority that controls and tracks data storage ownership, mapping, and protection. As volumes are created, the MDM provides the 

information application servers need to connect to the cluster’s virtualized storage. The MDM also ensures the even distribution of writes, to ensure balanced resource utilization. The MDM ensures data integrity by running periodic background processes to verify protection. All this is done through 

communication with the other PowerFlex processes. Despite the MDM's responsibilities, user data never passes through MDM. 

Typically, MDM nodes are deployed in a clustered configuration to ensure availability. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 22  |

![Image](images/admguide_Image506.png)

Software Architecture Storage Data Server (SDS) 

 The SDS is a software daemon that enables a server in the cluster to contribute local storage devices to an aggregated storage pool. It owns 

the contributing devices, and together with the other SDSs, forms a protected mesh from which storage pools are created. 

An instance of the SDS runs on every server that contributes some or all its local storage space, managing the capacity of that single server. The SDS is also responsible for performing requested SDC back-end I/O operations, and MDM rebuild and rebalance operations. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 23  |

![Image](images/admguide_Image509.jpeg)

Software Architecture Storage Data Client (SDC) 

 The SDC is a block device driver that exposes shared block volumes from the SDS to the operating system. The SDC runs on the same server as 

the application. In practice, an application issues an I/O request, and the SDC fulfills the request, regardless of what SDS the requested blocks physically reside on. 

The SDC communicates with other nodes (beyond its own local server) over TCP/IP-based protocol. The only I/O in the stack that the SDC intercepts are the I/O that are directed at the volumes that belong to PowerFlex. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 24  |

![Image](images/admguide_Image512.png)

Software Architecture File Services Node (FSN) 

 The FSN is a software component that allows the PowerFlex cluster to make data available over file-based protocols such as NAS. The FSN 

supports protocols such as: $•$  Server Message Block (SMB): A shared resource on an SMB server, typically Windows. The PowerFlex FSN component supports the SMB3 

protocol. $•$  Network File System **(**NFS**):** Most common protocol for sharing files between UNIX systems over a network. NFS exports can be either 

NFSv3 or NFSv4 protocol. $•$  File Transfer Protocol (FTP): A standard network protocol used to transfer files over a TCP-based network. Both secure and non-secure 

FTP are supported. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 25  |

![Image](images/admguide_Image515.png)

Software Architecture Storage Data Target (SDT) 

 PowerFlex uses a SDT to expose NVMe over TCP targets. The SDT is 

deployed with the SDS on each storage server and provides access to the volumes inside that protection domain. 

 PowerFlex Cluster To learn about how a PowerFlex cluster operates, watch the video below. 

# Video:

*The web version* *of this content contains a* *video.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 26  |

![Image](images/admguide_Image518.png)

Software Architecture Storage Data Replication (SDR) 

 Extending the PowerFlex cluster beyond a single site is accomplished through the Storage Data Replication component. It is an optional 

component responsible for managing all the aspects of PowerFlex replication. It is installed on an SDS node that contains storage media contributing to storage pools that have the potential to be replicated. 

The SDR intercepts write IOs from the SDC for replication and pass it to the SDS for storage. If the data written is being replicated, the I/O is redirected to the SDR instead of the SDS. In simple terms, to the SDC, SDR appears as an SDS. To the SDS, it appears to be an SDC. 

# Component Communication 

The software components that makeup PowerFlex (the SDCs, SDSs, and MDMs) converse with each other in predictable ways. Understanding these communication patterns enable administrators to make wise network layout decisions when designing a PowerFlex deployment. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 27  |

![Image](images/admguide_Image521.png)

Software Architecture Storage Data Client to Storage Data Server 

 Traffic between the SDCs and the SDSs forms the bulk of front-end 

storage traffic. Front-end storage traffic includes all read and write traffic arriving at or originating from a client. This network has a high throughput requirement. 

Storage Data Server to Storage Data Server 

 Traffic between SDSs forms the bulk of back-end storage traffic. Back-end 

storage traffic includes writes that are mirrored between SDSs, rebalance 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 28  |

![Image](images/admguide_Image524.png)

![Image](images/admguide_Image526.png)

traffic, rebuild traffic, and volume mig throughput requirement. 

Software Architecture n traffic. This network has a high ratio

Meta Data Manager to Meta Data Manager 

 MDMs are used to coordinate operations inside the cluster. They direct PowerFlex to manage, rebalance, rebuild, and redirect traffic. They also 

coordinate Replication Consistency Groups, determine replication journal interval closures, and maintain metadata synchronization with PowerFlex replica-peer systems. MDMs are redundant and must continuously communicate with each other to establish quorum and maintain a shared 

understanding of data layout. 

MDMs do not carry or directly interfere with I/O traffic. The data exchanged among them is relatively lightweight, and MDMs do not require the same level of throughput required for SDS or SDC traffic. 

MDM to MDM traffic requires a stable, reliable, low latency network. MDM to MDM traffic is considered back-end storage traffic. PowerFlex supports the use of one or more networks dedicated to traffic between MDMs. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 29  |

![Image](images/admguide_Image529.png)

Software Architecture Meta Data Manager to Storage Data Client 

 The Primary MDM must communicate with SDCs in the event that data layout changes. This can occur because the SDSs that host an SDC’s 

volume(s) storage for the SDCs are added, removed, placed in maintenance mode, or go offline. It may also happen if a volume is placed into a Replication Consistency Group. 

Communication between the Primary MDM and the SDCs is lazy and asynchronous but still requires a reliable, low latency network. MDM to SDC traffic is considered front-end storage traffic. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 30  |

![Image](images/admguide_Image532.png)

Software Architecture Meta Data Manager to Storage Data Server 

 The Primary MDM must communicate with SDSs to monitor SDS and device health and to issue rebalance and rebuild directives. MDM to SDS 

traffic requires a reliable, low latency network. MDM to SDS traffic is considered back-end storage traffic. 

Other Traffic There are many other types of low-volume traffic in a PowerFlex cluster. The other traffic includes infrequent management, installation, and 

reporting. This also includes traffic to the PowerFlex Gateway (REST API Gateway, Installation Manager, and SNMP trap sender), PowerFlex Manager, traffic to and from the Light Installation Agent (LIA), and reporting or management traffic to the MDMs (such as syslog for reporting 

and LDAP for administrator authentication). It also includes CHAP authentication traffic among the MDMs the SDSs and SDCs. 

**Tip**: Front-end and back-end storage traffic are logical distinctions and do not require physically distinct networks.  

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 31  |

![Image](images/admguide_Image535.png)

![Image](images/admguide_Image537.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Software Architecture 

# Deployment Options 

PowerFlex supports deploying nodes in three different configurations: Storage-only, Compute-only, and Hyperconverged. As the name implies, a Storage-only node handles strictly virtualizing and presenting storage. Compute-only nodes do not have storage responsibilities, only that of 

running applications. A Hyperconverged node performs both storage and compute responsibilities. 

As with the nodes, the options for deploying the nodes vary as well. Based on what nodes make up the deployment, there can be three different deployment models. 

Two-Layer 

 In a two-layer deployment, the SDS is installed on a separate host from the SDC. The front-end (client) is separated from the back end (storage) 

data traffic. 

Two-layer deployments allow compute and storage resources to grow independently. PowerFlex compute-only nodes host end-user applications. PowerFlex storage-only nodes contribute storage to the system pool. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 32  |

![Image](images/admguide_Image540.png)

Software Architecture Hyperconverged 

 In the hyperconverged configuration, both the SDC and the SDS can be 

installed on the same host. Hyperconverged enables the applications and storage to reside on the same host. 

Hyperconverged deployments ma infrastructure requirements. ize hardware utilization and reduces xim

Mixed 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 33  |

![Image](images/admguide_Image543.png)

![Image](images/admguide_Image544.png)

Software Architecture Hybrid hyperconverged deployment consists of hyperconverged, compute-only, and storage-only nodes. Some nodes contribute both compute 

resources and storage resources (hyperconverged nodes), some contribute only compute resources (compute-only nodes), and some contribute only storage resources (storage-only nodes). 

# Knowledge Check - Deployment Case Study 

Platform Needs Analysis A customer is planning their PowerFlex implementation, leveraging VMware by Broadcom, and is primarily interested in simplified 

deployment. They also want to use as many resources on as many nodes as possible. 

# 1.  Which PowerEdge deployment method best suits these requirements? 

a.  Hyperconverged b.  Two-layer 

c.  Mixed Answer Key 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 34  |

Management Interfaces 

# Management Interfaces 

## PowerFlex Manager 

PowerFlex Manager automates deploying, configuring, managing, and upgrading a PowerFlex system. The PowerFlex Manager Dashboard provides an overview of the PowerFlex system. It allows administrators to quickly view the status of the systems hardware and software 

components. 

 **1: Dashboard Menu** 

The Dashboard displays the gene of the PowerFlex system. verall health, status, and utilization ral o

**2: Block Menu** Manage block storage in the PowerFlex cluster. 

      

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 35  |

![Image](images/admguide_Image549.png)

Management Interfaces **3: File Menu** Create and configure storage resources as network file storage, and 

protect those file system resources using snapshots. 

Using File Nodes, customers cannot create File Storage without purchasing File Nodes and then deploying the File service. 

**4: Protection Menu** Use local and remote protection features to protect block storage within 

the cluster. **5: Lifecycle Menu** 

Displays tasks that are associated with deploying resource groups and building templates. 

**6: Resources Menu** Displays a list of discovered hardware and software components in 

PowerFlex Manager. **7: Monitoring Menu** 

Displays tasks for monitoring events and alerts from resources that are installed or automatically discovered. 

**8: Settings Menu** The Settings Menu displays a list of settings that can be configured in 

PowerFlex Manager. **9: Running Management Jobs** 

Lists running jobs that are related to management activities. **10: Running Storage** **Jobs** 

Lists running jobs that are related to storage activities. **11: Alerts Menu** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 36  |

This section lists the current ale categorized by severity level: 

$•$  Critical $•$  Major 

$•$  Minor **12: User Menu** 

Management Interfaces thin the system, which are rts wi

Displays the user who is curren provides a logout button. gged into PowerFlex Manager and tly lo

**13: Help Menu** Displays the PowerFlex Manager Online user guide. Administrators can 

search the online help for guidance on PowerFlex Manager tasks. **14: Overall Performance and Latency** 

This section displays graphs that show the overall performance (IOPS and bandwidth) and latency for block and file configurations within the system. 

**15: Alerts** This section lists the current ale

categorized by severity level: 

- thin the system, which are rts wi

$•$  Critical $•$  Major 

$•$  Minor **16: Usable Capacity** 

This section shows the total capacity, along with details about the physical, system, and free capacity. 

**17: Data Savings** This section provides details about the overall savings and thin 

provisioning savings. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 37  |

**18: Resources** **and Inventory** This section shows the number of resou

$•$  VM managers $•$  Nodes 

$•$  Switches $•$  Protection domains 

$•$  Storage pools $•$  Volumes 

$•$  Hosts $•$  File systems 

$•$  NAS servers **19: Resources** **Groups** 

Management Interfaces 

in the current inventory: rces 

This section displays a graphical representation of the resource groups deployed based on status. The number next to each icon indicates the number of resource groups in a particular state. 

# Explore PowerFlex Manager 

To explore PowerFlex Manager, use the simulation below.  

*The web version* *of this content contains an interactive activity.* 

# SCLI Management Interface 

The command-line interface to PowerFlex provides a command set with which many repeatable day-to-day administrative tasks can be automated through scripting. The CLI runs directly on a PowerFlex node by invoking the **scli** command followed by a set of actions. Commands should be 

run against the Primary MDM in the PowerFlex cluster. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 38  |

Management Interfaces The **scli** command is structured in a command | parameter | value format, for example: 

# #> scli --login --management_system_ip <MNO IP> --

# username <username> --password <password>

Where **login** is the command, **username**, **password**, **management_system_ip** are the parameters, and the values follow each parameter. All actions and parameters are preceded by a double-dash "**--**". For example, the command above with values would look like: 

 In order to help in formatting syntax, the PowerFlex CLI supports autocompletion. To complete a command or parameters, press the Tab 

key while typing CLI commands. A complete list of scli commands can be found in the[ PowerFlex 4.5.x CLI Reference Guide ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_cli_reference_guide_4.5/introduction?)located in dell.com/support. 

# Query Command Examples 

The query group of commands provides general information about the health and operations of the PowerFlex system. It also enables an examination of the detailed information about most components. Complete the exercise to view more details on the **query_cluster**, the 

# query_volume**, the **query_all_sds**, and the **query_all_sdc

commands. 

In the simulation, autocomplete is not available. Use the[ PowerFlex 4.5.x ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_cli_reference_guide_4.5/introduction?) [CLI Reference Guide ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_cli_reference_guide_4.5/introduction?)instead to look up the syntax for each of the commands. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 39  |

![Image](images/admguide_Image559.png)

Management Interfaces To try to query commands, use the simulation below. The username for the exercise is "admin", and the password is "Pow3rFlex". 

*The web version* *of this content contains an interactive activity.* 

# PowerAPI 

PowerAPI (the PowerFlex REST API) allows users to automate PowerFlex deployment, configuration, and management tasks. PowerAPI is divided into two distinct sets of calls, the PowerFlex Block API, and the PowerFlex Manager API. Each set of APIs interact with a different set of components 

within a PowerFlex cluster. 

PowerAPI is installed as part of the PowerFlex Gateway package, and adheres to the standardized REST API call/response structure. For example, to log on and begin sending API calls users would issue a command similar to the following using a REST client: 

curl --location --request POST 'https://powerflex.example.com/rest/auth/login' \ --header 'Accept: application/json' \ --header 'Content-Type: application/json' \ 

--data-raw '{ 

```
   "username": "Admin",  
   "password": "Pow3rFlex"  
```

}' 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 40  |

Management Interfaces 

 **1: **The **REST Client** is a framework that allows end users to interact directly with the RESTful API through the gateway. Some examples of 

REST Clients are cURL and Postman. Users use the REST Client to send an API request. 

**2: **The REST Client sends the **request** in JSON format. **3: **The REST API is served from the **PowerFlex Gateway** through the REST gateway. The PowerFlex Gateway sends the request to the Primary 

MDM (metadata manager) in the form of queries. When the Gateway receives the response from the MDM, it reformats the response in a RESTful manner back to a REST client.  

**4: **The PowerFlex Gateway leverages role-based security access to expose RESTful APIs. Users send an API **log in request** to the PowerFlex Gateway. The API sends back a token. This token is used to authorize API requests. 

**5: **Once the token is received, subsequent communication takes place directly between the REST client and the **MDM** while the token is valid. 

When an API request is sent, the **MDM** receives the query and returns a response. 

**6: **The REST Gateway returns the **response** in JSON format. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 41  |

![Image](images/admguide_Image565.png)

Management Interfaces Review more examples of Po Developers Portal I calls in the[ Dell Technology ](https://developer.dell.com/apis/4008/versions/4.5.1/docs/Introduction/Introduction_powerAPI.md)werAP

# PowerAPI Example 

PowerAPI can be used with various REST API Clients. For this simulation, you are using the cURL CLI tool. The basic structure for cURL commands is shown below. The –k flag is for skipping the SSL certificate verification. Use this structure for the API calls in the simulation. 

curl --location --request POST '<gateway IP address>:<port>/<api call>' \ --header 'Accept: application/json' \ --header 'Content-Type: application/json' \ 

--data-raw '{ 

```
  "username": "Admin", 
    "password": "Pow3rFlex" 
```

}' The web version of this content contains an interactive activity. 

**Deep Dive: **For more examples of PowerAPI calls, see the Dell Technology Developers Portal  

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 42  |

![Image](images/admguide_Image570.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path4.svg)

Deployment Process 

# Deployment Process 

## Outline of the Implementation Process 

The following image provides a high-level overview of the implementation process. 

 

## Implementation Process 

 The first phase of a PowerFlex Deployment is to validate that all 

prerequisites are met. Checklists can be found in the Field Logical Build 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 43  |

![Image](images/admguide_Image573.png)

![Image](images/admguide_Image575.pngimages/admguide_Path3.svg)

Deployment Process Guide. Next, the PowerFlex management virtual machine templates must be downloaded from the Dell support site. Finally, configuration of the 

physical networking infrastructure, access, and aggregation switches should be completed before continuing the deployment process. 

# Implementation Process 

 Next, the virtual environment that supports PowerFlex must be configured. Configuration includes creating and configuring the distributed switches 

that the management infrastructure run on. PowerFlex Management Platform (PFMP) VMs must also be configured for the production environment. Parameters include configuring networking features such as IP address, name resolution, and network time synchronization. Once 

network configuration is complete, the logical configuration for the PFMP must be defined using the **PFMP_config.json** file on the installer VM. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 44  |

![Image](images/admguide_Image577.png)

An example **PFMP_Co** sections of configuration data: 

- Deployment Process 

 file is provided below. Note the four **nfig.json**

# 1.  **Management Node Hostname and IP Address** - the IP addresses of 

the nodes where the PFMP is deployed and the hostnames that are assigned to these nodes. DNS should be set up to resolve these names before the execution of the installation script. Ensure that the hostnames are in lowercase. 

# 2.  **Subnet IP Addresses** - a private subnet that does not conflict with any 

of the subnets in the data center. The same subnets may be used for multiple PowerFlex management platform deployments. 

# 3.  **Management, Out-of-Band Management, and Data communication** 

- for the flex-node-mgmt and oob-mgmt VLANs, three IP addresses are 
- required for the ingress (web server), NFS, and SNMP (receiving 
- SNMP traps) services. The flex-node-mgmt requires an additional IP 
- address for the SupportAssist Remote service (only on the oob-mgmt 
- VLAN). Also, the data network IP addresses are defined for the 
- PowerFlex data communication. 

# 4.  **The** **PowerFlex FQDN hostname and IP Address** - the first 

parameter is the fully qualified domain name (FQDN) or hostname. A user can connect to the PowerFlex management platform UI through a browser using the FQDN. The FQDN must be resolvable using DNS or by a host file (where the browser is running). The second parameter 

specifies the IP address for ingress to the PowerFlex management platform UI by a browser. The Ingress IP needs to be in the same 

```
range as the flex-node-mgmt pool range, and be the first IP in that 
```

pool. 

# Implementation Process 

PowerFlex is installed by running the appropriate shell scripts from the installer VM. Once the PowerFlex Management Platform is successfully deployed, the remainder of the PowerFlex infrastructure is defined. Defining these resources is achieved by completing the "Getting Started" 

wizard. The wizard defines a software catalog, and physical infrastructure such as network switches and production nodes. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 45  |

Deployment Process 

 

# Implementation Process 

 Next, the production resources are defined. The deployment can be done using a comma-delimited configuration text file, or by using one of the 

integrated templates and customizing it to the environment. Once customized it can be deployed as a resource group. Templates and resource groups are covered later in this material. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 46  |

![Image](images/admguide_Image581.png)

![Image](images/admguide_Image582.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Deployment Process 

**Go to:**  the Dell PowerFlex 4.5.x Install and Upgrade Guide to review more information on[ Deploying a PowerFlex ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_install_upgrade_guide_4.5.x/deploying-a-powerflex-cluster-using-a-csv-topology-file?guid=guid-1e17cfb9-ca2c-4fd3-b9d7-422f393a7a3f&lang=en-us) [cluster using a CSV topology file.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_install_upgrade_guide_4.5.x/deploying-a-powerflex-cluster-using-a-csv-topology-file?guid=guid-1e17cfb9-ca2c-4fd3-b9d7-422f393a7a3f&lang=en-us)  

 

# Implementation Process 

The post-installation checklist is a list of configuration tasks to perform, some mandatory and others optional, after the installation of a PowerFlex system. All the mandatory configuration post-installation tasks must be completed in order for the PowerFlex system to function. 

# Post-Installation Checklist

| **Checklist item** |  **Mandatory ** | **Optional tasks**  |

# tasks

**Check port usage.**       

# Configure NVMe targets and

# hosts for block** **storage if they

# were not configured** **by the

**installation template.** 

# Add block storage if it was not

# configured by the installation

# template, in this order: Add

# protection domains > add

# storage pools > add acceleration

# pools (for fine granularity data

# layout only) > add devices to

# SDSs (including mapping them to

# storage pools) > add volumes >

**map volumes to hosts.** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 47  |

![Image](images/admguide_Image485.png)

![Image](images/admguide_Image586.png)

![Image](images/admguide_Image586.png)

![Image](images/admguide_Image586.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Post-Installation Tasks 

## Configure SDSs for optimal

**memory usage.**      

## Configure SDSs for NVDIMM

**health awareness.**      

## Configure SDC authentication for

**hardened access to the MDM.**      

## Enable** **OpenSSL Federal

## Information Processing

**Standards (FIPS) compliance.**  

## Configure local and Lightweight

## directory** **access protocol (LDAP)     

## users

**Install the system's license.**       

## Configure event and alert

**notifications to external systems.**      

## Install** **and configure** **Sanford

**Rose Associates (SRA) software.**       

# Post-Installation Tasks 

## Getting Started Wizard: Download a Compliance File 

The Getting Started wizard is provided once logged in to a new PowerFlex cluster, to walk administrators through the initial configuration tasks. The tasks covered in the wizard are also available elsewhere in the UI, the wizard simply consolidates them into one place. The first step in the 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 48  |

![Image](images/admguide_Image586.png)

![Image](images/admguide_Image586.png)

![Image](images/admguide_Image586.png)

![Image](images/admguide_Image586.png)

![Image](images/admguide_Image586.png)

![Image](images/admguide_Image586.png)

![Image](images/admguide_Image586.png)

![Image](images/admguide_Image586.png)

Post-Installation Tasks wizard is to download a compliance file. A compliance file contains a compilation of software that is certified to interact without issues. These 

code bases are what are used to deploy new hardware to the cluster and remediate existing hardware to supported levels. It is up to an administrator to download the compliance file from Dell support and specify it for use in PowerFlex Manager. 

To try uploading a new compliance file to a PowerFlex rack cluster, use the simulation below. 

For the simulation, you need the path to a previously downloaded PowerFlex 4.5.1 RCM file and credentials to access the file location: 

Path: \\127.0.0.1\share\Flex_RCM_3_7_3_1_r6_2.zip User: Administrator 

Password: P0werFl3x!  *The web version* *of this content contains an interactive activity.* 

# Getting Started Wizard: Define the Networks 

The next step in the Getting Started wizard is to define the networks PowerFlex will use for the various functions of the cluster. The variety and quantity of networks vary by environment, but at a minimum includes: 

- $•$  Hypervisor Management 
- $•$  PowerFlex Management 
- $•$  General Purpose LAN 
- $•$  Hypervisor Migration 
- $•$  PowerFlex Data 

To try defining one network for a PowerFlex rack cluster, use the simulation below. The same process would then be repeated for each network definition. The information about each network can be found in the customer site documentation. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 49  |

Post-Installation Tasks *The web version* *of this content contains an interactive activity.* 

# Getting Started Wizard: Discover Resources 

Once the compliance file is uploaded and the networks are defined, the next step is to discover resources in the PowerFlex cluster. This step allows the cluster to discover and usually manage the resources that are a part of the configuration. These resources include: 

- $•$  Nodes 
- $•$  Element Managers 
- $•$  Switches 
- $•$  Virtual Machine Managers (vCenter) 
- $•$  PowerFlex Gateway 
- $•$  PowerFlex System 

Once the resource discovery is started, the progress is monitored in the Resources tab of PowerFlex Manager. 

To try discovering managed switches and nodes to a PowerFlex rack cluster, use the simulation below. Other resources such as VMware vCenter and CloudLink Center may be required depending on the cluster configuration. 

Use the following values for the simulation:  

| **Managed Switches** |  **PowerFlex Nodes**  |

IP range: 10.10.101.2 - 10.10.101.9  IP range: 10.10.101.10 - 10.10.101.30 

| Credential Name: Switches |  Credential Name: iDRAC  |
| Username: root |  Username: root  |
| Password: P0werFl3x! |  Password: P0werFl3x!  |

       *The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 50  |

Post-Installation Tasks 

# Getting Started Wizard: Manage Deployed Resources 

The fourth step in the Getting Started wizard is to manage deployed resources. Managing deployed resources is optional. Included in managing deployed resources is an existing PowerFlex Cluster (3.x) that can be ingested into the new PFMP 2.0 environment. 

 

# Getting Started Wizard: Deploy Resources 

When deploying resources, there are a few options. The first is creating a template to deploy to discovered nodes. Another is cloning one of the existing templates and modifying it to fit the needs of the implementation. There is also an option to upload a pre-existing template from outside the 

cluster. Finally, administrators can use a comma delimited deployment file to define the resources and configuration of a cluster. Available sample templates consist of compute-only, storage-only, hyperconverged, PowerFlex-file, and CloudLink Center varieties. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 51  |

![Image](images/admguide_Image595.pngimages/admguide_Path3.svg)

Post-Installation Tasks 

 

# Deploy the Production Cluster - CSV Method 

For REST API deployments using a CSV topology file, see the[ Dell ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_install_upgrade_guide_4.5.x/deploying-a-powerflex-cluster-using-a-csv-topology-file?guid=guid-1e17cfb9-ca2c-4fd3-b9d7-422f393a7a3f&lang=en-us) [PowerFlex 4.5.x Install and Upgrade Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_install_upgrade_guide_4.5.x/deploying-a-powerflex-cluster-using-a-csv-topology-file?guid=guid-1e17cfb9-ca2c-4fd3-b9d7-422f393a7a3f&lang=en-us) 

 Deploying a PowerFlex cluster using a CSV topology file is an advanced deployment method for deploying block storage. A CSV topology file can 

be configured, uploaded, and used to deploy using either REST API commands or PowerFlex Manager. This section outlines the deployment using PowerFlex Manager. 

Details specific to the environment are entered into the CSV file, and then PowerFlex reads the configuration and builds it accordingly. Shown below is an example of a completed configuration file. 

 To reconfigure existing components, use a copy of the current CSV topology file, and add rows or columns, or edit existing rows and columns. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 52  |

![Image](images/admguide_Image599.png)

![Image](images/admguide_Image600.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path7.svg)

![Image](images/admguide_Path8.svg)

Post-Installation Tasks For new features, you can copy columns from the latest complete CSV template and paste them into the copy of the current CSV file. 

**Important**: Save a copy of the CSV file for future use. It is recommended to save the file in a secure encrypted folder on your operating system after deployment. When you add nodes to the system in the future, you need the CSV file that 

 you used for initial deployment of the system. 

# CSV File Formatting 

 Sample CSV files can be downloaded from the PowerFlex Manager UI when deploying from an installation file. A good practice is to download 

one or both files, and then modify them to suit the needs of the deployment. The table below lists all the fields available for configuration and their function. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 53  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Image604.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path14.svg)

Post-Installation Tasks During CSV file preparation there are important differences between hyperconverged and two-layer topologies, namely: 

$•$  In hyperconverged topologies, each row above the Storage Pool Configuration row in the CSV file represents a server that can host MDM, SDS, and SDC components. 

$•$  In two-layer topologies where front-end servers are Linux or Windows-based, SDCs are represented by separate, additional rows. The values for SDC-specific rows are described in CSV topology for two-layer deployment. Alternatively, administrators may omit the SDCs from the 

CSV file, and install them manually after deploying the backend. 

**Go to:**  To learn more about formatting the CSV file, refer to the[ Dell PowerFlex 4.5.x Install and Upgrade Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_install_upgrade_guide_4.5.x/prepare-the-csv-topology-file?guid=guid-928d51c0-c3b2-4bc2-9df8-0824798e7b65&lang=en-us)  

# Apply a PowerFlex License 

The final step in implementing a PowerFlex 4.x cluster is to apply a license file. By default, a new PowerFlex cluster runs for 90 days under the TRIAL license. After that, a license file must be obtained and applied either through the getting started wizard, or from Settings > License 

Management. The Product Version, Installation ID, and capacity must be provided to Dell to obtain a license. 

 When you upload a license file, PowerFlex Manager checks the license file to ensure that it is valid. After the upload is complete, PowerFlex 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 54  |

![Image](images/admguide_Image485.png)

![Image](images/admguide_Image607.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path11.svg)

Manager stores the license details and d Manager License page. 

Post-Installation Tasks ays them on the PowerFlex ispl

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 55  |

Upgrade a PowerFlex Cluster 

# Upgrade a PowerFlex Cluster 

## Upgrade a PowerFlex 4.x System 

An upgrade engineer upgrades the PFMP and PowerFlex Manager software when upgrading a PowerFlex system in 4.x. 

Backup PowerFlex Manager Before upgrading PowerFlex Manager, backup the current PowerFlex Manager configuration. Review the[ Dell PowerFlex 4.5.x Administration ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/backup-details?guid=guid-10d07c40-6ab5-499b-8481-2172a00d6c31&lang=en-us)

[Guide ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/backup-details?guid=guid-10d07c40-6ab5-499b-8481-2172a00d6c31&lang=en-us)for a comprehensive list of what files are backed up in this task. 

## Video:

*The web version* *of this content contains a* *Video.* Upload the Compatibility Management and Compliance file 

An updated Compatibility Management file is required before an upgrade of the PowerFlex PFMP and PowerFlex Manager can be started. The Compliance file will identify what PowerFlex System resources need upgrading to the target version. 

## Video:

*The web version* *of this content contains a* *Video.* Upgrade the PFMP, PowerFlex Manager software and Update the PowerFlex System Resource 

The PowerFlex Management Platform (PFMP) is key to the management of PowerFlex 4.x. Anytime an upgrade is performed for PowerFlex 4.x, the PFMP must be upgraded to match the same version. 

## Video:

*The web version* *of this content contains a* *Video.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 56  |

Upgrade a PowerFlex Cluster Upgrade the SDCs 

PowerFlex 4.x upgrades will require the engineer to reboot the nodes that host the SDCs. However, in a software-only co-resident environment, the engineer will use manual procedures to upgrade the SDCs, especially nodes running the PFMP. The SDC nodes in a co-resident environment 

must be properly prepared for the upgrade before they are rebooted. 

**Important**: While the example used is a PowerFlex 4.0.2 to a PowerFlex 4.5.1 system upgrade, the same process workflow is used for all 4.x.x to 4.x.x upgrades.  

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 57  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

PowerFlex Templates and Resource Groups 

# PowerFlex Templates and Resource Groups 

## Templates 

Templates in PowerFlex define the resources that are used by a production cluster. The number of nodes, node type, networks, and access privileges are all defined in the template. For most environments, you can clone one of the sample templates that are provided with 

PowerFlex Manager and edit as needed. 

Access Templates from the Lifecycle menu 

 From the dashboard, locate the **Lifecycle** menu item and select **Templates** from the dropdown menu. Click **Create** to begin. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 58  |

![Image](images/admguide_Image614.pngimages/admguide_Path2.svg)

PowerFlex Templates and Resource Groups Create a Template 

 Templates can be created in multiple ways: 

$•$  **Clone an existing PowerFlex Manager template -** Use one of the integrated sample templates as a base, and customize for the environment. $•$  **Upload External Template -** Create a template from an external file. 

$•$  **Create a new template** **-** Create a template from scratch, defining all properties individually. 

Provide Template information 

 In the Clone Template information page, provide the following information: 

$•$  Enter a **Template Name**. $•$  Select a **Template Category** from the list. 

$•$  Enter a **Template Description** (optional). 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 59  |

![Image](images/admguide_Image617.png)

![Image](images/admguide_Image618.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

PowerFlex Templates and Resource Groups $•$  Specify the **Firmware** **and Software** **Compliance** version. $•$  Specify **Who should have** **access to the Resource Group deployed **

**from this Template**. 

Complete the Template Settings 

 On the Additional Settings page, complete the remaining template 

settings: $•$  **Network Settings** 

$•$  **PowerFlex** **Gateway** **Settings** $•$  **Node Pool Settings** 

Click **Finish** to add the new template. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 60  |

![Image](images/admguide_Image620.pngimages/admguide_Path1.svg)

PowerFlex Templates and Resource Groups Edit the Draft Template 

 The template is created in draft mode. Additional settings to the new template can be made through the resource group diagram. 

$•$  Add nodes or clusters to your template. Adding nodes or clusters allows you to set the settings for each. $•$  Edit the nodes or clusters in the template. For example, sample templates may contain only the latest best practice node networking 

configurations. If customers have older configurations, the template may have to be edited to match their settings. $•$  Edit the template information, such as the name and description, template category, compliance file, or the administrative rights for the 

template. 

Publish the Template Once all template edits are completed, the template is ready to be published. From the template details page, click **Publish Template** to add 

it to the Templates list. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 61  |

![Image](images/admguide_Image622.pngimages/admguide_Path1.svg)

PowerFlex Templates and Resource Groups 

 

# Resource Groups 

Resource Groups are created by deploying a finished template. The deployment process configures the available nodes with an operating system, IP addresses, and PowerFlex software to create a production cluster. 

Deploy Template The first step in creating a Resource Group is to **Deploy** the published template: 

- PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 62  |

![Image](images/admguide_Image625.pngimages/admguide_Path2.svg)

PowerFlex Templates and Resource Groups 

 Provide Resource Group Information 

 In the **Deploy Resource Group** information page, provide the following 

information: $•$  Enter the template name in the **Select Published Template** text entry 

box. $•$  Provide a **Resource Group Name**. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 63  |

![Image](images/admguide_Image628.jpeg)

![Image](images/admguide_Image630.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

PowerFlex Templates and Resource Groups $•$  Enter a **Resource Group Description** (optional). $•$  Specify the version to use for compliance in the **Firmware and **

**Software Compliance** entry box. $•$  Specify **Who should have** **access to the to the Resource Group ** **from this template**. $•$  Click **Next** to continue. 

Validate the Deployment Settings 

 On the **Deployment Settings** page, verify the remaining deployment 

settings: $•$  **PowerFlex Cluster** 

$•$  **Physical Nodes** Click **Next** to go to deployment scheduling. 

Schedule the Deployment 

 In the **Schedule Deployments** section, deployments can be scheduled to 

start immediately following the wizard completion, or to start later. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 64  |

![Image](images/admguide_Image632.png)

![Image](images/admguide_Image633.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

PowerFlex Templates and Resource Groups Review the Summary 

 The **Summary** section of the Deploy Resource Group provides a recap of 

the deployment. Clicking **Finish** will start the deployment process (unless scheduled for a later time). 

# Working With Templates and Resource Groups: 

# Demonstration 

To see a demonstration of deploying a PowerFlex cluster as a resource group using a template, watch the video below. 

*Video:* *The web version* *of this content contains a* *Video.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 65  |

![Image](images/admguide_Image635.pngimages/admguide_Path3.svg)

Configure Protection Domains 

# Configure Protection Domains 

## About Protection Domains 

In PowerFlex, storage is organized into a storage pool which resides in a Protection Domain. A Protection Domain provides data isolation, security, and performance benefits for the pool. An SDS and the storage pool participate in one Protection Domain at a time. 

Logical Groups 

 Protection Domains contain SDSs that are configured as a separated logical group. The function of the logical group is to physically isolate 

specific data, such as production environment data, into a dataset. Isolation of the dataset provides a performance increase for the member SDSs (or Tenants) and limits the effects of a node device failure. 

When data is written to a PowerFlex cluster, the data is mirrored on two different SDSs within a Protection Domain. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 66  |

![Image](images/admguide_Image637.png)

Protection Domain Benefits Configure Protection Domains 

 Protection Domain benefits include protection against simultaneous 

failures, performance isolation, data location control, and network constraints. 

$•$  Improved system resilience - I/O is unaffected when a server or media device fails. $•$  Establish SLA tiers by separating volumes for performance planning. $•$  The Protection Domain SDS has efficient and secure data transfer. 

$•$  Spread out network workloads evenly on the vLANs. 

**Important**: Protection Domains can also contain SDTs and SDRs, as well as Acceleration Pools and Fault Sets.  

# Configure a Protection Domain with PowerFlex 

# Manager 

IT administrators can add, modify, activate, inactivate, or remove a Protection Domain in the PowerFlex system with PowerFlex Manager. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 67  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Image640.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

To learn more about configuring a Prote Manager, select each tab below. 

Add a Protection Domain 

Configure Protection Domains n Domain in PowerFlex 

 

ctio

After the Protection Domain is created, the administrator can add  SDSs, fault sets, storage pools, and acceleration pools to the Protection Domain. Replication can also be set up to ensure that the data is protected and saved to a remote cluster. 

Modify the Protection Domain 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 68  |

![Image](images/admguide_Image642.png)

![Image](images/admguide_Image643.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

Configure Protection Domains There are two options to modify a Protection Domain. 

Rename: Change the name of a Protection Domain. The administrator can start the Protection Domain with a test name. After validation, the administration 

renames the Protection Domain to match production environment naming conventions for the PowerFlex system. 

Network Throttling: Network throttling is configured separately for each protection domain and is used to limit the available bandwidth for backend data transfers. 

Inactivate or activate a Protection Domain

 Inactivate or activate a Protection Domain 

When a Protection Domain is created, it is automatically activated. However, IT administrators can also inactivate the Protection Domain at any time. When a Protection Domain is inactivated, the data remains on the SDSs and allows the user to conduct a graceful system shutdown. The 

user is then free to conduct any management activities on the Protection Domain, such as a reload of all the SDSs. When the activities are complete, the user can activate the Protection Domain to enable access to data. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 69  |

![Image](images/admguide_Image645.jpegimages/admguide_Path1.svg)

Remove a Protection Domain Configure Protection Domains  

When you inactivate a Protection Domain, the data remains on the SDSs. It is preferable to remove a Protection Domain when the organization no longer needs it. 

Prerequisites: Ensure that all SDSs, storage pools, acceleration pools, and fault sets 

have been removed from the Protection Domain before removing the Protection Domain from the system. 

**Caution**: Protection Domain configuration is an advanced user task. A nonadvanced user should contact Dell Technologies Support for help with Protection Domain 

configuration. Configuration of Protection Domains affects  system performance, so caution is warranted. 

 

# Configure a Protection Domain 

To try configuring a Protection Domain, use the simulation below. *The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 70  |

![Image](images/admguide_Image647.png)

![Image](images/admguide_Image649.jpeg)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path14.svg)

Configure Fault Sets 

# Configure Fault Sets 

## Fault Sets 

Fault Sets are logical entities that contain a group of SDSs within a protection domain. A Fault Set is defined for a set of servers that are likely to fail together, for example an entire rack full of servers. 

Fault Sets Overview 

 PowerFlex requires a minimum of three Fault Sets per protection domain, with at least two nodes in each Fault Set. Each Fault Set acts as a single 

fault unit. When defining fault sets, the term fault units refer to either a fault set, or an SDS not associated with a fault set. 

In PowerFlex, it is a requirement to have two nodes in each Fault Set to achieve high availability and fault tolerance. The two-node configuration provides redundancy, minimizing the risk of data unavailability and enhancing system reliability if there were to be a node failure. 

However, the actual number of nodes in a Fault Set can be adjusted according to specific requirements and the desired level of redundancy. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 71  |

![Image](images/admguide_Image651.png)

Configure Fault Sets Fault Sets Data Mirroring 

 PowerFlex maintains a copy of all data chunks within the Fault Set on 

SDSs outside of itself. 

Mirroring Fault Set data ensures that there is always another copy available even if all the servers within the defined Fault Set fail simultaneously. 

There must be enough capacity mirroring. in at least three fault units to enable  with

In the example Fault Set diagram, the combination of Fault units is: $•$  FS-1, FS-2, and FS-3 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 72  |

![Image](images/admguide_Image653.png)

Configure Fault Sets Add Fault Sets 

 Fault Sets are deployed automatically when a template is deployed to 

create a Resource Group. If manual deployment methods are used, the following requirements must be verified before a Fault Set can be added to a Protection Domain. 

$•$  Ensure that a Protection Domain exists or add a new Protection Domain. $•$  Ensure that a Storage Pool and Fault Sets (with a minimum of three Fault units) exist or add new ones. 

$•$  Add the SDS, designate a protection domain and the fault set and at the same time add the SDS devices into a storage pool. 

Steps: 

 1.  On the menu bar, click **Block** then **Fault Sets**. 

 2.  In the right pane, click **+Create Fault Set**. 

 3.  In the Create Fault Set dialog box, enter a name and select the 

protection domain, and click **Create**.       

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 73  |

![Image](images/admguide_Image655.pngimages/admguide_Path1.svg)

Configure Fault Sets 

 In PowerFlex 4.5.x or greater, when you create a new Fault Set, the new 

Fault Set can be added to an existing service. Also, existing nodes or newly discovered nodes can be added to the new Fault Set. 

Delete Fault Sets 

 Ensure that any configured S

before attempting to delete it.  have been removed from the Fault Set DSs Steps: 

 1.  On the menu bar, click **Block** then **Fault Sets**. 

 2.  From the list of Fault Sets, select the Fault Set to be deleted. 

 3.  Select **More Actions** then **Delete**. 

 4.  In the **Delete Fault Set** dialog box, verify that the desired fault set will 

be deleted, and click **Delete**. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 74  |

![Image](images/admguide_Image657.png)

![Image](images/admguide_Image658.pngimages/admguide_Path1.svg)

Configure Fault Sets 

**Important**: You can only create and configure Fault Sets before adding SDSs to the system, and configuring the Fault Sets incorrectly may prevent the creation of volumes. An SDS can be added to a Fault Set only during the creation of 

 the SDS. You can also add Fault Sets when adding SDS nodes after initial installation. 

 

# Create and Delete Fault Sets 

To try creating and deleting the Fault Sets, use the simulation below. *The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 75  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Configure Storage Pools 

# Configure Storage Pools 

## Introduction** **to** **Storage** **Pools

For more information, see the[ Add Storage Pools ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/add-storage-pools?guid=guid-76701bd1-85ef-41d5-a9de-1b0c2041ae4f&lang=en-us)section of the PowerFlex 4.5.x Administration Guide. 

 A Storage Pool is a subset of physical storage devices in a Protection Domain. Each storage device belongs to only one Storage Pool. When a 

PowerFlex volume is configured, the volume contents are distributed over all the devices residing in the same Storage Pool. Each block of data on a volume consists of two copies that are on different SDSs. The copies enable the system to maintain data availability following a single device, 

network, or server failure. The figure below shows two storage pools, spread across multiple SDSs. 

Storage Pools enable creating different storage tiers in the PowerFlex system. Devices belonging to a storage pool should have the same media type, size, and interface, ensuring each volume is distributed over devices of the same performance profile. PowerFlex provides two types of Storage 

Pools: Medium Granularity (MG) and Fine Granularity (FG). The granularity of a storage pool refers to its allocation unit size. 

$•$  In MG storage pools, volumes are divided into 1 MB allocation units, which are distributed and replicated across all disks contributing to a pool. MG storage pools are best suited for a performance-driven workload. 

$•$  FG storage pools are more space efficient, with an allocation unit of 4 KB. The physical data placement scheme is based on a Log Structure Array (LSA) architecture that is built on NVDIMMs. If administrators want to enable compression, FG storage pools are required. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 76  |

![Image](images/admguide_Path1.svg)

Configure Storage Pools 

 

# Acceleration Pools 

Acceleration Pools are used to provide write-caching using NVDIMM devices when configured with Fine Granularity storage pools. 

When deploying nodes, the check box **Enable Compression** is available under the template settings for the Protection Domain. If checked, when deploying the template, PowerFlex uses nodes that have at least two NVDIMMs installed. Compatible nodes also require SSD or NVMe devices 

and have persistent memory that is turned on. When these conditions are met, the resource group template adds fields that allow administrators to specify the acceleration pool name, and granularity setting. Deploying the template results in the Acceleration Pool being created. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 77  |

![Image](images/admguide_Image665.png)

Configure Storage Pools 

 

# Configure Storage Pool Settings 

Configure storage pool settings, compression. ding checksum, zero padding, and inclu

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 78  |

![Image](images/admguide_Image668.pngimages/admguide_Path2.svg)

Configure Storage Pools 

 **1: **Enable Rebuild/Rebalance: By default, the rebuild/rebalance features are enabled in the system because they are essential for system health, 

optimal performance, and data protection. **2: **Enable Inflight Checksum: Inflight checksum protection mode can be used to validate data reads and writes in storage pools in order to protect 

data from data corruption. **3: **Enable Persistent Checksum: This can be used to support the medium granularity layout in protecting the storage device from data corruption. 

Select **validate on read** to validate data reads in the storage pool. **4: **Enable Zero Padding: Use the zero-padded policy when the storage pool data layout is fine granularity. The zero-padded policy ensures that 

every read from an area that is previously not written to returns zeros. **5: **Enable Compression: For fine granularity storage pools, inline compression allows you to gain more effective capacity. 

For more information, see the[ Configure Storage Pool Settings ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/configure-storage-pool-settings?guid=guid-8d1d3b75-7885-4dc2-8f6b-c55bd1d7ab94&lang=en-us)section of the PowerFlex 4.5.x Administration Guide. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 79  |

![Image](images/admguide_Image672.jpegimages/admguide_Path1.svg)

Configure Storage Pools 

# Enable the Background Device Scanner 

Enable the background device in the specified storage pool. nner to check for errors on the devices  sca

 **1: Enable the Background Device Scanner** 

Check for errors on the devices in the specified storage pool. **2: Fix Local Device Errors** 

Automatically fixes device errors if they are found. **3: Compare Data** 

Compares primary and secondary copies of data. This setting is only available if zero padding is enabled. 

**4: Fix Comparison Errors** If Compare Data is selected, Fix Comparison Errors will enable the system 

to automatically fix local device errors. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 80  |

![Image](images/admguide_Image674.png)

**5: Bandwidth Limits** Configured in KB/s, the default Configure Storage Pools it is 3072 KB/S.  lim

**Note:** Higher bandwidth settings can affect the performance. 

For more information, see the[ Enable the Background Device Scanner ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/enable-the-background-device-scanner?guid=guid-42e5e8fb-8881-4e1d-a156-a286372947b5&lang=en-us) section of the PowerFlex 4.5.x Administration Guide. 

# Configure I/O Priorities (Rebuild/Rebalance) 

PowerFlex maintains user data in a distributed mesh mirrored layout. Each piece of data is stored on two different fault units. The copies are distributed over the storage devices according to an algorithm that ensures uniform load of each fault unit by terms of capacity and expected 

network load. Rebuild and rebalance processes are fully automated but are configurable. 

Rebuild When a failure occurs, such as on a server, device, or network failure, PowerFlex immediately initiates a process of protecting the data. This 

process is called rebuild, and comes in two flavors: $•$  **Forward rebuild** is the process of creating another copy of the data on a new server. In this process, all the devices in the storage pool work 

together, in a many-to-many fashion, to create new copies of all the failed storage blocks. This method ensures a fast rebuild. $•$  **Backward rebuild** is the process of re-synchronization of one of the copies. This re-synchronization is done by passing to the copy only 

changes made to the data while this copy was inaccessible. This process minimizes the amount of data that is transferred over the network during recovery. 

PowerFlex automatically selects the type of rebuild to perform. This fact implies that sometimes more data is transferred to minimize the time that the user data is not fully protected. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 81  |

Configure Storage Pools Rebuild Throttling 

Rebuild throttling sets the rebuild priority policy for a storage pool. The policy determines the priority between the rebuild I/O and the application I/O when accessing SDS devices. Regardless of which policy is chosen, application I/Os are continuously served. 

Applying rebuild throttling increases the time that the system is exposed with a single copy of some of the data but will reduce impact on the application. The right balance between the two must be determined. 

 The following priority policies may be applied: 

$•$  **No Limit:** No limit on rebuild I/Os. Any rebuild I/O is submitted to the device immediately, without further queuing. $•$  **Limit Concurrent I/O:** Limit the number of concurrent rebuild I/Os per SDS device (default). The rebuild I/Os are limited to a predefined 

number of concurrent I/Os. When the limit is reached, the next incoming rebuild I/O waits until the completion of a currently running rebuild I/O. This policy completes the Rebuild quickly for best reliability, however, there is a risk of host application impact. 

$•$  **Favor Application I/O:** Limit rebuild in both bandwidth and concurrent I/Os. The rebuild I/Os are limited both in bandwidth and in the number of concurrent I/Os. If the number of concurrent rebuild I/Os, and the 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 82  |

![Image](images/admguide_Image679.pngimages/admguide_Path1.svg)

Configure Storage Pools bandwidth they consume, do not exceed the predefined limits, rebuild I/Os will be served. Once either threshold is reached, the rebuild I/Os 

wait until both I/O and bandwidth are below their thresholds. For example, setting the value to **1** guarantees the device only has one concurrent rebuild I/O at any given moment, which ensures the application I/Os only wait for one rebuild I/O. This imposes bandwidth 

on top of the Limit Concurrent I/Os option, which is a prerequisite to using this policy. $•$  **Dynamic Bandwidth Throttling:** This policy is similar to Favor Application I/O but extends the interval in which application I/Os are 

considered to be flowing by defining a minimal quiet period. This quiet period is defined as a certain interval in which no application I/Os occurred. The limits on the rebuild bandwidth and concurrent I/Os are still imposed. 

The default policy for rebuild is Lim one concurrent I/O. oncurrent I/O with an I/O limit set to it C

Rebalance When PowerFlex detects that user data is not balanced across devices in the storage pool, it initiates a process to restore the balance. This process 

involves copying data from the most used devices to the least used. 

Rebalance is the process of moving one of the data copies to a different server. It occurs when PowerFlex detects that the user data is not evenly balanced across the fault units in a storage pool. This imbalance can occur because of several conditions such as: SDS addition or removal, 

device addition or removal, or following a recovery from a failure. PowerFlex moves copies of the data from the most used devices to the least used ones. 

Both rebuild and rebalance compete with the application I/O for the system resources, which include network, CPU and disks. PowerFlex provides a rich set of parameters that can control this resource consumption. While the system is factory-tuned for balancing between 

speedy rebuild or rebalance and minimization of the effect on the application I/O, Administrators have fine-grain control over the behavior. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 83  |

Configure Storage Pools Rebalance Throttling 

Rebalance throttling sets the rebalance priority policy for a storage pool. The policy determines the priority between the rebalance I/O and the application I/O when accessing SDS devices. Note that application I/Os are continuously served. Rebalance, unlike rebuild, does not impact the 

reliability of the system, and therefore reducing its impact is not risky. 

 The following possible priority policies may be applied: 

$•$  **No Limit:** No limit on rebalance I/Os. Any rebalance I/O is submitted to the device immediately, without further queuing. Note that rebalance I/Os are relatively large and hence setting this policy speeds up the rebalance, but has the most effect on the application I/O. 

$•$  **Limit Concurrent I/O:** Limit the number of concurrent rebalance I/Os per SDS device. The rebalance I/Os are limited to a predefined number of concurrent I/Os. Once the limit is reached, the next incoming rebalance I/O waits until the completion of a currently running 

rebalance I/O. For example, setting the value to **1** guarantee that the device has only one rebalance I/O at any given moment, which ensures that the application I/Os only wait for one rebalance I/O in the worst case. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 84  |

![Image](images/admguide_Image682.pngimages/admguide_Path1.svg)

Configure Storage Pools $•$  **Favor Application I/O:** Limit rebalance in both bandwidth and concurrent I/Os. The rebalance I/Os are limited both in bandwidth and 

in the number of concurrent I/Os. If the number of concurrent rebalance I/Os, and the bandwidth they consume, do not exceed the predefined limits, rebalance I/Os are served. Once either limiter is reached, the rebalance I/Os wait until such time that the limits are not 

met again. This policy imposes a bandwidth limit on top of the Limit Concurrent I/Os option. $•$  **Dynamic Bandwidth Throttling:** This policy is similar to Favor Application I/O, but extends the interval in which application I/Os are 

considered to be flowing by defining a minimal quiet period. This quiet period is defined as a certain interval in which no application I/Os occurred. The limits on the rebalance bandwidth and concurrent I/Os are still imposed. 

The default policy for rebalance is to favor application I/O, with a rebalance concurrent I/O Limit of one concurrent I/O per SDS device, and a rebalance bandwidth limit of 10240 KB/s. 

**Important**: Rebuild and rebalance throttling affects the performance of the system and should only be used by advanced users. For more information, see the[ Configure ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/configuring-io-priorities-and-bandwidth-use?guid=guid-5ee6b49f-ab70-49bb-b5ac-a7355df8d2da&lang=en-us) [IOPS and Bandwidth ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/configuring-io-priorities-and-bandwidth-use?guid=guid-5ee6b49f-ab70-49bb-b5ac-a7355df8d2da&lang=en-us)section of the PowerFlex 4.5.x 

 Administration Guide.  

# Configure Storage Pools 

# Important**: In the simulation exercise, the **Use Read RAM

**Cache** option defines whether the system stores the data of this storage pool's **writes** in the SDS RMcache, or not. The default is to store the **write** data in cache (cached). The Storage Pool RMcache features are advanced features, and 

 it is recommended to accept the default values. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 85  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Image602.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path14.svg)

![Image](images/admguide_Path15.svg)

![Image](images/admguide_Path16.svg)

Create, Modify, and Delete Volumes 

# Create, Modify, and Delete Volumes 

## Adding Volumes Using PowerFlex Manager 

After storage devices are added to the storage pool, a PowerFlex volume can be created. A PowerFlex volume is considered to be similar to a Logical Unit Number (LUN) from a physical storage array. To start allocating volumes, the system requires that there be at least three SDS 

nodes. Each SDS should be in a separate fault unit, and each device has a minimum of 240 GB free storage capacity. For users and applications on hosts to have access to a volume, the volume must be mapped to an SDC host. 

A volume name must contain fewer than 32 characters, contain only alphanumeric and punctuation characters, and be unique. The size of the volume that is created can range from a minimum of 8 GB to a maximum of 1 PB. 

To try adding volumes, use the simulation below. *The web version* *of this content contains an interactive activity.* 

## Increase Volume Size Using PowerFlex Manager 

Volume capacity can be increased (but not decreased) when there is available free capacity. Expansion can be performed on a single, or multiple volumes simultaneously. 

To try increasing the size of a volume, use the simulation below. *The web version* *of this content contains an interactive activity.* 

## Map Volumes to a Host Using PowerFlex Manager 

Mapping exposes the volume to the specified host, creating a block device on the host. You can map a volume to one or more hosts. Volumes can be mapped to either an SDC or NVMe host, but not both simultaneously. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 86  |

Create, Modify, and Delete Volumes Ensure that you know which type of hosts are being used for each volume, to avoid mixing host types. 

To try mapping volumes, use the simulation below. *The web version* *of this content contains an interactive activity.* 

# Set Volume Bandwidth and IOPS Limits Using 

# PowerFlex Manager 

Setting bandwidth and IOPS limits for volumes lets you control the quality of service. Bandwidth and IOPS limits are set on a per-host basis. 

Ensure that the volumes are mapped before you set the limits. 

To try setting volume limits, use the simulation below. *The web version* *of this content contains an interactive activity.* 

# vTree Migration Using PowerFlex Manager 

PowerFlex uses a volume tree (vTree) structure to manage snapshots of volumes. A vTree is a structure that spans from the source volume as the root, to its derivative snapshots (children and siblings). Snapshots are either copies of the base volume or descendants (snapshots of snapshots) 

of it. No matter how many copies exist, as long as they relate back to the source volume, they all belong to the same vTree. 

Migration of a vTree allows you to move a source volume, along with all snapshots of that volume, to a different storage pool. It frees up capacity in the source storage pool. For example, you can migrate a vTree from an HDD-based storage pool to an SSD-based storage pool. Or migrate a 

vTree from a storage pool with different attributes such as thin or thick. 

VTree migration is a long proce on the size of the vTree. nd can take days or weeks, depending ss a

To try performing a vTree migration, use the simulation below. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 87  |

Create, Modify, and Delete Volumes *The web version* *of this content contains an interactive activity.* 

# Unmap Volumes from a Host Using PowerFlex 

# Manager 

When volumes are no longer needed, they can be removed from PowerFlex. Before they can be removed, they must be unmapped from the host they are connected to. 

To try unmapping volumes, use the simulation below. *The web version* *of this content contains an interactive activity.* 

# Removing Volumes Using PowerFlex Manager 

When removing volumes in a PowerFlex system, always ensure that the volume that you are deleting is not mapped to any hosts. If it is mapped, unmap it before deleting the volume as all data is erased from a deleted volume. 

In addition, ensure that the volume is not the source volume of any snapshot policy. If it is, then you must remove the volume from the snapshot policy before you can remove the volume. 

To try removing a volume, use the simulation below. *The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 88  |

Create a NAS Server (NFS and SMB) 

# Create a NAS Server (NFS and SMB) 

## Overview of a File System Storage 

File System represents a set of storage resources that provide network file storage. The storage system establishes a file system that Windows users or Linux/UNIX hosts can connect to and use for file-based storage. Users access a file system through its shares, which draw from the total storage 

that is allocated to the file system. 

NAS Server $•$  NAS servers provide access to file systems. Each NAS server supports Windows (SMB) file systems, Linux/UNIX (NFS) exports, or 

both. $•$  A file server that is configured with its network interfaces and other settings exclusively exporting the set of specified file systems through 

mount points are called "shares". $•$  Client systems connect to a NAS server on the storage system to get access to the file system shares. 

$•$  A NAS server can have more than one file system, but each file system can only be associated with one NAS server. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 89  |

Create a NAS Server (NFS and SMB) 

 File System 

A File system is a manageable container for file-based storage that is associated with the following properties: 

$•$  A specific quantity of storage. $•$  A particular file access protocol (SMB, NFS, or multiprotocol). 

$•$  One or more shares (through which network hosts or users can access shared files or folders). 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 90  |

![Image](images/admguide_Image692.png)

Create a NAS Server (NFS and SMB) 

 Share or Export 

$•$  Share or export is a mountable access point to file system storage that network users and hosts can use for file-based storage. $•$  Shares represent mount points through which users or hosts can access file system resources. 

$•$  Each share is associated with a single file system and inherits the file system protocol (SMB or NFS) established for that file system. $•$  Shares of a multiprotocol file system can be either SMB or NFS. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 91  |

![Image](images/admguide_Image694.png)

Create a NAS Server (NFS and SMB) 

 Windows users or Linux/UNIX hosts 

Access to shares is determined depending on the type of file system: **Windows (SMB) shares**: Access is controlled by SMB share permissions and the ACLs on the shared directories and files. 

$•$  **Active Directory SMB servers**: Configure access for users and groups using Windows directory access controls. User/group authentication is performed through Active Directory. 

For Windows file systems, access to the share is based on share permissions and ACLs that assign privileges to objects defined in Active Directory. $•$  **Stand-alone SMB servers**: Manage a stand-alone SMB server within 

a workgroup from a Microsoft Windows host. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 92  |

![Image](images/admguide_Image696.png)

Create a NAS Server (NFS and SMB) **Linux/UNIX (NFS) exports**: Hosts access is defined by the NFS access control settings of the NFS export. 

$•$  For Linux/UNIX file systems, access is permitted based on NFS access settings. 

 

# PowerFlex 4.5.x File Service System Overview 

File management was introduced in PowerFlex appliance, rack, and software-only in the 4.0.x version. The File Service Node (FSN), or NAS node, acts like the SDC using a client to transfer files from a storage pool through a NAS server. In the PowerFlex 4.5.x release, the File Services 

System (Resource Group) functions have been changed to allow for easier file management and enhanced scalability. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 93  |

![Image](images/admguide_Image698.png)

Create a NAS Server (NFS and SMB) PowerFlex 4.5.x introduced the Global Namespace (GNS). In PowerFlex 4.5.x, managing collections of files becomes easier for administrators 

since all files fit neatly into a unified directory structure. 

 

# File Systems

File systems can be accessed through a wide range of protocols such as File Transfer Protocol (FTP and SFTP), SMB, and NFS. 

# File Controllers

The NAS Servers are created in the PowerFlex file controllers by the administrator before the PowerFlex File Services System is implemented. 

File controllers are dedicated physical nodes (FSNs) to host the NAS Servers. 

# NAS Servers

NAS Servers provide the file services to the client applications through the NAS container. The NAS Servers are responsible for namespaces, security policies, and serving file systems to the applications. NAS Servers can be discovered as a resource through PowerFlex Manager. 

      

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 94  |

![Image](images/admguide_Image700.png)

Create a NAS Server (NFS and SMB) 

# GNS

GNS simplifies the task of locating and accessing files. GNS also enforces a uniform naming convention of files and directories and provides a platform for multiple users to collaborate. With GNS, multiple NAS Servers or NAS systems can be created, all accessible through the global 

namespace. 

# Storage Pool Requirements

Beginning with version 4.5.x, a dedicated NAS Storage Pool is not required. File share and Block data share can be combined as one Storage Pool in a Protection Domain. 

# NAS Capabilities of a PowerFlex Cluster (System) 

Management of the PowerFlex file services can be performed by selecting the PowerFlex Manager File menu. 

 **1: ** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 95  |

![Image](images/admguide_Image702.pngimages/admguide_Path2.svg)

Create a NAS Server (NFS and SMB) 

 The NAS Servers page contains details about existing NAS servers and is 

used to: $•$  Create, modify, or remove a NAS Server. $•$  Configure existing NAS Server settings. 

$•$  Move a NAS Server to a different node. $•$  Swap NAS Servers between nodes. 

**2: ** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 96  |

![Image](images/admguide_Image704.png)

Create a NAS Server (NFS and SMB) 

 The File Systems page contains details about existing file systems and is 

used to: $•$  Create, modify, or remove a File System. $•$  Configure existing File System settings. 

$•$  Assign or Unassign a protection policy. $•$  Create or restore from a snapshot. 

$•$  Refresh quotas. **3: ** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 97  |

![Image](images/admguide_Image706.pngimages/admguide_Path1.svg)

Create a NAS Server (NFS and SMB) 

 The SMB Shares page contains details about existing SMB shares and is 

used to create, modify, or remove SMB shares. **4: ** 

 The NFS Exports page contains details about existing NFS exports and is 

used to create, modify, or remove NFS Exports. Host access is also managed from this page. 

**5: ** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 98  |

![Image](images/admguide_Image708.png)

![Image](images/admguide_Image709.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

Create a NAS Server (NFS and SMB) 

 The Global Namespaces page contains details about existing Global 

Namespace (GNS) deployments and is used to: $•$  Create, modify, or remove a GNS server $•$  Create, modify, or remove a link and target for a GNS server 

$•$  Restore a GNS server **6: **The File Protection page contains details about existing protection policies and snapshot rules. It is used to create, modify, or delete 

protection policies and snapshot rules. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 99  |

![Image](images/admguide_Image711.jpegimages/admguide_Path1.svg)

For more information, see the[ ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/managing-file-storage?guid=guid-27482e91-d7e4-457a-9ccd-faff1162955b&lang=en-us) and select Managing file storage. 

 

Create a NAS Server (NFS and SMB) 

 

owerFlex 4.5.x Administration Guide Dell P

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 100  |

![Image](images/admguide_Image714.png)

![Image](images/admguide_Image715.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path6.svg)

![Image](images/admguide_Path7.svg)

Create a NAS Server (NFS and SMB) 

# Create a NAS Server in PowerFlex Manager 

To create a NAS Server for the file node (controller), the administrator selects the File/NAS Servers option in PowerFlex Manager. In the NAS Servers page, the administrator then selects the Create NAS Server button. The first page in the Create NAS Server wizard is the Details page. 

The Details page is where a protection domain is selected and a NAS server name, description, and network information is entered. 

 **1: **When selecting the sharing protocol that the NAS server uses, the options are SMB, NFSv3, and NFSv4. If SMB and an NFS protocol are 

both selected, the NAS server is enabled to support multiprotocol. 

 **2: **With SMB selected, the wizard would look for the type of windows server. The two options are **Join to the Active Directory Domain** or 

**Standalone**. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 101  |

![Image](images/admguide_Image717.png)

![Image](images/admguide_Image718.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Create a NAS Server (NFS and SMB) 

 **3: **With the Unix Directory Services, the naming services can be 

configured with a combination of local files and NIS or LDAP. 

 If Secure NFS is enabled, the following requirements must be met: 

$•$  At least one NTP server must be configured (two NTP servers per domain is the recommendation). $•$  A Unix Directory Service (UDS) is configured. $•$  One or more DNS servers are configured. 

$•$  An Active Directory (AD) or custom realm must be added for Kerberos authentication. 

**4: **DNS information is mandatory when joining an AD domain or configuring Secure NFS. DNS is optional for a stand-alone NAS server. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 102  |

![Image](images/admguide_Image720.png)

![Image](images/admguide_Image721.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

Create a NAS Server (NFS and SMB) DNS can also be used to resolve hosts defined on NFS export access lists. 

 **5: **For user mapping a default account can be enabled for both a Windows 

and Linux user, or automatic user mapping can be selected. 

 For information on creating a NFS NAS Server, refer to the[ Dell ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/create-a-nas-server-for-nfs-linux-or-unix-only-file-systems?guid=guid-64ac5daf-c551-4fad-b22f-29fc1b14ab21&lang=en-us) [PowerFlex 4.5.x Administration Guide](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/create-a-nas-server-for-nfs-linux-or-unix-only-file-systems?guid=guid-64ac5daf-c551-4fad-b22f-29fc1b14ab21&lang=en-us). For information on creating a 

SMB NAS Server, refer to the[ Dell PowerFlex 4.5.x Administration Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/create-nas-server-for-smb-windows-only-file-systems?guid=guid-de5bd816-e251-488d-9c24-0f0233a775b1&lang=en-us) 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 103  |

![Image](images/admguide_Image726.png)

![Image](images/admguide_Image727.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path7.svg)

![Image](images/admguide_Path8.svg)

Export/Share Filesystems 

# Export/Share Filesystems 

## Configure the Settings of an Existing NAS Server 

For more information on a configuring NAS servers, see the[ Dell ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/overview-of-configuring-nas-servers?guid=guid-2a252a43-a706-4b14-a47d-ddaa09b09201&lang=en-us) [PowerFlex 4.5.x Administration Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/overview-of-configuring-nas-servers?guid=guid-2a252a43-a706-4b14-a47d-ddaa09b09201&lang=en-us) 

 Before administrators can provision file storage in PowerFlex, the NAS Servers must be running on the system. Administrators can modify or 

configure the settings for the newly created NAS Servers in the NAS Server View Details page in PowerFlex Manager. 

 NAS Server Settings Rules 

There are rules for modifying the NAS Server settings. $•$  Disabling multiprotocol file sharing for a NAS Server once a Global Namespace (GNS) file system is created on that NAS Server is not 

allowed. $•$  Disabling the NAS Server DNS for NAS Servers that support multiprotocol file sharing or NAS Servers that support SMB file sharing 

and that are joined to an Active Directory is not allowed. $•$  Enabling a UNIX directory service and DNS server to reconfigure a NAS Server that supports SMB-only or NFS-only file systems is 

required to support multiprotocol. $•$  Changing from an AD realm to a custom realm after the NAS Server is successfully created with Secure NFS will prohibit creation of any NFS 

exports until the following operations are performed: 

- PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 104  |

![Image](images/admguide_Image731.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path7.svg)

# 1.  Create a Keytab file. 

# 2.  Remove the AD realm from the

Export/Share Filesystems 

S Server.  NA

# 3.  Enter the username and password for the AD server. 

# 4.  Enter the custom realm. 

# 5.  Upload the Keytab. 

Creating a NAS Server Network Administrators create NAS Server networks in PowerFlex 4.x by adding production or backup file interfaces. When the additional interfaces are 

added, the administrator selects which interface is the preferred network to use for the NAS Server. PowerFlex assigns a preferred interface by default, but an administrator can set which interface to use first for production and backup, IPv4, and IPv6. 

To add the production and backup file interfaces to the NAS server and create routes to external services: 

# 1.  From the PowerFlex Manager UI, select the File menu. 

# 2.  From the options dropdown, select NAS Servers. 

# 3.  In the NAS Server page repository, select the NAS Server desired. 

# 4.  Select the Network option from the NAS Server View Details page. 

# 5.  Add more interfaces to the NAS Server as needed and then define the 

preferred interface to use with the NAS Server. 

Modifying NAS Server Naming Services In the NAS Server View Details page, administrators can configure the NAS Server naming services for the following: 

# DNS

$•$  Provide a DNS for Secure NFS. 

# UDS

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 105  |

Export/Share Filesystems $•$  UDS with NIS - a NIS domain name and the IP addresses is required for each of the NIS Servers. 

$•$  UDS with LDAP - LDAP must adhere to the IDMU, RFC2307, or RFC2307bis schemas. Administrator can also configure LDAP with SSL (LDAP Secure) and enforce the use of a Certificate Authority certificate for authentication. 

# Local Files

$•$  Local files can be used instead of, or in addition to DNS, LDAP, and NIS directory services. $•$  To use local files, configuration information must be provided through the files listed in PowerFlex Manager. 

$•$  To use local files for NFS, FTP access, the password file must include an encrypted password for the users. 

NAS Server Sharing Protocols NAS Servers are configured as NFS servers (Linux) or SMB servers (Windows). The setting for the connection File Transfer Protocol (FTP) or 

Secure FTP are configured after the NAS Server is created. 

# SMB servers

Configuration options for SMB servers: $•$  Select to Join to the Active Directory Domain if configuring SMB with Kerberos security. 

$•$  Make sure the DNS card is configured in the NAS server naming services if the Join to the Active Directory Domain is selected for SMB. $•$  When the Windows Server Type is set to Join to the Active Directory Domain, then Enable automatic mapping for unmapped Windows 

accounts/users must be selected in the User Mapping tab. 

# NFS

Configuration options for NFS servers: 

- PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 106  |

Export/Share Filesystems $•$  Select Enable extended Unix credentials when the NAS server uses the User ID (UID) to obtain the primary Group ID (GID) and all group 

GIDs to which it belongs. $•$  Clear Enable extended Unix credentials when the UNIX credential of the NFS request is directly extracted from the network information that 

is contained in the frame. $•$  In the Credential cache retention field, enter a time period (in minutes) for which access credentials are retained in the cache. The default 

value is 15 minutes. 

# FTP

FTP access can be authenticated using the same methods as NFS or SMB. Once authentication is complete, access is the same as SMB or NFS for security and permission purposes. 

# User Mapping

If you are configuring a NAS Server to support both types of protocols, SMB and NFS, you must configure the user mapping. When configured for both types of protocols, the user mapping requires that the NAS Server is joined with an AD domain. You can configure the SMB server with AD 

from the SMB Server card. 

If the Windows Server Type is set to Join to the Active Directory Domain, then the Enable automatic mapping for unmapped Windows accounts/users option must be selected. 

# Create a File System for NFS Export 

To create a File System for NFS Export, use the simulation below. *The web version* *of this content contains an interactive activity.* 

# Create a File System for SMB Share 

To create a File System for SMB Share, use the simulation below. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 107  |

Export/Share Filesystems *The web version* *of this content contains an interactive activity.* 

# PowerFlex File System Global Name Space 

 The Global Name Space (GNS) in PowerFlex Manager provides one IP address for the NAS File service. A single IP address means multiple NAS 

Servers and their file systems can be created and all are accessible through a single namespace. A single file node (controller) holds the information about all the NAS Servers and their file systems that are in separate file nodes. 

Any and all request from the NAS Servers are redirected to the single namespace host node. By adding the GNS function in PowerFlex for file storage, up to 8 times greater storage capacity is realized. The scalability in storage enables larger workloads and accommodates future growth for 

PowerFlex file storage. 

The operational improvement enables the addition of new NAS Servers and file systems to the GNS without having to map them to clients. 

 

      

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 108  |

![Image](images/admguide_Image736.png)

![Image](images/admguide_Image737.png)

Export/Share Filesystems 

# Create a Global Namespace in PowerFlex 4.5.x 

 Administrators use PowerFlex Manager to create a GNS to allow the NAS 

user to access a single namespace supported by the NAS cluster with a single export. 

Access the Global Namespaces page 

 In PowerFlex Manager, select the File menu. 

From the menu options, select Global Namespace. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 109  |

![Image](images/admguide_Image736.png)

![Image](images/admguide_Image739.pngimages/admguide_Path2.svg)

Create a Global Namespace 

In the Global Namespaces page, button. 

Select a Namespace Type 

Export/Share Filesystems 

 e + Create Global Namespace 

 

 click th

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 110  |

![Image](images/admguide_Image741.png)

![Image](images/admguide_Image742.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

Export/Share Filesystems In the Create Global Namespace wizard, select the radio button for the 

```
type of namespace wanted. The options are: 
```

$•$  Namespace for NFS $•$  Namespace for SMB 

$•$  Namespace for both NFS and SMB (this is the default selection) Click the Next button after the GNS type is selected. 

Select NAS Server 

 After selecting a Namespace type, choose a NAS Server on which to 

create the GNS. A NAS Server can host a NFS, SMB, or both namespaces. 

Click the **Next** button after selecting a NAS server. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 111  |

![Image](images/admguide_Image744.pngimages/admguide_Path1.svg)

Export/Share Filesystems Select the File System 

 In PowerFlex 4.5.x, the same storage pool used for Block can be used for File. Creating a new File System storage pool is recommended, however, 

there is an option to use an existing General Type File system. 

Provide a unique name for the GNS storage pool and optionally enter a description for the same. Provide a size for the GNS and then click the Next button to continue. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 112  |

![Image](images/admguide_Image746.pngimages/admguide_Path1.svg)

Export/Share Filesystems Namespace Details 

 Provide a unique name for the new GNS server and optionally provide a description for the same. Provide the Local Path in the entry field. 

Click the Next button to continue. 

Summary 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 113  |

![Image](images/admguide_Image748.png)

![Image](images/admguide_Image749.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

Export/Share Filesystems Verify the summary of the the GNS server. After verification, the administrator clicks the Create Namespace button. 

Result The namespace is now displayed in the Global Namespace page repository. From the Global Namespace page, the administrator can link 

file systems and NAS Servers to the new GNS server. 

 When the GNS is created, the root shares and exports are automatically created on the file system. Since this was a NFS-only GNS, the NFS 

exports is automatically created. 

 

**Deep Dive: **Review the[ Dell PowerFlex 4.5.x Administration ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/create-a-global-namespace?guid=guid-0e8bf6b2-5b35-4c03-8e09-2f76d13115a9&lang=en-us) [Guide ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/create-a-global-namespace?guid=guid-0e8bf6b2-5b35-4c03-8e09-2f76d13115a9&lang=en-us)to learn more about working with the PowerFlex GNS.  

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 114  |

![Image](images/admguide_Image753.png)

![Image](images/admguide_Image754.png)

![Image](images/admguide_Image755.png)

![Image](images/admguide_Image570.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path11.svg)

![Image](images/admguide_Path12.svg)

![Image](images/admguide_Path13.svg)

Protect NAS Servers 

# Protect NAS Servers 

## File System and NAS Server Protection Options 

PowerFlex uses Snapshots and NDMP to protect file system data. Depending on the recovery requirements, an organization can choose to use file system snapshots for a quicker recovery time objective (RTO), or an NDMP backup for long-term retention. NAS Servers can be provisioned 

with security to further protect the integrity of the data. 

File System Snapshots 

 Snapshots are point-in-time captures which save the state of the file system including all files and data within it. Snapshots can be used to 

restore the entire file system to a previous state. 

Up to 126 snapshots per file system can be taken. Manual snapshots that are created with PowerFlex Manager are retained for one week after creation (unless configured otherwise). 

The default access type for file snapshots is read-only. File Snapshot Access Types: 

- PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 115  |

![Image](images/admguide_Image757.pngimages/admguide_Path2.svg)

Protect NAS Servers $•$  **Protocol (read-only)**: Creates a read-only snapshot that can be mounted and accessed later through an NFS export or SMB share. 

$•$  **Snapshot**: Creates a read-only automounted snapshot accessible through the snapshot directory in the file system. $•$  **Protocol (read-write)**: Creates a read/write snapshot that can be mounted and accessed later through an NFS export or SMB share. 

For more information on file prote [Administration Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/file-protection?guid=guid-23311d3c-75a3-4c09-80b6-63c552de31c6&lang=en-us) n, see the[ Dell PowerFlex 4.5.x ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/file-protection?guid=guid-23311d3c-75a3-4c09-80b6-63c552de31c6&lang=en-us)ctio

NAS Sever Protection - NDMP 

 To access the NAS Server settings to provide protection, select PowerFlex Manager File > NAS Servers > [NAS Server Name View 

Details]>Protection. 

NDMP provides a standard for backing up file servers on a network. Once NDMP is enabled, a third-party Data Management Application (DMA), 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 116  |

![Image](images/admguide_Image761.jpegimages/admguide_Path1.svg)

Protect NAS Servers such as Dell NetWorker, can detect the PowerFlex NDMP using the NAS Server IP address. 

The NDMP Backup username is always **ndmp**. 

For more information on NAS Se [4.5.x Administration Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/nas-server-protection-and-events?guid=guid-bd17ef22-c500-4319-a7c1-d5ecb306d156&lang=en-us) rotection, see the[ Dell PowerFlex ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/nas-server-protection-and-events?guid=guid-bd17ef22-c500-4319-a7c1-d5ecb306d156&lang=en-us)rver p

For more information about using Dell NetWorker for the NDMP backups, see the[ NetWorker Administration Guide.](https://dl.dell.com/content/manual31623401-dell-emc-networker-19-8-administration-guide.pdf?language=en-us&ps=true) 

NAS Server Security 

 The security types for a NAS Server are: 

$•$  Kerberos - a distributed authentication service designed to provide strong authentication with secret-key cryptography. $•$  Common AntiVirus Agent (CAVA) - available for SMB servers, CAVA is an antivirus solution to clients using the NAS server. 

For more information on NAS Se [4.5.x Administration Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/nas-server-security?guid=guid-a9d327ea-5f1f-407b-b7db-0c808688a903&lang=en-us) ecurity, see the[ Dell PowerFlex ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/nas-server-security?guid=guid-a9d327ea-5f1f-407b-b7db-0c808688a903&lang=en-us)rver S

# Protect NAS Servers 

To create a protection policy, a snapshot rule, and enable NDMP, use the simulation below. 

*The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 117  |

![Image](images/admguide_Image769.pngimages/admguide_Path2.svg)

Configure NAS Quotas 

# Configure NAS Quotas 

## File System Quotas in a PowerFlex Cluster 

PowerFlex file systems include quota support. Quota is the method administrators use to ensure no individual user consumes all the available storage from a file system. Quotas are supported on SMB, NFS, FTP, NDMP, and multiprotocol file systems. 

PowerFlex file system supports user quotas, quota trees, and user quotas on tree quotas. All three types of quotas can co-exist on the same file system and can be used together to achieve fine-grained control over storage usage. If the limits from multiple quota types are reached 

simultaneously, whichever threshold was met first is considered the effective policy. 

Quotas are disabled by default. Administrators set quotas on a file system from the **File > File System > [selected file** **system] > Details > Quotas** tab. A description of the available quota types is listed below: 

| **TYPE** |  **DESCRIPTION**  |

User quotas  User quotas are set at a file system level to limit the amount of space a user may consume on a file system. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 118  |

Configure NAS Quotas 

Tree quota  Tree quotas limit the maximum size of a directory on a  file system. Unlike user quotas, which are applied and  tracked on a user-by-user basis, quota trees are applied 

 to directories within the file system. Tree quotas can be  applied to new or existing directories. 

 Administrators can use tree quotas to:     $•$  Set storage limits on a project basis. For example,      tree quotas can be established for a project directory that has multiple users sharing and creating files in it. 

$•$  Track directory usage by setting the tree quota hard and soft limits to zero. 

User quota on  User quotas on tree quotas occur when a tree quota is a quota tree  created, and the administrator creates additional user quotas within the directory and enforces user quotas. User quotas on a quota tree limit the amount of storage 

that an individual user consumes to store data on the quota tree. 

 

# Quota Limits for a File System 

A soft limit can be passed temporarily until the grace period expires. A hard limit is an absolute limit on storage usage. If both the hard and soft limits are set to zero, this effectively disables the quotas. 

The administrators can enable or disable quotas at any time. It is recommended that they are enabled or disabled during nonpeak production hours to avoid impacting file system operations. The administrator can set grace periods on a quota that determines how long 

the soft limit can be exceeded. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 119  |

Quota reaching soft limit 

In this scenario, the user is writing on it. The file system usage is in limit. 

Configure NAS Quotas 

 he data to a directory with a tree quota sing and climbing towards the soft  t crea

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 120  |

![Image](images/admguide_Image773.png)

Configure NAS Quotas Quota exceeding soft limit 

 Once the soft limit is reached, the user is notified but will still be allowed to 

use space until the grace period has elapsed. $•$  The grace period counts down the time between the soft and hard limits. It alerts the user about the time remaining before the hard limit is 

met. $•$  The quota grace period is a specific amount of time for each tree quota on a file system. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 121  |

![Image](images/admguide_Image776.png)

Configure NAS Quotas Quota grace period expires 

 Once the grace period expires the user cannot write to the file system until 

more space has been added, even if the hard limit has not been met. 

Administrators can set an expiration date for the grace period. The default is seven days. Alternatively, the expiration can be set to an infinite amount of time so the grace period does not expire. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 122  |

![Image](images/admguide_Image779.png)

Configure NAS Quotas Quota reaching hard limit 

 Finally, when the hard limit is reached for a quota tree, no user can write 

data to the tree or the file system until more space becomes available. 

**Important**: Administrators cannot create quotas for read-only file systems.  

# Configure NAS Quotas for File Systems in PowerFlex 

# Manager 

To learn how to configure file system NAS quotas in PowerFlex Manager, use the simulation below. 

*The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 123  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Image782.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Enter and Exit Maintenance Modes 

# Enter and Exit Maintenance Modes 

## Introduction** **to** **Maintenance** **Modes

Maintenance Mode is a feature that is used to safely streamline system operation when maintenance or a planned restart of an SDS is required. 

Maintenance Mode takes a node offline to repair, replace, upgrade hardware components, or due to rolling PowerFlex upgrades. All while allowing the data residing on that node to remain available. 

There are some general criteria that must be met when using maintenance mode: 

 $•$  Two nodes from the same Protection Domain cannot be put into maintenance mode simultaneously. 

$•$  Protected and Instant Maintenance Modes (IMM) cannot be simultaneously active within a single protection domain. $•$  If fault sets are in use for the containing Protection Domain, all SDSs concurrently in Protected Maintenance Mode (PMM) must belong to 

the same fault set. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 124  |

![Image](images/admguide_Image785.pngimages/admguide_Path1.svg)

Enter and Exit Maintenance Modes 

# Types of Maintenance Modes 

PowerFlex provides two options for putting a node into maintenance mode. 

Instant Maintenance Mode 

# Video:

*The web version* *of this content contains a* *Video.* In IMM, an SDS node is immediately and temporarily removed from active participation without building a new copy of the data on other nodes. 

Changes are tracked and resynched when the node is available again. 

Advantages of IMM are the speed at resync of changes when exiting. h the mode is entered, and quick whic

A disadvantage of IMM is single copy exposure, which could lead to data unavailability or data loss. 

For more information, see the[ Instant Maintenance Mode ](https://www.dell.com/support/manuals/en-us/scaleio/flex-software-to-45x/instant-maintenance-mode-imm?guid=guid-9dfdec3a-7f7c-4b5e-a4a5-72c7b98eaf71&lang=en-us)section of the Dell PowerFlex 4.5.x Technical Overview. 

Protected Maintenance Mode 

# Video:

*The web version* *of this content contains a* *Video.* PMM creates a third copy of the data before entering maintenance mode. Data is mirrored on two nodes. Changes are tracked and resynched when 

the node is available again. 

Advantages: $•$  Two copies of the data are always available. $•$  Many-to-one rebuilds are avoided. 

$•$  Performance is better when exiting PMM than IMM. 

Disadvantages: 

- PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 125  |

Enter and Exit Maintenance Modes $•$  PMM requires more spare capacity than IMM. $•$  It takes longer to enter PMM than IMM. 

For more information, see the[ Protected](https://www.dell.com/support/manuals/en-us/scaleio/flex-software-to-45x/protected-maintenance-mode-pmm?guid=guid-2b71174b-9413-4b0e-944d-994e0ade92b7&lang=en-us) PowerFlex 4.5.x Technical Overview. [intenance Mode ](https://www.dell.com/support/manuals/en-us/scaleio/flex-software-to-45x/protected-maintenance-mode-pmm?guid=guid-2b71174b-9413-4b0e-944d-994e0ade92b7&lang=en-us)section of the [ Ma](https://www.dell.com/support/manuals/en-us/scaleio/flex-software-to-45x/protected-maintenance-mode-pmm?guid=guid-2b71174b-9413-4b0e-944d-994e0ade92b7&lang=en-us)

# Abort Protected Maintenance Mode 

Protected Maintenance Mode can be aborted both manually by the user or automatically by the system. 

Manual Abort PMM Users can manually abort PMM before entering PMM is complete, by selecting the **Abort entering Protected Maintenance Mode** option in the 

**More Actions** menu. The extra data copies are cleaned up, and the SDS returns to its normal state. 

 Auto Abort PMM 

If there is a system failure while the node is entering PMM the system aborts entering PMM automatically. If there is not enough spare capacity 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 126  |

![Image](images/admguide_Image790.jpegimages/admguide_Path2.svg)

Enter and Exit Maintenance Modes available to copy data from the target SDS node, the system aborts entering PMM automatically also. Auto-Abort PMM frees up the system 

spare capacity for higher priority tasks. 

When the system Auto-Aborts PMM, an alert is generated, and the nodes state is changed to **Maintenance Aborted by System**. This state is changed only after the system finishes aborting the Enter PMM phase. Users can exit the **Maintenance Aborted by** **System** state by starting a 

new PMM, IMM, or selecting the **Exit Maintenance Mode** option in the **More Actions** menu. 

 

# Put an SDS into Maintenance Mode 

To try putting an SDS into maintenance mode, use the simulation below. *The web version* *of this content contains an interactive activity.* 

# Put a Node into Service Mode 

PowerFlex Manager enables you to put a node in service mode when you must perform maintenance operations on the node. When you put a node in service mode, you can specify whether you are performing short-term maintenance or long-term maintenance work. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 127  |

![Image](images/admguide_Image792.jpegimages/admguide_Path3.svg)

Perform a Node Expansion To try using maintenance modes, use the simulation below. 

*The web version* *of this content contains an interactive activity.*  

# Perform a Node Expansion 

## Expansion Process Overview 

Implementation engineers perform the tasks needed to complete the PowerFlex rack expansion. Each task is completed in a four-phase process. 

 

## PowerFlex 4.x rack Expansion Types 

Implementation engineers are responsible for deploying the different types of PowerFlex rack expansions. There are four types that the solution architect (SA) or the customer design to expand a PowerFlex 4.x rack system. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 128  |

![Image](images/admguide_Image794.png)

Perform a Node Expansion Single Node Expansion 

 In a single node expansion, the implementation engineer deploys a single node to an existing PowerFlex rack cluster. The node may be an 

additional storage node to expand the storage capacity of the cluster. Or the node may be required to increase the compute capabilities of the cluster. The storage-only node is an SDS node while the compute-only node is an SDC node. On occasion, the node that is added to the cluster 

acts as a combination of storage and compute. The storage and compute single node is known as an HCI node. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 129  |

![Image](images/admguide_Image797.pngimages/admguide_Path1.svg)

Perform a Node Expansion Node-Only Expansion 

 In a Node-only expansion, implementation engineers are required to expand the node capacity of a PowerFlex rack cluster. The requirement is 

to expand the PowerFlex rack cluster to include more nodes for compute, storage, or both. The maximum number of nodes that can be added per rack cabinet are 24 units. The new nodes connect to the existing PowerFlex rack system switches. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 130  |

![Image](images/admguide_Image799.png)

Perform a Node Expansion Switch Expansion 

 A switch expansion is when the implementation engineer is required to add more ToR switches, such as management and access switches. 

Adding more switches increases the performance capabilities of a PowerFlex rack system that expanded the node capacity and uses multicabinets. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 131  |

![Image](images/admguide_Image801.png)

Perform a Node Expansion Cabinet Expansion 

 Adding cabinets to an existing PowerFlex rack system scales the capacity of storage and compute nodes. Implementation engineers ensure that the 

aggregation switches are placed appropriately in the new cabinet setup to provide resiliency for the system. 

# Validate System Health 

Before beginning the expansion, ensure that the components in the system are in a healthy state. Ensure that any critical issues are remediated. Resolving critical issues before an expansion avoids most deployment failures. If there are issues, reach out to support as needed. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 132  |

![Image](images/admguide_Image803.png)

Perform a Node Expansion PowerFlex Manager Health In the PowerFlex Manager Dashboard, verify the health of the Services 

and Resources, and the overall health of PowerFlex Manager. 

 The Resources in PowerFlex Manager must be in Managed mode to 

perform the expansion. Check the Resources page to verify the management state of the resources. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 133  |

![Image](images/admguide_Image806.png)

![Image](images/admguide_Image807.jpegimages/admguide_Path1.svg)

Perform a Node Expansion Element Managers Health All element managers, including vCenter, iDRAC, and CloudLink, must be 

in a healthy state for the expansion. Log in to each element manager and verify the health of the environment before beginning the expansion. 

  Hardware Components Health 

PowerFlex Manager and iDRAC are used to monitor the hardware components which include nodes and switches. 

Log in to all switches and nodes, and verify they are healthy. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 134  |

![Image](images/admguide_Image809.png)

![Image](images/admguide_Image810.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

Perform a Node Expansion 

# Verify Cluster Configuration Details 

To learn how to verify cluster configuration details, select each tab below. 

Logical Network Configuration 

 Before beginning the expansion, verify the details in the LCS for the 

PowerFlex rack system and what type of expansion is needed. Check the details to verify if the network configuration is a single network (one VLAN) or multinetwork (multiple VLANs). The logical network configuration (types of switches) for the expansion system must match the existing network 

configuration. Ensure nothing in the existing configuration changes during the expansion other than the addition of new nodes and switches. 

To check the LCS details on the expansion:   Open the existing LCS for the expansion.   Expand the **Dell Converged System Information** category of the 

LCS.   Scroll through the content to verify the details match the configuration details of the PowerFlex rack expansion.   Check the categories for **Storage Layout Information**, **System VLAN** 

# Information**, **Fast Track IP Address Information**, **Device

# Information**, **Cluster Configuration**, and **Core Environment

**Information**. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 135  |

![Image](images/admguide_Image812.jpegimages/admguide_Path2.svg)

Perform a Node Expansion Licenses 

 Verify in the SCR and in the individual component management interface 

that there are enough available licenses for the expansion. Ensure that there are licenses for the following: 

  ESXI licenses   Licenses for each installed node 

  Storage and capacity licenses   Licenses for each switch port, particularly for Dell switches 

  Cloudlink (if applicable) licenses Number of Ports on Node 

For the expansion projects, the customer requests more storage, compute, or hyperconverged nodes added to the existing PowerFlex rack. 

The typical node ports for the expansion nodes are four 25GB ports in two NIC slots. 

However, there are some nodes that come out of the factory with two additional 25 GB or 100 GB ports, or one extra NIC. Nodes with the extra ports are nodes that are tailored for a network configuration, node type and CPU, or management. Check the number of ports and speed of ports 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 136  |

![Image](images/admguide_Image815.pngimages/admguide_Path1.svg)

Perform a Node Expansion on the nodes before performing the expansion. The nodes should always match what is designated in the LCS. If a field engineer is installing a six-

port node, any network configurations that are implemented during the expansion should match the existing node network configuration. 

 Download the Dell PowerFlex Rack with PowerFlex 4.x Cabling and Connectivity Guide from[ SolVe Online.](https://solve.dell.com/solve/home/PowerFlex%20Family) Go to the PowerFlex Family > 

PowerFlex rack > Install > Build Procedures > PowerFlex appliance and rack with PowerFlex 4.x Cabling and Connectivity Guide section to download the guide. 

**Important**: Verify that all required licenses are available before going on site. Not having the required licenses available delays the expansion. Licenses must be validated as being on the correct site ID of the original build, or the 

 system cannot be expanded.  

# Cabling Pre-Requisites 

Before beginning the switch expansion, complete the following: Check cabling 

  Double and triple check that the cabling is done correctly and is following best practices. Errors in cabling can cause long delays in completing a successful expansion. 

  Ensure you have all the cables and that they are the correct cables. 

Verify system health and check alerts 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 137  |

![Image](images/admguide_Image818.png)

![Image](images/admguide_Image602.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Perform a Node Expansion   Ensure all components, and element managers are healthy and at the correct version for a successful expansion. 

  If there are any issues with hardware components or element managers, resolve issues before starting the expansion. Failure to resolve any issues before beginning the work delays the expansion. 

Verify that enough IP Addresses are available   If needed, expand the IP address range on the existing networks so that enough IPs are available for the new switches. 

  The new switches must have the same VLANs as existing switches and be able to communicate across all the switches. 

# PowerFlex rack Node Port Locations 

When cabling the PowerFlex rack nodes, be sure to use the proper NICs provided and ensure that the cables are plugged into the appropriate ports. 

To learn about the NIC ports for each node type, select the red outlined NIC ports below. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 138  |

![Image](images/admguide_Image821.jpeg)

**1: ** 

# Slot:Port**  **Description

| **00:01** |  **Trunk** **1 to switch A ** |

**(management, **

# data1, data3, rep1)

 **2: **NIC Y, Port 1 - Management 

# Slot:Port**  **Description

Perform a Node Expansion 

| **VMW** |  **dvSwitch**  |

# Label

# vmnic2**  **cust_dvswitch

# VMW Label**  **dvSwitch

| **00:02** |  **Trunk** **2 to switch  vmnic3**  **flex_dvswitch**  |

**B (data2, data4, **

# rep2)

 **3: ** 

| **Slot:Port** |  **Description** |  **VMW Label**  **dvSwitch**  |
| **03:02** |  **Trunk** **2 to switch  vmnic5**  **flex_dswitch**  |

**A (data2, data4, **

# rep2)

 **4: ** 

| **Slot:Port**  **Description** |  **VMW**  |

# dvSwitch

# Label

| **03:01** |  **Trunk** **1 to switch B  vmnic4**  **cust_dvswitch**  |

**(management, **

# data1, data3, rep1)

 **5: ** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 139  |

Perform a Node Expansion 

| **Slot:Port**  **Description** |  **VMW**  |

# dvSwitch

# Label

| **06:01** |  **Trunk** **1 to switch  vmnic4**  **cust_dvswitch**  |

**B**  

**6: ** 

| **Slot:Port**  **Description** |  **VMW Label**  **dvSwitch**  |
| **06:02** |  **Trunk** **2 to switch  vmnic3**  **flex_dvswitch**  |

# B (data)

 **7: ** 

| **Slot:Port**  **Description** |  **VMW Label**  **dvSwitch**  |
| **00:02** |  **Trunk** **2 to switch  vmnic1**  **flex_dvswitch**  |

# A (data)

 **8: ** 

| **Slot:Port**  **Description** |  **VMW**  |

# dvSwitch

# Label

| **00:01** |  **Trunk** **1 to switch  vmnic2**  **cust_dvswitch**  |

**A**  

**9: ** 

| **Slot:Port**  **Description** |  **VMW**  |

# dvSwitch

# Label

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 140  |

Perform a Node Expansion 

| **00:01** |  **Trunk** **1 to switch A  vmnic0**  **cust_dvswitch**  |

**(management, **

# data1, data3, rep1)

 **10: ** 

| **Slot:Port**  **Description** |  **VMW Label**  **dvSwitch**  |
| **00:02** |  **Trunk** **2 to switch  vmnic1**  **flex_dvswitch**  |

**B (data2, data4, **

# rep2)

 **11: ** 

| **Slot:Port**  **Description** |  **VMW**  |

# dvSwitch

# Label

| **04:01** |  **Trunk** **1 to switch B  vmnic2**  **cust_dvswitch**  |

**(management, **

# data1, data3, rep1)

 **12: ** 

| **Slot:Port**  **Description** |  **VMW Label**  **dvSwitch**  |
| **04:02** |  **Trunk** **2 to switch  vmnic3**  **flex_dvswitch**  |

**A (data2, data4, **

# rep2)

 **13: **iDRAC Management Port 

**14: **iDRAC Management Port **15: **iDRAC Management Port 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 141  |

Perform a Node Expansion 

**Deep Dive:**  Download the Dell PowerFlex Rack with PowerFlex 4.x Cabling and Connectivity Guide from[ SolVe ](https://solve.dell.com/solve/home/PowerFlex%20Family) [Online](https://solve.dell.com/solve/home/PowerFlex%20Family). Go to the PowerFlex Family > PowerFlex rack > Install > Build Procedures > PowerFlex appliance and rack 

  with PowerFlex 4.x Cabling and Connectivity Guide section to download the guide. 

 

# Switch Cabling 

Pre-Expansion The following is an Access/Aggregation networking model. To keep the graphic readable, the image shows connections of a node only expansion 

to the switches.   The two PowerFlex data network connections for each node are on different NICs and are connected to each access switch. 

  The management network switch has a connection for each node from a different NIC and is connected to each access switch.   Each access switch is connected to both aggregate switches with connections from two different ports. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 142  |

![Image](images/admguide_Image570.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Perform a Node Expansion 

 Post-Expansion 

The cabling diagram represents a multi-cabinet PowerFlex rack expansion with switch cabling. 

  PowerFlex data network connections for each node are on different NICs and are connected to each access switch by two different ports.   A third connection is installed between the nodes and the management switch. 

  The access switches are connected to both aggregation switches across the cabinets, with cables from two different ports. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 143  |

![Image](images/admguide_Image829.jpeg)

Perform a Node Expansion 

 

# Verify Cabling 

Once the cabling of the rack is completed, it is a best practice to verify that the cabling has been done correctly. The most reliable method is to verify the cabling through iDRAC, because iDRAC reads the MAC address directly from the NIC. 

*Video:* *The web version* *of this content contains a video.* 

A video demonstration on how to verify the cabling of the node ports using iDRAC. 

**Important**: The iDRAC verification method only works if LLDP is enabled on the switches. Ensure to enable LLDP when you configure the switches.  

 PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 144  |

![Image](images/admguide_Image831.jpeg)

![Image](images/admguide_Image602.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path4.svg)

Perform a Node Expansion 

# Perform Expansion Using PowerFlex Manager 

Review the steps that are involved in the Discovering Resources section from the[ Dell PowerFlex 4.5.x Administration Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/discover-a-resource?guid=guid-50b5731c-0904-40c9-b3f2-e509e7422862&lang=en-us) 

 Before you start discovering a resource, gather the IP addresses and credentials that are associated with the resources. Ensure that there is 

network connectivity. PowerFlex Manager automates the node and switch hardware expansion and configuration. 

To learn about discovering, adding, and configuring resources, select each tab below. 

Adding Expansion Node Resources PowerFlex Manager provides the process to discover the node resources and add the node resources to an existing resource group. 

To learn about discovering and adding the expansion node hardware using PowerFlex Manager, watch the video below. Click the download icon on the navigation bar to download the video transcript. 

# Video:

*The web version* *of this content contains a video.*  

A video demonstration on PowerFlex rack node discovery and expansion using PowerFlex Manager. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 145  |

![Image](images/admguide_Path2.svg)

Perform a Node Expansion 

**Important**: It is a pre-requisite that any upgrade to the system is performed prior to adding an expansion node. 

Upgrades to the system may be installation of drivers for Broadcom cards, CPLD drives, and other system requirements. The remediation tasks performed by the customer ensures there are no issues during the expansion 

 activity. View the course page **Identify Existing Product** Versions to learn how to identify elements of the PowerFlex system that require attention. 

 Configure Expansion Switches 

 Once the node components are cabled to the network, expansion engineers must configure the switches. At times, new cabinets come with 

new access switches, but sometimes existing access and aggregation switches need configuration in PowerFlex Manager. New switches must be discovered as a resource and ensure that they are compliant to the PowerFlex 4.x RCM/IC. 

Download the Dell PowerFlex Rack Field Logical Build Guide from[ SolVe ](https://solve.dell.com/solve/home/PowerFlex%20Family) [Online ](https://solve.dell.com/solve/home/PowerFlex%20Family)to learn how to configure the new switches. Chapter 4 of the guide covers the process to configure Cisco Nexus and Dell Networking switches. Go to the SolVe Online PowerFlex Family > PowerFlex rack > 

Miscellaneous > Expansion procedures section to download the guide. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 146  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Image838.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path17.svg)

Perform a Node Expansion Ensure that licenses for all the required virtual networking ports are available before you arrive on site for the hardware expansion. 

Upgrade the Resource Group The newly added expansion resources must meet the current PowerFlex rack RCM/IC. Ideally, code upgrades are manually completed before 

adding the expanded components to the cluster. If after adding new components, there are code inconsistencies found, the resource group must be upgraded. To complete the upgrade of the new resources to 4.x RCM/IC compliance, the engineer goes to the Resource Group in 

PowerFlex Manager. The engineer proceeds with an upgrade of the Resource Group. 

 

 1.  The engineer selects the **Lifecycle** menu in PowerFlex Manager. 

 2.  Next, the engineer selects the **Resource** **Group** option from the menu. 

 3.  In the **Resource** **Groups** page, the engineer selects the **Resource** 

**Group**. 

 4.  The engineer selects the **Update Resources** to initiate the upgrade of 

all the noncompliant resources. 

Typical estimates to complete the Resource Group firmware upgrade is 1 to 2 hours per node in the cluster. 

      

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 147  |

![Image](images/admguide_Image840.pngimages/admguide_Path1.svg)

Perform a Node Expansion 

# Performing Expansion Using CSV File 

The process of expansion using a CSV file is actually very easy. Before expanding, ensure the prerequisites have been met, just as with a PowerFlex Manager deployed expansion. Open the deployment file used for initial cluster configuration and add lines according to the number of 

nodes being added. Once the deployment file is updated, rerun the PowerFlex Installer wizard. The wizard will recognize the existing nodes and only configure the new (expansion) nodes. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 148  |

![Image](images/admguide_Image842.pngimages/admguide_Path2.svg)

Manage Users and Passwords in a PowerFlex Cluster 

# Manage Users and Passwords in a PowerFlex 

# Cluster 

## Change PowerFlex Cluster System Passwords 

IT administrators use PowerFle manage the PowerFlex cluster. ager and other applications to x Man

The external passwords that are configured with PowerFlex must be kept in sync with those configured on the components themselves. The components include iDRAC, Server BIOS, vCenter (and ESXi), and host passwords. 

When the administrator first logs into PowerFlex Manager, they must set their password. Administrators can also change their password at any time after the first login. 

 Steps to update passwords external to PowerFlex: 

- $•$  iDRAC password 
- $•$  Server BIOS password 
- $•$  vCenter and ESXi password 
- $•$  Embedded Linux management password 

 Steps to change the PowerFlex Manager login password are: 

# 1.  Click the **user icon** in the upper right corner of PowerFlex Manager. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 149  |

![Image](images/admguide_Image848.jpeg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path10.svg)

Manage Users and Passwords in a PowerFlex Cluster 

# 2.  Click **Change password**. 

# 3.  Type the password in the **New Password **field. 

# 4.  Type the password again in the **Verify Password** field. 

# 5.  Click **Apply.** 

# Manage User Roles 

PowerFlex Manager allows an administrator to set User Roles for the PowerFlex 4.x.x cluster. The different type of User Roles dictates what the specific user can do while managing the PowerFlex cluster. The administrator provides the User Roles in the PowerFlex Manager User 

Management page. The User Management page allows the administrator to manage: 

$•$  **Local Users** - create, modify, delete, or reset the password for a local user. Review the User Roles available in PowerFlex Manager. $•$  **LDAP** **Users** **and Groups** - Add LDAP users or groups or modify LDAP users or groups. 

$•$  **Directory Services** - Add, modify, and remove a directory service. The PowerFlex Manager accesses the[ directory services ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/directory-services?guid=guid-59435ddd-4b6d-47d5-84da-d82d5d055225&lang=en-us)to authenticate users. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 150  |

Manage Users and Passwords in a PowerFlex Cluster 

 To access the User Management page in PowerFlex Manager, click the **Settings > User Management** section and then select the option that is 

wanted. Use the Dell PowerFlex Manager 4.x.x application Online Help to review how to accomplish adding users as an administrator. 

# Change Passwords

# PowerFlex Manager 

# ources Discovered in  for Res

To try changing passwords for discovered resources in PowerFlex Manager, use the simulation below. 

*The web version* *of this content contains an interactive activity.* 

# Manage the Secadmin User in the CloudLink Center 

# and PowerFlex 

The security administrator (**secadmin**) is an integrated role that CloudLink provides. The **secadmin** has full access to all CloudLink Center functionality, including managing user accounts, configuring key stores, and accessing event logs. If the CloudLink **secadmin** user password is 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 151  |

![Image](images/admguide_Image856.pngimages/admguide_Path5.svg)

Manage Users and Passwords in a PowerFlex Cluster changed after deployment, administrators must also change the CloudLink **secadmin** user password in PowerFlex Manager. 

 The steps to change the password of the CloudLink secadmin user using the CloudLink Center and PowerFlex Manager are: 

# 1.  Open a web browser and log in to the CloudLink VM. 

# 2.  Log in with **secadmin** username. 

# 3.  On the upper right corner, click **secadmin**, and click **Change **

**Password**. 

# 4.  On the **CHANGE PASSWORD** screen, type the **Current Password** 

and **New Password** into the respective fields. Click  **Change **to complete the password change. 

# 5.  On the upper right corner click **secadmin,** and then select **Logout.** 

# 6.  Log in with the **secadmin** username and the new password. 

# 7.  Change the CloudLink password in PowerFlex Manager. 

i.  In PowerFlex Manager, go to **Settings > Security > Resource** **Credentials > Credentials Management**, select the CloudLink credential, click **Edit,** change the **Password,** and click **Save**. 

ii.  Test the changes, In the PowerFlex Manager GUI, go to the 

# Resources** page, select the CloudLink center, and click **Run

**Inventory.** iii.  To confirm that the proce **Settings> Logs**. mpletes with no errors, check ss co

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 152  |

![Image](images/admguide_Image859.pngimages/admguide_Path1.svg)

Protect Volumes using SnapShots 

# Protect Volumes using SnapShots 

## PowerFlex Snapshot 

Snapshots are block images in the form of a storage volume or LUN. Snapshots are used to instantaneously capture the state of a volume at a specific point in time. Once a snapshot is taken, it becomes a new unmapped volume in the system. It can be manipulated to be in a state 

such as mapped, unmapped, renamed, and resized like other volumes in the cluster. Each volume in a PowerFlex cluster can have up to 128 snapshots. 

 **1: **A snapshot is an instantaneous, point-in-time copy of a volume. The snapshot records changes to the volume after the snapshot is taken. 

Snapshots can also be made read/write and modifications to the snapshot are recorded. Snapshots give administrators the ability to revert changes should the original copy be needed. 

**2: **PowerFlex snapshots can also have relationships across different vTrees. When you take a snapshot, multiple volumes can be selected simultaneously. All snapshots taken together this way form a Consistency Group. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 153  |

![Image](images/admguide_Image861.png)

# Difference Between Regul

# Snapshot 

Regular Snapshot PowerFlex enables you to create snapshots and some of 

the key features of regular snapshots include: 

$•$  They are thinly provisioned and writable, regardless of the original volumes. 

$•$  When a volume or snapshot is created, snapshot capacity and volume capacity can be viewed in 

the Capacity view of the dashboard. 

Protect Volumes using SnapShots 

# apshot and Secure ar Sn

$•$  PowerFlex snapshots can be initiated manually through      any of the clients including PowerFlex Manager UI, CLI, or REST API. 

$•$  Snapshots and their source volume are organized into a V-Tree. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 154  |

![Image](images/admguide_Image864.png)

Secure Snapshot Secure snapshots support compliance applications and 

can be configured with predefined retention periods (or policies). Snapshots with active retention cannot be 

deleted. 

Some of the key features of Secure snapshots include: 

$•$  Secure snapshots provide two key properties, secured flag and expiration time. 

$•$  Secure snapshots are read-only. $•$  Secure snapshots cannot be overwritten. 

$•$  Dell Technical Support remove strict dual-signature policy. 

Protect Volumes using SnapShots 

 

cure snapshots while following a s se $•$  Pausing, altering, or deleting the policy do not delete the snapshots that are marked as secure. 

# Create, Modify, and Delete Snapshots 

For more information, see the[ Dell Po](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/snapshots?guid=guid-a744b69f-834a-4206-a89f-02f30f5c216c&lang=en-us) and select the Snapshots section. ex 4.5.x Administration Guide werFl

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 155  |

![Image](images/admguide_Image871.png)

![Image](images/admguide_Path2.svg)

Protect Volumes using SnapShots Create a volume Snapshot with PowerFlex Manager 

 Using PowerFlex Manager, you can create instantaneous snapshots of 

one or more volumes. From the menu bar, select the relevant volume and create snapshots by configuring the given parameters. 

The **Use secure** **snapshots** option prohibits deletion of the snapshots until the defined expiration period has elapsed. 

Set Bandwidth and IOPS limits for a Snapshot 

 Setting bandwidth and IOPS limits for snapshots lets you control the 

quality of service. Bandwidth and IOPS limits are set on a per host basis. 

Ensure that the snapshots are mapped before you set the limits. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 156  |

![Image](images/admguide_Image874.jpeg)

![Image](images/admguide_Image875.jpeg)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

Lock and Unlock a Snapshot Protect Volumes using SnapShots 

 A snapshot lock prevents a snapshot from being deleted. You can lock 

auto snapshots (snapshots that are created through a snapshot policy automatically) through the auto-removal process to avoid deletion. You can unlock the snapshots later so that they are automatically removed. **Note:** If a snapshot policy is displayed for the snapshot in the Snapshot 

Policy column of the snapshots list, it is an auto snapshot. 

Delete a Snapshot 

 Ensure that the snapshot that you are removing is not mapped to any 

hosts. If the snapshot is mapped, unmap it before removing it. In addition, ensure that the snapshot is not the source volume of any snapshot policy. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 157  |

![Image](images/admguide_Image877.jpeg)

![Image](images/admguide_Image878.jpeg)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

Protect Volumes using SnapShots You must remove the volume from the snapshot policy before you can remove the snapshot. 

**Note:** Removing or deleting a snapshot erases all the data in the corresponding snapshot. 

Migrate vTree Snapshot 

 Snapshots and their source volume are organized into a Volume Tree or vTree. A vTree includes the root volume and all the descendant snapshots 

resulting from that volume. You can migrate vTree for a snapshot to a different storage pool. Volumes undergoing migration remain available for I/O. 

**Note:** vTree migration is a long proce depending on the size of the vTree. nd can take days or weeks, ss a

# Create a Volume Snapshot 

To try creating snapshots and setting IOPS limits for the snapshot, use the simulation below. 

*The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 158  |

![Image](images/admguide_Image880.jpegimages/admguide_Path2.svg)

Protect Volumes using SnapShots 

# Managing Volume Snapshots 

To try locking, unlocking, mapping, and unmapping snapshots, use the simulation below. 

*The web version* *of this content contains an interactive activity.* 

# Snapshot Policies 

To review detailed steps for creating snapshot policies using SCLI, see the[ Dell Technologies Infohub Site.](https://infohub.delltechnologies.com/l/dell-powerflex-snapshots/snapshot-policies-19) For more information, see the Snapshot Policies section of the[ Dell ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/snapshot-policies?guid=guid-67fe6287-552e-43e9-8e7a-f3d711a44503&lang=en-us) [PowerFlex 4.5.x Administration Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/snapshot-policies?guid=guid-67fe6287-552e-43e9-8e7a-f3d711a44503&lang=en-us) 

 Snapshot Policies contain a few attributes and elements, which allows the system to automatically run snapshots for specified volumes based on 

specified retention schedules. 

Steps to create a Snapshot Policy with PowerFlex Manager: 

# 1.  On the menu bar, click **Protection > Snapshot Policies**. 

# 2.  On the Snapshot Policies page, click **Create Snapshot Policy > Add **

**Snapshot Policy parameters > Create and Activate**. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 159  |

![Image](images/admguide_Image885.png)

![Image](images/admguide_Path3.svg)

Protect Volumes using SnapShots 

# Managing Snapshot Policies 

To try creating, activating, assigning, and modifying a Snapshot Policy, use the simulation below. 

*The web version* *of this content contains an interactive activity.* 

# Working With Snapshot Policies 

To try pausing, unassigning, and deleting Snapshot Policies, use the simulation below. 

*The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 160  |

Replicate Volumes Between Clusters 

# Replicate Volumes Between Clusters 

## Native Asynchronous Replication 

Native asynchronous replication enables disaster recovery and data protection of database or application data. It can also be used for data migration or distribution of workloads to secondary environments. 

With native asynchronous replication, volumes are replicated from one cluster to another cluster. The cluster can be a replication Source, Target, or both. Asynchronous replication defines a point in time, and ensures that all writes carried out before that point are copied to the destination. 

The main components of native asynchronous replication are: 

- $•$  [Storage Data Replicator ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/sdrs?guid=guid-f0dbae54-b53a-484e-8af7-911fdda5d449&lang=en-us)
- $•$  [**Replication Consistency Group (RCG)** ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/replication-consistency-group?guid=guid-71d462f6-805e-48e1-a2ae-7b0f77d91d7f&lang=en-us)

 

## Summary of Remote** **Replication Setup

Install the remote (target) PowerFlex system In order for replication to take place, there must be a second PowerFlex cluster (Target) present. While not covered here, the PowerFlex 

Implementation class walks learners through this process. 

If not installed during the initial deployment, the Storage Data Replicator must be installed and configured on each side of the replication pair. 

Exchange Root Certificates For the two clusters to communicate securely, each PowerFlex system must export its root certificate, upload it to its peer, and import the peer 

certificate. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 161  |

![Image](images/admguide_Path2.svg)

Configure journal capacity Replicate Volumes Between Clusters PowerFlex uses the journal capacity to hold the replication data. The source system accumulates data changes in the journal and sends them 

to the target. The target system accumulates received data in the target journal until a complete consistent image is received and can be applied to the target volumes. 

Add a Replication Peer Once the source and target systems are configured and authentication is established, the replication target system can be defined from the source 

cluster. 

Configure replication volumes The final step in the process is to define the volumes for replication. The volumes to be replicated must be the same size on the source and target 

systems. If the network is up, the systems should be connected. 

# Configure Replication 

Using the SCLI to Extract and Import Root Certificates Swapping root certificates is required to allow communication and data transfer between PowerFlex clusters. The process consists of extracting 

the root certificate from both source and target replication systems, and copying it to the peer system. 

The following process assumes admin user level credentials for command line access to both systems: 

# 1.  Log in to the source system using the SCLI: 

# scli --login_certificate --p12_path

# /opt/emc/scaleio/mdm/cfg/cli_certificate.p12

# 2.  Extract the root certificate: 

# scli --extract_root_ca --certificate_file

# <FILE_PATH>

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 162  |

Replicate Volumes Between Clusters 

# 3.  Copy the certificate file to the Target system (using SCP or similar 

tool). 

# 4.  On the Target system, perform steps 1 and 2. 

# 5.  Copy the Target system's certificate file to the source system. 

# 6.  On the source system, add the certificate for the Target system using 

the SCLI: 

# scli --add_trusted_ca

# <PATH_TO_LOCAL_COPY_OF_TARGET_CERT>

# 7.  On the Target system, add the source system's certificate using the 

SCLI: 

# scli --add_trusted_ca

# <PATH_TO_LOCAL_COPY_OF_SOURCE_CERT>

# Swap Root Certificates - Exercise 

Use the steps that you just learned to complete the simulation. The following environment details help you complete the activity: 

| **SCLI User (both systems):** |  **admin**  |
| **Administrator password:** |  **P0werFl3x**  |
| **Source IP:** |  **10.220.20.6**  |
| **Target IP:** |  **10.220.21.7**  |
| **Local Path for Certificates:** |  **/tmp**  |

 *The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 163  |

# Configure Replication 

How to Calculate Journal Size? 

licate Volumes Between Clusters Rep

There are several factors to consider when allocating journal capacity. 

Journal capacity is defined as a percentage of the total storage capacity in the storage pool. At a minimum, the journal capacity must equal 10% of the total application pool. It is important to assign enough storage capacity for the replication journal. 

The amount of capacity that following factors:  eeded for the journal is based on the  is n $•$  Before configuring journal capacity, ensure that there is enough space in the storage pool. 

$•$  When the total storage capacity in the system increases, a small percentage is needed for the journal capacity. $•$  As application workload increases, more journal capacity must be added accordingly. 

$•$  The journal capacity depends on the change rate of the dataset, and the Recovery Point Objective (RPO). Data writes are accumulated in the journal until half the RPO time is reached, to ensure a consistent copy is maintained between the volumes. 

$•$  The journal capacity must sustain an outage, as determined by the application WAN bandwidth requirement multiplied by the expected WAN outage. If an application has a heavy I/O load, larger capacity should be used. Similarly, if a longer outage is expected, a greater 

capacity should be allocated. 

Example: The following is an example of how to calculate journal capacity allocation: 

# 1.  An application generates 1 GB/s of writes. 

# 2.  The maximal supported outage is 3 hours (3 hours x 3600 seconds = 

10800 seconds). 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 164  |

Replicate Volumes Between Clusters 

# 3.  The journal capacity that is needed for this application is 1 GB/s x 

10800 s = ~10.547 TB. 

# 4.  Since the journal capacity is expressed as a percentage of the storage 

pool capacity, divide the 10.547 TB by the size of the storage pool, which is 200 TB: 100 x 10.547 TB/200 TB = 5.27%. Round up to 6%. 

# 5.  Repeat the process for each application being replicated. 

If there are replicated volumes in more than one storage pool in the protection domain, this calculation should be repeated for each storage pool. The allocated journal capacity in the protection domain must at least equal the sum of the size per application pool. 

# Journal Configuration - Knowledge Check 

Customer Scenario:  

# 1.  A customer has an application that generates about 2 GB/s of writes. 

Their maximum expected outage time is no more than two hours. The Storage Pool that is used by the application is 150 TB. Given these parameters, what percentage of the Storage Pool should be allocated for the Replication Journal? 

a.  5% b.  7% 

c.  10% d.  12% 

Answer Key 

# Configure Replication 

On the source system, add the connection information of the target system (remote site) to enable replication between the peer systems. Prior to performing the steps below, the system ID of both source and target systems must be obtained. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 165  |

Replicate Volumes Between Clusters The system ID is displayed immediately after login to the SCLI. It can also be obtained by running the command **scli --query_all**. 

 Steps:  

 1.  On the source system, on the menu bar, click **Protection** > **Peer **  

**Systems**. 

 2.  Click **Add Peer System**. 

 3.  In the Add Peer System dialog box, enter the connection information of 

the peer system: a.  Enter the peer system's name. Enter the system ID of the remote site. b.  Accept the default or ente

connect the systems. e port number that will be used to r th c.  Enter the MDM IP address of the remote site. 

d.  Either enter the remote system's virtual IP address, or enter both primary and secondary MDM IP addresses, using the **Add IP** option. e.  Click **Add Peer** to initiate a connection with the peer system. 

 4.  Repeat steps 1–3 on the target system. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 166  |

![Image](images/admguide_Image899.png)

Replicate Volumes Between Clusters 

# Configure Replication 

The web version of this content contains an interactive activity. 

The final step in configuring replication is to create a Replication Consistency Group (RCG) and add it to the system. Replication occurs between volumes, and RCGs maintain consistency between volume pairs in an RCG. 

**Important**: 

- Volumes in an RCG pair must be exactly the same size.

 

- Protection Domains must be configured on both source 

  and target systems.  

When configuring an RCG pair, a Recovery Point Object is set. This RPO defines the maximum amount of time during which data can be lost. 

Setting a low RPO ensures that minimal data is lost should the data transfer from source to target be interrupted. 

**Tip**: The data loss exposure is half the RPO value. If one minute is set as the RPO, no more than 30 seconds of data will be lost. Dell highly recommends that the RPO is set low to ensure minimal data loss. The minimum amount of time 

 this feature allows is 15 seconds.  

Steps to configure an RCG: 

# 1.  On the menu bar, click Protection > RCGs. 

# 2.  Click **Create RCG**. 

# 3.  In the Create RCG wizard, enter the information for the RCG: 

- a.  On the Properties page: 
  - $•$  Enter the RCG Name. 
  - $•$  Enter the number of RPO seconds or minutes. 
    - PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 167  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Image535.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path4.svg)

![Image](images/admguide_Path13.svg)

![Image](images/admguide_Path14.svg)

![Image](images/admguide_Path15.svg)

Replicate Volumes Between Clusters $•$  Select the Source Protection Domain. $•$  Select the Target Protection Domain. 

b.  Select the **Provisioning Type**: 

- $•$  Select **Auto Provisioning** if target replication volumes have not 
  - defined on the remote system. 

$•$  Select **Manual Provisioning** if target replication volumes have already been defined on the remote system. c.  For **Auto Provisioning** perform the following: 

- $•$  Click the desired volume in the Source column. 
- $•$  Select a type, Thick or Thin. 
- $•$  Select a storage pool from the target system. 
- $•$  Click **Add Pair** and continue to the Map Volumes page. 
- $•$  Select a volume on the target side. 
- $•$  Select a host to which to map the volume. Repeat this step until 
  - at least three hosts are mapped to each volume. 

$•$  Continue to step 3(e). d.  For **Manual Provisioning** perform the following: 

- $•$  Select the desired volume in the Source column, and then in the 
  - Target column, select the target volume. 

$•$  Click **Add Pair**. $•$  Continue to step 3(e). 

e.  On the Summary page, ensure that the correct source and volume pair are defined, and then click **Create and Activate** or **Create**. 

**Go to:** To view a list of the operations available to manage an RCG, refer to the[ Dell PowerFlex 4.5.x Administration ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/replication-consistency-group?guid=guid-71d462f6-805e-48e1-a2ae-7b0f77d91d7f&lang=en-us) [Guide/Replication consistency group ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/replication-consistency-group?guid=guid-71d462f6-805e-48e1-a2ae-7b0f77d91d7f&lang=en-us)section  

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 168  |

![Image](images/admguide_Image485.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Replicate Volumes Between Clusters 

# Replication I/O Flow 

PowerFlex asynchronous replication uses a journal-based architecture. Journals reside as volumes in a storage pool. A replication journal is copied from the source journal buffer by the SDR to the target journal buffer. Once the journals are secured at the destination, they are deleted 

from the source, making space for new journals. 

To learn about the replication I/O workflow, watch the video below. 

For more information about Asynchronous Replication in PowerFlex, see the Remote Protection section of the[ Dell PowerFlex 4.5.x ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/remote-protection?guid=guid-7df0ed65-dfa1-482e-b48d-6642cd92043d&lang=en-us) [Administration Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_sw_admin_guide_45/remote-protection?guid=guid-7df0ed65-dfa1-482e-b48d-6642cd92043d&lang=en-us) 

 *Video:* 

*The web version* *of this content contains a* *video.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 169  |

![Image](images/admguide_Path2.svg)

Configure Storage Data Servers 

# Configure Storage Data Servers 

## Prerequisites for Adding a Storage Data Server 

 Before adding an SDS, ensure that the following prerequisites are met: 

$•$  At least one storage pool is defined in the required protection domain. $•$  All devices in the storage pool are the same media type, and the storage pool is configured to receive that media type. $•$  At least one acceleration pool is defined if adding acceleration devices. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 170  |

![Image](images/admguide_Image909.pngimages/admguide_Path2.svg)

Configure Storage Data Servers 

# Add a Storage Data Server (SDS) 

 Administrators can add an SDS to a PowerFlex system using PowerFlex 

Manager. The steps to add an SDS are as follows: 

# 1.  On the PowerFlex Manager menu bar, click **Block > SDSs**. 

# 2.  Click **+ Add SDS**. 

# 3.  Configure the following settings: 

- a.  Enter an SDS name. 
- b.  Select a protection domain. 
- c.  Select a fault set. 
- d.  Enter SDS port used for communication. 
- e.  Enter the IP address for SDC, SDS or both and click **Add IP**. 

# 4.  For additional IP addresses, enter the IP address, select the 

communication role and click **Add IP**. 

# 5.  Expand **Advanced** for more options. Configure the following options 

(for advanced users): a.  To enable RMcache, select **Use Read RAM** **Cache** and enter the size in MB. b.  Click one of the options for **Performance Profile**: Compact or High. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 171  |

![Image](images/admguide_Image912.pngimages/admguide_Path2.svg)

c.  To force clean a node, se

# 6.  Click **Add SDS**. 

# Add an SDS - Exercise 

Configure Storage Data Servers  **Force Clean SDS**. lect

To try adding an SDS, use the simulation below. *The web version* *of this content contains an interactive activity.* 

# Add an SDS - Advanced Options 

There are configuration settings that can be applied when adding an SDS. These settings are applied through the Modify drop-down menu in the PowerFlex Manager SDS page. 

 **1: **Use Read RAM Cache is configured to improve the performance of 

systems using hard drives. **2: **Performance profiles are set by default to HIGH but can be changed to 

Compact. The compact setting may impact system performance. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 172  |

![Image](images/admguide_Image915.pngimages/admguide_Path3.svg)

Configure Storage Data Servers **3: **When devices are added to an SDS, PowerFlex checks that the device is clean before adding it. If the device is not clean, an error message is 

returned, and the command fails for that device. If you want to overwrite existing data on the device by forcing the command, set Force Clean SDS to YES. Select YES with caution, because all data on the device is destroyed. 

# Add a Device - Advanced Options 

By default, PowerFlex tests the performance of the device being added before its capacity can be used and saves the results. Two tests are performed: random writes and random reads. When the tests are complete, the device capacity is added automatically to the storage pool 

used by the MDM. 

 **1: **A read and write test is run on the device before its capacity is used when Test and Activate Device is selected. 

**2: **Devices are tested, but not used when Test Only is selected. **3: **The device capacity is used Without Test is selected. out any device testing when Activate  with

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 173  |

![Image](images/admguide_Image917.pngimages/admguide_Path2.svg)

Configure Storage Data Servers **4: **This value is the maximum test run time in seconds. The test stops when it reaches either this limit, or the time it takes to complete 128 MB of 

data read/write, whichever is first. When **Activate without test** is selected, this timeout is ignored. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 174  |

Reconfiguring MDMs 

# Reconfiguring MDMs 

## Reconfiguring MDM Roles 

PowerFlex Manager enables you to change the MDM role for a node in a PowerFlex cluster. For example, if you add a node to the cluster, you might want to switch the MDM role from an existing node to a new node. The MDM Reconfiguration wizard can be used to change MDM role 

assignments after deployment. 

 

## Reconfigure MDM Roles in a PowerFlex Cluster 

To try reconfiguring the MDM roles in a PowerFlex cluster, use the simulation below. 

*The web version* *of this content contains an interactive activity.* 

## A PowerFlex Cluster Refresh 

A PowerFlex customer wants to perform a server refresh, for instance moving from 13G servers to 15G. The administrator must be aware of the current large configuration issues to avoid problems. If the PowerFlex large-scale configuration are node racks with their own subnets, then the 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 175  |

![Image](images/admguide_Image922.jpegimages/admguide_Path4.svg)

Reconfiguring MDMs refresh might add or remove racks. In such a case, the MDM IP addresses change due to the MDMs transferring to new racks with new subnets. 

Automated MDM IP address updates for the SDCs lower the amount of effort in supporting a large-scale server technology refresh. Automation helps to avoid user errors in updates to the SDCs and does not require root access to the compute servers. 

 **Important**: Virtual IP is an effective solution only when the new MDM IP addresses are in the same network as the 

existing MDM IP addresses.   

# Enable and Disable SDC Automatic Update with MDM 

# IP Addresses 

An administrator uses the PowerFlex SDC management SCLI commands to enable or disable the SDC automatic update with MDM IP addresses. 

Specifically, the set_auto_update_sdc_with_mdm_ip_addresses (--enable | --disable) command. 

 SCLI Command Parameters 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 176  |

![Image](images/admguide_Image925.png)

![Image](images/admguide_Image602.png)

![Image](images/admguide_Image927.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path13.svg)

Reconfiguring MDMs 

| **Parameter** |  **Function**  |
| **--enable** |  **Enable** **SDC automatic update with ** |

**MDM IP addresses.** 

| **--disable** |  **Disable** **SDC automatic update ** |

**with MDM IP addresses.**  

**Go to:**  Review the NVMe Host and SDC SCLI Commands for PowerFlex in the[ Dell PowerFlex Manager 4.6.x CLI ](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_cli_reference_guide_46/nvme-host-and-sdc-commands?guid=guid-8ce2e2d2-6dbe-4a05-b93d-66ffaf2842b9&lang=en-us) [Reference Guide.](https://www.dell.com/support/manuals/en-us/scaleio/powerflex_cli_reference_guide_46/nvme-host-and-sdc-commands?guid=guid-8ce2e2d2-6dbe-4a05-b93d-66ffaf2842b9&lang=en-us)  

 

# How the Automated Solution Works 

 The following is a review of how the automated solution works in 

PowerFlex: 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 177  |

![Image](images/admguide_Image485.png)

![Image](images/admguide_Image931.png)

![Image](images/admguide_Image932.png)

![Image](images/admguide_Path9.svg)

![Image](images/admguide_Path10.svg)

![Image](images/admguide_Path11.svg)

Reconfiguring MDMs 

# 1.  The SDC requests updated MDM IP addresses using the existing 

periodic SDC query. The MDM replies to the query. 

# 2.  The SDC sends the new MDM IP addresses to the daemon/agent. 

# 3.  The daemon/agent updates the SDC persistent configuration. 

# 4.  The MDM queries all the SDCs and then creates a report of which 

SDCs were updated and which were not. The administrator queries all SDC reports created by the MDM of all the SDCs. 

After the SDC reboots, the SDC accesses the persistent configuration and uses the new MDM IPs from the persistent configuration to connect to the MDM. If the MDM IPs were modified when the SDC was disconnected, then the SDC requires a manual update after the reboot. 

**Note**: As a kernel driver, the SDC cannot update its persistent configuration in ESXi. A user agent is required to update ESXi. The agent does not open an external port.  

 

# Automated SDC Update with MDM IP Addresses 

# Scenarios 

 MDM IP Change 

When there is an MDM change in a PowerFlex cluster, the MDM carries out the operation through the **replace_cluster_mdm** command. 

The administrator that performed the MDM IP change (technical refresh) adds the new MDMs first, including the standby. The administrator then waits for all the SDCs to update through the query operation. After the SDCs are updated, the administrator removes the old MDMs. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 178  |

![Image](images/admguide_Image934.png)

![Image](images/admguide_Image931.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Reconfiguring MDMs 

 Deployment of a New SDC Including an Agent 

If the SDC and SDC daemon/agent a Installation Bundles (VIBs):  stalled in two separate vSphere re in

$•$  Deployment of the SDC daemon/agent is optional. $•$  The installation of the SDC daemon/agent is manual. 

$•$  The MDM IP is configured on the SDC first installation of the VIBs. 

The behavior of the SDC is the same, in that the SDC queries the MDM for any updates to the MDM IP addresses. However, the SDC sends the query only to the MDMs that support automatic updates. 

Agent Deployment and Upgrade The deployment of an SDC daemon/agent on a server with an existing SDC is supported. Once the daemon/agent is deployed, the SDC begins 

to query the MDM for updates of the MDM IP addresses. The MDMs must support automatic updates to complete the reply for the requested updated MDM IP addresses. 

Any upgrade to the SDC daemon/agent is a manual operation. $•$  Linux: rpm upgrade 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 179  |

![Image](images/admguide_Image937.png)

Reconfiguring MDMs $•$  ESXi: Use the VIB upgrade that is supported through the VMware Update Manager (VUM). 

$•$  Future Orchestration: The PowerFlex road map includes full orchestration through the PowerFlex Installation Manager and PowerFlex Gateway. 

**Important**: Only approved SDCs receive an automatic update.  

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 180  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

NVMe/TCP 

# NVMe/TCP 

## Configure Storage Target (SDT) Service 

Storage Data Targets (or NVMe targets) must be configured on the PowerFlex system in order to use NVMe over Fabric technology. 

The NVMe target is a front-end component that translates NVMe over Fabric protocol into internal PowerFlex protocols. 

The NVMe target provides I/O and discovery services to NVMe hosts configured on the PowerFlex system. 

Before a protection domain can serve NVMe hosts, a minimum of two NVMe targets must be assigned to the protection domain for minimal path resiliency. 

TCP ports, IP addresses, and IP address roles must be configured for each NVMe target. There is the option to assign both storage and host roles to the same target IP addresses. Alternatively, there is the option to assign the storage role to one target IP address and add another target IP 

address for the host role. Both roles must be configured on each NVMe target. The host port listens for incoming connections from hosts over the NVMe protocol. The storage port listens for connections from the MDM. 

Once the NVMe targets have been configured, add hosts to PowerFlex, and then map volumes to the hosts. Connect hosts to NVMe targets, preferably using the discovery feature. 

On the operating system of the compute nodes, NVMe initiators must be configured. Network connectivity is required between the NVMe targets and the NVMe initiators, and between NVMe targets and SDSs. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 181  |

NVMe/TCP 

 

# Add an NVMe Target Using PowerFlex Manager 

In the simulation below, you will be performing the steps to add an NVMe target to PowerFlex. 

*The web version* *of this content contains an interactive activity.* 

# Register NVMe Host Initiator 

Hosts are entities that consume PowerFlex storage for application usage. 

There are two methods of consuming PowerFlex block storage: using the SDC kernel driver or using NVMe over TCP connectivity. A host is either an SDC or an NVMe host. 

Before configuring the Hosts: $•$  Ensure that you have the host NVMe Qualified Name (NQN). If you do not know the NQN, see the host operating system documentation. 

$•$  Ensure that the host is connected to the ethernet switch. $•$  Ensure that the host is configured with the correct VLAN ID and routing rules. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 182  |

![Image](images/admguide_Image941.pngimages/admguide_Path3.svg)

NVMe/TCP 

 

# Add an NVMe Host Using PowerFlex Manager 

In the simulation below, you will be performing the steps to add an NVMe host using PowerFlex Manager. 

# Host** **configuration details:

| **Hostname:** |  Pflex-nvme-1 |  **Number of Paths:**  4  |
| **NQN:** |  nqn.2014-08.com.dell.lab:nvme:plex-nvme-1  **Number of Ports:**  10  |

 *The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 183  |

![Image](images/admguide_Image945.pngimages/admguide_Path2.svg)

Add CloudLink to a PowerFlex Cluster 

# Add CloudLink to a PowerFlex Cluster 

## Role of Cloudlink in a PowerFlex Cluster 

CloudLink is software that secures sensitive information within virtual machines across both private and public clouds. CloudLink features policy-based agent encryption for multiple levels of the data center, including direct integration with PowerFlex to provide software-based 

device encryption. 

PowerFlex does not encrypt pre-existing data on SDS devices. CloudLink is installed to encrypt storage devices before they are added as available PowerFlex storage. 

CloudLink encrypts the SDS devices with unique keys that enterprise security administrators’  control. CloudLink Center provides centralized, policy-based management for these keys. 

       

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 184  |

![Image](images/admguide_Image947.jpeg)

Add CloudLink to a PowerFlex Cluster 

# Add a CloudLink License to the PowerFlex Cluster 

Step 1 On the menu bar, click Settings and then click License Management and 

select other Software Licenses. 

 Step 2 

From the licenses menu, click Add. 

 Step 3 

To upload the PowerFlex license start by clicking the Upload License button in the Add Software License page. Select CloudLink as the license Type. Upload the file and Click Save. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 185  |

![Image](images/admguide_Image949.png)

![Image](images/admguide_Image950.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

# Discover and Deploy Cl

Discover CloudLink Center 

Add CloudLink to a PowerFlex Cluster 

 

# nk oudLi

Before deploying CloudLink, you must discover at least two CloudLink Centers in PowerFlex Manager. 

 Deploy CloudLink Center 

CloudLink can be deployed using the sample template available in PowerFlex Manager. Before cloning the template, ensure that all the VM networks are added and the required Cloudlink resources are discovered in PowerFlex Manager. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 186  |

![Image](images/admguide_Image952.png)

![Image](images/admguide_Image953.png)

![Image](images/admguide_Path2.svg)

Add CloudLink to a PowerFlex Cluster 

 

# Deploy a CloudLink Cluster Using PowerFlex Manager 

To try deploying a CloudLink Cluster using PowerFlex Manager, use the simulation below. 

*The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 187  |

![Image](images/admguide_Image956.pngimages/admguide_Path2.svg)

Configure PowerFlex Alerting 

# Configure PowerFlex Alerting 

## Monitoring Menu 

 Events and Alerts monitoring provide logging for activities happening in the cluster. These logging options can be found, along with the Jobs 

listing, under the Monitoring menu. 

Events 

 An event is a notification that something happened in the system. Events monitor configuration changes, activities of the system, faults, and errors 

in the system and send notifications accordingly to inform the administrator. An event that requires attention generates an alert as well. 

Other events can update or clear an alert when the system detects a change in the condition that needs attention. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 188  |

![Image](images/admguide_Image958.png)

![Image](images/admguide_Image959.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Configure PowerFlex Alerting Alerts 

 An alert is a state in the system which is usually on or off. Alerts monitor 

serious events that require user attention or action. Although the events may be interesting for troubleshooting purposes, it is not necessary to monitor events. Whereas an alert is a summation of one or more events that need or needed attention. 

Jobs 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 189  |

![Image](images/admguide_Image961.jpeg)

![Image](images/admguide_Image962.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

Configure PowerFlex Alerting When user-initiated processes are running on the system, they are listed under the Jobs menu. The entries under Jobs are transient, in that they 

appear only while the process is running and disappear when complete. 

# Enabling SupportAssist 

 PowerFlex Manager can connect directly or through a secure connect gateway to access SupportAssist. Enabling SupportAssist ensures better 

alignment with Dell Technologies services initiatives to enhance the user experience. SupportAssist is the functionality that enables this connection. 

When you enable SupportAssist, you can take advantage of the following benefits, depending on the service agreement on your device: 

$•$  Automated issue detection - SupportAssist monitors your Dell Technologies devices and automatically detects hardware issues, both proactively and predictively. 

$•$  Automated case creation - When an issue is detected, SupportAssist automatically opens a support case with Dell Technologies Support. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 190  |

![Image](images/admguide_Image965.jpegimages/admguide_Path2.svg)

Configure PowerFlex Alerting $•$  Automated diagnostic collection - SupportAssist automatically collects system state information from your devices and uploads it securely to 

Dell Technologies. Dell Technologies Support uses this information to troubleshoot the issue. $•$  Proactive contact - A Dell Technologies Support agent contacts you about the support case and helps you resolve the issue. 

# Configure SupportAssist 

SupportAssist can be configured during the initial configuration of a PowerFlex cluster or skipped and completed after the initial configuration. 

To configure SupportAssist after the initial configuration, browse to Settings, Events and Alerts, then Notification Policies. The steps to configure SupportAssist are the same whichever method is used. 

Enable and Accept EULA The first step in configuring Suppo End-User License Agreements (EULA). is to enable it and accept the rtAssist 

 PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 191  |

![Image](images/admguide_Image967.jpeg)

Configure PowerFlex Alerting Option 1: Direct Connect 

Next, two options are presented for connecting to Dell Support, Connect Directly and Connect using Gateway Server. Connect directly if there is not a secure connect gateway on site to connect through. An assigned access key and PIN are required to make the connection. 

 Option 2: Connect using Gateway 

If there is a secure connect gateway on site, configure SupportAssist to use it as a method for contacting support. The correct IP and configured port are required. An access key and PIN are still required. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 192  |

![Image](images/admguide_Image969.jpeg)

Configure PowerFlex Alerting 

 Support Contact 

Finally, local point-of-contact details are input. Multiple contacts can be designated by clicking the Add Contact link. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 193  |

![Image](images/admguide_Image971.png)

# Configure SupportAssist 

To try configuring SupportAssist u simulation below. 

Configure PowerFlex Alerting 

 

 PowerFlex Manager, use the sing

Use the following Configuration Details for the simulation exercise: 

| Access Key |  AQTPN12345  |
| PIN |  12345  |
| Software Unique ID |  eSWUID2468  |
| Solution Serial Number |  908060  |
| Site ID |  0000  |

 *The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 194  |

![Image](images/admguide_Image973.jpegimages/admguide_Path2.svg)

Configure PowerFlex Alerting 

# Connecting PowerFlex to CloudIQ 

 CloudIQ is a cloud-based application that leverages machine learning to proactively monitor and measure the overall health of Dell IT infrastructure 

systems through intelligent, comprehensive, and predictive analytics. Integrating CloudIQ with PowerFlex is as simple as selecting that option when setting up SupportAssist. The image to the right highlights the option to enable CloudIQ. Once connected, log in to the CloudIQ web console to 

see the status of your PowerFlex system. 

To learn more about the CloudIQ interface, and how it reports on a PowerFlex system, watch the video below: 

- *Video:* 

*The web version* *of this content contains a* *video.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 195  |

![Image](images/admguide_Image975.pngimages/admguide_Path2.svg)

Configure PowerFlex Alerting 

# Configure Events and Alert Notifications to External 

# Systems Overview 

To review additional details about configuring an external alert, see[ Dell ](https://www.dell.com/support/manuals/en-us/powerflex-rack-hw/flex-rack-admin-guide-4x/managing-events-and-alerts?guid=guid-c130858a-3ab0-410a-a5fa-63f6672bfa43&lang=en-us) [PowerFlex Rack with PowerFlex 4.x Administration Guide.](https://www.dell.com/support/manuals/en-us/powerflex-rack-hw/flex-rack-admin-guide-4x/managing-events-and-alerts?guid=guid-c130858a-3ab0-410a-a5fa-63f6672bfa43&lang=en-us) 

 While SupportAssist is technically an external alerting system, this section deals with non-SupportAssist external alerting. PowerFlex is capable of 

both receiving SNMP or Syslog alerts from external sources and sending alerts to external sources such as messaging servers and system management ecosystems. In both cases, configuration is done from the Events and Alerts section of PowerFlex Manager Settings page. 

Configure an External Source A source is used to configure the receiving of external events and syslog content. 

SNMP sources are not automatically d to receive events about these sources. ered and must be configured iscov

Resources in the PowerFlex systems are automatically discovered. Any future resources, for example, switch replacements or additional nodes, are considered external and must be added manually as sources. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 196  |

![Image](images/admguide_Image979.png)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path8.svg)

Configure PowerFlex Alerting Configure Destination A destination is used to configure the ability to send events and alerts 

information out. A destination is always external. 

Notification policies define what information is sent to each destination. 

SupportAssist, email, SNMP, and remote syslog are considered as possible external destinations. 

 

# Configure External Source and Destination 

To try configuring External Sou below and Destination, use the simulation rce 

*The web version* *of this content contains an interactive activity.* 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 197  |

![Image](images/admguide_Image982.pngimages/admguide_Path2.svg)

PowerFlex in the Cloud (APEX) 

# PowerFlex in the Cloud (APEX) 

## What Is Dell APEX? 

 Dell APEX is a portfolio of cloud and consumption experience offerings 

that brings a cloud operating model to devices, applications, and data. $•$  Dell APEX provides the following benefits: 

- –  Delivers simplified cloud experiences by providing choice, 
  - consistency, and freedom without fear of being locked into a 
  - solution. 

–  Improves agility by providing the flexibility to deploy resources quickly and predictably in a transparent pay as you go model. –  Allows customers to maintain control of their data to both minimize risk and maximize performance. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 198  |

![Image](images/admguide_Image985.jpegimages/admguide_Path2.svg)

PowerFlex in the Cloud (APEX) 

# The Dell APEX as-a-Service Portfolio Offerings 

 Dell APEX is an end-to-end as-a-Service offering. Customers can choose which service fits their needs for cloud deployments and related services. 

The service categories with offerings are: $•$  Storage: Customers can choose the cloud storage service offerings for either a private or public cloud deployment. 

$•$  Cyber and Data Protection: Safeguarding business operations with APEX are provided in this service.  The offerings are backup services, data storage services backup target, and cyber recovery services. 

$•$  Cloud Platforms: Customers can modernize their cloud deployments with the platform choices that are offered in APEX. The offerings include Cloud Platform for Microsoft Azure and Cloud Platform for Red Hat OpenShift. 

$•$  Custom Solutions: A customer can implement usage-based consumption on the technology they choose through Dell APEX On Demand and Dell APEX Data Center Utility. 

**Go to:** Review the[ Dell APEX site ](https://www.dell.com/en-us/dt/apex/index.htm?gacd=9650523-1137-5761040-266691960-0&dgc=ST&SA360CID=71700000083414790&gad_source=1)to learn more about APEX features and services.  

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 199  |

![Image](images/admguide_Image485.png)

![Image](images/admguide_Image988.png)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path4.svg)

![Image](images/admguide_Path11.svg)

# APEX Subscription Model 

PowerFlex in the Cloud (APEX) 

Customers choose the APEX services through a subscription model, where customers pay only for the services consumed. The subscription model is a simplified 1-2-3 approach to get customers the services they need from APEX. 

Step 1: Choose Your Technology 

 To start, customers configure the technology capabilities to align with their technology requirements. For example, customers may need Storage 

Infrastructure technology support for a PowerFlex block storage implementation. The customer configuration for each technology dictates the deployment requirements for APEX. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 200  |

![Image](images/admguide_Image991.png)

Step 2: Choose the APEX Service PowerFlex in the Cloud (APEX) 

 Next, customers choose the APEX as-a-Service to deploy and support the APEX technology subscription. Customers add additional infrastructure 

and managed services to support people, processes, and workloads to meet their business needs. In the example PowerFlex block storage implementation, the cloud deployment service is a Dell APEX Block Storage for Public Cloud offering. The selected deployment service 

adheres to the configuration previously implemented. Other services customers may choose for the cloud deployment might be a combination of data protection services and consumption utilities for the data center. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 201  |

![Image](images/admguide_Image994.png)

PowerFlex in the Cloud (APEX) Step 3: Pay Based on Consumption 

 The expense for a customer is directly proportional to the number of service features that are used. Usage includes maintenance and operation 

of the APEX cloud deployments. 

Summary Overview To review the Dell APEX Subscription model features, watch the video below. 

# Video:

*The web version* *of this content contains a* *video.* When customers choose APEX for their technology requirements, they are choosing both functional and cost savings advantages. In the APEX 

Block Storage for Cloud, for example, the inclusion of PowerFlex with APEX provides functional advantages over public cloud-native block storage offerings. With Dell APEX Block Storage for AWS, customers realize up to 87% cost savings over AWS EBS volume storage. With Dell 

APEX Block Storage for Azure, customers can realize up to 82% cost savings over Microsoft Azure managed virtual disk storage. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 202  |

![Image](images/admguide_Image997.png)

PowerFlex in the Cloud (APEX) 

**Go to:** learn more about the APEX subscription model from the[ Dell APEX Subscription site.](https://www.dell.com/en-us/dt/apex/subscriptions.htm?hve=dell+apex+subscriptions)  

# PowerFlex Cloud 

 

PowerFlex offers three on-premises deployment models: $•$  PowerFlex software-only $•$  PowerFlex appliance 

$•$  PowerFlex rack PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 203  |

![Image](images/admguide_Image485.png)

![Image](images/admguide_Image1001.png)

![Image](images/admguide_Image1002.jpeg)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

PowerFlex in the Cloud (APEX) With PowerFlex 4.5.x, customers realize a new deployment model: PowerFlex Cloud. Customers experience the same benefits of enterprise-

class storage services in the cloud as with on-premises. PowerFlex accomplishes this cloud deployment model through one of two ways: 

$•$  PowerFlex cloud-based block storage cluster deployed through the Dell APEX Block Storage for Public Cloud service. 

$•$  Administrators manually deploy PowerFlex Cloud for compliant public cloud environments. 

PowerFlex Cloud, also known as Dell APEX Block Storage for Public Cloud, is a cloud-hosted version of the PowerFlex software-defined storage (SDS). PowerFlex Cloud provides higher performance, larger volume sizes, and improved resilience than what is available on the public 

cloud. 

Review the capabilities and benefits of the PowerFlex Cloud deployment model in the[ PowerFlex 4.6.x Technical Overview.](https://www.dell.com/support/manuals/en-us/scaleio/flex-software-to-46x/dell-apex-block-storage?guid=guid-cbf6ab42-d831-4356-8605-3aa6fe675902&lang=en-us) 

**Important**: The following features in PowerFlex 4.5.x and 4.6.x are NOT supported with the APEX Block Storage for Public Cloud: compression, fine granularity storage pools in the public cloud, SDNAS, PowerFlex file services, or NVMe 

 TCP.  

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 204  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

PowerFlex in the Cloud (APEX) 

# PowerFlex Public Cloud Environments 

 Dell APEX Block Storage for Public Cloud brings many PowerFlex on-premises features to Amazon Web Services™  (AWS™ ) and Microsoft® 

Azure™  public cloud environments. The features include linear scaling of performance and capacity. 

$•$  Dell APEX Block Storage for AWS empowers enterprises to run diverse workloads in the public cloud while ensuring extreme performance, scalability, and a simplified cloud experience. To learn more about Dell APEX Block Storage for AWS, review the[ Dell APEX ](https://infohub.delltechnologies.com/en-us/t/dell-apex-block-storage-for-aws/)

[Block Storage for AWS white paper.](https://infohub.delltechnologies.com/en-us/t/dell-apex-block-storage-for-aws/) $•$  Dell APEX Block Storage for Azure can be deployed in two configurations. The deployments use Azure Managed Disks or virtual 

machines with attached NVMe SSDs based on the use case. To learn more about Dell APEX Block Storage for Microsoft Azure, review the [Dell APEX Block Storage for Microsoft Azure white paper.](https://infohub.delltechnologies.com/en-us/t/dell-apex-block-storage-for-microsoft-azure/) 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 205  |

![Image](images/admguide_Image1010.png)

PowerFlex in the Cloud (APEX) 

**Go to:**  The[ Dell APEX Block Storage page ](https://www.dell.com/en-us/dt/apex/storage/public-cloud/block.htm?hve=explore+block)to learn more about the PowerFlex enabled cloud deployment. Go to the [Dell APEX Block Storage Technical Documentation ](https://www.dell.com/support/kbdoc/en-ie/000226210)page to review the Getting Started, Deployment, and the 

 Management and Operations Guides.  

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 206  |

![Image](images/admguide_Image485.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

 

# Appendix 

## Lightweight Installer Agent (LIA) 

The PowerFlex Lightweight Installation Agent is installed on every node during PowerFlex deployment. LIA is used to upgrade the component on which it is installed and is required for many maintenance operations. LIA can be configured to use LDAP authentication following PowerFlex 

deployment. To learn more about setting the LIA authentication, review the[ PowerFlex REST API Guide ](https://www.dell.com/support/manuals/en-ie/scaleio/vxf_p_rest_api_reference_guide_3_x_sw/set-the-lia-authentication-method?guid=guid-246542a2-daa7-4d3b-ad48-6a3e376b4b6b&lang=en-us)located in dell.com/support. 

## Nodes 

Nodes are the basic hardware units, hypervisor and/or PowerFlex software.  are used to install and run a which

Dell PowerFlex custom node, PowerFlex appliance, and PowerFlex rack bring together Dell PowerEdge servers and Dell PowerFlex software into an integrated solution. The combination is done in a reliable, quick, and easy to deploy building block. These building blocks are ideal for server 

SAN, heterogeneous virtualized environments, and high-performance databases. 

## Storage Media 

In PowerFlex, the storage media is categorized based on the type, location, and memory modules as shown in the diagram. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 207  |

Appendix 

 

# MDM Cluster 

The MDM cluster consists of a combination of primary MDM, secondary MDMs, and tiebreaker MDMs. An MDM is assigned a Primary, Secondary, or a Tie-breaker role, during deployment. 

$•$  A three-node MDM cluster contains a primary MDM, a secondary MDM, and a tie-breaker MDM. An additional node is used as a Standby. 

$•$  A five-node MDM cluster contains a primary MDM, two secondary MDMs, and two tie-breaker MDMs. An additional node is used as a Standby. 

$•$  Primary MDM, secondary MDMs, and tie-breakers can be distributed across multiple physical cabinets and access switch pairs to ensure maximum availability of the cluster. 

$•$  MDM nodes can be Storage-only or Hyperconverged servers. $•$  MDM nodes should be placed on SO nodes first when available. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 208  |

![Image](images/admguide_Image1018.pngimages/admguide_Path2.svg)

Appendix 

 

# Compliance Files 

 To upgrade a cluster, the upgrade engineer must first upload a specific compliance file bundle to establish the software and firmware standards. 

$•$  Release Certification Matrix (RCM) compliance files are used to provide the software and firmware standards for a PowerFlex rack environment. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 209  |

![Image](images/admguide_Image1021.png)

![Image](images/admguide_Image1022.png)

Appendix $•$  Intelligent Catalog (IC) compliance files are used to provide the software and firmware standards for a PowerFlex appliance 

environment. $•$  Software-only compliance files are used to provide software standards for a software only PowerFlex environment. 

# Storage Pool Granularity Comparison 

# Medium Granularity** **Pools**  **Fine Granularity

# Pools

| Volumes |  Supports thick or thin- | Supports only thin- |
| provisioned volumes and |  provisioned, "zero- |

supports Zero padded and Non-padded" volumes zero padded Storage pools 

| Space |  1 MB units |  4 KB units  |

Allocation 

| Device Type  SSD Media |  Require SSD/NVMe  |

Media with NVDIMM-N persistent memory modules. 

| Compression  Not supported |  Supported  |
| Usage |  Recommended for workloads  Ideal for workloads  |
| with high-performance |  where space  |
| requirements |  efficiency is more  |

valuable than raw I/O performance  

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 210  |

Appendix 

# Zero-Padding

Each Storage Pool can work in one of the following modes: **Zero-padding enabled:** Ensures that every read from an area previously not written to returns zeros. Some 

applications might depend on this behavior. Zero padding also ensures that reading from a volume will not return information that was previously deleted from the volume. This behavior incurs some performance overhead on the 

first write to every area of the volume since the area needs to be filled with zeros first. FG is always zero padded.  **Zero-padding disabled (default only for MG):** A read from 

an area previously not written to will return unknown content. This content might change on subsequent reads. 

Zero padding must be enabled if you plan to use any other application that assumes that when reading from areas not written to before, the storage will return zeros or consistent data. 

 

# NVDIMM Acceleration Pool 

In an Acceleration Pool, the cached data resides on NVDIMM devices. These devices are mapped to a logical DAX device. DAX devices are a specialized file system that are created on NVDIMMs for caching use on fine granularity storage pools. Writes are assembled and buffered in the 

Acceleration Pool. 

The NVDIMM device allows a faster write acknowledgment and improves the performance of operations such as compression with minimal addition of latency. The acknowledgment is sent for an I/O when the I/O is written to the battery protected NVDIMM. When the compression operation has 

finished, the blocks of I/O are written to the SSD device. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 211  |

![Image](images/admguide_Image570.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

Appendix 

 

# Change the iDRAC Password 

The Integrated Dell Remote Access Controller is a piece of hardware that allows administrators to update and manage Dell systems. IT administrators deploy, update, and monitor PowerEdge servers anywhere and anytime with the iDRAC secure local and remote server management. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 212  |

![Image](images/admguide_Image1026.png)

Appendix 

 The steps to change iDRAC password using the iDRAC web interface are: 

# 1.  In the iDRAC Web Interface, go to **iDRAC Settings > User.** 

# 2.  In the User ID column, select **user ID 2** and click Edit. 

# 3.  Modify the user settings as needed. 

# 4.  Click **Save**. 

# Change the Server BIOS Password 

The server BIOS, also called System Setup, manages data flow between the operating system and attached devices. Sometimes security policies require that BIOS passwords be changed at regular intervals. Whenever changing the BIOS password, ensure that the PowerFlex cluster is 

updated to reflect the new password. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 213  |

![Image](images/admguide_Image1029.jpegimages/admguide_Path2.svg)

Appendix 

 The steps to change Dell Server's BIOS password are: 

 1.  Enter System Setup by pressing **F2** immediately after turning on or 

restarting your system. 

 2.  On the **System Setup Main Menu**, click **System BIOS > System** 

**Security.** 

 3.  On the System Security screen, ensure that the **Password Status** is 

set to **Unlocked**. Note that, if the password status is set to **locked**, you cannot change or delete the system or set up a password. 

 4.  In the **System Password** field, change or delete the existing system 

password and press **Enter.** 

 5.  In the **Setup Password** field, change or delete the existing setup 

password and press **Enter**. If you change the system or setup password, a message prompts you to re-enter the password. If you delete a password, a message prompts you to confirm the deletion. 

 6.  Press **Esc** to return to the System BIOS screen. Press **Esc** again. A 

message prompts you to save the changes. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 214  |

![Image](images/admguide_Image1031.pngimages/admguide_Path1.svg)

Appendix 

# Change the vCenter and ESXi Password 

During the deployment of the PowerFlex appliance, the installation user sets the VMware vCenter password and VMware ESXi operating system password in PowerFlex Manager. 

The steps to change vCenter and ESXi passwords are: 

 

# 1.  Change the PowerFlex Manager **VMware vCenter** or **VMware** **ESXi** 

**operating system** password by completing the following: $−$  In PowerFlex Manager, go to **Settings > Security > Resource** 

# Credentials > Credentials Management**, select the **VMware

**vCenter** or **VMWare ESXi operating system** credential, click **Modify**, change the **Password** to the **<NEW_PASSWORD>**, and click **Save**. 

# 2.  Change the **VMware** **vCenter** password by completing the following: 

- $−$  Log in to the VMware vCenter web interface using the 
  - **<OLD_PASSWORD>.** 

$−$  Click the username in upper right of page and select **Change ** **password**. 

$−$  Type the **<OLD_PASSWORD>** and the **<NEW_PASSWORD>** and click **OK.** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 215  |

![Image](images/admguide_Image1033.jpegimages/admguide_Path2.svg)

Appendix 

# 3.  Change the **VMware** **ESXi operating system root** password on every 

hyperconverged or PowerFlex compute-only node, by completing the following: 

$−$  Log in to the VMware ESXi web interface on the PowerFlex node using **root** and the **<OLD_PASSWORD>.** 

$−$  In the upper right of page, click the **root@<ip address>**, and select **Change password**. 

$−$  Type the **<NEW_PASSWORD>** twice and click **Change password.** 

# 4.  Test the changes. Even though the cluster is operating properly, 

because of the time between changing the password in PowerFlex Manager and changing the password in the ESXi OS, nodes may show a critical error on the **Services** page in PowerFlex Manager. The following steps will return the nodes to a healthy state. 

$−$  In the PowerFlex Manager GUI, go to **Resources** page, select **vCenter or ESXi nodes** and click **Run Inventory.** 

$−$  To confirm that the process completes with no errors, check **Settings > Logs.** 

$−$  In the PowerFlex Manager GUI, go to **Services** page for ESXi nodes and click **Update Service Details.** 

$−$  After **Update Service** **Details** completes the process, confirm that all cluster objects report as healthy (green check mark). 

# Change the Embedded Linux Management VM 

# Password 

An administrator sets the embedded operating system password in PowerFlex Manager during the deployment of the PowerFlex appliance. If the embedded operating system password is changed after deployment, administrators must also change it in PowerFlex Manager. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 216  |

The Steps to change Embedded L

Appendix 

 x Management VM password are: inu

 1.  Change the PowerFlex Manager embedded operating system 

password by completing the following: $−$  In PowerFlex Manager, go to **Settings > Security > Resource** 

**Credentials. **The Credentials Management page opens, select the **embedded Linux operating system** credential, click **Modify**, change the **Password** to the **<NEW_PASSWORD>**, and click **Save**. 

 2.  Change the **embedded Linux operating system root** password on 

every PowerFlex storage-only node, by completing the following: $−$  Use an SSH client program like **PuTTY** to log in as a user with 

administrative privileges to the embedded operating system console using the **<OLD_PASSWORD>.** 

$−$  Change the embedded operating system root password using **passwd** command: o  **[root@node1 ~] # passwd** o  **Changing password for user root.** 

o  **New password: <NEW_PASSWORD>** o  **Retype new password: <NEW_PASSWORD>** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 217  |

![Image](images/admguide_Image1036.pngimages/admguide_Path1.svg)

Appendix o  **passwd: all authentication tokens updated ** **successfully.** 

# 3.  Test the changes. Nodes may show a critical error on the PowerFlex 

Manager **Services** page. The error may appear even though the cluster is operating properly. The error occurs because of the time lapse between changing the password between the Linux OS and in PowerFlex Manager. The following steps will return the nodes to a 

healthy state. $−$  In the PowerFlex Manager GUI, go to **Resources** page, select **embedded operating system nodes** and click **Run Inventory.** 

$−$  To confirm that the process completes with no errors, check **Settings > Logs.** 

$−$  In the PowerFlex Manager GUI, go to **Services** page for Linux nodes and click **Update Service Details.** 

$−$  After **Update Service** **Details** completes the process, confirm that all cluster objects report as healthy (green check mark). 

# PowerFlex Manager User Roles 

User roles control the activities that can be performed by different types of users, depending on the activities that they perform when using PowerFlex Manager. 

The table summarizes the activities that are performed by each user role.       

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 218  |

 

| **Role** |  **Description**  |

# Super User**  **Performs all** **system

# operations

Appendix 

# Activities

$•$  Manage storage resources $•$  Manage lifecycle operations, resource groups, templates, deployment, backend 

operations $•$  Manage replication operations, peer systems, RCGs 

$•$  Manage snapshots, snapshot policies $•$  Manage users, certificates $•$  Replace drives 

$•$  Hardware operations $•$  View storage configurations, resource details $•$  View platform configuration, 

resource details $•$  System monitoring (events, alerts) $•$  Perform serviceability 

operations $•$  Update system settings 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 219  |

Appendix 

| **Super** |  **Perform all **   $•$  Manage storage resources  |

# Admin**  **operations,** **except for

$•$  Manage lifecycle operations, 

# user management** **and

resource groups, templates, 

# security

deployment, backend operations $•$  Manage replication operations, peer systems, 

RCGs $•$  Manage snapshots, snapshot policies $•$  Replace drives 

$•$  Hardware operations $•$  View storage configurations, resource details $•$  View platform configuration, 

resource details $•$  System monitoring (events, alerts) $•$  Perform serviceability 

operations $•$  Update system settings 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 220  |

Appendix 

**Storage**  **Perform all storage-** $•$  Manage storage resources 

# Admin**  **related front-end

$•$  Manage lifecycle operations, 

# operations including

resource groups, templates, 

# element management

deployment, backend 

# of already setup NAS

operations **and block** **systems.** $•$  Manage replication 

# For example: create

operations, peer systems, 

# volume, create file

RCGs 

# system, manage file-

**server user quotas.**  $•$  Manage snapshots, snapshot policies 

$•$  Replace drives $•$  Hardware operations 

$•$  View storage configurations, resource details $•$  View platform configuration, resource details 

$•$  System monitoring (events, alerts) 

**Lifecycle**  **Manage the life cycle**  $•$  Manage lifecycle operations, 

| **Admin**  **of hardware** **and ** | resource groups, templates,  |
| **PowerFlex** **systems.** |  deployment, backend  |

> **NOTE: Operations **    operations 

# such as create

$•$  Replace drives 

# storage pool,** **create

$•$  Hardware operations 

# file-server, and add

**NAS node cannot be**  $•$  View resource groups and 

| **performed by** **a** |  templates  |

# Storage Admin, but

$•$  System monitoring (events, 

# can be performed by

alerts) 

# the Lifecycle Admin

**role.** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 221  |

# Replication**  **A subset of the

**Manager**  **Storage Admin role, **

# for work on existing

# systems for setup an

# management of

# replication and

**snapshots.** 

Appendix 

$•$  Manage replication operations, peer systems, RCGs 

$•$  Manage snapshots, snapshot policies $•$  View storage configurations, resource details (volume, 

snapshot, replication views) $•$  System monitoring (events, alerts) 

**d **

**Snapshot**  **A subset of Storage**  $•$  Manage snapshots, snapshot 

| **Manager**  **Admin, working only** |  policies  |

**on existing systems. ** $•$  View storage configurations, 

# This role** **includes all

resource details 

# operations required to

$•$  System monitoring (events, 

# set up and manage

alerts) **snapshots.** 

| **Security**  **Manages** | $•$  Manage users, certificates  |

# Admin**  **PowerFlex** **role-based

$•$  System monitoring (events, 

# access** **control

alerts) 

# (RBAC), and LDAP

# user federation. It

# includes all security

# aspects of the

**system.** 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 222  |

Appendix 

**Technician**  **Do all hardware field ** $•$  Replace drives 

# replacement unit

$•$  Hardware operations 

# (FRU) operations on

$•$  System monitoring (events, 

# the system. A

alerts) 

# Technician also

**performs the relevant**  $•$  Perform serviceability **commands for proper ** operations 

# maintenance, such as

# entering a node into

**maintenance mode.** 

| **Drive** |  **A subset of the **   $•$  Replace drives  |

# Replacer**  **Technician role. The

$•$  System monitoring (events, 

# Drive Replacer is a

alerts) 

# user who is only

# allowed to do

# operations required

**for drive replacement. **

# For example: life

# cycle operations on

# the node and

# evacuating a block

**system device.** **Monitor**  **Read-only access to ** $•$  View storage configurations, **the system, including**  resource details 

**topology, alerts, ** $•$  View platform configuration, **events, and metrics.** resource details 

$•$  System monitoring (events, alerts) 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 223  |

Appendix 

**Support**  **A special kind of **   $•$  Manage storage resources 

# System Admin (all

$•$  Manage lifecycle operations, 

# activities** **except for

resource groups, templates, 

# user/security

deployment, backend 

# management

operations 

# operations) to be used

$•$  Manage replication 

# only by Dell

operations, peer systems, 

# Technologies** **support

RCGs **staff and developers.** **This user role has**  $•$  Manage snapshots, snapshot 

| **access to ** | policies  |

**undocumented, ** $•$  Replace drives 

# special operations

# and options for

$•$  Hardware operations **common operations,** $•$  View storage configurations, 

# required only for

resource details **support purposes.** $•$  View platform configuration, 

> **NOTE: This special **

resource details 

# role should be used

| **only by** **Dell** | $•$  System monitoring (events,  |

**Technologies** **support. ** alerts) 

# It opens special, often

$•$  Perform serviceability **dangerous,** operations 

# commands for

$•$  Special Dell Technologies 

# advanced trouble

Support operations **shooting.**  

# Add LDAP Users or Groups 

You can add and modify LDAP users and groups in PowerFlex Manager and assign roles to them. These roles control access permissions for the corresponding LDAP user or group. These users and groups must also be configured in the directory service. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 224  |

Steps to add LDAP Users or Groups are: 

# 1.  On the menu bar, click **Settings.** 

# 2.  In the left pane, click **User Manageme**

**LDAP Users**. 

# 3.  Click **Add.** 

Appendix 

 

, then in the right pane, click **nt**

# 4.  In the **Add LDAP User/Group** dialog box, select a **Type** option. 

$•$  User - a user definition will be configured for an individual user. $•$  Group - a group definition will be configured for a specific group. 

# 5.  In the **User Name** box, enter the user or group name. 

# 6.  In the **User Role** box, select the role to be assigned to the user or 

group. 

# 7.  Click **Apply.** 

# Modify LDAP Users or Groups 

You can add and modify LDAP users and groups in PowerFlex Manager and assign roles to them. These roles control access permissions for the corresponding LDAP user or group. These users and groups must also be configured in the directory service. 

Steps to modify LDAP Users or Groups are: 

- PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 225  |

![Image](images/admguide_Image1046.jpegimages/admguide_Path2.svg)

Appendix 

# 1.  On the menu bar, click **Settings > LDAP Users.** 

# 2.  Click **Modify.** 

# 3.  In the **Modify LDAP User/Group** dialog box, change the user role by 

selecting one of the **User Role** options: $•$  SuperUser 

$•$  SystemAdmin $•$  StorageAdmin 

$•$  LifecycleAdmin $•$  ReplicationManager 

$•$  SnapshotManager $•$  SecurityAdmin 

$•$  DriveReplacer $•$  Technician 

$•$  Monitor $•$  Support 

# 4.  Click **Apply.** 

# Add Storage Data Replicator 

A minimum of two SDRs are required on each replication system. Each SDR must be configured with one or more IP addresses and roles. 

The SDR communicates with several components, including: SDC (application), SDS (storage), and remote SDR (external). When an IP address is added to an SDR, the role or roles of the IP address must be defined. The IP address role determines the component with which that IP 

address communicates. For example, the application role means that the associated IP address is used for SDR-SDC communication. By default, all the roles are selected for an IP address. SDR components must be deployed as resources before you can add them using this procedure. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 226  |

Appendix 

 

 1.  On the menu bar, click **Protection** > **SDRs**. 

 2.  Click **Add SDR**. 

 3.  In the Add SDR dialog box, enter the connection information of the 

SDR: a.  Enter the SDR name. 

b.  If necessary, modify the SDR port number. c.  Select the relevant protection domain. 

d.  Enter the IP address of the SDR. e.  Select one or more roles, for example, default: All roles are selected. f.  If the SDR has more than one IP address, click **Add IP** to add more 

IP addresses and their roles. g.  Click **Add SDR** to initiate a connection with the peer system. 

 4.  Verify that the operation has finished and was successful, and click 

**Dismiss**. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 227  |

![Image](images/admguide_Image1050.png)

Appendix 

**Important**: SDR components (Storage-only or Hyperconverged "-replication" Resource Group) must be deployed before you can successfully add SDRs using this  procedure. 

 

# Add Peer System Dialog Box 

 **1: **A name for the replication Peer PowerFlex System (displays in PowerFlex Manager) 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 228  |

![Image](images/admguide_Image602.png)

![Image](images/admguide_Image1053.png)

![Image](images/admguide_Path1.svg)

![Image](images/admguide_Path2.svg)

![Image](images/admguide_Path3.svg)

![Image](images/admguide_Path13.svg)

Appendix **2: **System ID of the remote replication Peer. This can be gathered by issuing the **SCLI --query_all** command on the target system. 

**3: **IP registration box - this is where the MDMs from the target system are entered. 

**4: **IP registration list - this is where the registered MDMs from the target system are listed. 

# Modify an SDS 

 

# PowerFlex Cluster Refresh Problem Definition 

$•$  Some large-scale PowerFlex cluster hardware refreshes can involve moving the MDMs to new networks. $•$  In a refresh, the current PowerFlex Cluster SDCs and their persistent configuration must be updated with the new MDM IP addresses. 

$•$  In PowerFlex 3.x (not including 3.6.1) and 4.0 clusters, the SDC persistent configuration update is done manually. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 229  |

![Image](images/admguide_Image1055.pngimages/admguide_Path3.svg)

Appendix $•$  Manual SDC configuration requires a host reboot and can be a challenge in a large PowerFlex Cluster configuration. 

$•$  **Drv_cfg** can be used to update the SDC configuration without reboot but does require root access. 

# Automated MDM IP Address Update for SDC 

# Configurations Requirements and Limitations 

Automation Requirements: $•$  Backward Compatibility - existing update methods are fully supported $•$  Deployment - optional 

$•$  Security - No change $•$  PowerFlex releases - 3.6.1 and 4.5 and higher 

$•$  Operating Systems - ESXi 7.0 (PowerFlex 3.6.1), Linux (PowerFlex 3.6.1) 

Automation Limitations: $•$  The automation solution is not supported in PowerFlex 4.0. $•$  The automated solution is supported in ESXi version 7.0 or higher that 

support Daemon-SDK. $•$  SDCs that are disconnected from the MDMs for the length of the technical refresh require a manual update. A list of the SDCs that require manual update are provided by the system. 

# Add SDR System Dialog Box 

Enter the remote SDR details on the "Add SDR" page. To learn about the fields in the dialog below, select each red outlined hotspot. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 230  |

Appendix 

 **1: **Name for this Storage Data Replicator. 

**2: **Storage data replicator communications port (pre-populated) can be updated here. 

**3: **Local protection domain that contains the storage pool being replicated. **4: **Remote IP address of the peer SDR. 

**5: **The IP address role determines the component with which that IP address communicates. For example, the application role means that the associated IP address is used for SDR-SDC communication. By default, all the roles are selected for an IP address. 

**6: **List of currently registered peer SDR IP addresses. 

# Identify Existing Product Versions 

PowerFlex Manager, Element Managers, and PowerFlex rack components must be at the supported version for a successful expansion. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 231  |

![Image](images/admguide_Image1058.png)

Appendix To verify the version of the components, review the compliance report from the Resources page in PowerFlex Manager. 

If any of the components are not at the documented version, they must be upgraded in PowerFlex Manager before the expansion begins. This is a remediation task for the customer to perform before the expansion. Upgrades may require customer change windows or downtime. 

 New nodes being added must be at the same operating system level as the current nodes or be remediated as soon as being added to the cluster. 

 

PowerFlex 4.x Administration-SSP 

| © Copyright 2023 Dell Inc |  Page 232  |

![Image](images/admguide_Image1060.jpeg)

 

# Glossary 

Bring your own hypervisor (BYO) Option to allow customer to supply a hypervisor to host the Kubernetes VM/Env, and enable setup of the PFMP system for management 

purposes. Can be ESXi or Kubernetes VM. 

Business Continuity and Disaster Recovery (BCDR) A business continuity and disaster recovery (BCDR) plan protect an organization from disruption, downtime, and data loss if a disaster occurs. 

Such as ransomware attack, human error, accidental or malicious deletion, and natural disasters. 

Co-residency The PowerFlex storage-only nodes in a protection domain runs the management of the system. 

Dual Signature Policy There is a multistep business-approval workflow with specific documentation requirements (including two business leader signatures) 

before a support ticket is created. The required documentation is reviewed, and a login session with support is facilitated in which the snapshot retention period is reset, and the snapshot is deleted. 

Dual Signature Policy There is a multistep business-approval workflow with specific documentation requirements (including two business leader signatures) 

before a support ticket is created. The required documentation is reviewed, and a login session with support is facilitated in which the snapshot retention period is reset, and the snapshot is deleted. 

Enterprise Project Services (EPS) EPS is an end-to-end deployment project platform that contains the configuration details for the installation and implementation of a PowerFlex 

appliance system. 

Expiration Time 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 233  |

 

The expiration time determines the retention interval. 

Fault Set A Fault Set is a logical entity that contains a group of SDSs within a Protection Domain (PD). The SDSs have a higher chance of failing 

together. By grouping the SDSs into a Fault Set, PowerFlex mirrors the data for a Fault Set on separate SDSs that are outside the Fault Set. In this way, availability is assured even if all the servers within one Fault Set fail simultaneously. 

File Services Node (FSN) The File Services Node (FSN) is a software component that allows the PowerFlex cluster to make data available over file-based protocols (NAS). 

The FSN supports protocols such as SMB, NFS, and FTP. 

Fine Granularity Pool In a Fine Granularity (FG) Pool, the volumes are divided into 4 KB allocation units. The FG layout requires both Flash media (SSD or NVMe) 

and NVDIMM to create an FG pool. FG layout is thin-provisioned and zero-padded natively, and enables support for inline compression, more efficient snapshots, and persistent checksums. 

Glossary Term PowerFlex terms defined from acronyms or new features and functions. 

Hyperconverged Deployment In the hyperconverged (HCI) configuration, both the SDC and the SDS can be installed on the same host. HCI deployment maximizes hardware 

utilization and reduces infrastructure requirements. 

Intelligent Catalog (IC) The PowerFlex Intelligent Catalog (IC) is a definition of all compatible versions of hardware and software that are tested and validated together 

for a PowerFlex appliance deployment. 

IOPS 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 234  |

 

IOPS (input/output operations per second) is the standard unit of measurement for the maximum number of reads and writes. 

Lightweight Installer Agent (LIA) The PowerFlex Lightweight Installation Agent (LIA) is used to upgrade the component on which it is installed and is required for many maintenance 

operations. LIA can be configured to use LDAP authentication following the PowerFlex deployment. 

Management data store (MDS) The PowerFlex datastore on the dedicated multi-node controller. It is not the same as the datastore used in a co-resident option. Kubernetes VMs 

that host the management containers reside on this datastore. 

Management virtual machine (MVM) The PowerFlex Manager UI is hosted in a Kubernetes environment running across a series of Linux machines. The VM version is termed the 

management VM. The OS is a customized eSLES version that is maintained by PowerFlex Engineering. 

Medium Granularity Pool In a Medium Granularity (MG) Pool, the volumes are divided into 1 MB allocation units, which are distributed and replicated across all disks which 

are contributing to the pool. 

Meta Data Manager (MDM) The MDM is the authority that controls and tracks data storage ownership, mapping, and protection. As volumes are created, the MDM provides the 

information application servers need to connect to the cluster’s virtualized storage. 

Mixed Deployment Hybrid hyperconverged deployment consists of hyperconverged, compute-only, and storage-only nodes. Some nodes contribute both compute 

resources and storage resources (hyperconverged nodes), some contribute only compute resources (compute-only nodes), and some contribute only storage resources (storage-only nodes). 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 235  |

 

Network Time Protocol (NTP) Network Time Protocol (NTP) is an Internet protocol that is used to synchronize with computer clock time sources in a network. 

Nodes Nodes are the basic hardware units that are used to install and run a hypervisor and PowerFlex software. 

PowerAPI PowerAPI (the PowerFlex REST API) allows users to automate PowerFlex deployment, configuration, and management tasks.  

PowerEdge R650 The PowerEdge R650 platform is a 15G single rack unit. As such, PowerEdge R650 provides great CPU core density for usage in a compute-only, or hyperconverged node use-case. PowerEdge R650 

provides ten drive slots for an overall maximum capacity of 76.8 TB. While more expensive than two rack unit systems, two R650 systems can provide a similar storage density to the R750 and R850. 

PowerFlex PowerFlex is an enterprise-class, software-defined block, and file storage solution that is deployed, managed, and supported as a single system. 

PowerFlex appliance PowerFlex appliance has a smaller starting point than PowerFlex rack, but scales to hundreds of nodes. PowerFlex appliance allows customers to use a broad set of supported networking options and can be added to 

existing networking infrastructures. PowerFlex appliance comes with licensing for a PowerFlex system and a unified management platform. 

PowerFlex Cluster A collection of multiple nodes, either SDS, SDC, or HCI, that communicates with each other to perform set of operations. 

PowerFlex custom node 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 236  |

 

PowerFlex custom nodes are servers that are configured to support a PowerFlex system. The servers are tested and certified to deliver predictable performance for a wide variety of workloads. The PowerFlex custom node offering is ideal for customers who prefer to build their own 

environments and have their own management services. PowerFlex custom node also allows broader networking options than the appliance. 

PowerFlex custom node comes with licensing for a PowerFlex system and a unified management platform. 

PowerFlex Device Local, direct attached block storage (DAS) in a node that is managed by an SDS and is contributed to a storage pool. 

PowerFlex large-scale configuration PowerFlex large-scale configuration can include tens to hundreds to thousands of nodes working together in a cluster to provide software-

defined storage. PowerFlex can scale to over 2,000 nodes that can drive up to 240 million IOPS of performance. 

PowerFlex management controller (PFMC) The infrastructure underlying the PowerFlex management controller in a dedicated multi-node configuration. It is mandatory in the PowerFlex rack, 

and optional for PowerFlex appliance. The physical controllers with ESXi hosting PFMP can be a three-node or a five-node ESXi cluster with vSAN (PFMC 1.0) or PowerFlex as shared storage (PFMC 2.0). 

PowerFlex management platform (PFMP) The management software stack for PowerFlex. It includes the Kubernetes/RKE environment running on physical or virtual Linux 

instances, and the containers that provide services. This management solution can be implemented on various infrastructure options. 

PowerFlex Manager 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 237  |

 

PowerFlex Manager automates deploying, configuring, managing, and upgrading a PowerFlex system. The PowerFlex Manager Dashboard provides an overview of the PowerFlex system. PowerFlex Manager allows administrators to quickly view the status of the system hardware 

and software components. 

PowerFlex rack A fully engineered rack-scale system with integrated networking. PowerFlex rack has a larger starting point, and scales to hundreds of 

nodes and racks. PowerFlex rack comes with licensing for a PowerFlex system and a unified management platform. 

PowerFlex System A PowerFlex system is the colle Management (MDM) cluster. n of entities managed by the Metadata ctio

Protection Domain A Protection Domain (PD) is a group of nodes or SDSs that provide data isolation, security, and performance benefits. A node can only participate 

in one PD. Separate PDs can be created for different node types with unequal performance profiles. 

Recovery Point Objective (RPO) The maximum amount of data, as measured by time that can be lost after a recovery from a disaster failure, or comparable event before data loss 

will exceed what is acceptable to an organization. 

Recovery Time Objective (RTO) The goal your organization sets for the maximum length of time it should take to restore normal operations following an outage or data loss. 

Release Certification Matrix (RCM) The PowerFlex Release Certification Matrix (RCM) is a definition of all compatible versions of hardware and software that are tested and 

validated together for a PowerFlex rack deployment. 

Secured Flag PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 238  |

 

The secure flag determines if the expiration field is active. 

Software Defined Storage Physical storage from aggregate nodes are pooled together to create volumes. Volumes are mapped to targets, allowing an application to have 

access to the storage. 

Storage Data Client (SDC) The Storage Data Client (SDC) is a block device driver that exposes shared block volumes from the SDS to the operating system. The SDC 

runs on the same server as the application. In practice, an application issues an I/O request, and the SDC fulfills the request regardless of what SDS the requested blocks physically reside on. 

Storage Data Replication (SDR) Extending the PowerFlex cluster beyond a single site is accomplished through the Storage Data Replication (SDR) component. SDR is an 

optional component responsible for managing all the aspects of PowerFlex replication. SDR is installed on an SDS node that contains storage media contributing to storage pools that have the potential to be replicated. 

Storage Data Server (SDS) The Storage Data Server (SDS) is a software daemon that enables a server in the cluster to contribute local storage devices to an aggregated 

storage pool. SDS owns the contributing devices, and together with the other SDSs, forms a protected mesh from which storage pools are created. 

Storage Data Target (SDT) The Storage Data Target (SDT) manages host connections and controllers that are connected over NVMe/TCP. SDT sends administrative 

and I/O commands forward to the SDS programmatically such that the SDS is oblivious to the source of the I/O. This makes traffic from an SDT look like it came from an SDC. 

Storage Pools PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 239  |

 

The Storage Pools are a subset of physical storage devices within a Protection Domain (PD). Each storage device belongs to a Storage Pool. 

Dell Technologies recommends having the same type of storage devices within a Storage Pool to ensure that the volumes are distributed across the same type of storage within the PD. 

Storage Virtualization Storage virtualization is the process of presenting a logical view of the physical storage resources to a host computer system. 

SupportAssist SupportAssist is a proactive monitoring software with automatic failure detection and notifications for Dell PCs, tablets, and servers. 

SupportAssist is free of charge, secure, and streamlines traditional support routines. 

SupportAssist SupportAssist is a proactive monitoring software with automatic failure detection and notifications for Dell PCs, tablets, and servers. 

SupportAssist is free of charge, secure, and streamlines traditional support routines. 

Two-Layer Deployment In a two-layer deployment, the SDS is installed on a separate host from the SDC. The front-end (client) is separated from the back-end (storage) 

data traffic.Two-layer deployments allow compute and storage resources to grow independently. PowerFlex compute-only nodes host end-user applications. PowerFlex storage-only nodes contribute storage to the system pool. 

VCSA The vCenter Server Appliance (VCSA) is a preconfigured Linux virtual machine, which is optimized for running VMware vCenter Server and the 

associated services on Linux. Using LogicMonitor's VMware VCSA package, you can monitor CPU usage, file system capacity, disk performance, memory, and much more. 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 240  |

 

Volume Analogous to a LUN, a volume is a subset of a storage pool’s capacity presented by an SDC as a local block device. A volume’s data is evenly 

distributed across all disks comprising a storage pool, according to the data layout selected for that storage pool. 

      

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 241  |

 

# Knowledge Check Answer Key 

## Knowledge Check - PowerFlex Offerings Case Study 1 

A customer is expanding their data center to include large scale virtual storage and compute capabilities. They have recently upgraded their networking infrastructure with the latest Dell switches. The customer wants to increase their compute, storage, and networking capacity with a 

simple to deploy, scalable multi-node solution that includes a unified management platform. They also want to be able to maintain their existing network infrastructure. 

Which PowerFlex offering best suits these requirements? a.  PowerFlex rack 

## b.**  **PowerFlex appliance

c.  PowerFlex custom node d.  PowerFlex software only 

*PowerFlex appliance allows the customer to get the benefits of the* *PowerFlex multi-node* *virtual storage and compute features and maintain* *their own network infrastructure.* 

## Knowledge Check - PowerFlex Offerings Case Study 2 

A customer has a large data center that has been recently updated with hundreds of new servers. They have learned about PowerFlex and want to leverage the scalability of software-defined storage on their existing servers without investing in new hardware. 

Which PowerFlex offering best suits these requirements? a.  PowerFlex rack 

b.  PowerFlex appliance c.  PowerFlex custom node 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 242  |

 

# d.**  **PowerFlex software** **only

*A PowerFlex software-only solution* *allows the customer to* *get the benefits* *of PowerFlex and maintain their existing infrastructure.* 

# Knowledge Check - PowerFlex Model Case Study 3 

A customer wants to build their own virtualized storage clusters. They are not interested in upgrading their network infrastructure or managing it through a shared interface. They do want hardware that Dell supports and is configured for PowerFlex. 

Which PowerFlex offering best suits these requirements? a.  PowerFlex rack 

b.  PowerFlex appliance 

# c.**  **PowerFlex custom node

d.  PowerFlex software only *PowerFlex custom node allows the customer to get the benefits of a* *PowerFlex solution* *on* *validated and supported hardware.* 

# Knowledge Check - PowerFlex Model Case Study 4 

A customer wants to switch their current data center over to virtual storage and compute capabilities. The current bare-metal data center infrastructure consists of stand-alone servers with internal storage. The data center has led to underperformance, limited scalability, and 

management complexity. The customer is looking for a software-defined storage and compute solution that has Dell provided fully integrated network devices and a unified management platform. The customer also wants the solution to give them room to scale as the company grows. 

Which PowerFlex offering best suits these requirements? 

# a.**  **PowerFlex rack

b.  PowerFlex appliance PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 243  |

 

c.  PowerFlex custom node d.  PowerFlex software only 

*PowerFlex rack allows the customer to get the benefits of a PowerFlex * *solution on* *a Dell fully integrated, configured, and supported system* *that * *can also* *scale or expand as needed.* 

# Knowledge Check - Hardware Case Study 1 

Platform Needs Analysis A customer is planning a two-layer PowerFlex implementation, with dedicated storage and compute nodes. Their applications run exclusively 

on AMD processors. 

Which PowerEdge platform best suits these requirements? a.  R650 

b.  R750 c.  R860 

# d.**  **R6525

*The R6525* *features AMD processors.* 

# Knowledge Check - Hardware Case Study 2 

Platform Needs Analysis A customer is planning a hyperconverged PowerFlex solution, where all nodes will perform both compute and storage roles. They want a single 

hardware platform for all nodes that can provide the greatest amount of computing power. The customer also wants the greatest storage-density and expansion capabilities. Cost is not a factor. 

Which PowerEdge platform best suits these requirements? a.  R650 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 244  |

 b.  R750 

# c.**  **R860

d.  R6525 *With four Intel Xeon®* *Scalable P* *compute power.* 

 

*ssors, the R860 will provide the most roce*

# Knowledge Check - Hardware Case Study 3 

Platform Needs Analysis A customer requires a compute-dense solution that minimizes the amount of rack space for their PowerFlex cluster. The nodes should be capable of 

providing hyperconverged functionality. Cost is not a factor. 

Which PowerEdge platform best suits these requirements? 

# a.**  **R650

b.  R750 c.  R860 

d.  R6525 *As a single rack server, the R650 takes up less space than the* *two rack * *servers.* 

# Knowledge Check - Hardware Case Study 4 

Platform Needs Analysis A customer is planning a two-layer PowerFlex implementation, with dedicated storage and compute nodes. They are looking for the most cost-

effective solution for their storage nodes, with room for expansion. Cost to storage density is their greatest driving factor in choosing a hardware platform for the storage nodes. 

Which PowerEdge platform best suits these requirements? 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 245  |

 

a.  R650 

# b.**  **R750

c.  R860 d.  R6525 

*The cost-to-storage density ratio of the R750* *makes this the* *best option* *for the customer.* 

# Knowledge Check - Deployment Case Study  

Platform Needs Analysis A customer is planning their PowerFlex implementation, leveraging VMware by Broadcom, and is primarily interested in simplified 

deployment. They also want to use as many resources on as many nodes as possible. 

Which PowerEdge deployment method best suits these requirements? 

# a.**  **Hyperconverged

b.  Mixed c.  Two-layer 

*Using hyperconverged simplifies the cluster architecture, while also* *using* *storage and compute resources on* *all nodes.* 

# Journal Configuration - Knowledge Check  

Customer Scenario: A customer has an application that generates about 2 GB/s of writes. Their maximum expected outage time is no more than two hours. The 

Storage Pool that is used by the application is 150 TB. 

Given these parameters, what percentage of the Storage Pool should be allocated for the Replication Journal? 

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 246  |

 

# a.**  **10%

b.  5% c.  12% 

d.  7% *The formula to calculate the journal capacity is: 2 GB/s * 14,400 s =* *14.4* *TB | (100 * 14.4 TB) / 150 TB = 9.6% (rounded up to 10%)*  

PowerFlex 4.x Administration-SSP 

| © Copyright 2024 Dell Inc |  Page 247  |

