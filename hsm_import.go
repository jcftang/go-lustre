package lustre

//
// #cgo LDFLAGS: -llustreapi
// #include <sys/types.h>
// #include <sys/stat.h>
// #include <unistd.h>
// #include <lustre/lustreapi.h>
// #include <stdlib.h>
//
import "C"

import (
	"errors"
	"log"
	"os"
	"syscall"
)

var errStatError = errors.New("stat failure")

func stat_to_cstat(fi os.FileInfo) *C.struct_stat {
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		log.Printf("no stat info")
		return nil
	}

	var st C.struct_stat
	st.st_uid = C.__uid_t(stat.Uid)
	st.st_gid = C.__gid_t(stat.Gid)
	st.st_mode = C.__mode_t(stat.Mode)
	st.st_size = C.__off_t(stat.Size)
	st.st_mtim.tv_sec = C.__time_t(stat.Mtim.Sec)
	st.st_mtim.tv_nsec = C.long(stat.Mtim.Nsec)
	st.st_atim.tv_sec = C.__time_t(stat.Atim.Sec)
	st.st_atim.tv_nsec = C.long(stat.Atim.Nsec)

	return &st
}

func HsmImport(
	f string,
	archive uint,
	fi os.FileInfo,
	stripe_size uint64,
	stripe_offset int,
	stripe_count int,
	stripe_pattern int,
	pool_name string) (*Fid, error) {

	var fid Fid

	st := stat_to_cstat(fi)
	if st == nil {
		return nil, errStatError
	}

	rc, err := C.llapi_hsm_import(
		C.CString(f),
		C.int(archive),
		st,
		C.ulonglong(stripe_size),
		C.int(stripe_offset),
		C.int(stripe_count),
		C.int(stripe_pattern),
		nil,
		(*_Ctype_lustre_fid)(&fid),
	)
	if rc < 0 {
		return nil, err
	}
	return &fid, nil
}