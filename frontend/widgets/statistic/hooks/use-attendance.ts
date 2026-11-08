import { useMemo, useState } from 'react';
import { useGetAbcenses } from '@/features/statistic';

type AttendanceFilter = 'all' | 'absent' | 'present';

export const useAttendance = () => {
  const { data, isError, isLoading } = useGetAbcenses();
  const [filter, setFilter] = useState<AttendanceFilter>('all');
  const [search, setSearch] = useState('');

  const students = data?.allStudents.length ? data.allStudents : (data?.abcensesStudents ?? []);
  const absentIds = useMemo(
    () => new Set(data?.abcensesStudents.map((student) => student.id) ?? []),
    [data?.abcensesStudents]
  );
  const absentCount = students.filter((student) => absentIds.has(student.id)).length;
  const presentCount = students.length - absentCount;
  const attendancePercent = students.length
    ? Math.round((presentCount / students.length) * 100)
    : 0;
  const filteredStudents = students.filter((student) => {
    const isAbsent = absentIds.has(student.id);
    const matchesFilter =
      filter === 'all' || (filter === 'absent' && isAbsent) || (filter === 'present' && !isAbsent);
    const matchesSearch = student.name
      .toLocaleLowerCase('ru-RU')
      .includes(search.trim().toLocaleLowerCase('ru-RU'));

    return matchesFilter && matchesSearch;
  });

  return {
    isError,
    isLoading,
    students,
    absentIds,
    absentCount,
    presentCount,
    attendancePercent,
    filteredStudents,
    filter,
    setFilter,
    search,
    setSearch,
  };
};

export const getStudentInitials = (name: string) =>
  name
    .trim()
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
    .toLocaleUpperCase('ru-RU');
