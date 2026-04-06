package db

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"rules-resolution-service/internal/domain"
	"rules-resolution-service/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedDefaults(ctx context.Context, db *pgxpool.Pool) {
	data, err := os.ReadFile("seed/defaults.json")
	if err != nil {
		log.Fatal(err)
	}

	var raw map[string]map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		log.Fatal("unmarshal defaults:", err)
	}

	for step, traits := range raw {
		for trait, value := range traits {

			valJSON, _ := json.Marshal(value)

			_, err := db.Exec(ctx,`
				INSERT INTO defaults (step_key, trait_key, value)
				VALUES ($1,$2,$3)
				ON CONFLICT (step_key, trait_key) DO NOTHING
			`, step, trait, valJSON)

			if err != nil {
				log.Println("default seed error:", err)
			}
		}
	}

	log.Println("defaults seeded")
}

func SeedSteps(ctx context.Context, db *pgxpool.Pool) {
	data, err := os.ReadFile("seed/steps.json")
	if err != nil {
		log.Fatal(err)
	}

	var steps []map[string]interface{}
	if err := json.Unmarshal(data, &steps); err != nil {
		log.Fatal(err)
	}

	for _, s := range steps {
		_, err := db.Exec(ctx, `
			INSERT INTO steps (key, name, description, position)
			VALUES ($1,$2,$3,$4)
			ON CONFLICT (key) DO NOTHING
		`,
			s["key"],
			s["name"],
			s["description"],
			int(s["position"].(float64)),
		)

		if err != nil {
			log.Println("step seed error:", err)
		}
	}

	log.Println("steps seeded")
}

func SeedOverrides(ctx context.Context, repo repository.OverrideRepository) {
	data, err := os.ReadFile("seed/overrides.json")
	if err != nil {
		log.Fatal("read overrides:", err)
	}

	var raw []map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		log.Fatal("unmarshal overrides:", err)
	}

	for _, item := range raw {

		effStr, ok := item["effectiveDate"].(string)
		if !ok {
			log.Println("invalid effectiveDate")
			continue
		}

		eff, err := time.Parse("2006-01-02", effStr)
		if err != nil {
			log.Println("date parse error:", err)
			continue
		}

		var exp *time.Time
		if expVal, ok := item["expiresDate"]; ok && expVal != nil {
			expStr, ok := expVal.(string)
			if ok {
				t, err := time.Parse("2006-01-02", expStr)
				if err == nil {
					exp = &t
				}
			}
		}

		selectorMap, ok := item["selector"].(map[string]interface{})
		if !ok {
			log.Println("invalid selector")
			continue
		}

		o := domain.Override{
			ID:            safeString(item["id"]),
			StepKey:       safeString(item["stepKey"]),
			TraitKey:      safeString(item["traitKey"]),
			Selector:      toStringMap(selectorMap),
			Value:         item["value"],
			Specificity:   len(selectorMap),
			EffectiveDate: eff,
			ExpiresDate:   exp,
			Status:        safeString(item["status"]),
		}

		if err := repo.Create(ctx, o); err != nil {
			log.Println("override seed error:", err)
		}
	}

	log.Println("overrides seeded")
}

func toStringMap(m interface{}) map[string]string {
	res := make(map[string]string)

	for k, v := range m.(map[string]interface{}) {
		res[k] = v.(string)
	}
	return res
}

func safeString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}