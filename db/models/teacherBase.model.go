package models

// structure des ensignant temporaire
type TeacherBase struct {
	UserBase

	//nom de la matière qu'il veut enseigner
	SubjectName string `gorm:"column:subject_name;not null" validate:"required,name"`
}
