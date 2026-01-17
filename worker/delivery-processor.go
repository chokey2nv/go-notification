package worker

import (
	"context"
	"time"

	"github.com/chokey2nv/go-notification/dispatcher"
	"github.com/chokey2nv/go-notification/domain"
	"github.com/chokey2nv/go-notification/repository"
	"github.com/google/uuid"
)

//	type DeliveryProcessor struct {
//		deliveries    repository.DeliveryRepository
//		notifications repository.NotificationRepository
//		dispatchers   map[domain.Channel]dispatcher.Dispatcher
//		maxRetries    int
//	}
type DeliveryProcessor struct {
	deliveries    repository.DeliveryRepository
	notifications repository.NotificationRepository
	dlq           repository.DeadLetterRepository
	dispatchers   map[domain.Channel]dispatcher.Dispatcher
	maxRetries    int
}

func NewDeliveryProcessor(
	deliveries repository.DeliveryRepository,
	notifications repository.NotificationRepository,
	dlq repository.DeadLetterRepository,
	dispatchers []dispatcher.Dispatcher,
) *DeliveryProcessor {

	m := make(map[domain.Channel]dispatcher.Dispatcher)
	for _, d := range dispatchers {
		m[d.Channel()] = d
	}

	return &DeliveryProcessor{
		deliveries:    deliveries,
		notifications: notifications,
		dlq:           dlq,
		dispatchers:   m,
		maxRetries:    5,
	}
}
func (p *DeliveryProcessor) RunOnce(ctx context.Context) error {
	pending, err := p.deliveries.GetPending(ctx, 50)
	if err != nil {
		return err
	}

	for _, d := range pending {
		p.processOne(ctx, d)
	}

	return nil
}
func (p *DeliveryProcessor) processOne(
	ctx context.Context,
	d *domain.NotificationDelivery,
) {

	// Guard: respect scheduling
	if time.Now().Before(d.NextAttemptAt) {
		return
	}

	// Load notification
	n, err := p.notifications.GetByID(ctx, d.NotificationID)
	if err != nil {
		d.Status = domain.DeliveryFailed
		d.LastError = err.Error()
		_ = p.deliveries.Update(ctx, d)

		p.sendToDLQ(
			ctx,
			d,
			"notification not found",
			nil,
		)
		return
	}

	// Dispatcher lookup (NON-RETRYABLE)
	dispatcher, ok := p.dispatchers[d.Channel]
	if !ok {
		d.Status = domain.DeliveryFailed
		d.LastError = "dispatcher not registered"
		_ = p.deliveries.Update(ctx, d)

		p.sendToDLQ(
			ctx,
			d,
			"dispatcher missing",
			n,
		)
		return
	}

	// Attempt send
	err = dispatcher.Send(ctx, n)
	if err == nil {
		d.Status = domain.DeliverySent
		d.UpdatedAt = time.Now()
		_ = p.deliveries.Update(ctx, d)
		return
	}

	// Failure path (RETRYABLE)
	d.Attempts++
	d.LastError = err.Error()
	d.UpdatedAt = time.Now()

	// MAX RETRIES EXCEEDED → DLQ
	if d.Attempts >= p.maxRetries {
		d.Status = domain.DeliveryFailed
		_ = p.deliveries.Update(ctx, d)

		p.sendToDLQ(
			ctx,
			d,
			"max retries exceeded: "+d.LastError,
			n,
		)
		return
	}

	// Schedule retry
	d.Status = domain.DeliveryRetrying
	d.NextAttemptAt = backoff(d.Attempts)
	_ = p.deliveries.Update(ctx, d)
}

func (p *DeliveryProcessor) sendToDLQ(
	ctx context.Context,
	d *domain.NotificationDelivery,
	reason string,
	n *domain.Notification,
) {

	dlq := &domain.DeadLetter{
		ID:             uuid.NewString(),
		DeliveryID:     d.ID,
		NotificationID: d.NotificationID,
		Channel:        d.Channel,
		Reason:         reason,
		Payload:        n.Metadata,
		CreatedAt:      time.Now(),
	}

	_ = p.dlq.Create(ctx, dlq)
}

func backoff(attempt int) time.Time {
	delay := time.Duration(1<<attempt) * time.Minute
	if delay > 30*time.Minute {
		delay = 30 * time.Minute
	}
	return time.Now().Add(delay)
}

func ReplayDLQ(
	ctx context.Context,
	dlq *domain.DeadLetter,
	deliveryRepo repository.DeliveryRepository,
	dlqRepo repository.DeadLetterRepository,
) error {

	newDelivery := &domain.NotificationDelivery{
		ID:             uuid.NewString(),
		NotificationID: dlq.NotificationID,
		Channel:        dlq.Channel,
		Status:         domain.DeliveryPending,
		Attempts:       0,
		NextAttemptAt:  time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := deliveryRepo.Create(ctx, newDelivery); err != nil {
		return err
	}

	return dlqRepo.Delete(ctx, dlq.ID)
}
