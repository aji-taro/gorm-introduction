package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// データベースに接続する
// https://gorm.io/ja_JP/docs/connecting_to_the_database.html#MySQL
func connectDb() error {
	dataSourceName := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") +
		"@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" +
		os.Getenv("DB_NAME") +
		"?charset=utf8mb4&parseTime=True&loc=Asia%2FTokyo"
	var err error
	db, err = gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("gorm.Open error", err)
	}

	if false {
		// 取得テスト
		var userTest User
		errTest := db.First(&userTest, "id = ?", 1).Error
		if errTest != nil {
			if errors.Is(errTest, gorm.ErrRecordNotFound) {
				fmt.Println("err:", errTest)
			} else {
				fmt.Println("err:", errTest)
			}
			return nil
		}
		fmt.Println("# db test", userTest)
	}
		
	return nil
}

// DB接続終了
func closeDb() {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}

	fmt.Println("# db closed")
}

// User テーブルのモデル
// ※参照：モデルを宣言する（https://gorm.io/ja_JP/docs/models.html）
// ※参照：Belongs To（https://gorm.io/ja_JP/docs/belongs_to.html#Belongs-To）
// ※ここでは、上記ページの構造体をコピペしたけどgorm.Modelを使ってもよいかも
type User struct {
	ID           uint           // Standard field for the primary key
	Name         string         // A regular string field
	Email        *string        // A pointer to a string, allowing for null values
	Age          uint8          // An unsigned 8-bit integer
	Birthday     *time.Time     // A pointer to time.Time, can be null
	MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields
	CompanyId    uint
	Belong       Company `gorm:"foreignKey:CompanyId"`
	CreatedAt    time.Time      // Automatically managed by GORM for creation time
	UpdatedAt    time.Time      // Automatically managed by GORM for update time
	ignored      string         // fields that aren't exported are ignored
}

type Company struct {
	ID   uint
	Name string
}

func main() {
	fmt.Println("# gorm")

	// DB接続処理を実行
	if err := connectDb(); err != nil {
		fmt.Println("connectDb error", err)
		return
	}
	defer closeDb()

	// レコードの作成
	// https://gorm.io/ja_JP/docs/create.html
	fmt.Println("- Create")
	now := time.Now()
	newUser := User{Name: "Jinzhu", Age: 18, Birthday: &now}
	createResult := db.Create(&newUser);
	if createResult.Error != nil {
		fmt.Println(createResult.Error)
		return
	}
	fmt.Println("inserted data's primary key", newUser.ID)
	fmt.Println("inserted records count", createResult.RowsAffected)

	// レコードの取得 その 1
	// https://gorm.io/ja_JP/docs/query.html
	fmt.Println("- First")
	user01 := User{}
	if result := db.First(&user01, newUser.ID); result.Error != nil {
		fmt.Println(result.Error)
		return
	}
	fmt.Println("user01", user01)

	// レコードの取得 その 2
	fmt.Println("- Where")
	user02 := User{}
	if result := db.Where("name = ?", "jinzhu").First(&user02); result.Error != nil {
		fmt.Println(result.Error)
		return
	}
	fmt.Println("user02", user02)

	// トランザクション
	// https://gorm.io/ja_JP/docs/transactions.html#%E3%83%88%E3%83%A9%E3%83%B3%E3%82%B6%E3%82%AF%E3%82%B7%E3%83%A7%E3%83%B3
	fmt.Println("- Transaction")
	db.Transaction(func(tx *gorm.DB) error {
		// レコードの更新
		// https://gorm.io/ja_JP/docs/update.html
		fmt.Println("- Save")
		user02.Name = "jinzhu 2"
		user02.Age = 100
		if result := tx.Save(&user02); result.Error != nil {
			fmt.Println(result.Error)
			return errors.New("rollback user2")
		}
		fmt.Println("update user02")
		return nil
	})

	// Eager Loading(Preload)
	// https://gorm.io/ja_JP/docs/preload.html#Preload
	fmt.Println("- Preload")
	preloadUser := User{}
	if result := db.Preload("Belong").First(&preloadUser, 1); result.Error != nil {
		fmt.Println(result.Error)
		return
	}
	fmt.Println("preloadUser.Name", preloadUser.Name )
	fmt.Println("preloadUser.Belong", preloadUser.Belong )

	// ジョイン
	// https://gorm.io/ja_JP/docs/query.html#Joins
	fmt.Println("- Joins")
	type AliasModel struct {
		UserName string
		CompanyName string
	}
	aliasModel := AliasModel{}
	if result := db.Model(&User{}).Select("users.name as user_name, companies.name as company_name").Joins("JOIN companies ON users.company_id = companies.id").Scan(&aliasModel); result.Error != nil {
		fmt.Println("error", result.Error)
		return
	}

	fmt.Println("aliasModel", aliasModel )

}
