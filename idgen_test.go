package idgen

import (
	convey "github.com/smartystreets/goconvey/convey"
	"github.com/sumory/baseN4go"
	"testing"
	"time"
)

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func TestNotErr(t *testing.T) {
	convey.Convey("should not err", t, func() {
		err, idWorker := NewIdWorker(1)
		if err != nil {
			t.Fatal("can not initialize")
		}
		for i := 0; i < 10; i++ {
			err, _ := idWorker.NextId()
			convey.So(err, convey.ShouldBeNil)
		}
	})

}

func TestRabaseShortRadix(t *testing.T) {
	convey.Convey("test method 'RabaseShortRadix', functions should be right after rebase", t, func() {
			err, idWorker := NewIdWorker(1)
			if err != nil {
				t.Fatal("can not initialize")
			}
			err, baseN := baseN4go.NewBaseN(int8(62))
			convey.So(err, convey.ShouldBeNil)
			convey.So(baseN, convey.ShouldNotBeNil)

			err, newId := idWorker.NextId()
			convey.So(err, convey.ShouldBeNil)
			convey.So(newId, convey.ShouldNotBeNil)
			err, shortId := idWorker.ShortenId(newId)
			convey.So(err, convey.ShouldBeNil)

			_, baseNId := baseN.Encode(newId)
			convey.So(shortId, convey.ShouldEqual, baseNId)

			idWorker.RabaseShortRadix(int8(16))
			err, newId = idWorker.NextId()
			err, shortId = idWorker.ShortenId(newId)

			err, baseN = baseN4go.NewBaseN(int8(16))
			_, baseNId = baseN.Encode(newId)
			convey.So(shortId, convey.ShouldEqual, baseNId)
		})
}

func TestContainsWorkerId(t *testing.T) {
	convey.Convey("generated id should contains the `workerId`", t, func() {
		var workerId int64
		for workerId = 0; workerId < 1024; workerId++ {
			err, idWorker := NewIdWorker(workerId)
			if err != nil {
				t.Fatal("can not initialize")
			}
			for i := 0; i < 1; i++ {
				_, newId := idWorker.NextId()
				newId2 := uint(newId << 42)
				newId3 := uint(newId2 >> 54)
				//fmt.Printf("%d~%b~%b~%b\n",workerId,workerId,newId, newId3);
				convey.So(newId3, convey.ShouldEqual, workerId)
			}
		}
	})
}

func TestWorkerId(t *testing.T) {
	convey.Convey("method 'WorkerId' should return the right workerId", t, func() {
		var workerId int64
		for workerId = 0; workerId < 1024; workerId++ {
			err, idWorker := NewIdWorker(workerId)
			if err != nil {
				t.Fatal("can not initialize")
			}
			for i := 0; i < 1; i++ {
				_, newId := idWorker.NextId()
				wId := idWorker.WorkerId(newId)
				//fmt.Println(workerId,wId)
				convey.So(wId, convey.ShouldEqual, workerId)
			}
		}
	})
}

func TestShortenId(t *testing.T) {
	convey.Convey("method 'ShortenId' should return the right shorten id", t, func() {
		err, idWorker := NewIdWorker(1)
		if err != nil {
			t.Fatal("can not initialize")
		}
		err, baseN := baseN4go.NewBaseN(int8(62))
		convey.So(err, convey.ShouldBeNil)
		convey.So(baseN, convey.ShouldNotBeNil)

		err, newId := idWorker.NextId()
		convey.So(err, convey.ShouldBeNil)
		convey.So(newId, convey.ShouldNotBeNil)
		err, shortId := idWorker.ShortenId(newId)
		convey.So(err, convey.ShouldBeNil)

		_, baseNId := baseN.Encode(newId)
		convey.So(shortId, convey.ShouldEqual, baseNId)
	})
}

func TestShortId(t *testing.T) {
	convey.Convey("method 'ShortId' should return the right result", t, func() {
			err, idWorker := NewIdWorker(1)
			if err != nil {
				t.Fatal("can not initialize")
			}

			err, newIdStr := idWorker.ShortId()
			convey.So(err, convey.ShouldBeNil)
			convey.So(newIdStr, convey.ShouldHaveSameTypeAs,"abc")

			err, baseN := baseN4go.NewBaseN(int8(62))
			_,newId:=baseN.Decode(newIdStr)
			convey.So(newId, convey.ShouldHaveSameTypeAs,int64(1))
		})
}
