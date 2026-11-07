import { type CreateApplicationRequest } from '@/entities/application';
import { yupResolver } from '@hookform/resolvers/yup';
import React from 'react';
import { useForm } from 'react-hook-form';
import { applicationCreateRequestSchema } from '../config';
import { useCreateApplication } from './use-create-application';

export const useApplicationForm = (userId: string, name: string) => {
  const [isVisibleForm, setIsVisibleForm] = React.useState<boolean>(false);
  const { create, isLoading, isError, isSuccess } = useCreateApplication();

  const { control, watch, setValue, register, handleSubmit } = useForm<CreateApplicationRequest>({
    resolver: yupResolver(applicationCreateRequestSchema),

    defaultValues: {
      type: '',
      name: '',
      organizationId: '',
      groupId: '',
      userId: '',
    },
  });

  React.useEffect(() => {
    if (userId) {
      setValue('userId', userId, { shouldValidate: true });
    }
    if (name) {
      setValue('name', name, { shouldValidate: true });
    }
  }, [setValue, userId]);

  const selectRole = (role: string) => {
    setValue('type', role, { shouldDirty: true, shouldValidate: true });
    setIsVisibleForm(true);
  };

  const selectedType = watch('type');
  const isSelected =
    (selectedType === 'teacher' || (watch('groupId')?.length ?? 0) > 5) &&
    watch('name').length >= 10 &&
    watch('organizationId').length > 5 &&
    (selectedType === 'student' || selectedType === 'teacher') &&
    userId.length > 0;

  const onSubmit = (data: CreateApplicationRequest) => create(data);

  return {
    control,
    watch,
    register,
    setValue,
    handleSubmit,
    onSubmit,
    selectRole,
    isVisibleForm,
    isSelected,
    isLoading,
    isError,
    isSuccess,
  };
};
