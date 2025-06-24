package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
)

func main() {
	role := flag.String("role", "", "Role name (e.g. admin, teacher)")
	perm := flag.String("permission", "", "Permission string (e.g. user_v1_user_service_register)")
	flag.Parse()

	if *role == "" || *perm == "" {
		log.Fatal("Both --role and --permission are required")
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres passwordUtils=postgres dbname=attendly sslmode=disable")
	if err != nil {
		log.Fatal("db connection failed:", err)
	}
	defer db.Close()

	var roleID, permID int

	err = db.QueryRow(`SELECT id FROM roles WHERE name=$1`, *role).Scan(&roleID)
	if err != nil {
		log.Fatalf("Role '%s' not found: %v", *role, err)
	}

	err = db.QueryRow(`INSERT INTO permissions (action) VALUES ($1)
		ON CONFLICT (action) DO UPDATE SET action=EXCLUDED.action
		RETURNING id`, *perm).Scan(&permID)
	if err != nil {
		log.Fatalf("Permission insert failed: %v", err)
	}

	_, err = db.Exec(`INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING`, roleID, permID)
	if err != nil {
		log.Fatalf("Failed to grant permission: %v", err)
	}

	fmt.Printf("✅ Granted permission '%s' to role '%s'\n", *perm, *role)
}
