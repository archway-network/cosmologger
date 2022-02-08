package dbinit

import (
	"log"

	"strings"

	"github.com/archway-network/cosmologger/database"
)

/**
* This function initializes the database
* It checks if DB is not ready and then creates tabels, indices, ...
 */
func DatabaseInit(db *database.Database) {

	if !NeedToInitDB(db) {
		return
	}

	/*--------------*/
	log.Printf("Database initialization started.")
	log.Printf("\tCreating Tables and Indices...")

	err := CreateTables(db)
	if err != nil {
		panic(err)
	}
	log.Printf("Done")

	log.Printf("Database initialization Done.\n\n")
}

/*--------------------------------*/

func NeedToInitDB(db *database.Database) bool {

	SQL := `SELECT * FROM "tx_events" LIMIT 1;`
	_, err := db.Query(SQL, database.QueryParams{})
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return true
		}
		panic(err)
	}
	return false
}

/*--------------------------------*/

func CreateTables(db *database.Database) error {
	SQList := []string{
		`CREATE TABLE IF NOT EXISTS public.tx_events
			(
				"txHash" character(64) COLLATE pg_catalog."default" NOT NULL,
				height bigint DEFAULT 0,
				module character varying(20) COLLATE pg_catalog."default" DEFAULT ''::character varying,
				sender character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
				receiver character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
				validator character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
				action character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
				amount character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
				"txAccSeq" character varying(255) COLLATE pg_catalog."default" DEFAULT ''::character varying,
				"txSignature" character varying(255) COLLATE pg_catalog."default" DEFAULT ''::character varying,
				"proposalId" bigint DEFAULT 0,
				"txMemo" text COLLATE pg_catalog."default",
				json text COLLATE pg_catalog."default",
				"logTime" timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
				CONSTRAINT tx_events_pkey PRIMARY KEY ("txHash")
			)
		TABLESPACE pg_default;`,

		// `ALTER TABLE IF EXISTS public.tx_events OWNER to root;`,

		`CREATE INDEX IF NOT EXISTS height
			ON public.tx_events USING btree
			(height ASC NULLS LAST)
			INCLUDE(height)
		TABLESPACE pg_default;`,

		`CREATE TABLE IF NOT EXISTS public.block_signers
			(
				"blockHeight" bigint NOT NULL,
				"valConsAddr" character varying(150) COLLATE pg_catalog."default" NOT NULL,
				"time" timestamp without time zone,
				signature text COLLATE pg_catalog."default",
				CONSTRAINT block_signers_pkey PRIMARY KEY ("blockHeight", "valConsAddr")
			)
			
			TABLESPACE pg_default;`,
		`CREATE TABLE IF NOT EXISTS public.blocks
		(
			"blockHash" character varying(255) COLLATE pg_catalog."default" NOT NULL,
			height bigint NOT NULL,
			"numOfTxs" bigint DEFAULT 0,
			"time" timestamp without time zone,
			CONSTRAINT blocks_pkey PRIMARY KEY ("blockHash")
		)
			TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS height_index
		ON public.blocks USING btree
		(height ASC NULLS LAST)
		TABLESPACE pg_default;`,
		`CREATE TABLE public.validators
		(
			"oprAddr" character varying(255) NOT NULL,
			"consAddr" character varying(255) NOT NULL,
			PRIMARY KEY ("oprAddr", "consAddr")
		)
			TABLESPACE pg_default;`,
	}

	for _, SQL := range SQList {
		_, err := db.Exec(SQL, database.QueryParams{})
		if err != nil {
			// fmt.Printf("\n\tError in SQL: %+v\n", SQL)
			return err
		}
	}

	return nil
}

/*--------------------------------*/
