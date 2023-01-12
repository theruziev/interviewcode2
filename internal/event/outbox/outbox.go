package outbox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/theruziev/oson_auth/internal/db"
	"github.com/theruziev/oson_auth/internal/event/constants"
	"github.com/theruziev/oson_auth/internal/model"
	"github.com/theruziev/oson_auth/internal/pkg/logging"
	"github.com/theruziev/oson_auth/internal/pkg/rabbitmqx"
	"github.com/wagslane/go-rabbitmq"
)

const ticker = 50 * time.Millisecond

type Outbox struct {
	producer *rabbitmq.Publisher
	store    *db.OutBoxStore
}

func NewOutBox(producer *rabbitmq.Publisher, store *db.OutBoxStore) *Outbox {
	return &Outbox{
		producer: producer,
		store:    store,
	}
}

func (o *Outbox) Serve(ctx context.Context) {
	logger := logging.FromContext(ctx)
	ticker := time.NewTicker(ticker)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := o.store.ProcessMessage(ctx, o.processMessage); err != nil {
				logger.Errorf("failed to process message: %s", err)
			}
		}
	}
}

func (o *Outbox) processMessage(_ context.Context, m *model.OutBox) error {
	msgBytes, err := json.Marshal(m.Data)
	if err != nil {
		return err
	}
	err = o.producer.Publish(
		msgBytes,
		[]string{m.Topic},
		rabbitmq.WithPublishOptionsMessageID(fmt.Sprintf("%d", m.ID)),
		rabbitmq.WithPublishOptionsExchange(constants.ExchangeUser),
		rabbitmqx.WithPublishJSONContentType(),
	)

	return err
}
