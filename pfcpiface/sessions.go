// SPDX-License-Identifier: Apache-2.0
// Copyright(c) 2020 Intel Corporation

package main

import (
	"sync"

	"github.com/omec-project/upf-epc/pfcpiface/metrics"
)

type notifyFlag struct {
	flag bool
	mux  sync.Mutex
}

// PFCPSession implements one PFCP session.
type PFCPSession struct {
	localSEID        uint64
	remoteSEID       uint64
	notificationFlag notifyFlag
	pdrs             []pdr
	fars             []far
	qers             []qer
	metrics          *metrics.Session
}

// NewPFCPSession allocates an session with ID.
func (pConn *PFCPConn) NewPFCPSession(rseid uint64) uint64 {
	for i := 0; i < pConn.maxRetries; i++ {
		lseid := pConn.rng.Uint64()
		// Check if it already exists
		if _, ok := pConn.sessions[lseid]; ok {
			continue
		}

		s := PFCPSession{
			localSEID:  lseid,
			remoteSEID: rseid,
			pdrs:       make([]pdr, 0, MaxItems),
			fars:       make([]far, 0, MaxItems),
			qers:       make([]qer, 0, MaxItems),
		}
		pConn.sessions[lseid] = &s

		// Metrics update
		s.metrics = metrics.NewSession(pConn.nodeID.remote)
		pConn.SaveSessions(s.metrics)

		return lseid
	}

	return 0
}

// RemoveSession removes session using lseid.
func (pConn *PFCPConn) RemoveSession(lseid uint64) {
	s, ok := pConn.sessions[lseid]
	if !ok {
		return
	}

	// Metrics update
	s.metrics.Delete()
	pConn.SaveSessions(s.metrics)

	delete(pConn.sessions, lseid)
}
