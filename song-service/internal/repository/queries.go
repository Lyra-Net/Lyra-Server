package repository

import (
	"context"
)

func (q *Queries) GetEnumValues(ctx context.Context, enumType string) ([]string, error) {
	query := `SELECT unnest(enum_range(NULL::` + enumType + `))::text`
	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	return values, nil
}
