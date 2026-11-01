import { Controller } from 'react-hook-form';
import { AnimatePresence, motion } from 'motion/react';
import { Building2, Check, ChevronDown, GraduationCap, UserRound } from 'lucide-react';

import { InputGroup, InputGroupInput } from '@/shared/components/ui/input-group';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/shared/components/ui/select';

import type { useInitialForm } from '../../hooks';
import { Icon } from '@/shared/components';
import { Button } from '@/shared/components/ui/button';
import { getMax } from '@/shared/lib/platform/max';

const organizations = [
  {
    id: '20000000-0000-0000-0000-000000000001',
    name: 'ГАПОУ ОНТ им. В. А. Сорокина',
  },
];

const groupsByOrganization: Record<string, { id: string; name: string }[]> = {
  '20000000-0000-0000-0000-000000000001': [
    { id: '30000000-0000-0000-0000-000000000001', name: 'РУПО-1' },
    { id: '30000000-0000-0000-0000-000000000002', name: 'РУПО-2' },
    { id: '30000000-0000-0000-0000-000000000003', name: 'РУПО-3' },
  ],
};

const itemVariants = {
  hidden: { opacity: 0, y: 10 },
  visible: { opacity: 1, y: 0 },
};

const fieldClassName =
  'w-full text-[1rem] py-5 rounded-lg border border-[#e5eaf2] bg-[#f5f7fb] px-3.5 shadow-none transition-all duration-200 hover:border-[#c8d7f5] focus-visible:border-[#3984f4] focus-visible:bg-white focus-visible:ring-[3px] focus-visible:ring-[#3984f4]/10';

const iconClassName = 'size-[18px] shrink-0 text-[#7b8ba5]';

export const InitialForm = ({ form }: { form: ReturnType<typeof useInitialForm> }) => {
  const { control, register, setValue, watch, isSelected } = form;

  const max = getMax();

  const organizationId = watch('organizationId');
  const groupId = watch('groupId');

  const groups = groupsByOrganization[organizationId] ?? [];
  const selectedOrganization = organizations.find(
    (organization) => organization.id === organizationId
  );
  const selectedGroup = groups.find((group) => group.id === groupId);

  return (
    <motion.section
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      exit={{ opacity: 0, y: -8 }}
      transition={{ duration: 0.3, ease: 'easeOut' }}
      className="mx-auto w-full max-w-110 text-left"
    >
      <div className="overflow-hidden rounded-[28px] border border-[#e8edf5] bg-white shadow-[0_12px_48px_-20px_rgba(32,65,115,0.18)]">
        <div className="relative overflow-hidden bg-[#3478f6] px-6 pb-7 pt-6 sm:px-7">
          <div className="pointer-events-none absolute -right-12 -top-16 size-44 rounded-full bg-white/8" />
          <div className="pointer-events-none absolute -bottom-20 right-20 size-36 rounded-full bg-white/[0.07]" />

          <div className="relative flex items-center justify-between">
            <div className="flex items-center gap-2.5">
              {watch('role') == 'student' ? (
                <>
                  <div className="flex size-10 items-center justify-center rounded-[14px] bg-white text-[#3478f6] shadow-sm">
                    <Icon icon={'graduation-cap'} size={50} className="text-[33px]" />
                  </div>

                  <div>
                    <p className="text-[15px] font-bold tracking-tight text-white">Учёба</p>
                    <p className="text-md font-medium text-white/75">Личный кабинет студента</p>
                  </div>
                </>
              ) : (
                <>
                  <div className="flex size-10 items-center justify-center rounded-[14px] bg-white text-[#3478f6] shadow-sm">
                    <Icon icon={'user'} size={50} className="text-[33px]" />
                  </div>

                  <div>
                    <p className="text-[15px] font-bold tracking-tight text-white">Управление</p>
                    <p className="text-md font-medium text-white/75">
                      Личный кабинет преподавателя
                    </p>
                  </div>
                </>
              )}
            </div>
          </div>

          <div className="relative mt-7">
            <h2 className="text-[27px] font-bold leading-tight tracking-[-0.8px] text-white">
              Давайте
              <br />
              знакомиться!
            </h2>

            <p className="mt-2 max-w-72.5 text-[13px] leading-5 text-white/80">
              Заполните данные, чтобы настроить ваш профиль.
            </p>
          </div>
        </div>

        <div className="px-5 py-6 sm:px-7 sm:py-7">
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
            className="flex flex-col gap-5"
          >
            <motion.div variants={itemVariants} className="flex flex-col gap-2">
              <label htmlFor="initial-name" className="text-[13px] font-semibold text-[#26344b]">
                Как вас зовут?
              </label>

              <InputGroup className="h-12 rounded-xl border border-[#e5eaf2] bg-[#f5f7fb] px-3.5 shadow-none transition-all duration-200 hover:border-[#c8d7f5] focus-within:border-[#3984f4] focus-within:bg-white focus-within:ring-[3px] focus-within:ring-[#3984f4]/10">
                <UserRound className={iconClassName} />

                <InputGroupInput
                  {...register('name')}
                  id="initial-name"
                  autoComplete="name"
                  placeholder="Иванов Иван Иванович"
                  className="h-full border-0 bg-transparent px-2 text-[14px] text-[#26344b] shadow-none placeholder:text-[#9aa7ba] focus-visible:ring-0"
                />
              </InputGroup>

              <p className="text-[11px] leading-4 text-[#8b98ad]">
                Укажите фамилию, имя и отчество полностью.
              </p>
            </motion.div>

            <motion.div variants={itemVariants} className="flex flex-col gap-2">
              <label id="organization-label" className="text-[13px] font-semibold text-[#26344b]">
                Где вы учитесь?
              </label>

              <Controller
                control={control}
                name="organizationId"
                render={({ field }) => (
                  <Select
                    value={field.value}
                    onValueChange={(value) => {
                      field.onChange(value);

                      setValue('groupId', '', {
                        shouldDirty: true,
                        shouldValidate: true,
                      });
                    }}
                  >
                    <SelectTrigger aria-labelledby="organization-label" className={fieldClassName}>
                      <div className="flex min-w-0 items-center gap-2.5">
                        <Building2 className={iconClassName} />

                        <SelectValue placeholder="Выберите организацию" />
                      </div>
                    </SelectTrigger>

                    <SelectContent className="rounded-xl border-[#e8edf5] bg-white p-1.5 shadow-xl">
                      {organizations.map((organization) => (
                        <SelectItem
                          key={organization.id}
                          value={organization.id}
                          className="min-h-11 rounded-lg py-2 pl-3 pr-8 text-[13px] text-[#26344b] focus:bg-[#eef4ff] focus:text-[#2468e8]"
                        >
                          {organization.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                )}
              />

              {selectedOrganization && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: 'auto' }}
                  className="flex items-center gap-1.5 overflow-hidden px-0.5 text-[11px] text-[#3478f6]"
                >
                  <Check className="size-3.5 shrink-0" />
                  <span className="truncate">{selectedOrganization.name}</span>
                </motion.div>
              )}
            </motion.div>

            <motion.div variants={itemVariants} className="flex flex-col gap-2">
              <label id="group-label" className="text-[13px] font-semibold text-[#26344b]">
                Ваша группа
              </label>
              <Controller
                control={control}
                name="groupId"
                render={({ field }) => (
                  <Select
                    value={field.value}
                    onValueChange={field.onChange}
                    disabled={!organizationId}
                  >
                    <SelectTrigger
                      aria-labelledby="group-label"
                      className={`${fieldClassName} disabled:cursor-not-allowed disabled:opacity-50`}
                    >
                      <div className="flex min-w-0 items-center gap-2.5">
                        <GraduationCap className={iconClassName} />

                        <SelectValue
                          placeholder={
                            organizationId
                              ? 'Выберите учебную группу'
                              : 'Сначала выберите организацию'
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
                    </SelectContent>
                  </Select>
                )}
              />
              {selectedGroup && (
                <motion.div
                  initial={{ opacity: 0, height: 0 }}
                  animate={{ opacity: 1, height: 'auto' }}
                  className="flex items-center gap-1.5 overflow-hidden px-0.5 text-[11px] text-[#3478f6]"
                >
                  <Check className="size-3.5 shrink-0" />
                  <span>Группа {selectedGroup.name} выбрана</span>
                  <span>{max.initData}</span>
                </motion.div>
              )}
              <AnimatePresence mode="wait">
                {isSelected && (
                  <motion.div
                    className="w-full"
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    transition={{ duration: 0.3 }}
                  >
                    <Button className="w-full h-10 text-lg">Отправить заявку</Button>
                  </motion.div>
                )}
              </AnimatePresence>
            </motion.div>
          </motion.div>
        </div>
      </div>
    </motion.section>
  );
};
