package dbinit

// import (
// 	"fmt"
// 	"log"

// 	// "sensor-data-simulator/database"
// 	// "sensor-data-simulator/global"
// 	"strings"
// )

// /**
// * This function initializes the database
// * It checks if DB is not ready and then creates tabels, indices, ...
//  */
// func DatabaseInit() {

// 	if !NeedToInitDB() {
// 		return
// 	}

// 	/*--------------*/
// 	log.Printf("Database initialization started.")
// 	log.Printf("\tCreating Tables and Indices...")

// 	err := CreateTables()
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Printf("Done")

// 	log.Printf("Database initialization Done.\n\n")
// }

// /*--------------------------------*/

// func NeedToInitDB() bool {

// 	SQL := `SELECT * FROM "tx_events" LIMIT 1;`
// 	_, err := global.DB.Query(SQL, database.QueryParams{})
// 	if err != nil {
// 		if strings.Contains(err.Error(), "does not exist") {
// 			return true
// 		}
// 		panic(err)
// 	}
// 	return false
// }

// /*--------------------------------*/

// func CreateTables() error {
// 	SQList := []string{
// 		`CREATE TABLE IF NOT EXISTS public.tx_events
// 		(
// 			"txHash" character(64) COLLATE pg_catalog."default" NOT NULL,
// 			height bigint DEFAULT 0,
// 			module character varying(20) COLLATE pg_catalog."default" DEFAULT ''::character varying,
// 			sender character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
// 			receiver character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
// 			validator character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
// 			action character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
// 			amount character varying(100) COLLATE pg_catalog."default" DEFAULT ''::character varying,
// 			"txAccSeq" character varying(255) COLLATE pg_catalog."default" DEFAULT ''::character varying,
// 			"txSignature" character varying(255) COLLATE pg_catalog."default" DEFAULT ''::character varying,
// 			"proposalId" bigint DEFAULT 0,
// 			"txMemo" text COLLATE pg_catalog."default",
// 			json text COLLATE pg_catalog."default",
// 			CONSTRAINT tx_events_pkey PRIMARY KEY ("txHash")
// 		)
// 		TABLESPACE pg_default;`,

// 		// `ALTER TABLE IF EXISTS public.tx_events OWNER to root;`,

// 		`CREATE INDEX IF NOT EXISTS height
// 			ON public.tx_events USING btree
// 			(height ASC NULLS LAST)
// 			INCLUDE(height)
// 			TABLESPACE pg_default;`,
// 	}

// 	for _, SQL := range SQList {
// 		_, err := global.DB.Exec(SQL, database.QueryParams{})
// 		if err != nil {
// 			fmt.Printf("\n\tError in SQL: %+v\n", SQL)
// 			return err
// 		}
// 	}

// 	return nil
// }

// /*--------------------------------*/
