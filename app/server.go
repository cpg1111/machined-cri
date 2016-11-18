package app

import (
	"errors"
	"fmt"
	"log"

	"k8s.io/kubernetes/cmd/kubelet/app/options"
	"k8s.io/kubernetes/pkg/util/flock"

	"github.com/cpg1111/nspawnlet/nspawnlet"
)

// Run runs the nspawnlet
func Run(s *options.KubeletServer, kubeDeps *nspawnlet.Deps) error {
	err := run(s, kubeDeps)
	if err != nil {
		return fmt.Errorf("failed torun nspawnlet: %v", err)
	}
	return nil
}

func run(s *options.KubeletServer, kubeDeps *nspawnlet.Deps) error {
	if s.ExitOnLockContention && s.LockFilePath == "" {
		return errors.New("cannot exit on lock file contention: no lock file specified")
	}
	done := make(chan struct{})
	if s.LockFilePath != "" {
		log.Printf("acquiring file lock on %q\n", s.LockFilePath)
		lockErr := flock.Acquire(s.LockFilePath)
		if lockErr != nil {
			return fmt.Errorf("unable to acquire file lock on %q: %v", s.LockFilePath, lockErr)
		}
		if s.ExitOnLockContention {
			log.Printf("watching for inotify events for: %v\n", s.LockFilePath)
			watchErr := watchForLockfileContention(s.LockFilePath, done)
			if watchErr != nil {
				return watchErr
			}
		}
	}
	ftErr := utilconfig.DefaultFeatureGate.Set(s.KubeletConfiguration.FeatureGates)
	if ftErr != nil {
		return ftErr
	}
	cfgz, cfgzErr := initConfigz(&s.KubeletConfugration)
	if utilconfig.DefaultFeatureGate.DynamicKubeletConfig() {
		if s.RunOnce == false {
			remoteKC, err := initKubeletConfigSync(s)
			if err == nil {
				s.KubeletConfiguration = *remoteKC
				if cfgzErr != nil {
					log.Println(fmt.Errorf("was unable to register configz before due to %s, will not be able to set now", cfgzErr))
				} else {
					setConfigz(cfgz, &s.KubeletConfiguration)
				}
				err = utilconfig.DefaultFeatureGate.Set(s.KubeletConfiguration.FeatureGates)
				if err != nil {
					return err
				}
			}
		}
	}
	if s.StandAlone {

	}
	return nil
}
