// Package fastid is a distributed, k-ordered unique ID generator.
//  Under 64 bits (Long Integer)
//  Lock-free (using atomic CAS)
//  Decentralized and no coordination needed
//  Docker friendly
package fastid

import (
	"errors"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

const (
	//StartTimeEnvName is the env key for ID generating start time
	StartTimeEnvName = "FASTID_START_TIME"
	//MachineIDEnvName is the env key for machine id
	MachineIDEnvName           = "FASTID_MACHINE_ID"
	defaultStartTimeStr        = "2018-06-01T00:00:00.000Z"
	defaultStartTimeNano int64 = 1527811200000000000
)

//Config maintains the settings for id generating
type Config struct {
	timeBits      uint
	seqBits       uint
	machineBits   uint
	timeMask      int64
	seqMask       int64
	machineID     int64
	machineIDMask int64
	lastID        int64
}

//ConstructConfig creates an instance of FastIDConfig with the given settings
func ConstructConfig(timeBits, seqBits, machineBits uint) *Config {
	return ConstructConfigWithMachineID(timeBits, seqBits, machineBits, getMachineID())
}

//ConstructConfigWithMachineID creates an config with machine id, in case you don't want to use the lower 16 bits of the IP address.
func ConstructConfigWithMachineID(timeBits, seqBits, machineBits uint, machineID int64) *Config {
	machineIDMask := ^(int64(-1) << machineBits)
	return &Config{
		timeBits:      timeBits,
		seqBits:       seqBits,
		machineBits:   machineBits,
		timeMask:      ^(int64(-1) << timeBits),
		seqMask:       ^(int64(-1) << seqBits),
		machineIDMask: machineIDMask,
		machineID:     machineID & machineIDMask,
		lastID:        0,
	}
}

// BenchmarkConfig is a high performance setting for benchmark
//  40 bits timestamp
//  15 bits seq
//  8  bits machine id
var BenchmarkConfig = ConstructConfig(40, 15, 8)

// CommonConfig is the recommended setting for most applications
//  40 bits timestamp
//  7  bits seq
//  16 bits machine id
var CommonConfig = ConstructConfig(40, 2, 16)

var startEpochNano = getStartEpochFromEnv()

func (c *Config) getCurrentTimestamp() int64 {
	//devided by 2^20 (~10^6, nano to milliseconds)
	return (time.Now().UnixNano() - startEpochNano) >> 20 & c.timeMask
}

//GenInt64ID generates unique int64 IDs with the setting in the methond owner
func (c *Config) GenInt64ID() int64 {
	for {
		localLastID := atomic.LoadInt64(&c.lastID)
		seq := c.GetSeqFromID(localLastID)
		lastIDTime := c.GetTimeFromID(localLastID)
		now := c.getCurrentTimestamp()
		if now > lastIDTime {
			seq = 0
		} else if seq >= c.seqMask {
			time.Sleep(time.Duration(0xFFFFF - (time.Now().UnixNano() & 0xFFFFF)))
			continue
		} else {
			seq++
		}

		newID := now<<(c.machineBits+c.seqBits) + seq<<c.machineBits + c.machineID
		if atomic.CompareAndSwapInt64(&c.lastID, localLastID, newID) {
			return newID
		}
		time.Sleep(time.Duration(20))
	}
}

//GetSeqFromID extracts seq number from an existing ID
func (c *Config) GetSeqFromID(id int64) int64 {
	return (id >> c.machineBits) & c.seqMask
}

//GetTimeFromID extracts timestamp from an existing ID
func (c *Config) GetTimeFromID(id int64) int64 {
	return id >> (c.machineBits + c.seqBits)
}

func getMachineID() int64 {
	//getting machine from env
	if machineIDStr, ok := os.LookupEnv(MachineIDEnvName); ok {
		if machineID, err := strconv.ParseInt(machineIDStr, 10, 64); err == nil {
			return machineID
		}
	}
	//take the lower 16bits of IP address as Machine ID
	if ip, err := getIP(); err == nil {
		return (int64(ip[2]) << 8) + int64(ip[3])
	}
	return 0
}

func getStartEpochFromEnv() int64 {
	startTimeStr := getEnv(StartTimeEnvName, defaultStartTimeStr)
	var startEpochTime, err = time.Parse(time.RFC3339, startTimeStr)

	if err == nil {
		return defaultStartTimeNano
	}

	return startEpochTime.UnixNano()
}

func getIP() (net.IP, error) {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					ip := ipNet.IP.To4()

					if ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168 {
						return ip, nil
					}
				}
			}
		}
	}
	return nil, errors.New("Failed to get ip address")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
