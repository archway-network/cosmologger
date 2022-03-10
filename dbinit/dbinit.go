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

		`CREATE INDEX IF NOT EXISTS unjail_tx
		ON public.tx_events USING btree
		(action COLLATE pg_catalog."default" varchar_ops ASC NULLS LAST, sender COLLATE pg_catalog."default" varchar_ops ASC NULLS LAST)
		TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS action_indx
		ON public.tx_events USING btree
		(action COLLATE pg_catalog."default" ASC NULLS LAST)
		TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS sender
		ON public.tx_events USING btree
		(sender COLLATE pg_catalog."default" ASC NULLS LAST)
		TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS validator
		ON public.tx_events USING btree
		(validator COLLATE pg_catalog."default" ASC NULLS LAST)
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

		`CREATE TABLE IF NOT EXISTS public.validators
		(
			"oprAddr" character varying(255) COLLATE pg_catalog."default" NOT NULL,
			"consAddr" character varying(255) COLLATE pg_catalog."default" NOT NULL,
			"accountAddr" character varying(255) COLLATE pg_catalog."default" NOT NULL,
			CONSTRAINT validators_pkey PRIMARY KEY ("oprAddr", "consAddr", "accountAddr")
		)
		
		TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS "blockHeight"
			ON public.block_signers USING btree
			("blockHeight" ASC NULLS LAST)
			TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS "valConsAddr"
			ON public.block_signers USING btree
			("valConsAddr" COLLATE pg_catalog."default" ASC NULLS LAST)
			TABLESPACE pg_default;`,
		`CREATE TABLE IF NOT EXISTS public.participants
		(
			"accountAddress" character varying(200) COLLATE pg_catalog."default" NOT NULL,
			"fullLegalName" character varying(255) COLLATE pg_catalog."default",
			"githubHandle" character varying(255) COLLATE pg_catalog."default",
			"emailAddress" character varying(255) COLLATE pg_catalog."default" NOT NULL,
			pubkey character varying(255) COLLATE pg_catalog."default",
			"kycSessionId" character varying(255) COLLATE pg_catalog."default",
			"kycVerified" boolean DEFAULT false,
			CONSTRAINT participants_pkey PRIMARY KEY ("accountAddress")
		)
		
		TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS "emailAddress"
		ON public.participants USING btree
		("emailAddress" COLLATE pg_catalog."default" ASC NULLS LAST)
		TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS "kycVerified"
		ON public.participants USING btree
		("kycVerified" ASC NULLS LAST)
		TABLESPACE pg_default;`,

		`CREATE TABLE IF NOT EXISTS public.contracts
		(
			"contractAddress" character varying(255) COLLATE pg_catalog."default",
			"rewardAddress" character varying(255) COLLATE pg_catalog."default",
			"developerAddress" character varying(255) COLLATE pg_catalog."default",
			"blockHeight" bigint,
			"gasConsumed" bigint,
			"rewardsDenom" character varying(255) COLLATE pg_catalog."default",
			"contractRewardsAmount" double precision,
			"inflationRewardsAmount" double precision,
			"leftoverRewardsAmount" double precision,
			"collectPremium" boolean,
			"gasRebateToUser" boolean,
			"premiumPercentageCharged" bigint,
			"metadataJson" text COLLATE pg_catalog."default",
			"incId" bigint NOT NULL DEFAULT nextval('"contracts_incId_seq"'::regclass),
			CONSTRAINT contracts_pkey PRIMARY KEY ("incId")
		)
		TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS "contractAddress"
		ON public.contracts USING btree
		("contractAddress" COLLATE pg_catalog."default" ASC NULLS LAST)
		TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS "developerAddress"
		ON public.contracts USING btree
		("developerAddress" COLLATE pg_catalog."default" ASC NULLS LAST)
		TABLESPACE pg_default;`,

		`CREATE INDEX IF NOT EXISTS "rewardAddress"
		ON public.contracts USING btree
		("rewardAddress" COLLATE pg_catalog."default" ASC NULLS LAST)
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
