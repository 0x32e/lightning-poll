package votes

import (
	"context"
	"encoding/hex"
	"log"
	"time"

	votes_db "lightning-poll/votes/internal/db/votes"
	"lightning-poll/votes/internal/types"
)

func Start(b Backends) {
	go expireVotesForever(b)
}

func expireVotesForever(b Backends) {
	for {
		if err := expireVotes(context.Background(), b); err != nil {
			log.Printf("votes/ops: expireVotesForever error: %v", err)
		}
		time.Sleep(time.Minute * 5)
	}
}

func expireVotes(ctx context.Context, b Backends) error {
	expired, err := votes_db.ListExpired(ctx, b.GetDB())
	if err != nil {
		return err
	}

	if len(expired) == 0 {
		return nil
	}

	for _, exp := range expired {
		// if the invoice has been settled, update vote to paid
		inv, err := b.GetLND().LookupInvoice(ctx, exp.PayReq)
		if err != nil {
			return err
		}

		if inv.SettleIndex > 0 {
			if err := markInvoicePaid(ctx, b, hex.EncodeToString(inv.RHash),
				inv.AmtPaidSat, inv.SettleIndex); err != nil {
				return err
			}
		}

		// if the invoice has not been paid, expire it
		if err := votes_db.UpdateStatus(ctx, b.GetDB(), exp.ID, types.VoteStatusCreated,
			types.VoteStatusExpired); err != nil {
			return err
		}
	}
	return nil
}

// markInvoicePaid marks an invoice as paid, so that it can be settled or released in future
func markInvoicePaid(ctx context.Context, b Backends, payHash string, settledAmount int64, settleIndex uint64) error {
	vote, err := votes_db.LookupByHash(ctx, b.GetDB(), payHash)
	if err != nil {
		return err
	}

	return votes_db.MarkPaid(ctx, b.GetDB(), vote.ID, settledAmount, settleIndex)
}
