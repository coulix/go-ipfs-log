package test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	idp "github.com/berty/go-ipfs-log/identityprovider"
	io "github.com/berty/go-ipfs-log/io"
	ks "github.com/berty/go-ipfs-log/keystore"
	"github.com/berty/go-ipfs-log/log"
	ds "github.com/ipfs/go-datastore"
	dssync "github.com/ipfs/go-datastore/sync"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLog(t *testing.T) {
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ipfs := io.NewMemoryServices()

	datastore := dssync.MutexWrap(ds.NewMapDatastore())
	keystore, err := ks.NewKeystore(datastore)
	if err != nil {
		panic(err)
	}

	idProvider := idp.NewOrbitDBIdentityProvider(keystore)

	var identities []*idp.Identity

	for i := 0; i < 4; i++ {
		identity, err := idProvider.GetID(fmt.Sprintf("User%d", i))
		if err != nil {
			panic(err)
		}

		identities = append(identities, identity)
	}

	Convey("Log", t, FailureHalts, func(c C) {
		c.Convey("sets an id and a clock id", FailureHalts, func(c C) {
			log1 := log.NewLog(ipfs, identities[0], &log.NewLogOptions{ID: "A"})
			c.So(log1.ID, ShouldEqual, "A")
			c.So(log1.Clock.ID.Equals(identities[0].PublicKey), ShouldBeTrue)
		})

		c.Convey("set time.now as id string if id is not passed as an argument", FailureHalts, func(c C) {
			before := time.Now().Unix() / 1000
			log1 := log.NewLog(ipfs, identities[0], &log.NewLogOptions{})
			after := time.Now().Unix() / 1000

			logid, err := strconv.ParseInt(log1.ID, 10, 64)
			c.So(err, ShouldBeNil)
			c.So(logid, ShouldBeGreaterThanOrEqualTo, before)
			c.So(logid, ShouldBeLessThanOrEqualTo, after)
		})
	})
}
