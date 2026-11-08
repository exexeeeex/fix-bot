import { CircleAlert, Search, UserRoundCheck, UserRoundX, UsersRound } from 'lucide-react';
import type { Student } from '@/entities/student';
import { getStudentInitials, useAttendance } from '../hooks';

export const Abcenses = () => {
  const {
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
  } = useAttendance();

  if (isLoading) {
    return (
      <section
        className="mx-auto w-full max-w-6xl px-4 py-8 sm:px-6 lg:px-8"
        aria-busy="true"
        aria-label="Загрузка статистики посещаемости"
      >
        <div className="h-52 animate-pulse rounded-3xl bg-white" />
        <div className="mt-5 grid gap-4 sm:grid-cols-3">
          {[0, 1, 2].map((item) => (
            <div key={item} className="h-36 animate-pulse rounded-3xl bg-white" />
          ))}
        </div>
      </section>
    );
  }

  if (isError) {
    return (
      <section className="mx-auto w-full max-w-6xl px-4 py-8 sm:px-6 lg:px-8">
        <div className="rounded-3xl border border-red-200 bg-red-50 p-6 text-red-900" role="alert">
          <div className="flex items-start gap-3">
            <CircleAlert className="mt-0.5 size-5 shrink-0" aria-hidden="true" />
            <div>
              <h2 className="font-semibold">Не удалось загрузить посещаемость</h2>
              <p className="mt-1 text-sm text-red-800">
                Проверьте подключение и попробуйте открыть страницу позже.
              </p>
            </div>
          </div>
        </div>
      </section>
    );
  }

  return (
    <section className="mx-auto w-full max-w-6xl px-4 py-4 pb-28 sm:px-6 lg:px-8">
      <h1 className="text-foreground font-bold text-2xl">Посещаемость</h1>
      <div className="mt-4 grid gap-4 sm:grid-cols-3">
        <SummaryCard
          label="Всего студентов"
          value={students.length}
          description="В списке организации"
          icon={<UsersRound aria-hidden="true" />}
          color="blue"
        />
        <SummaryCard
          label="Отсутствуют"
          value={absentCount}
          description="Нужно уточнить причину"
          icon={<UserRoundX aria-hidden="true" />}
          color="red"
        />
        <SummaryCard
          label="Присутствуют"
          value={presentCount}
          description="На занятиях сейчас"
          icon={<UserRoundCheck aria-hidden="true" />}
          color="green"
        />
      </div>

      <div className="mt-8">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h2 className="text-xl font-semibold tracking-tight text-slate-950">
              Список студентов
            </h2>
            <p className="mt-1 text-sm text-slate-500">
              {absentCount ? `Отсутствуют: ${absentCount}` : 'Все студенты на месте'}
            </p>
          </div>
          <label className="relative block w-full sm:max-w-xs">
            <span className="sr-only">Поиск по имени студента</span>
            <Search
              className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-slate-400"
              aria-hidden="true"
            />
            <input
              className="h-11 w-full rounded-xl border border-slate-200 bg-white pl-10 pr-4 text-sm outline-none transition placeholder:text-slate-400 focus-visible:border-blue-500 focus-visible:ring-4 text-foreground focus-visible:ring-blue-500/15"
              type="search"
              placeholder="Найти студента"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
          </label>
        </div>

        <div className="mt-5 flex flex-wrap gap-2" role="group" aria-label="Фильтр студентов">
          <FilterButton active={filter === 'all'} onClick={() => setFilter('all')}>
            Все <span className="tabular-nums">{students.length}</span>
          </FilterButton>
          <FilterButton active={filter === 'absent'} onClick={() => setFilter('absent')}>
            Отсутствуют <span className="tabular-nums">{absentCount}</span>
          </FilterButton>
          <FilterButton active={filter === 'present'} onClick={() => setFilter('present')}>
            Присутствуют <span className="tabular-nums">{presentCount}</span>
          </FilterButton>
        </div>

        {filteredStudents.length ? (
          <ul className="mt-4 grid gap-3 md:grid-cols-2">
            {filteredStudents.map((student) => (
              <StudentRow key={student.id} student={student} isAbsent={absentIds.has(student.id)} />
            ))}
          </ul>
        ) : (
          <div className="mt-4 rounded-2xl border border-dashed border-slate-300 bg-white/70 px-6 py-10 text-center">
            <p className="font-medium text-slate-800">
              {students.length
                ? 'По вашему запросу ничего не найдено'
                : 'Пока нет данных о студентах'}
            </p>
            <p className="mt-1 text-sm text-slate-500">
              {students.length
                ? 'Попробуйте изменить имя или фильтр.'
                : 'Список появится, когда данные о посещаемости будут загружены.'}
            </p>
          </div>
        )}
      </div>
    </section>
  );
};

function SummaryCard({
  label,
  value,
  description,
  icon,
  color,
}: {
  label: string;
  value: number;
  description: string;
  icon: React.ReactNode;
  color: 'blue' | 'red' | 'green';
}) {
  const colors = {
    blue: 'bg-blue-50 text-blue-700',
    red: 'bg-red-50 text-red-700',
    green: 'bg-emerald-50 text-emerald-700',
  };

  return (
    <article className="rounded-3xl border border-slate-200/80 bg-white p-5 shadow-sm sm:p-6">
      <div className="flex items-start justify-between gap-4">
        <div>
          <p className="text-sm font-medium text-slate-500">{label}</p>
          <p className="mt-3 text-3xl font-semibold tracking-tight text-slate-950 tabular-nums">
            {value}
          </p>
          <p className="mt-1 text-xs text-slate-500">{description}</p>
        </div>
        <span className={`flex size-11 items-center justify-center rounded-2xl ${colors[color]}`}>
          {icon}
        </span>
      </div>
    </article>
  );
}

function FilterButton({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      type="button"
      aria-pressed={active}
      onClick={onClick}
      className={`inline-flex min-h-10 items-center gap-2 rounded-full px-4 text-sm font-medium transition focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-blue-500/20 ${
        active
          ? 'bg-slate-950 text-white'
          : 'border border-slate-200 bg-white text-slate-600 hover:border-slate-300 hover:bg-slate-50'
      }`}
    >
      {children}
    </button>
  );
}

function StudentRow({ student, isAbsent }: { student: Student; isAbsent: boolean }) {
  return (
    <li
      className={`flex min-w-0 items-center gap-3 rounded-2xl border p-4 transition ${
        isAbsent ? 'border-red-200 bg-red-50/80' : 'border-slate-200/80 bg-white'
      }`}
    >
      <span
        className={`flex size-11 shrink-0 items-center justify-center rounded-full text-sm font-semibold ${
          isAbsent ? 'bg-red-100 text-red-700' : 'bg-blue-50 text-blue-700'
        }`}
        aria-hidden="true"
      >
        {getStudentInitials(student.name)}
      </span>
      <span className="min-w-0 flex-1 truncate font-medium text-slate-900">{student.name}</span>
      <span
        className={`inline-flex shrink-0 items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold ${
          isAbsent ? 'bg-red-100 text-red-700' : 'bg-emerald-100 text-emerald-700'
        }`}
      >
        <span
          className={`size-1.5 rounded-full ${isAbsent ? 'bg-red-600' : 'bg-emerald-600'}`}
          aria-hidden="true"
        />
        {isAbsent ? 'Отсутствует' : 'На месте'}
      </span>
    </li>
  );
}
