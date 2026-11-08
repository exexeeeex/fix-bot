import type { Student } from '@/entities/student';

export interface Abcense {
  allStudents: Student[];
  abcensesStudents: Student[];
  lateStudents: Student[];
}
