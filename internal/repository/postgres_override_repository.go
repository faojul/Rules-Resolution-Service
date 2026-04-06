package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"rules-resolution-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOverrideRepository struct {
    db *pgxpool.Pool
}

func NewPostgresOverrideRepository(db *pgxpool.Pool) *PostgresOverrideRepository {
    return &PostgresOverrideRepository{db: db}
}

func (r *PostgresOverrideRepository) FindByStepAndTrait(ctx context.Context, step, trait string) ([]domain.Override, error) {

    query := `
    SELECT id, step_key, trait_key, selector,
           value, specificity, effective_date, expires_date, status
    FROM overrides
    WHERE step_key=$1 AND trait_key=$2
    `

    rows, err := r.db.Query(ctx, query, step, trait)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []domain.Override

    for rows.Next() {
        var o domain.Override
        err := rows.Scan(
            &o.ID,
            &o.StepKey,
            &o.TraitKey,
            &o.Selector,
            &o.Value,
            &o.Specificity,
            &o.EffectiveDate,
            &o.ExpiresDate,
            &o.Status,
        )
        if err != nil {
            return nil, err
        }
        result = append(result, o)
    }

    return result, nil
}

func (r *PostgresOverrideRepository) FindAllOverrides(ctx context.Context) ([]domain.Override, error) {

    query := `
    SELECT id, step_key, trait_key, selector,
           value, specificity, effective_date, expires_date, status
    FROM overrides
    WHERE status = 'active'
    `

    rows, err := r.db.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var result []domain.Override

    for rows.Next() {
        var o domain.Override
        err := rows.Scan(
            &o.ID,
            &o.StepKey,
            &o.TraitKey,
            &o.Selector,
            &o.Value,
            &o.Specificity,
            &o.EffectiveDate,
            &o.ExpiresDate,
            &o.Status,
        )
        if err != nil {
            return nil, err
        }
        result = append(result, o)
    }

    return result, nil
}

func scanOverride(rows Scanner) (domain.Override, error) {
	var o domain.Override

	err := rows.Scan(
		&o.ID,
		&o.StepKey,
		&o.TraitKey,
		&o.Selector,
		&o.Value,
		&o.Specificity,
		&o.EffectiveDate,
		&o.ExpiresDate,
		&o.Status,
		&o.CreatedAt,
		&o.UpdatedAt,
	)

	return o, err
}

type Scanner interface {
	Scan(dest ...any) error
}

func (r *PostgresOverrideRepository) List(
	ctx context.Context,
	filter domain.OverrideFilter,
) ([]domain.Override, error) {

	query := `
	SELECT id, step_key, trait_key, selector,
	       value, specificity, effective_date, expires_date, status,
	       created_at, updated_at
	FROM overrides
	WHERE 1=1
	`

	var args []any
	argIdx := 1

	add := func(condition string, value any) {
		query += fmt.Sprintf(" AND %s = $%d", condition, argIdx)
		args = append(args, value)
		argIdx++
	}

	if filter.StepKey != nil {
		add("step_key", *filter.StepKey)
	}
	if filter.TraitKey != nil {
		add("trait_key", *filter.TraitKey)
	}
	if filter.State != nil {
		add("state", *filter.State)
	}
	if filter.Client != nil {
		add("client", *filter.Client)
	}
	if filter.Investor != nil {
		add("investor", *filter.Investor)
	}
	if filter.CaseType != nil {
		add("case_type", *filter.CaseType)
	}
	if filter.Status != nil {
		add("status", *filter.Status)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Override

	for rows.Next() {
		o, err := scanOverride(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, o)
	}

	return result, nil
}

func (r *PostgresOverrideRepository) Create(
	ctx context.Context,
	o domain.Override,
) error {

	query := `
	INSERT INTO overrides (
		id, step_key, trait_key,
		selector,
		value, specificity,
		effective_date, expires_date,
		status,
		created_at, updated_at
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`

	valueJSON, _ := json.Marshal(o.Value)

	_, err := r.db.Exec(ctx, query,
		o.ID,
		o.StepKey,
		o.TraitKey,
		o.Selector,
		valueJSON,
		o.Specificity,
		o.EffectiveDate,
		o.ExpiresDate,
		o.Status,
		o.CreatedAt,
		o.UpdatedAt,
	)

	return err
}

func (r *PostgresOverrideRepository) Update(
	ctx context.Context,
	o domain.Override,
) error {

	query := `
	UPDATE overrides
	SET step_key=$2,
	    trait_key=$3,
	    selector=$4,
	    value=$5,
	    specificity=$6,
	    effective_date=$7,
	    expires_date=$8,
	    status=$9,
	    updated_at=$10
	WHERE id=$1
	`

	valueJSON, _ := json.Marshal(o.Value)

	_, err := r.db.Exec(ctx, query,
		o.ID,
		o.StepKey,
		o.TraitKey,
		o.Selector,
		valueJSON,
		o.Specificity,
		o.EffectiveDate,
		o.ExpiresDate,
		o.Status,
		o.UpdatedAt,
	)

	return err
}

func (r *PostgresOverrideRepository) UpdateStatus(
	ctx context.Context,
	id string,
	status string,
) error {

	query := `
	UPDATE overrides
	SET status=$2, updated_at=NOW()
	WHERE id=$1
	`

	_, err := r.db.Exec(ctx, query, id, status)
	return err
}

func (r *PostgresOverrideRepository) InsertHistory(
	ctx context.Context,
	before domain.Override,
	after domain.Override,
) error {

	query := `
	INSERT INTO override_history (
		override_id, before, after, changed_at, changed_by
	)
	VALUES ($1,$2,$3,NOW(),'system')
	`

	beforeJSON, _ := json.Marshal(before)
	afterJSON, _ := json.Marshal(after)

	_, err := r.db.Exec(ctx, query,
		before.ID,
		beforeJSON,
		afterJSON,
	)

	return err
}

func (r *PostgresOverrideRepository) GetHistory(
	ctx context.Context,
	id string,
) ([]domain.OverrideHistory, error) {

	query := `
	SELECT id, override_id, before, after, changed_at, changed_by
	FROM override_history
	WHERE override_id = $1
	ORDER BY changed_at DESC
	`

	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.OverrideHistory

	for rows.Next() {
		var h domain.OverrideHistory
		var beforeJSON, afterJSON []byte

		err := rows.Scan(
			&h.ID,
			&h.OverrideID,
			&beforeJSON,
			&afterJSON,
			&h.ChangedAt,
			&h.ChangedBy,
		)
		if err != nil {
			return nil, err
		}

		_ = json.Unmarshal(beforeJSON, &h.Before)
		_ = json.Unmarshal(afterJSON, &h.After)

		result = append(result, h)
	}

	return result, nil
}

func (r *PostgresOverrideRepository) GetByID(
	ctx context.Context,
	id string,
) (*domain.Override, error) {

	query := `
	SELECT id, step_key, trait_key, selector,
	       value, specificity, effective_date, expires_date, status,
	       created_at, updated_at
	FROM overrides
	WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	o, err := scanOverride(row)
	if err != nil {
		return nil, err
	}

	return &o, nil
}