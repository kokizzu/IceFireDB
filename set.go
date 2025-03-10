package main

import (
	"time"

	"github.com/ledisdb/ledisdb/ledis"
	"github.com/tidwall/redcon"
	"github.com/tidwall/uhaha"
)

func init() {
	conf.AddWriteCommand("SADD", cmdSADD)
	conf.AddReadCommand("SCARD", cmdSCARD)
	conf.AddReadCommand("SDIFF", cmdSDIFF)
	conf.AddWriteCommand("SDIFFSTORE", cmdSDIFFSTORE)
	conf.AddReadCommand("SINTER", cmdSINTER)
	conf.AddWriteCommand("SINTERSTORE", cmdSINTERSTORE)
	conf.AddReadCommand("SISMEMBER", cmdSISMEMBER)
	conf.AddReadCommand("SMEMBERS", cmdSMEMBERS)
	conf.AddWriteCommand("SREM", cmdSREM)
	conf.AddReadCommand("SUNION", cmdSUNION)
	conf.AddWriteCommand("SUNIONSTORE", cmdSUNIONSTORE)
	conf.AddWriteCommand("SCLEAR", cmdSCLEAR)
	conf.AddWriteCommand("SMCLEAR", cmdSMCLEAR)
	conf.AddWriteCommand("SEXPIRE", cmdSEXPIRE)
	conf.AddWriteCommand("SEXPIREAT", cmdSEXPIREAT)
	conf.AddReadCommand("STTL", cmdSTTL)
	conf.AddWriteCommand("SPERSIST", cmdSPERSIST)
	conf.AddReadCommand("SKEYEXISTS", cmdSKEYEXISTS)
}

func cmdSADD(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) < 3 {
		return nil, uhaha.ErrWrongNumArgs
	}
	n, err := ldb.SAdd([]byte(args[1]), stringSliceToBytes(args[2:])...)
	return redcon.SimpleInt(n), err
}

func soptGeneric(args [][]byte, optType byte) ([][]byte, error) {
	if len(args) < 2 {
		return nil, uhaha.ErrWrongNumArgs
	}

	var v [][]byte
	var err error
	switch optType {
	case ledis.UnionType:
		v, err = ldb.SUnion(args[1:]...)
	case ledis.DiffType:
		v, err = ldb.SDiff(args[1:]...)
	case ledis.InterType:
		v, err = ldb.SInter(args[1:]...)
	}
	if err != nil {
		return nil, err
	}
	return v, nil
}

func soptStoreGeneric(args [][]byte, optType byte) (interface{}, error) {
	if len(args) < 3 {
		return 0, uhaha.ErrWrongNumArgs
	}

	var n int64
	var err error

	switch optType {
	case ledis.UnionType:
		n, err = ldb.SUnionStore(args[1], args[2:]...)
	case ledis.DiffType:
		n, err = ldb.SDiffStore(args[1], args[2:]...)
	case ledis.InterType:
		n, err = ldb.SInterStore(args[1], args[2:]...)
	}
	return redcon.SimpleInt(n), err
}

func cmdSCARD(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 2 {
		return nil, uhaha.ErrWrongNumArgs
	}

	n, err := ldb.SCard([]byte(args[1]))
	return redcon.SimpleInt(n), err
}

func cmdSDIFF(m uhaha.Machine, args []string) (interface{}, error) {
	return soptGeneric(stringSliceToBytes(args), ledis.DiffType)
}

func cmdSDIFFSTORE(m uhaha.Machine, args []string) (interface{}, error) {
	return soptStoreGeneric(stringSliceToBytes(args), ledis.DiffType)
}

func cmdSINTER(m uhaha.Machine, args []string) (interface{}, error) {
	return soptGeneric(stringSliceToBytes(args), ledis.InterType)
}

func cmdSINTERSTORE(m uhaha.Machine, args []string) (interface{}, error) {
	return soptStoreGeneric(stringSliceToBytes(args), ledis.InterType)
}

func cmdSISMEMBER(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 3 {
		return nil, uhaha.ErrWrongNumArgs
	}
	n, err := ldb.SIsMember([]byte(args[1]), []byte(args[2]))

	return redcon.SimpleInt(n), err
}

func cmdSMEMBERS(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 2 {
		return nil, uhaha.ErrWrongNumArgs
	}
	return ldb.SMembers([]byte(args[1]))
}

func cmdSREM(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) < 3 {
		return nil, uhaha.ErrWrongNumArgs
	}

	n, err := ldb.SRem([]byte(args[1]), stringSliceToBytes(args[2:])...)
	return redcon.SimpleInt(n), err
}

func cmdSUNION(m uhaha.Machine, args []string) (interface{}, error) {
	return soptGeneric(stringSliceToBytes(args), ledis.UnionType)
}

func cmdSUNIONSTORE(m uhaha.Machine, args []string) (interface{}, error) {
	return soptStoreGeneric(stringSliceToBytes(args), ledis.UnionType)
}

func cmdSCLEAR(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 2 {
		return nil, uhaha.ErrWrongNumArgs
	}

	n, err := ldb.SClear([]byte(args[1]))

	return redcon.SimpleInt(n), err
}

func cmdSMCLEAR(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) < 2 {
		return nil, uhaha.ErrWrongNumArgs
	}

	n, err := ldb.SMclear(stringSliceToBytes(args[1:])...)

	return redcon.SimpleInt(n), err
}

func cmdSEXPIRE(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 3 {
		return nil, uhaha.ErrWrongNumArgs
	}

	duration, err := ledis.StrInt64([]byte(args[2]), nil)
	if err != nil {
		return nil, uhaha.ErrInvalid
	}

	timestamp := m.Now().Unix() + duration
	if timestamp < time.Now().Unix() {
		_, err := ldb.SClear([]byte(args[1]))
		if err != nil {
			return nil, err
		}
		return redcon.SimpleInt(0), nil
	}

	v, err := ldb.SExpireAt([]byte(args[1]), timestamp)
	return redcon.SimpleInt(v), err
}

func cmdSEXPIREAT(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 3 {
		return nil, uhaha.ErrWrongNumArgs
	}

	timestamp, err := ledis.StrInt64([]byte(args[2]), nil)
	if err != nil {
		return nil, uhaha.ErrInvalid
	}
	if timestamp < time.Now().Unix() {
		_, err := ldb.SClear([]byte(args[1]))
		if err != nil {
			return nil, err
		}
		return redcon.SimpleInt(0), nil
	}

	v, err := ldb.SExpireAt([]byte(args[1]), timestamp)
	return redcon.SimpleInt(v), err
}

func cmdSTTL(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 2 {
		return nil, uhaha.ErrWrongNumArgs
	}

	v, err := ldb.STTL([]byte(args[1]))
	return redcon.SimpleInt(v), err
}

func cmdSPERSIST(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 2 {
		return nil, uhaha.ErrWrongNumArgs
	}

	n, err := ldb.SPersist([]byte(args[1]))
	return redcon.SimpleInt(n), err
}

func cmdSKEYEXISTS(m uhaha.Machine, args []string) (interface{}, error) {
	if len(args) != 2 {
		return nil, uhaha.ErrWrongNumArgs
	}

	n, err := ldb.SKeyExists([]byte(args[1]))
	return redcon.SimpleInt(n), err
}

func stringSliceToBytes(args []string) [][]byte {
	bs := make([][]byte, len(args))
	for k, v := range args {
		bs[k] = []byte(v)
	}
	return bs
}
