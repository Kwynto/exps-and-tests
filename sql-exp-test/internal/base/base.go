package base

import (
	"context"
	"fmt"
	"log"
	"sql-exp-test/internal/storage"
)

func Run(ctx context.Context, db storage.Storage) {

	db.Init(ctx)

	ent01 := storage.Entities{
		Name:        "Mister First",
		Value:       1.0,
		Description: "This man was a first man in the cinima.",
		Flag:        false,
	}
	{
		id, err := db.Create(ctx, &ent01)
		if err != nil {
			log.Fatalln("Create failed")
		}
		ent01.Id = id
	}
	fmt.Println("Created entity:", "№", ent01.Id)

	ent02, err := db.Read(ctx, ent01.Id)
	if err != nil {
		log.Fatalln("Not read id:", ent01.Id, "Error:", err)
	}
	fmt.Println(*ent02)

	ent02.Flag = true
	ent02.Description = "This is changed string."
	if err := db.Update(ctx, ent02); err != nil {
		log.Fatalln("Not update entity:", err)
	}
	fmt.Println("New value of entity: ", *ent02)

	exist, err := db.IsExists(ctx, ent02)
	if err != nil {
		log.Fatalln("Not run IfExists: ", err)
	}
	if exist {
		fmt.Println("Entity ", *ent02, "is exist.")
	} else {
		fmt.Println("Entity ", *ent02, "is'nt exist.")
	}

	exist, err = db.IsExistsById(ctx, ent02.Id)
	if err != nil {
		log.Fatalln("Not run IfExists: ", err)
	}
	if exist {
		fmt.Println("Entity ", *ent02, "is exist.")
	} else {
		fmt.Println("Entity ", *ent02, "is'nt exist.")
	}

	if err := db.Delete(ctx, ent02); err != nil {
		log.Fatalln("Not delete entity:", ent02)
	}
	fmt.Println("Entity was deleted.")

	exist, err = db.IsExists(ctx, ent02)
	if err != nil {
		log.Fatalln("Not run IfExists: ", err)
	}
	if exist {
		fmt.Println("Entity ", *ent02, "is exist.")
	} else {
		fmt.Println("Entity ", *ent02, "is'nt exist.")
	}

	{
		id, err := db.Create(ctx, &ent01)
		if err != nil {
			log.Fatalln("Create failed")
		}
		ent01.Id = id
	}
	fmt.Println("Created entity:", "№", ent01.Id)

	if err := db.DeleteId(context.Background(), ent01.Id); err != nil {
		log.Fatal("Not deleted entity by Id:", ent01.Id)
	}
	fmt.Println("Entity was deleted.")

}
