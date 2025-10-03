package models

//structure des enseignant en attente de validation
type TeacherWaiting struct {
	TeacherBase
}

func (TeacherWaiting) TableName() string {
	return "teacher_waiting"
}
