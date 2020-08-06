package config

import "syscall"

// SIGUSR1 linux SIGUSR1
const SIGUSR1 = syscall.Signal(0xa)

// SIGUSR2 linux SIGUSR2
const SIGUSR2 = syscall.Signal(0xc)

const (
	FileStateType = iota + 1
	FileReadType
	FileSyncAttr
)
const ConnPath = "/v1/conn"

const ClientHeartTime = 4 * 60

const ClientRetryConnectTime = 10

const ServerHeartTime = 4 * 60

const ServerRetryConnectTime = 10
