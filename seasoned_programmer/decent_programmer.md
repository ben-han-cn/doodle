# Decent Programmer
* [Language](#language)
* [OS](#os)
* [Storage](#storage)
* [Network](#network)
* [Architecture](#architecture)


## Language<a id="language"></a>
1. Paradigm
1. Type System
1. Data structure
1. Memory management
1. Concurrency model

## OS<a id="os"></a>
### Linux
1. Linux API
  It composed out of the system call interface of the linux kernel and subroutines
  in GNU c library.
1. Linux ABI
  Application binary interface refers to the compiled binaries in machine code, therefore
  ABI is bound to the intruction set. It's a task for linux distributions. An ABI has 
  to defined for very instruction set such as x86, x86-64, MIPS, ARMv7 with endianness.

## Storage<a id="storage"></a>
### Block device
### File system
### Database
1. SQL
1. NOSQL
1. Data Structure
 * B+ tree
 * Log-Structured Merge Tree(LSM)
  + RocksDB
    - Memtable, Transaction log, Static Sorted Table(SST)
    - Key, Value are sorted and has associated timestamp
    - Always append new operation(update or delete), without update old record.
    - Read will go through Memtable and SST then merged based on timestamp
    - Use bloom filter to decide whether a SST file contains a given key
    - Backend compact process to merge the SST files then free the space
 * Log-Structured Hash-table
## Network<a id="network"></a>
* TCP/IP
* Async programming
* Virtualization
* Mobile network

## Architecture<a id="architecture"></a>
