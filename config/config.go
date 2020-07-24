package config

import "syscall"

// SIGUSR1 linux SIGUSR1
const SIGUSR1 = syscall.Signal(0xa)

// SIGUSR2 linux SIGUSR2
const SIGUSR2 = syscall.Signal(0xc)