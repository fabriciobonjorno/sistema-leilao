package auction

import (
	"context"
	"sync"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type fakeAuctionCollection struct {
	mu     sync.Mutex
	status auction_entity.AuctionStatus
	closed chan struct{}
}

func (c *fakeAuctionCollection) InsertOne(_ context.Context, document interface{}, _ ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status = document.(*AuctionEntityMongo).Status
	return &mongo.InsertOneResult{}, nil
}

func (c *fakeAuctionCollection) UpdateOne(_ context.Context, filter interface{}, update interface{}, _ ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	query := filter.(bson.M)
	if query["status"] == c.status {
		c.status = update.(bson.M)["$set"].(bson.M)["status"].(auction_entity.AuctionStatus)
		close(c.closed)
	}
	return &mongo.UpdateResult{}, nil
}

func (*fakeAuctionCollection) FindOne(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult {
	return nil
}

func (*fakeAuctionCollection) Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error) {
	return nil, nil
}

func TestCreateAuctionClosesAfterConfiguredDuration(t *testing.T) {
	collection := &fakeAuctionCollection{closed: make(chan struct{})}
	repository := &AuctionRepository{Collection: collection, auctionDuration: 20 * time.Millisecond}
	entity := &auction_entity.Auction{Id: "auction-1", Status: auction_entity.Active, Timestamp: time.Now()}

	if err := repository.CreateAuction(context.Background(), entity); err != nil {
		t.Fatalf("CreateAuction() error = %v", err)
	}
	collection.mu.Lock()
	status := collection.status
	collection.mu.Unlock()
	if status != auction_entity.Active {
		t.Fatalf("status immediately after creation = %v, want Active", status)
	}

	select {
	case <-collection.closed:
	case <-time.After(time.Second):
		t.Fatal("auction was not closed after configured duration")
	}
	collection.mu.Lock()
	status = collection.status
	collection.mu.Unlock()
	if status != auction_entity.Completed {
		t.Fatalf("status after duration = %v, want Completed", status)
	}
}
