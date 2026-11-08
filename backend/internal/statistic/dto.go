package statistic

import "github.com/exexeeeex/fiks-api/internal/domain"

type Absence struct {
	AllStudents      []domain.Student    `json:"all_students"`
	AbsencesStudents []*domain.Student `json:"abcenses_students"`
}

type Statistics struct {
	AllStudents      []domain.Student    `json:"all_students"`
	AbsencesStudents []*domain.Student `json:"abcenses_students"`
	LateStudents     []*domain.Student `json:"late_students"`
}