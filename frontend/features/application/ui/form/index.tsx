import { Controller } from 'react-hook-form';
import { AnimatePresence, motion } from 'motion/react';
import { Building2, Check, GraduationCap, Search, UserRound } from 'lucide-react';
import { useEffect, useState } from 'react';

import { InputGroup, InputGroupInput } from '@/shared/components/ui/input-group';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';

import { Icon } from '@/shared/components';
import { Button } from '@/shared/components/ui/button';
import { useGroupsByOrganizationQuery } from '@/entities/group';
import { useOrganizationsByNameQuery, useOrganizationsQuery } from '@/entities/organization';
import type { useApplicationForm } from '../../hooks';

const itemVariants = {
  hidden: { opacity: 0, y: 10 },
  visible: { opacity: 1, y: 0 },
};

const fieldClassName =
  '!h-12 w-full rounded-lg border border-[#e5eaf2] bg-[#f5f7fb] px-3.5 py-0 text-base shadow-none transition-colors duration-200 hover:border-[#c8d7f5] focus-visible:border-[#3984f4] focus-visible:bg-white focus-visible:ring-[3px] focus-visible:ring-[#3984f4]/10';

const iconClassName = 'size-[18px] shrink-0 text-[#7b8ba5]';

export const ApplicationForm = ({ form }: { form: ReturnType<typeof useApplicationForm> }) => {
  const [organizationSearch, setOrganizationSearch] = useState('');
  const [debouncedOrganizationSearch, setDebouncedOrganizationSearch] = useState('');
  const {
    control,
    register,
    setValue,
    watch,
    isSelected,
    handleSubmit,
    onSubmit,
    isLoading,
    isError,
    isSuccess,
  } = form;

  const organizationId = watch('organizationId');
  const groupId = watch('groupId');
  const applicationType = watch('type');

  useEffect(() => {
    const timeoutId = window.setTimeout(
      () => setDebouncedOrganizationSearch(organizationSearch.trim()),
      300
    );
    return () => window.clearTimeout(timeoutId);
  }, [organizationSearch]);

  const isOrganizationSearchEnabled = debouncedOrganizationSearch.length >= 3;
  const isOrganizationSearchDebouncing =
    organizationSearch.trim().length >= 3 &&
    organizationSearch.trim() !== debouncedOrganizationSearch;
  const {
    data: organizations = [],
    isLoading: isOrganizationsLoading,
    isError: isOrganizationsError,
    refetch: refetchOrganizations,
  } = useOrganizationsQuery();
  const {
    currentData: searchedOrganizations = [],
    isFetching: isOrganizationSearchLoading,
    isError: isOrganizationSearchError,
    refetch: refetchOrganizationsByName,
  } = useOrganizationsByNameQuery(debouncedOrganizationSearch, {
    skip: !isOrganizationSearchEnabled,
  });
  const {
    currentData: groups = [],
    isFetching: isGroupsLoading,
    isError: isGroupsError,
    refetch: refetchGroups,
  } = useGroupsByOrganizationQuery(organizationId, { skip: !organizationId });
  const visibleOrganizations = isOrganizationSearchEnabled ? searchedOrganizations : organizations;
  const isOrganizationOptionsLoading =
    isOrganizationsLoading ||
    isOrganizationSearchDebouncing ||
    (isOrganizationSearchEnabled && isOrganizationSearchLoading);
  const isOrganizationOptionsError = isOrganizationSearchEnabled
    ? isOrganizationSearchError
    : isOrganizationsError;
  const selectedOrganization =
    organizations.find((organization) => organization.id === organizationId) ??
    searchedOrganizations.find((organization) => organization.id === organizationId);
  const selectedGroup = groups.find((group) => group.id === groupId);

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <motion.section
        initial={{ opacity: 0, y: 12 }}
        animate={{ opacity: 1, y: 0 }}
        exit={{ opacity: 0, y: -8 }}
        transition={{ duration: 0.3, ease: 'easeOut' }}
        className="mx-auto w-full max-w-110 text-left"
      >
        <div className="overflow-hidden rounded-2xl border border-[#e8edf5] bg-white shadow-[0_12px_48px_-20px_rgba(32,65,115,0.18)] sm:rounded-[28px]">
          <div className="relative overflow-hidden bg-[#3478f6] px-5 pb-6 pt-5 sm:px-7 sm:pb-7 sm:pt-6">
            <div className="pointer-events-none absolute -right-12 -top-16 size-44 rounded-full bg-white/8" />
            <div className="pointer-events-none absolute -bottom-20 right-20 size-36 rounded-full bg-white/[0.07]" />

            <div className="relative flex items-center justify-between">
              <div className="flex items-center gap-2.5">
                {watch('type') == 'student' ? (
                  <>
                    <div className="flex size-10 items-center justify-center rounded-[14px] bg-white text-[#3478f6] shadow-sm">
                      <Icon icon={'graduation-cap'} size={50} className="text-[33px]" />
                    </div>

                    <div>
                      <p className="text-sm font-bold tracking-tight text-white">Учёба</p>
                      <p className="text-xs font-medium text-white/75 sm:text-sm">
                        Личный кабинет студента
                      </p>
                    </div>
                  </>
                ) : (
                  <>
                    <div className="flex size-10 items-center justify-center rounded-[14px] bg-white text-[#3478f6] shadow-sm">
                      <Icon icon={'user'} size={50} className="text-[33px]" />
                    </div>

                    <div>
                      <p className="text-sm font-bold tracking-tight text-white">Управление</p>
                      <p className="text-xs font-medium text-white/75 sm:text-sm">
                        Личный кабинет преподавателя
                      </p>
                    </div>
                  </>
                )}
              </div>
            </div>

            <div className="relative mt-6 sm:mt-7">
              <h2 className="text-2xl font-bold leading-tight tracking-tight text-white sm:text-[27px] sm:tracking-[-0.8px]">
                Давайте
                <br />
                знакомиться!
              </h2>

              <p className="mt-2 max-w-72.5 text-sm leading-5 text-white/80">
                Заполните данные, чтобы настроить ваш профиль.
              </p>
            </div>
          </div>

          <div className="px-4 py-5 sm:px-7 sm:py-7">
            <motion.div
              initial="hidden"
              animate="visible"
              variants={{
                hidden: {},
                visible: {
                  transition: {
                    staggerChildren: 0.07,
                    delayChildren: 0.05,
                  },
                },
              }}
              className="flex flex-col gap-4 sm:gap-5"
            >
              <motion.div variants={itemVariants} className="flex flex-col gap-2">
                <label htmlFor="initial-name" className="text-sm font-semibold text-[#26344b]">
                  Как вас зовут?
                </label>

                <InputGroup className="h-12 rounded-xl border border-[#e5eaf2] bg-[#f5f7fb] px-3.5 shadow-none transition-colors duration-200 hover:border-[#c8d7f5] focus-within:border-[#3984f4] focus-within:bg-white focus-within:ring-[3px] focus-within:ring-[#3984f4]/10">
                  <UserRound className={iconClassName} />

                  <InputGroupInput
                    {...register('name')}
                    id="initial-name"
                    autoComplete="name"
                    placeholder="Иванов Иван Иванович"
                    className="h-full border-0 bg-transparent px-2 text-base text-[#26344b] shadow-none placeholder:text-[#9aa7ba] focus-visible:ring-0"
                  />
                </InputGroup>

                <p className="text-xs leading-4 text-[#66758c]">
                  Укажите фамилию, имя и отчество полностью.
                </p>
              </motion.div>

              <motion.div variants={itemVariants} className="flex flex-col gap-2">
                <label id="organization-label" className="text-sm font-semibold text-[#26344b]">
                  Где вы учитесь?
                </label>

                <InputGroup className="h-12 rounded-xl border border-[#e5eaf2] bg-[#f5f7fb] px-3.5 shadow-none transition-colors duration-200 hover:border-[#c8d7f5] focus-within:border-[#3984f4] focus-within:bg-white focus-within:ring-[3px] focus-within:ring-[#3984f4]/10">
                  <Search className={iconClassName} aria-hidden="true" />
                  <InputGroupInput
                    id="organization-search"
                    type="search"
                    value={organizationSearch}
                    onChange={(event) => setOrganizationSearch(event.target.value)}
                    aria-label="Поиск организации по названию"
                    placeholder="Поиск организации по названию"
                    className="h-full border-0 bg-transparent px-2 text-base text-[#26344b] shadow-none placeholder:text-[#9aa7ba] focus-visible:ring-0"
                  />
                </InputGroup>
                <p className="text-xs leading-4 text-[#66758c]">
                  Введите не менее 3 символов для поиска.
                </p>

                <Controller
                  control={control}
                  name="organizationId"
                  render={({ field }) => (
                    <Select
                      value={field.value}
                      disabled={isOrganizationOptionsLoading || isOrganizationOptionsError}
                      onValueChange={(value) => {
                        field.onChange(value);
                        setOrganizationSearch('');

                        setValue('groupId', '', {
                          shouldDirty: true,
                          shouldValidate: true,
                        });
                      }}
                    >
                      <SelectTrigger
                        aria-labelledby="organization-label"
                        className={fieldClassName}
                      >
                        <div className="flex min-w-0 items-center gap-2.5">
                          <Building2 className={iconClassName} />

                          <SelectValue
                            placeholder={
                              isOrganizationOptionsLoading
                                ? 'Загружаем организации...'
                                : isOrganizationOptionsError
                                  ? 'Не удалось загрузить организации'
                                  : 'Выберите организацию'
                            }
                          />
                        </div>
                      </SelectTrigger>

                      <SelectContent className="rounded-xl border-[#e8edf5] bg-white p-1.5 shadow-xl">
                        {visibleOrganizations.map((organization) => (
                          <SelectItem
                            key={organization.id}
                            value={organization.id}
                            className="min-h-11 rounded-lg py-2 pl-3 pr-8 text-[13px] text-[#26344b] focus:bg-[#eef4ff] focus:text-[#2468e8]"
                          >
                            {organization.name}
                          </SelectItem>
                        ))}
                        {visibleOrganizations.length === 0 && (
                          <SelectItem value="__no-organizations" disabled>
                            {isOrganizationOptionsLoading
                              ? 'Загружаем организации...'
                              : 'Организации не найдены'}
                          </SelectItem>
                        )}
                      </SelectContent>
                    </Select>
                  )}
                />
                {isOrganizationOptionsError && (
                  <div className="flex items-center justify-between gap-3">
                    <p role="alert" className="text-sm text-red-600">
                      Не удалось загрузить организации.
                    </p>
                    <Button
                      type="button"
                      variant="link"
                      onClick={() =>
                        isOrganizationSearchEnabled
                          ? refetchOrganizationsByName()
                          : refetchOrganizations()
                      }
                    >
                      Повторить
                    </Button>
                  </div>
                )}

                {selectedOrganization && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: 'auto' }}
                    className="flex items-center gap-1.5 overflow-hidden px-0.5 text-xs text-[#3478f6]"
                  >
                    <Check className="size-3.5 shrink-0" />
                    <span className="truncate">{selectedOrganization.name}</span>
                  </motion.div>
                )}
              </motion.div>

              {applicationType === 'student' && (
                <motion.div variants={itemVariants} className="flex flex-col gap-2">
                  <label id="group-label" className="text-sm font-semibold text-[#26344b]">
                    Ваша группа
                  </label>
                  <Controller
                    control={control}
                    name="groupId"
                    render={({ field }) => (
                      <Select
                        value={field.value}
                        onValueChange={field.onChange}
                        disabled={!organizationId || isGroupsLoading || isGroupsError}
                      >
                        <SelectTrigger
                          aria-labelledby="group-label"
                          className={`${fieldClassName} disabled:cursor-not-allowed disabled:opacity-50`}
                        >
                          <div className="flex min-w-0 items-center gap-2.5">
                            <GraduationCap className={iconClassName} />

                            <SelectValue
                              placeholder={
                                !organizationId
                                  ? 'Выберите организацию'
                                  : isGroupsLoading
                                    ? 'Загружаем группы...'
                                    : isGroupsError
                                      ? 'Не удалось загрузить группы'
                                      : groups.length === 0
                                        ? 'Группы не найдены'
                                        : 'Выберите учебную группу'
                              }
                            />
                          </div>
                        </SelectTrigger>

                        <SelectContent className="rounded-xl border-[#e8edf5] bg-white p-1.5 shadow-xl">
                          {groups.map((group) => (
                            <SelectItem
                              key={group.id}
                              value={group.id}
                              className="min-h-11 rounded-lg py-2 pl-3 pr-8 text-[13px] text-[#26344b] focus:bg-[#eef4ff] focus:text-[#2468e8]"
                            >
                              {group.name}
                            </SelectItem>
                          ))}
                          {organizationId && groups.length === 0 && (
                            <SelectItem value="__no-groups" disabled>
                              {isGroupsLoading ? 'Загружаем группы...' : 'Группы не найдены'}
                            </SelectItem>
                          )}
                        </SelectContent>
                      </Select>
                    )}
                  />
                  {isGroupsError && (
                    <div className="flex items-center justify-between gap-3">
                      <p role="alert" className="text-sm text-red-600">
                        Не удалось загрузить группы.
                      </p>
                      <Button type="button" variant="link" onClick={() => refetchGroups()}>
                        Повторить
                      </Button>
                    </div>
                  )}
                  {selectedGroup && (
                    <motion.div
                      initial={{ opacity: 0, height: 0 }}
                      animate={{ opacity: 1, height: 'auto' }}
                      className="flex items-center gap-1.5 overflow-hidden px-0.5 text-xs text-[#3478f6]"
                    >
                      <Check className="size-3.5 shrink-0" />
                      <span>Группа {selectedGroup.name} выбрана</span>
                    </motion.div>
                  )}
                </motion.div>
              )}
              <motion.div variants={itemVariants} className="flex flex-col gap-2">
                <AnimatePresence mode="wait">
                  {isSelected && (
                    <motion.div
                      className="w-full"
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      exit={{ opacity: 0 }}
                      transition={{ duration: 0.3 }}
                    >
                      <Button
                        type="submit"
                        className="h-12 w-full text-base"
                        disabled={isLoading || isSuccess}
                      >
                        {isLoading
                          ? 'Отправляем...'
                          : isSuccess
                            ? 'Заявка отправлена'
                            : 'Отправить заявку'}
                      </Button>
                      {isError && (
                        <p role="alert" className="mt-2 text-sm text-red-600">
                          Не удалось отправить заявку. Проверьте соединение и попробуйте ещё раз.
                        </p>
                      )}
                      {isSuccess && (
                        <p role="status" className="mt-2 text-sm text-green-700">
                          Заявка успешно отправлена.
                        </p>
                      )}
                    </motion.div>
                  )}
                </AnimatePresence>
              </motion.div>
            </motion.div>
          </div>
        </div>
      </motion.section>
    </form>
  );
};
