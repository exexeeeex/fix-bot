import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';

type InitialForm = {
  role: string;
  name: string;
  organizationId: string;
  groupId: string;
};

export const useInitialForm = () => {
  const [isVisibleForm, setIsVisibleForm] = useState<boolean>(false);

  const schema = yup.object({
    role: yup.string().required('Выберите свою роль!'),
    name: yup.string().min(10).max(50).required('Укажите полное ФИО!'),
    organizationId: yup.string().required('Организация не выбрана!'),
    groupId: yup.string().required('Группа не выбрана!'),
  });

  const { control, watch, setValue, register } = useForm<InitialForm>({
    resolver: yupResolver(schema),

    defaultValues: {
      role: '',
      name: '',
      organizationId: '',
      groupId: '',
    },
  });

  const selectRole = (role: string) => {
    setValue('role', role, { shouldDirty: true, shouldValidate: true });
    setIsVisibleForm(true);
  };

  const isSelected =
    watch('groupId').length > 5 &&
    watch('name').length >= 10 &&
    watch('organizationId').length > 5 &&
    watch('role').length > 5;

  return {
    control,
    watch,
    register,
    setValue,
    selectRole,
    isVisibleForm,
    isSelected,
  };
};
