package models

import "database/sql"

type Models struct {
	Songs              SongModel
	Users              UserModel
	TILs               TILModel
	Pages              PageModel
	Blogs              BlogModel
	BellevueActivities BellevueActivityModel
}

func New(db *sql.DB) Models {
	return Models{
		Songs:              SongModel{DB: db},
		Users:              UserModel{DB: db},
		TILs:               TILModel{DB: db},
		Pages:              PageModel{DB: db},
		Blogs:              BlogModel{DB: db},
		BellevueActivities: BellevueActivityModel{DB: db},
	}
}
